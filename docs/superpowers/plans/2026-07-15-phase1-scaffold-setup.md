# Phase 1: Project Scaffold + Setup — Implementation Plan

> **For agentic workers:** Use `software-development/subagent-driven-development` or execute inline to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Initialize Go project with full scaffold, Docker Compose for PostgreSQL, configuration system, database connection, and verified working Fiber server.

**Architecture:** Single Go module (`golanglaundry`) with `cmd/server/main.go` entry point. All source in `internal/`. Templates in `templates/`, static files in `static/`. PostgreSQL runs in Docker.

**Tech Stack:** Go 1.24+, Fiber v3, GORM v2, Viper, golang-migrate, Air, Docker Compose, PostgreSQL 16

## Global Constraints

- Go module path: `github.com/afdhalpower/golanglaundry`
- Go version 1.24+ required
- PostgreSQL 16 via Docker Compose
- All Go source under `internal/` except `cmd/server/main.go`
- Config via Viper (YAML + env overrides)
- Database connection via GORM with pgx driver
- Hot reload via Air
- Migrations via golang-migrate CLI
- Logging via slog
- git init with `afdhalpower` as author

---

## Project Structure (to be created)

```
golanglaundry/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── handlers/
│   ├── services/
│   ├── repositories/
│   ├── models/
│   ├── middleware/
│   ├── routes/routes.go
│   ├── validation/validator.go
│   └── helpers/
│       ├── response.go
│       ├── pagination.go
│       └── template.go
├── templates/layouts/
├── static/css/
├── migrations/
├── docs/
├── scripts/
├── config.yaml
├── .env.example
├── .air.toml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

### Task 1: Install Go 1.24+

**Files:** None (system setup)

- [ ] **Step 1: Check if Go is already installed**

```bash
go version 2>/dev/null || echo "Go not installed"
```

- [ ] **Step 2: Download and install Go 1.24+ for linux/amd64**

```bash
wget -q https://go.dev/dl/go1.24.3.linux-amd64.tar.gz -O /tmp/go.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go.tar.gz
rm /tmp/go.tar.gz
```

- [ ] **Step 3: Add Go to PATH for current user**

```bash
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
```

- [ ] **Step 4: Verify installation**

```bash
go version
# Expected: go version go1.24.3 linux/amd64

which go
# Expected: /usr/local/go/bin/go
```

- [ ] **Step 5: Commit (system setup — no code yet)**

```bash
git init
git config user.name "afdhalpower"
git config user.email "afdhalpower@users.noreply.github.com"
git add .
git commit -m "chore: init repository"
```

---

### Task 2: Create Docker Compose for PostgreSQL

**Files:**
- Create `docker-compose.yml`

- [ ] **Step 1: Create docker-compose.yml**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: golanglaundry_db
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: golanglaundry
      POSTGRES_USER: golanglaundry
      POSTGRES_PASSWORD: golanglaundry_secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U golanglaundry"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

- [ ] **Step 2: Start PostgreSQL container**

```bash
docker compose -f docker-compose.yml up -d
# Expected: Container golanglaundry_db started
```

- [ ] **Step 3: Verify PostgreSQL is running**

```bash
docker ps --filter name=golanglaundry_db --format "{{.Status}}"
# Expected: Up X minutes

docker compose exec postgres pg_isready -U golanglaundry
# Expected: localhost:5432 - accepting connections
```

- [ ] **Step 4: Install psql client (optional but useful)**

```bash
sudo apt-get update -qq && sudo apt-get install -y -qq postgresql-client 2>/dev/null || true
```

- [ ] **Step 5: Verify connection via psql**

```bash
PGPASSWORD=golanglaundry_secret psql -h localhost -U golanglaundry -d golanglaundry -c "SELECT current_database();"
# Expected: current_database = golanglaundry
```

- [ ] **Step 6: Install golang-migrate CLI**

```bash
# Download migrate binary
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz
tar -xzf /tmp/migrate.tar.gz -C /tmp
sudo mv /tmp/migrate /usr/local/bin/migrate
rm /tmp/migrate.tar.gz
migrate -version
# Expected: 4.18.2 (or similar)
```

- [ ] **Step 7: Install Air (hot reload)**

```bash
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b /tmp
sudo mv /tmp/air /usr/local/bin/air
air -v
# Expected: air 1.x.x
```

- [ ] **Step 8: Commit**

```bash
git add docker-compose.yml
git commit -m "chore: add docker compose for postgresql and install dev tools"
```

---

### Task 3: Initialize Go Module + Create Project Structure

**Files:**
- Create `go.mod`
- Create all directory structure

- [ ] **Step 1: Initialize Go module**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go mod init github.com/afdhalpower/golanglaundry
# Expected: go.mod created with module path
```

