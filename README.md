# fast-gin

A scaffolder to initialize a project based on [Gin](https://gin-gonic.com/) web framework.

## Features

- **Initialization**: Configuration, Database connection, Redis, Logging, Misc.
- **Command line**: DB initialization, Data import and export.
- **Routes**: Static, Grouping
- **Middleware**: Authentication, Rate limit.
- **JWT**: Login, Logout
- **Common**: File upload, Captcha, List query
- **Deployment**: Dockerfile, docker-compose

## Configuration

In a program, some values that don’t change frequently are typically stored in a configuration file. For example, things like the database address, username and password, JWT expiration time, file upload paths, and so on.

If you don’t use a configuration file and want to change a specific configuration, you’d have to recompile your program.

Configuration files often use [`YAML`](https://yaml.org/) format. You could also use `TOML`, `INI`, or `JSON`, just parse the corresponding file with Go. However, `YAML` is a bit more flexible and easier to use, plus it supports comments.

Library:

```bash
go get gopkg.in/yaml.v3
```

Parse yaml file.

```go
package main

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

func main() {
    var data map[string]any
    
    byteData, _ := os.ReadFile("/path/to/settings.yaml") // Ignoring error for brevity
    err := yaml.Unmarshal(byteData, &data)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(data)
}
```

Improvements:

- Use `struct` which provides compile-time type checking and avoids runtime assertions over `map`.
- Use `os.Open` with a decoder for streaming or larger files, avoiding the need to load everything into memory at once.
- Encapsulate parsing into a reusable function for better modularity.
- Handle missing files or fields with defaults.

```yaml
# Example configuration
database:
  host:     localhost
  port:     5432
  user:     admin
  password: admin
```

```go
package main

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

func LoadConfig(filename string) (Config, error) {
    config := Config{
        Database: struct {
            Host string `yaml:"host"`
            Port int    `yaml:"port"`
        }{Host: "localhost", Port: 5432}, // Default
        JWTExpiry: "24h", // Default value
    }

    file, err := os.Open(filename)
    if err != nil {
        return config, nil // Return defaults on error
    }
    defer file.Close()

    decoder := yaml.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return config, fmt.Errorf("decoding YAML: %w", err)
    }
    return config, nil
}

func main() {
    config, err := LoadConfig("settings.yaml")
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Printf("Config: %+v\n", config)
}
```

Read the configuration file from a command-line flag instead of hardcoding it.

```go
func LoadConfig() (cfg *config.Config, err error) {  
    cfg = new(config.Config)  
    file, err := os.Open(flags.Options.File)
    ...
}
```

😖 If the configuration file is modified, the program needs to be restarted to retrieve the new values.

😕 Is there a way to dynamically modify the configuration without restarting the container?

💡

1. Store configuration directly in memory (small or medium project).
2. Access via APIs of configuration management system, such as etcd.

```go
func DumpConfig() error {  
    byteData, err := yaml.Marshal(global.Config)  
    if err != nil {  
       return fmt.Errorf("error when dumping configuration: %w", err)  
    }  
    err = os.WriteFile(flags.Options.File, byteData, 0666)  
    if err != nil {  
       return fmt.Errorf("error when dumping configuration: %w", err)  
    }  
    fmt.Println("Configuration dumped successfully")  
    return nil  
}
```

## Flags

| Option | Type     | Description               | Default                  |
| ------ | -------- | ------------------------- | ------------------------ |
| `-f`   | `string` | Configuration file        | `./config/settings.yaml` |
| `-v`   | `bool`   | Print version information | `false`                  |
| `-db`  | `bool`   | Database migration        | `false`                  |

## Logging

Logging is a very important aspect. It is highly recommended that everyone logs extensively when working on projects.

- Where the log is printed?
- When the log is printed?
- Log segmentation: by time, by size?
- Log level?
- ⚠ Fatal when error occurs in loading configurations. Do not fatal during runtime.

This also brings up the question of whether backend errors should be returned to the frontend.

- If it’s for internal company use, just return the errors directly. That way, when an error occurs later, you can immediately know the reason and fix it.
- But if it’s for external use, directly returning backend errors makes your product seem unprofessional. It’s better to standardize the responses, such as "network error" or "system error," and then display the specific error details in the logs.

To choose:

- **[zap](https://github.com/uber-go/zap)** for new projects, especially those requiring high performance and active maintenance, given its benchmarks and community backing.
- **[logrus](https://github.com/sirupsen/logrus)** for existing projects where its feature set is already integrated, but be aware of its maintenance mode and potential need for migration in the future.

Here we choose logrus.

```go
go get github.com/sirupsen/logrus
```

### Format

Implement `Format(entry *logrus.Entry) ([]byte, error)`.

```go
type MyLog struct {}

func (MyLog) Format(entry *logrus.Entry) ([]byte, error) {  
    // Color  
    var color int  
    switch entry.Level {  
    case logrus.DebugLevel, logrus.TraceLevel:  
       color = gray  
    case logrus.WarnLevel:  
       color = yellow  
    case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:  
       color = red  
    default:  
       color = blue  
    }  
  
    // Buffer is required for formatting log messages before outputting them.  
    var buf *bytes.Buffer  
    if entry.Buffer != nil {  
       buf = entry.Buffer  
    } else {  
       buf = &bytes.Buffer{}  
    }  
    
    // Time format  
    timeFormat := entry.Time.Format("2006-01-02T15:04:05Z0700")  
  
    if entry.HasCaller() {  
       // Custom file path and line  
       funcVal := entry.Caller.Function  
       fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)  
       // Custom format  
       _, err := fmt.Fprintf(buf, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timeFormat, color, entry.Level, fileVal, funcVal, entry.Message)  
       if err != nil {  
          return nil, err  
       }  
    }    return buf.Bytes(), nil  
}
```

```go
func InitLogger() {  
    logrus.SetLevel(logrus.DebugLevel)  
    logrus.SetReportCaller(true)  
    logrus.SetFormatter(MyLog{})  
    //logrus.SetFormatter(&logrus.JSONFormatter{})  // To external
}
```

### Hook

Hooks are called whenever a log entry is created. 

```go
type MyHook struct {  
    file     *os.File   // Log file  
    errFile  *os.File   // Error log file  
    fileDate string     // Date of log file  
    logPath  string     // Path of log file  
    mu       sync.Mutex // Mutex lock  
}
```

```go
func InitLogger() {  
    ...
    logrus.AddHook(&MyHook{  
       logPath: "logs",  
    })  
}
```

```go
func (hook *MyHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

    date := entry.Time.Format("2006-01-02")  
    if hook.fileDate != date {  
       // Rotate if day is passed  
       if err := hook.rotate(date); err != nil {  
          return err  
       }  
    }  
    
    // Dump logs to file  
    entryStr, err := entry.String()  
    if err != nil {  
       return fmt.Errorf("failed to get log entry: %v", err)  
    }  
    if _, err := hook.file.Write([]byte(entryStr)); err != nil {  
       return fmt.Errorf("failed to write to log file: %v", err)  
    }  
  
    // Dump error logs to file  
    if entry.Level <= logrus.ErrorLevel {  
       if _, err := hook.errFile.Write([]byte(entryStr)); err != nil {  
          return fmt.Errorf("failed to write to error log file: %v", err)  
       }  
    }  
    return nil  
}
```

```go
func (hook *MyHook) rotate(date string) error {  
    if hook.file != nil {  
       // Close the old one  
       if err := hook.file.Close(); err != nil {  
          return fmt.Errorf("failed to close the old log file when rotation: %v", err)  
       }  
    }    if hook.errFile != nil {  
       // Close the old one  
       if err := hook.errFile.Close(); err != nil {  
          return fmt.Errorf("failed to close the old error log file when rotation: %v", err)  
       }  
    }  
    
    // Log file directory  
    dir := fmt.Sprintf("%s/%s", hook.logPath, date)  
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {  
       return fmt.Errorf("failed to create log directory: %v", err)  
    }  
  
    infoLog := fmt.Sprintf("%s/info.log", dir)  
    errLog := fmt.Sprintf("%s/err.log", dir)  
  
    // Create new log files  
    var err error  
    hook.file, err = os.OpenFile(infoLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)  
    if err != nil {  
       return fmt.Errorf("failed to open log file: %v", err)  
    }  
    hook.errFile, err = os.OpenFile(errLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)  
    if err != nil {  
       return fmt.Errorf("failed to open error log file: %v", err)  
    }  
  
    // Update file date  
    hook.fileDate = date  
    return nil  
}
```

## [GORM](https://gorm.io/docs/connecting_to_the_database.html)

GORM officially supports the databases MySQL, PostgreSQL, SQLite, SQL Server, and TiDB.

This scaffolder supports MySQL, PostgreSQL, SQLite.

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
go get github.com/glebarez/sqlite
```

⚠ For SQLite in CGO.

```bash
go get gorm.io/driver/sqlite
```

```yaml
db:  
    mode: mysql # Supports: mysql pgsql sqlite  
    db_name:  
    host:   
    port: 3306  
    user:  
    password:
```

Use simple factory pattern to initialize `gorm.DB` given `mode`.

```go
type DB struct {  
    Mode     DBMode `yaml:"mode"` // Supports: mysql pgsql sqlite  
    DBName   string `yaml:"db_name"`  
    Host     string `yaml:"host"`  
    Port     int    `yaml:"port"`  
    User     string `yaml:"user"`  
    Password string `yaml:"password"`  
}  
  
func (db DB) GetDSN() gorm.Dialector {  
    switch db.Mode {  
    case MYSQL:  
       dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",  
          db.User,  
          db.Password,  
          db.Host,  
          db.Port,  
          db.DBName,  
       )  
       return mysql.Open(dsn)  
    case PG:  
       dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",  
          db.User,  
          db.Password,  
          db.Host,  
          db.Port,  
          db.DBName,  
       )  
       return postgres.Open(dsn)  
    case SQLITE:  
       return sqlite.Open(db.DBName)  
    case "":  
       logrus.Warnf("Database mode not specified")  
       return nil  
    default:  
       logrus.Fatalf("Database is not supported")  
       return nil  
    }  
}
```

```go
func InitGorm() (db *gorm.DB) {  
    cfg := global.Config.DB  
  
    dialector := cfg.GetDSN()  
    if dialector == nil {  
       return  
    }  
  
    // Open initialize db session based on dialector  
    database, err := gorm.Open(dialector, &gorm.Config{  
       DisableForeignKeyConstraintWhenMigrating: true,  
    })  
    if err != nil {  
       logrus.Fatalf("Failed to connect to database: %v", err)  
    }  
  
    // Get DB connection pool  
    sqlDB, err := database.DB()  
    if err != nil {  
       logrus.Fatalf("Failed to get database connection pool: %s", err)  
       return  
    }  
    err = sqlDB.Ping()  
    if err != nil {  
       logrus.Fatalf("Failed to probe database connection pool liveness: %s", err)  
       return  
    }  
  
    // Configure DB connection pool  
    // TODO: Add to configuration file  
    sqlDB.SetMaxIdleConns(10)  
    sqlDB.SetMaxOpenConns(100)  
    sqlDB.SetConnMaxLifetime(time.Hour)  
  
    logrus.Infof("DB initialized successfully")  
    return  
}
```

## Redis

```bash
go get github.com/redis/go-redis/v9
```

```yaml
redis:  
    addr: "127.0.0.1:6379"  
    password: ""  
    db: 1
```

```go
type Redis struct {  
    Addr     string `yaml:"addr"`  
    Password string `yaml:"password"`  
    DB       int    `yaml:"db"`  
}
```

```go
func InitRedis() *redis.Client {  
    cfg := global.Config  
    rdb := redis.NewClient(&redis.Options{  
       Addr:     cfg.Redis.Addr,  
       Password: cfg.Redis.Password,  
       DB:       cfg.Redis.DB,  
    })  
  
    _, err := rdb.Ping(context.Background()).Result()  
    if err != nil {  
       logrus.Errorf("Failed to connect to redis: %s", err)  
       return nil  
    }  
    logrus.Infof("Connect to redis successfully")  
    return rdb  
}
```

```go
func main() {  
    ... 
    // Redis  
    global.Redis = core.InitRedis()  
  
}
```

## Database migration

```go
type Model struct {  
    ID        uint `gorm:"primaryKey"`  
    CreatedAt time.Time  
    UpdatedAt time.Time  
}
```

```go
type UserModel struct {  
    Model           // Base  
    Username string `gorm:"size:16" json:"username"`  
    Nickname string `gorm:"size:32" json:"nickname"`  
    Password string `gorm:"size:64" json:"password"`  
    RoleID   int8   `json:"roleID"` // 1: admin, 2: normal  
  
    // TODO: Email, Phone, UUID, OpenID...  
}
```

```go
func MigrateDB() {  
    err := global.DB.AutoMigrate(&models.UserModel{})  
    if err != nil {  
       logrus.Errorf("Failed to migrate database: %s", err)  
       return  
    }  
    logrus.Infof("Migrate database successfully")  
}
```

## User

```go
type User struct {}
```

Create a user.

```go
func (User) Create() {  
    var user models.UserModel  
  
    // Role  
    fmt.Println("Please select a role for user (1 (admin) 2 (normal)): ")  
    _, err := fmt.Scanln(&user.RoleID)  
    if err != nil {  
       fmt.Println("Input error:", err)  
       return  
    }  
    if user.RoleID != 1 && user.RoleID != 2 {  
       fmt.Println("Role err:", err)  
       return  
    }  
  
    // Username  
    for {  
       fmt.Println("Please input username: ")  
       _, err = fmt.Scanln(&user.Username)  
       if err != nil {  
          fmt.Println("Input error:", err)  
          return  
       }  
       var u models.UserModel  
       err = global.DB.Take(&u, "username = ?", user.Username).Error  
       if err == nil {  
          fmt.Println("User already exists")  
          continue  
       }  
       break  
    }  
  
    // Password  
    fmt.Println("Please input password: ")  
    password, err := terminal.ReadPassword(int(os.Stdin.Fd()))  
    if err != nil {  
       fmt.Println("Failed to read password:", err)  
       return  
    }  
    fmt.Println("Please input password again: ")  
    rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))  
    if err != nil {  
       fmt.Println("Failed to read password:", err)  
       return  
    }  
    if string(password) != string(rePassword) {  
       fmt.Println("Password mismatched")  
       return  
    }  
  
    // Persist  
    encryptedPassword, err := pwd.Encrypt(string(password))  
    if err != nil {  
       fmt.Println("Failed to encrypt password:", err)  
    }  
    err = global.DB.Create(&models.UserModel{  
       Username: user.Username,  
       Password: encryptedPassword,  
       RoleID:   user.RoleID,  
    }).Error  
    if err != nil {  
       logrus.Errorf("Failed to create user: %s", err)  
       return  
    }  
    logrus.Infof("Create user [%s] successfully", user.Username)  
}
```

```bash
go run main.go -res user -op create
```

Encrypt password by `bcrypt`.

```bash
go get golang.org/x/crypto/bcrypt
```

```go
func Encrypt(password string) (string, error) {  
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)  
    if err != nil {  
       logrus.Errorf("Failed to encrypt password: %s", err)  
       return "", err  
    }  
    return string(hashedPassword), nil  
}
```

List users.

```go
func (User) List() {  
    var userList []models.UserModel  
    global.DB.Order("created_at desc").Limit(10).Find(&userList)  
    for _, model := range userList {  
       fmt.Printf("UserID: %d  Username: %s Nickname: %s Role: %d CreatedAt: %s\n",  
          model.ID,  
          model.Username,  
          model.Nickname,  
          model.RoleID,  
          model.CreatedAt.Format("2006-01-02 15:04:05"),  
       )  
    }  
}
```

```bash
go run main.go -res user -op list
```

Remove a user.

```go
func (User) Remove() {  
    var username string  
  
    // Username  
    for {  
       fmt.Println("Please input username of user to be deleted: ")  
       _, err := fmt.Scanln(&username)  
       if err != nil {  
          fmt.Println("Input error:", err)  
          return  
       }  
       var u models.UserModel  
       err = global.DB.Take(&u, "username = ?", username).Error  
       if err != nil {  
          fmt.Println("User does not exist")  
          continue  
       }  
       break  
    }  
  
    err := global.DB.  
       Where("username = ?", username).  
       Delete(&models.UserModel{}).Error  
    if err != nil {  
       logrus.Errorf("Failed to delete user: %s", err)  
       return  
    }  
    logrus.Infof("Delete user [%s] successfully", username)  
}
```

```bash
go run main.go -res user -op remove
```