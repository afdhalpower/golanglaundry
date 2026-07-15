# Laundry Management System

A modern, production-ready laundry management application built with Go (Fiber v3), PostgreSQL, and HTMX.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Backend** | Go 1.25+, Fiber v3 |
| **Database** | PostgreSQL 16 |
| **ORM** | GORM v2 |
| **Frontend** | HTML Templates, HTMX, Alpine.js, Tailwind CSS |
| **Auth** | Session-based |
| **Validation** | go-playground/validator |
| **Migrations** | golang-migrate |
| **Config** | Viper (YAML + env) |
| **Logging** | slog (stdlib) |
| **Dev Tools** | Air (hot reload) |

## Features

- 🔐 **Authentication** — Login, logout, role-based (Admin, Kasir, Pegawai)
- 📊 **Dashboard** — Real-time stats, charts, recent activity
- 👥 **Customers** — CRUD with search & pagination
- 🧺 **Laundry Services** — Manage service types & pricing
- 📦 **Orders** — Full order lifecycle with status tracking
- 💳 **Payments** — Multiple payment methods with receipt printing
- 📋 **Expenses** — Track operational costs
- 📦 **Inventory** — Stock management with low-stock alerts
- 📈 **Reports** — Revenue, expenses, profit with PDF/Excel/CSV export
- 👤 **User Management** — Multi-role users with permissions
- ⚙️ **Settings** — Store profile, working hours, receipt config

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Air (hot reload) — optional for development

## Quick Start

### 1. Clone & Setup

```bash
git clone <repo-url> golanglaundry
cd golanglaundry
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

### 3. Run Database Migrations

```bash
migrate -database "postgres://golanglaundry:golanglaundry_secret@localhost:5432/golanglaundry?sslmode=disable" -path migrations up
```

### 4. Start the Application

```bash
# Development (with hot reload)
./scripts/dev.sh

# Or manually
go run ./cmd/server/
```

### 5. Open Browser

Navigate to [http://localhost:3000](http://localhost:3000)

Default admin credentials:
- **Email:** admin@laundry.com
- **Password:** admin123

## Project Structure

```
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── config/                  # Configuration & database connection
│   ├── handlers/                # HTTP handlers (Fiber controllers)
│   ├── services/                # Business logic layer
│   ├── repositories/            # Database operations (GORM)
│   ├── models/                  # GORM models
│   ├── middleware/               # HTTP middleware (auth, logger)
│   ├── routes/                  # Route registration
│   ├── validation/              # Input validation
│   └── helpers/                 # Utility functions
├── templates/                   # HTML templates
│   ├── layouts/                 # Main layout (sidebar + navbar)
│   ├── partials/                # Reusable components
│   └── ...                      # Page templates per module
├── static/                      # Static assets (CSS, JS)
├── migrations/                  # SQL migration files
├── docs/                        # Documentation
├── scripts/                     # Development scripts
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # PostgreSQL container
├── config.yaml                  # Application configuration
├── .env.example                 # Environment variables template
├── .air.toml                    # Air hot reload configuration
└── README.md
```

## Environment Variables

See `.env.example` for all available variables.

Config is loaded from `config.yaml`. Environment variables override with `LAUNDRY_` prefix.

| Variable | Default | Description |
|----------|---------|-------------|
| `LAUNDRY_SERVER_PORT` | `3000` | HTTP server port |
| `LAUNDRY_DATABASE_HOST` | `localhost` | PostgreSQL host |
| `LAUNDRY_DATABASE_PORT` | `5432` | PostgreSQL port |
| `LAUNDRY_DATABASE_USER` | `golanglaundry` | Database user |
| `LAUNDRY_DATABASE_PASSWORD` | `golanglaundry_secret` | Database password |
| `LAUNDRY_DATABASE_NAME` | `golanglaundry` | Database name |
| `LAUNDRY_SESSION_SECRET` | `change-me-in-production` | Session encryption key |
| `LAUNDRY_APP_ENVIRONMENT` | `development` | App environment |
| `LAUNDRY_APP_DEBUG` | `true` | Debug mode |

## Database Migrations

```bash
# Apply all migrations
migrate -database "postgres://golanglaundry:golanglaundry_secret@localhost:5432/golanglaundry?sslmode=disable" -path migrations up

# Rollback last migration
migrate -database "$DATABASE_URL" -path migrations down 1

# Create new migration
migrate create -ext sql -dir migrations -seq create_table_name
```

## Docker Deployment

### Build & Run

```bash
# Build the Docker image
docker build -t golanglaundry .

# Run with PostgreSQL
docker compose up -d
```

## Development

### Hot Reload with Air

Air watches for file changes and automatically rebuilds/restarts the server.

```bash
./scripts/dev.sh
```

### Code Structure

- **Handlers** handle HTTP concerns only (parse request, call service, render response)
- **Services** contain business logic and orchestration
- **Repositories** handle database queries via GORM
- **Models** define database schemas and validation rules

## API Documentation

API documentation is available via Swagger at `/docs/api/` when running in development mode.

## Security

- Password hashing with bcrypt
- Session-based authentication with secure cookies
- CSRF protection on all forms
- Input validation on all endpoints
- Role-based access control (Admin, Kasir, Pegawai)
- XSS prevention via Go template auto-escaping
- SQL injection prevention via GORM parameterized queries

## License

MIT

## Author

Afdhal RZ (@afdhalpower)