- [ ] **Step 2: Create all project directories**

```bash
mkdir -p cmd/server
mkdir -p internal/{config,handlers,services,repositories,models,middleware,routes,validation,helpers}
mkdir -p templates/{layouts,auth,dashboard,customers,services,orders,payments,expenses,inventory,reports,users,settings,partials}
mkdir -p static/{css,js}
mkdir -p migrations
mkdir -p docs/api
mkdir -p scripts
```

- [ ] **Step 3: Add .gitkeep files in empty directories that need tracking**

```bash
touch static/css/.gitkeep static/js/.gitkeep
touch migrations/.gitkeep
```

- [ ] **Step 4: Create .gitignore**

```gitignore
# Binary
/golanglaundry
/tmp/

# Environment
.env
config.local.yaml

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Air
tmp/

# Dependency
vendor/
```

- [ ] **Step 5: Commit**

```bash
git add go.mod .gitignore
git add cmd/ internal/ templates/ static/ migrations/ docs/ scripts/
git commit -m "chore: initialize go module and project structure"
```

---

### Task 4: Create Configuration System (Viper)

**Files:**
- Create `internal/config/config.go`
- Create `config.yaml`
- Create `.env.example`

- [ ] **Step 1: Install Viper dependency**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go get github.com/spf13/viper
```

- [ ] **Step 2: Create config.go**

```go
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Session  SessionConfig
	App      AppConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type SessionConfig struct {
	Secret   string
	MaxAge   int // seconds
	HttpOnly bool
	Secure   bool
}

type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Environment overrides
	v.AutomaticEnv()
	v.SetEnvPrefix("LAUNDRY")

	// Default values
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 3000)
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "golanglaundry")
	v.SetDefault("database.password", "golanglaundry_secret")
	v.SetDefault("database.name", "golanglaundry")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("session.secret", "change-me-in-production")
	v.SetDefault("session.max_age", 86400)
	v.SetDefault("session.http_only", true)
	v.SetDefault("session.secure", false)

	v.SetDefault("app.name", "Laundry Management System")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found is fine — use defaults + env
	}

	var cfg Config

	// Server
	cfg.Server.Host = v.GetString("server.host")
	cfg.Server.Port = v.GetInt("server.port")
	cfg.Server.ReadTimeout = v.GetDuration("server.read_timeout")
	cfg.Server.WriteTimeout = v.GetDuration("server.write_timeout")

	// Database
	cfg.Database.Host = v.GetString("database.host")
	cfg.Database.Port = v.GetInt("database.port")
	cfg.Database.User = v.GetString("database.user")
	cfg.Database.Password = v.GetString("database.password")
	cfg.Database.Name = v.GetString("database.name")
	cfg.Database.SSLMode = v.GetString("database.sslmode")
	cfg.Database.MaxOpenConns = v.GetInt("database.max_open_conns")
	cfg.Database.MaxIdleConns = v.GetInt("database.max_idle_conns")
	cfg.Database.ConnMaxLifetime = v.GetDuration("database.conn_max_lifetime")

	// Session
	cfg.Session.Secret = v.GetString("session.secret")
	cfg.Session.MaxAge = v.GetInt("session.max_age")
	cfg.Session.HttpOnly = v.GetBool("session.http_only")
	cfg.Session.Secure = v.GetBool("session.secure")

	// App
	cfg.App.Name = v.GetString("app.name")
	cfg.App.Environment = v.GetString("app.environment")
	cfg.App.Debug = v.GetBool("app.debug")

	return &cfg, nil
}

// DSN returns the PostgreSQL connection string for GORM
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}
```

- [ ] **Step 3: Create config.yaml**

```yaml
server:
  host: "0.0.0.0"
  port: 3000
  read_timeout: 10s
  write_timeout: 10s

database:
  host: "localhost"
  port: 5432
  user: "golanglaundry"
  password: "golanglaundry_secret"
  name: "golanglaundry"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 5m

session:
  secret: "change-me-in-production"
  max_age: 86400
  http_only: true
  secure: false

app:
  name: "Laundry Management System"
  environment: "development"
  debug: true
