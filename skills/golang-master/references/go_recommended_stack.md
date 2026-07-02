# Go Recommended Tech Stack

> Đây là danh sách các library/tool được khuyến nghị sử dụng trong dự án Go.
> Agent **phải ưu tiên** dùng các library này khi tạo project mới hoặc thêm feature.
> Không tự ý thay thế bằng library khác trừ khi có lý do chính đáng và được xác nhận.

## Core Toolkit

| Mục đích | Library | Import Path |
|----------|---------|-------------|
| **Shared Utilities** | gotoolkit | `github.com/ducconit/gotoolkit` |

> **gotoolkit** là bộ toolkit nội bộ, ưu tiên dùng các utilities từ đây trước khi tìm library bên ngoài. Luôn kiểm tra gotoolkit có function/helper cần thiết không trước khi thêm dependency mới.

## Web Framework & API

| Mục đích | Library | Import Path | Ghi chú |
|----------|---------|-------------|---------|
| **HTTP Framework** | Gin | `github.com/gin-gonic/gin` | Framework chính cho REST API |
| **Validation** | validator v10 | `github.com/go-playground/validator/v10` | Validation structs, request binding |

### Gin Usage Guidelines

```go
// Router setup
r := gin.New()
r.Use(gin.Recovery())
r.Use(gin.Logger())

// Group routes
api := r.Group("/api/v1")
{
    api.GET("/users", handler.ListUsers)
    api.POST("/users", handler.CreateUser)
}

// Request binding + validation
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,gte=0,lte=150"`
}

func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // ...
}
```

## CLI

| Mục đích | Library | Import Path | Ghi chú |
|----------|---------|-------------|---------|
| **CLI Framework** | Cobra | `github.com/spf13/cobra` | Commands, flags, subcommands |

### Cobra Structure

```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "A brief description",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

// cmd/serve.go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the HTTP server",
    RunE: func(cmd *cobra.Command, args []string) error {
        return startServer()
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
    serveCmd.Flags().IntP("port", "p", 8080, "server port")
}
```

## Configuration

| Mục đích | Library | Import Path | Ghi chú |
|----------|---------|-------------|---------|
| **Config Management** | Viper | `github.com/spf13/viper` | Env, YAML, TOML, flags |

### Viper Usage

```go
func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./configs")

    // Environment variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.max_conns", 25)

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            log.Fatal(err)
        }
    }
}
```

## Database & Migration

| Mục đích | Library | Import Path | Ghi chú |
|----------|---------|-------------|---------|
| **Migration** | Goose | `github.com/pressly/goose/v3` | SQL migrations, up/down |

### Goose Usage

```bash
# Tạo migration mới
goose -dir migrations create add_users_table sql

# Chạy migrations
goose -dir migrations postgres "postgres://..." up

# Rollback
goose -dir migrations postgres "postgres://..." down
```

```go
// Embed migrations trong Go
import "github.com/pressly/goose/v3"

//go:embed migrations/*.sql
var embedMigrations embed.FS

func runMigrations(db *sql.DB) error {
    goose.SetBaseFS(embedMigrations)
    return goose.Up(db, "migrations")
}
```

## Logging

| Mục đích | Library | Import Path | Ghi chú |
|----------|---------|-------------|---------|
| **Logger** | zerolog | `github.com/rs/zerolog` | Structured JSON logging, zero alloc |

### Zerolog Usage

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func initLogger() {
    // Development: human-readable
    if os.Getenv("APP_ENV") == "development" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }

    // Set global level
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// Structured logging
log.Info().
    Str("user_id", user.ID).
    Str("action", "login").
    Msg("user logged in")

// Error logging
log.Error().
    Err(err).
    Str("request_id", reqID).
    Msg("failed to process request")

// Context-aware logging
logger := log.With().Str("request_id", reqID).Logger()
logger.Info().Msg("processing request")
```

## Summary: go.mod Dependencies

Khi tạo project Go mới, đây là danh sách dependencies chuẩn:

```
github.com/ducconit/gotoolkit       # Shared utilities
github.com/gin-gonic/gin            # HTTP framework
github.com/spf13/cobra              # CLI
github.com/spf13/viper              # Config
github.com/rs/zerolog               # Logging
github.com/pressly/goose/v3         # Migration
github.com/go-playground/validator/v10  # Validation
```

> **Lưu ý**: Không phải project nào cũng cần tất cả. Chỉ thêm khi thực sự cần.
> Ví dụ: CLI tool không cần Gin, library không cần Cobra.
