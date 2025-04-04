# fast-gin

A scaffolder to initialize a project based on [Gin](https://gin-gonic.com/) web framework.

## Features

- **Initialization**: Configuration, Database connection, Redis, Logging, Misc.
- **Command line**: DB initialization, Data import and export.
- **Routes**: Static, Grouping.
- **Middleware**: Authentication, Rate limit.
- **JWT**: Login, Logout.
- **Common**: File upload, Captcha, List query.
- **Deployment**: Dockerfile, docker-compose.

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

- Use `struct` which provides compile-time checking and avoids runtime assertions over `map`.
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

## Logging