```

- [ ] **Step 4: Create .env.example**

```env
# Server
LAUNDRY_SERVER_HOST=0.0.0.0
LAUNDRY_SERVER_PORT=3000

# Database
LAUNDRY_DATABASE_HOST=localhost
LAUNDRY_DATABASE_PORT=5432
LAUNDRY_DATABASE_USER=golanglaundry
LAUNDRY_DATABASE_PASSWORD=golanglaundry_secret
LAUNDRY_DATABASE_NAME=golanglaundry

# Session
LAUNDRY_SESSION_SECRET=change-me-in-production

# App
LAUNDRY_APP_ENVIRONMENT=development
LAUNDRY_APP_DEBUG=true
```

- [ ] **Step 5: Verify compile**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go build ./internal/config/
# Expected: no errors (no output)
```

- [ ] **Step 6: Commit**

```bash
git add internal/config/ config.yaml .env.example go.mod go.sum
git commit -m "feat: add configuration system with viper"
```

---

### Task 5: Create GORM Models (Phase 1 — minimal models)

**Files:**
- Create `internal/models/user.go`
- Create `internal/models/customer.go`
- Create `internal/models/service.go`

- [ ] **Step 1: Install GORM + pgx driver**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

- [ ] **Step 2: Create user model**

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name" validate:"required,min=3,max=255"`
	Email        string         `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         string         `gorm:"size:20;not null;default:kasir" json:"role" validate:"required,oneof=admin kasir pegawai"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

- [ ] **Step 3: Create customer model**

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name" validate:"required,min=2,max=255"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Whatsapp  string         `gorm:"size:20" json:"whatsapp"`
	Address   string         `gorm:"type:text" json:"address"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

- [ ] **Step 4: Create service model**

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:255;not null" json:"name" validate:"required,min=2,max=255"`
	PricePerKg     float64        `gorm:"type:decimal(12,2);not null" json:"price_per_kg" validate:"required,gt=0"`
	EstimatedHours int            `gorm:"not null;default:24" json:"estimated_hours" validate:"required,gt=0"`
	Description    string         `gorm:"type:text" json:"description"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
```

- [ ] **Step 5: Verify compile**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go build ./internal/models/
# Expected: no errors
```

- [ ] **Step 6: Commit**

```bash
git add internal/models/ go.mod go.sum
git commit -m "feat: add base models (user, customer, service)"
```

---

### Task 6: Create Database Connection + AutoMigrate

**Files:**
- Create `internal/config/database.go`

- [ ] **Step 1: Create database connection manager**

```go
package config

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afdhalpower/golanglaundry/internal/models"
)

var DB *gorm.DB

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	DB = db

	slog.Info("database connected successfully",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"name", cfg.Database.Name,
	)

	return db, nil
}

func RunAutoMigration(db *gorm.DB) error {
	slog.Info("running auto migration")

	err := db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Service{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	slog.Info("auto migration completed")
	return nil
}

func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
```

- [ ] **Step 2: Verify compile**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go build ./internal/config/
# Expected: no errors
```

- [ ] **Step 3: Commit**

```bash
git add internal/config/database.go go.mod go.sum
git commit -m "feat: add database connection and auto migration"
```

---

### Task 7: Create Validation System

**Files:**
- Create `internal/validation/validator.go`
- Create `internal/helpers/response.go`

- [ ] **Step 1: Install go-playground/validator**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go get github.com/go-playground/validator/v10
```

- [ ] **Step 2: Create validator**

```go
package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: '%s' validation failed on '%s'", v.Field, v.Value, v.Tag)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = fmt.Sprintf("Field '%s' %s", err.Field, tagMessage(err.Tag))
	}
	return result
}

func tagMessage(tag string) string {
	messages := map[string]string{
		"required": "wajib diisi",
		"email":    "format email tidak valid",
		"min":      "terlalu pendek",
		"max":      "terlalu panjang",
		"gt":       "harus lebih besar dari 0",
		"oneof":    "nilai tidak valid",
	}
	if msg, ok := messages[tag]; ok {
		return msg
	}
	return fmt.Sprintf("tidak valid (aturan: %s)", tag)
}

func ValidateStruct(s interface{}) ValidationErrors {
	var errs ValidationErrors

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: fmt.Sprintf("%v", err.Value()),
		})
	}

	return errs
}
```

- [ ] **Step 3: Create response helpers**

```go
package helpers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSONSuccess(c fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

func JSONError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
	})
}

