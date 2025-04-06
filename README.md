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

| Option | Type   | Description        | Default                  |
| ------ | ------ | ------------------ | ------------------------ |
| `-f`   | string | configuration file | `./config/settings.yaml` |

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