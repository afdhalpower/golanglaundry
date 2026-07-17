# Laundry Management System

A modern, production-ready laundry management application built with Go (Fiber v3), PostgreSQL, and Alpine.js.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Backend** | Go 1.25+, Fiber v3 |
| **Database** | PostgreSQL 16 |
| **ORM** | GORM v2 |
| **Frontend** | HTML Templates (Go `html/template`), Alpine.js, Tailwind CSS |
| **Auth** | Session-based (PostgreSQL-backed) |
| **Validation** | go-playground/validator |
| **Migrations** | golang-migrate |
| **Config** | Viper (YAML + env) |
| **Logging** | slog (stdlib) |
| **Charts** | Chart.js (Dashboard) |
| **Dev Tools** | Air (hot reload) |

## Features

### Core Management
- 🔐 **Authentication** — Login, logout, role-based access (Admin, Kasir, Pegawai)
- 👥 **Customers** — Full CRUD with search, pagination, and detail page
- 🧺 **Laundry Services** — Manage service types, unit pricing, and estimates
- 📦 **Orders** — Complete order lifecycle with status tracking, quick actions, and receipt printing
- 💳 **Payments** — Multiple payment methods with receipt printing
- 📋 **Expenses** — Track operational costs by category
- 📦 **Inventory** — Stock management with low-stock alerts and movement log

### Dashboard & Analytics
- 📊 **Real-time Dashboard** — Stats cards (orders today, in progress, completed), revenue (today & monthly), expenses, total customers
- 📈 **Interactive Charts** — 7-day revenue trend (line chart) and order volume (bar chart) via Chart.js
- 📋 **Activity Log** — Real-time feed of the 10 most recent order tracking events with status icons and relative timestamps
- ⏰ **Overdue Alerts** — Automatic detection of orders past their estimated completion time, with a prominent alert card and direct links to each order
- 📈 **Reports** — Revenue, expenses, profit summary with CSV export and status filtering

### Enhanced UX
- ⚡ **Quick Actions** — Update order status directly from the orders table via dropdown (no need to open detail page), powered by Alpine.js + fetch API
- 🔍 **Global Search** — Search across all pages (customers, orders, services) with real-time filtering
- 🗑️ **Modal Confirmation** — All delete actions use a proper Alpine.js modal with blur backdrop instead of native `confirm()` dialogs
- 🔔 **Toast Notifications** — Success/error toast alerts with auto-hide and slide-in animation, accessible via Alpine store

### Admin & Settings
- 👤 **User Management** — Multi-role users with permissions management
- ⚙️ **Settings** — Store profile, working hours, receipt configuration
- 🔔 **Toast Notifications** — Real-time feedback on all successful/error actions

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Air (hot reload) — optional, for development

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

Migrations run automatically on application startup via GORM AutoMigrate.
For manual migrations:

```bash
migrate -database "postgres://golanglaundry:***@localhost:5432/golanglaundry?sslmode=disable" -path migrations up
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
├── cmd/server/main.go           # Application entry point (94 handlers)
├── internal/
│   ├── config/                  # Configuration, database connection & seeder
│   ├── handlers/                # HTTP handlers (Fiber controllers)
│   ├── services/                # Business logic layer
│   ├── repositories/            # Database operations (GORM queries)
│   ├── models/                  # GORM models & enums
│   ├── middleware/               # HTTP middleware (auth, session, logger)
│   ├── routes/                  # Route registration & template functions
│   ├── validation/              # Input validation rules
│   └── helpers/                 # Utility functions (JSON response, formatters)
├── templates/                   # HTML templates (Go html/template)
│   ├── layouts/                 # Main layout (sidebar + navbar + toast)
│   ├── partials/                # Reusable components (pagination, modals, alerts)
│   └── ...                      # Page templates per module
├── static/                      # Static assets (CSS, JS, app.js)
├── migrations/                  # SQL migration files
├── docs/                        # Documentation & implementation plans
│   └── plans/                   # Feature implementation plans
├── scripts/                     # Development scripts (dev.sh, etc.)
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
migrate -database "postgres://golanglaundry:***@localhost:5432/golanglaundry?sslmode=disable" -path migrations up

# Rollback last migration
migrate -database "$DATABASE_URL" -path migrations down 1

# Create new migration
migrate create -ext sql -dir migrations -seq create_table_name
```

> Migration also runs automatically on every startup via GORM AutoMigrate.

## Quick Actions API

The order quick status update is available via:

```bash
POST /orders/:id/quick-status
Content-Type: application/json

{"status": "dicuci"}
```

Returns JSON `{"success": true}` or error with message.

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

### Code Architecture

- **Handlers** — HTTP concerns only (parse request, call service, render response)
- **Services** — Business logic and orchestration between handlers & repositories
- **Repositories** — Database queries via GORM (raw SQL when needed)
- **Models** — Database schemas, validation rules, and type definitions

### Template Functions

The following custom functions are available in all templates:

| Function | Description |
|----------|-------------|
| `add`, `sub` | Integer arithmetic |
| `mul` | Float multiplication |
| `slice` | String slicing |
| `upper` | Uppercase first letter |
| `loop` | Generate numeric range for iteration |
| `default` | Default value if input is zero/nil |
| `nowDate` | Current date (YYYY-MM-DD) |
| `timeAgo` | Relative time in Indonesian ("2 jam lalu") |
| `statusLabel` | Order status in Indonesian ("Dicuci") |
| `statusIcon` | Emoji icon for order status (🫧, ✅, etc.) |

## Security

- Password hashing with bcrypt
- Session-based authentication with PostgreSQL-backed secure cookies
- CSRF protection on all forms
- Input validation on all endpoints
- Role-based access control (Admin, Kasir, Pegawai)
- XSS prevention via Go template auto-escaping
- SQL injection prevention via GORM parameterized queries

## License

MIT

## Author

Afdhal RZ ([@afdhalpower](https://github.com/afdhalpower))