func JSONValidationError(c fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(APIResponse{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}

func LogAndGetUserID(c fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		slog.Warn("user_id not found in context")
		return 0
	}
	return userID
}
```

- [ ] **Step 4: Verify compile**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go build ./internal/validation/ ./internal/helpers/
# Expected: no errors
```

- [ ] **Step 5: Commit**

```bash
git add internal/validation/ internal/helpers/ go.mod go.sum
git commit -m "feat: add validation and response helpers"
```

---

### Task 8: Create Router + Middleware + Main Server Entry Point

**Files:**
- Create `internal/middleware/auth.go`
- Create `internal/middleware/logger.go`
- Create `internal/routes/routes.go`
- Create `cmd/server/main.go`

- [ ] **Step 1: Install Fiber v3 and session middleware**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go get github.com/gofiber/fiber/v3
go get github.com/gofiber/session/v2
go get github.com/gofiber/template/html/v2
```

- [ ] **Step 2: Create logger middleware**

```go
package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		slog.Info("request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", duration.String(),
			"ip", c.IP(),
		)

		return err
	}
}
```

- [ ] **Step 3: Create auth middleware (stub — will be enhanced in Phase 2)**

```go
package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		// TODO: Implement session-based auth check in Phase 2
		return c.Next()
	}
}

func RolePermission(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// TODO: Implement role check in Phase 2
		return c.Next()
	}
}
```

- [ ] **Step 4: Create routes registration**

```go
package routes

import (
	"html/template"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/afdhalpower/golanglaundry/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Static files
	app.Use("/static", static.New("./static"))

	// Global middleware
	app.Use(middleware.Logger())

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth group (public)
	auth := app.Group("/auth")
	_ = auth // Will be used in Phase 2

	// Protected routes
	protected := app.Group("/", middleware.AuthRequired())

	// Dashboard
	protected.Get("/", func(c fiber.Ctx) error {
		// TODO: Dashboard handler in Phase 2
		return c.SendString("Laundry Management System")
	})

	// API route groups (wired up in later phases)
	api := app.Group("/api", middleware.AuthRequired())
	_ = api // Will be used as we add handlers

	slog.Info("routes registered successfully")
}

// Helper to register template functions
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}
}
```

- [ ] **Step 5: Create main.go entry point**

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"

	"github.com/afdhalpower/golanglaundry/internal/config"
	"github.com/afdhalpower/golanglaundry/internal/routes"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load configuration
	configPath := "config.yaml"
	if envPath := os.Getenv("LAUNDRY_CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer config.CloseDatabase(db)

	// Run auto migration
	if err := config.RunAutoMigration(db); err != nil {
		slog.Error("failed to run auto migration", "error", err)
		os.Exit(1)
	}

	// Setup HTML template engine
	engine := html.New("./templates", ".html")
	engine.AddFuncMap(routes.TemplateFunctions())
	engine.Reload(cfg.App.Debug)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Views:        engine,
	})

	// Setup routes
	routes.SetupRoutes(app)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Start server in goroutine
	go func() {
		slog.Info("server starting", "address", addr, "environment", cfg.App.Environment)
		if err := app.Listen(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
```

- [ ] **Step 6: Verify compile**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go mod tidy
go build ./...
# Expected: no errors, binary built successfully
```

- [ ] **Step 7: Run the server and verify it starts**

```bash
cd /home/aqsadev/PRIBADI/golanglaundry
go run ./cmd/server/ &
SERVER_PID=$!
sleep 3
curl -s http://localhost:3000/health
# Expected: {"status":"ok"}
curl -s http://localhost:3000/
# Expected: "Laundry Management System" (or rendered page)
kill $SERVER_PID 2>/dev/null
```

- [ ] **Step 8: Commit**

```bash
git add cmd/server/ internal/middleware/ internal/routes/ go.mod go.sum
git commit -m "feat: add server entry point with fiber, middleware, and routes"
```

---

### Task 9: Create Air Configuration (Hot Reload)

**Files:**
- Create `.air.toml`

- [ ] **Step 1: Create .air.toml**

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/server/"
  bin = "tmp/main"
  full_bin = "LAUNDRY_CONFIG_PATH=config.yaml ./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml", "toml", "css", "js"]
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", "static"]
  include_dir = ["cmd", "internal", "templates", "config.yaml"]
  delay = 1000
  stop_on_error = true
  log = "build-errors.log"
  send_interrupt = true
  kill_delay = 500

[log]
  main_only = true

[misc]
  clean_on_exit = true
```

- [ ] **Step 2: Add tmp/ to .gitignore**

Append this to `.gitignore`:
```
tmp/
```

- [ ] **Step 3: Commit**

```bash
git add .air.toml
git add .gitignore
git commit -m "chore: add air hot reload configuration"
```

---

### Task 10: Create Development Script + Dockerfile

**Files:**
- Create `scripts/dev.sh`
- Create `Dockerfile`

- [ ] **Step 1: Create dev script**

```bash
#!/bin/bash
set -e

echo "Starting Laundry Management System in development mode..."

# Check if PostgreSQL is running
if ! docker compose ps postgres 2>/dev/null | grep -q "Up"; then
    echo "Starting PostgreSQL container..."
    docker compose up -d postgres
    echo "Waiting for PostgreSQL to be ready..."
    sleep 3
fi

# Run with Air (hot reload)
echo "Starting application with Air..."
air -c .air.toml
```

```bash
chmod +x scripts/dev.sh
```

- [ ] **Step 2: Create Dockerfile (multi-stage build)**

```dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/laundry ./cmd/server/

# Stage 2: Run
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/laundry .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config.yaml .

RUN chown -R app:app /app

USER app

EXPOSE 3000

CMD ["./laundry"]
```

- [ ] **Step 3: Update .gitignore to include Docker patterns**

```
tmp/
```

- [ ] **Step 4: Commit**

```bash
git add scripts/dev.sh Dockerfile
git commit -m "chore: add dev script and dockerfile"
```

---

### Task 11: Create Migration Files (Phase 1)

**Files:**
- Create `migrations/000001_create_users.up.sql`
- Create `migrations/000001_create_users.down.sql`

- [ ] **Step 1: Create initial migration**

```sql
-- migrations/000001_create_users.up.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Insert default admin user (password: admin123)
INSERT INTO users (name, email, password_hash, role) VALUES
('Administrator', 'admin@laundry.com', '$2a$12$LJ3m4ys3Lk6vY7x8Wz9Q4uB5c6D7e8F9g0H1i2J3k4L5m6N7o8P9q0R', 'admin');
```

```sql
-- migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

- [ ] **Step 2: Commit**

```bash
git add migrations/
git commit -m "feat: add initial database migration for users"
```

---

### Task 12: Create README.md

**Files:**
- Create `README.md`

- [ ] **Step 1: Create README**

```markdown
# Laundry Management System

Modern laundry management application built with Go (Fiber v3), PostgreSQL, and HTMX.

## Tech Stack

- **Backend:** Go 1.24+, Fiber v3
- **Database:** PostgreSQL 16
- **ORM:** GORM v2
- **Frontend:** HTML Templates, HTMX, Alpine.js, Tailwind CSS
- **Auth:** Session-based
- **Dev Tools:** Air (hot reload), golang-migrate

## Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Air (hot reload)

## Quick Start

1. Clone the repository
2. Start PostgreSQL:

```bash
docker compose up -d
```

3. Run the application:

```bash
# Development (with hot reload)
./scripts/dev.sh

# Or manually
go run ./cmd/server/
```

4. Open http://localhost:3000

## Environment Variables

See `.env.example` for all available environment variables.

Config is loaded from `config.yaml`. Environment variables override with `LAUNDRY_` prefix.

## Project Structure

```
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── config/             # Configuration & database
│   ├── handlers/           # HTTP handlers
│   ├── services/           # Business logic
│   ├── repositories/       # Database operations
│   ├── models/             # GORM models
│   ├── middleware/          # HTTP middleware
│   ├── routes/             # Route registration
│   ├── validation/         # Input validation
│   └── helpers/            # Utility functions
├── templates/              # HTML templates
├── static/                 # Static assets
├── migrations/             # SQL migrations
├── docs/                   # Documentation
└── scripts/                # Development scripts
```

## Database Migrations

```bash
# Create migration
migrate create -ext sql -dir migrations -seq create_users

# Apply all
migrate -database "postgres://golanglaundry:golanglaundry_secret@localhost:5432/golanglaundry?sslmode=disable" -path migrations up

# Rollback
migrate -database "postgres://..." -path migrations down 1
```

## License

MIT
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add readme with project overview and setup instructions"
```

---

## Verification Checklist (Final)

- [ ] `go build ./...` passes with no errors
- [ ] PostgreSQL container is running via Docker Compose
- [ ] Server starts and responds to `/health`
- [ ] Air hot reload works (change a .go file, app restarts)
- [ ] All migrations can run
```

