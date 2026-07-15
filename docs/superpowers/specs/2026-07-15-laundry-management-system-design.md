# Laundry Management System — Design Spec

> **Supersedes:** N/A (new project)
> **Status:** Approved
> **Date:** 2026-07-15
> **Author:** Afdhal RZ

## 1. Project Overview

Modern, production-ready Laundry Management System built with Go (Fiber v3) targeting small-to-medium laundry businesses. Single-binary deployment with embedded HTML templates.

## 2. Architecture

**Pattern:** Monolith Modular (Clean Architecture layers within a single binary)

```
┌─────────────────────────────────────────────────┐
│                  HTTP Router (Fiber)             │
├─────────────────────────────────────────────────┤
│  Middleware (Session, CSRF, Auth, Permission)    │
├─────────────────────────────────────────────────┤
│  Handlers  ←→  Services  ←→  Repositories       │
│                      ↕                           │
│                   GORM ORM                       │
│                      ↕                           │
│                PostgreSQL (Docker)               │
├─────────────────────────────────────────────────┤
│  Templates (Go html/template + HTMX + Alpine)    │
│  Static (Tailwind CSS output)                    │
└─────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility | Notes |
|-------|---------------|-------|
| **Router** | Route registration, middleware binding | Fiber v3 |
| **Middleware** | Session auth, CSRF, permission gates, logging | |
| **Handlers** | HTTP request parsing, response formatting, template rendering | No business logic |
| **Services** | Business rules, validation orchestration, transaction boundaries | |
| **Repositories** | GORM queries, DB operations, soft-delete filters | One file per model |
| **Models** | GORM structs, DB table mapping, validation tags | |
| **Templates** | HTML with HTMX for dynamic UI, Alpine.js for interactions | |
| **Config** | Viper-based env/config file loading | |

## 3. Tech Stack

| Component | Choice | Justification |
|-----------|--------|---------------|
| Language | Go 1.24+ | Performance, single binary, great concurrency |
| HTTP Framework | Fiber v3 | Fast, Express-like, built-in middleware |
| ORM | GORM v2 | Mature, auto-migration, soft-delete, hooks |
| Database | PostgreSQL 16 | Transactions, JSON, indexing, production-grade |
| Templates | Go html/template | No extra dependency, secure (auto-escape) |
| HTMX | htmx 2.x | Minimal JS, server-driven UI |
| Alpine.js | Alpine 3.x | Lightweight interactions (modals, toggles) |
| Tailwind CSS | Tailwind 3.x (CDN for dev, standalone CLI for prod) | Utility-first, consistent design |
| Auth | Session-based (gofiber/session + gorilla/sessions) | Simple, secure, good for kasir/employee |
| Validation | go-playground/validator v10 | Struct tags, i18n-ready |
| Migrations | golang-migrate/migrate v4 | Versioned, up/down, CI-friendly |
| Config | Viper | YAML, env, defaults |
| Logging | slog (stdlib) | Structured, zero-dependency |
| API Docs | Swagger (swaggo) | OpenAPI 3.0 |
| Testing | Testify | Assertions, mocking, suite |
| Hot Reload | Air | File watcher, auto-restart |
| Container | Docker Compose | PostgreSQL + app |

## 4. Database Schema

### Entity Relationship (Core)

```
users ──1:N── orders ──N:1── customers
              │
              ├──1:N── order_services
              │
              ├──1:1── payments
              │
              └──1:N── order_tracking

expenses ──N:1── expense_categories

inventory ──N:1── inventory_categories
```

### Tables

**users** — Authentication & role management
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| name | VARCHAR(255) | |
| email | VARCHAR(255) UNIQUE | Login credential |
| password_hash | VARCHAR(255) | bcrypt |
| role | VARCHAR(20) | 'admin', 'kasir', 'pegawai' |
| is_active | BOOLEAN | DEFAULT true |
| created_at, updated_at, deleted_at | TIMESTAMPTZ | GORM soft-delete |

**customers** — Pelanggan
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| name | VARCHAR(255) NOT NULL | |
| phone | VARCHAR(20) | |
| whatsapp | VARCHAR(20) | |
| address | TEXT | |
| notes | TEXT | |
| created_at, updated_at, deleted_at | TIMESTAMPTZ | |

**services** — Layanan Laundry
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| name | VARCHAR(255) NOT NULL | Cuci, Setrika, dll |
| price_per_kg | DECIMAL(12,2) | |
| estimated_hours | INT | |
| description | TEXT | |
| is_active | BOOLEAN | DEFAULT true |
| created_at, updated_at, deleted_at | TIMESTAMPTZ | |

**orders** — Pesanan Laundry
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| order_number | VARCHAR(50) UNIQUE | Auto-generated (INV/YYYYMMDD/XXXX) |
| customer_id | BIGINT FK | |
| user_id | BIGINT FK | Created by |
| weight_kg | DECIMAL(10,2) | |
| price_per_kg | DECIMAL(12,2) | Snapshot from service |
| discount | DECIMAL(12,2) | DEFAULT 0 |
| extra_cost | DECIMAL(12,2) | DEFAULT 0 |
| total | DECIMAL(12,2) | Computed |
| entry_date | DATE | |
| estimated_done_date | DATE | |
| status | VARCHAR(20) | 'menunggu', 'dicuci', 'dikeringkan', 'disetrika', 'siap_diambil', 'sudah_diambil', 'dibatalkan' |
| notes | TEXT | |
| created_at, updated_at, deleted_at | TIMESTAMPTZ | |

**order_tracking** — Riwayat status pesanan
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| order_id | BIGINT FK | |
| status | VARCHAR(20) | |
| note | TEXT | |
| created_by | BIGINT FK users | |
| created_at | TIMESTAMPTZ | |

**payments** — Pembayaran
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| order_id | BIGINT FK UNIQUE | One payment per order |
| amount | DECIMAL(12,2) | |
| method | VARCHAR(20) | 'tunai', 'qris', 'transfer' |
| status | VARCHAR(20) | 'lunas', 'belum_lunas' |
| payment_date | TIMESTAMPTZ | |
| note | TEXT | |
| created_by | BIGINT FK users | |
| created_at, updated_at | TIMESTAMPTZ | |

**expense_categories** — Kategori Pengeluaran
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| name | VARCHAR(255) | Sabun, Pewangi, dll |
| description | TEXT | |
| created_at, updated_at | TIMESTAMPTZ | |

**expenses** — Pengeluaran
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| expense_category_id | BIGINT FK | |
| amount | DECIMAL(12,2) | |
| description | TEXT | |
| date | DATE | |
| created_by | BIGINT FK users | |
| created_at, updated_at | TIMESTAMPTZ | |

**inventory_items** — Stok Inventaris
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| name | VARCHAR(255) | |
| category | VARCHAR(50) | 'sabun', 'pewangi', 'plastik', 'hanger', 'nota', 'lainnya' |
| quantity | INT | |
| min_stock | INT | Alert threshold |
| unit | VARCHAR(50) | 'pcs', 'liter', 'kg', 'pack' |
| notes | TEXT | |
| created_at, updated_at, deleted_at | TIMESTAMPTZ | |

**settings** — Pengaturan
| Column | Type | Notes |
|--------|------|-------|
| id | BIGSERIAL PK | |
| key | VARCHAR(255) UNIQUE | |
| value | TEXT | |
| updated_at | TIMESTAMPTZ | |

## 5. Route Design

```
# Public
GET  /login                  -> AuthHandler.Login
POST /login                  -> AuthHandler.DoLogin
GET  /logout                 -> AuthHandler.Logout

# Dashboard
GET  /                       -> DashboardHandler.Index
GET  /dashboard              -> DashboardHandler.Index

# Customers
GET  /customers              -> CustomerHandler.Index
GET  /customers/new          -> CustomerHandler.New
POST /customers              -> CustomerHandler.Create
GET  /customers/:id          -> CustomerHandler.Show
GET  /customers/:id/edit     -> CustomerHandler.Edit
PUT  /customers/:id          -> CustomerHandler.Update
DELETE /customers/:id        -> CustomerHandler.Delete

# Services (Layanan)
GET  /services               -> ServiceHandler.Index
GET  /services/new           -> ServiceHandler.New
POST /services               -> ServiceHandler.Create
GET  /services/:id/edit      -> ServiceHandler.Edit
PUT  /services/:id           -> ServiceHandler.Update
DELETE /services/:id         -> ServiceHandler.Delete

# Orders
GET  /orders                 -> OrderHandler.Index
GET  /orders/new             -> OrderHandler.New
POST /orders                 -> OrderHandler.Create
GET  /orders/:id             -> OrderHandler.Show
GET  /orders/:id/edit        -> OrderHandler.Edit
PUT  /orders/:id             -> OrderHandler.Update
POST /orders/:id/status      -> OrderHandler.UpdateStatus
DELETE /orders/:id           -> OrderHandler.Delete

# Payments
GET  /orders/:id/payment     -> PaymentHandler.Create
POST /orders/:id/payment     -> PaymentHandler.Store
GET  /payments               -> PaymentHandler.Index

# Expenses
GET  /expenses               -> ExpenseHandler.Index
GET  /expenses/new           -> ExpenseHandler.New
POST /expenses               -> ExpenseHandler.Create
GET  /expenses/:id/edit      -> ExpenseHandler.Edit
PUT  /expenses/:id           -> ExpenseHandler.Update
DELETE /expenses/:id         -> ExpenseHandler.Delete

# Inventory
GET  /inventory              -> InventoryHandler.Index
GET  /inventory/new          -> InventoryHandler.New
POST /inventory              -> InventoryHandler.Create
GET  /inventory/:id/edit     -> InventoryHandler.Edit
PUT  /inventory/:id          -> InventoryHandler.Update
DELETE /inventory/:id        -> InventoryHandler.Delete

# Reports
GET  /reports                -> ReportHandler.Index
GET  /reports/revenue        -> ReportHandler.Revenue
GET  /reports/expenses       -> ReportHandler.Expenses
GET  /reports/profit         -> ReportHandler.Profit
GET  /reports/export         -> ReportHandler.Export

# Users (Admin only)
GET  /users                  -> UserHandler.Index
GET  /users/new              -> UserHandler.New
POST /users                  -> UserHandler.Create
GET  /users/:id/edit         -> UserHandler.Edit
PUT  /users/:id              -> UserHandler.Update
DELETE /users/:id            -> UserHandler.Delete

# Settings
GET  /settings               -> SettingHandler.Index
PUT  /settings               -> SettingHandler.Update
GET  /profile                -> ProfileHandler.Index
PUT  /profile                -> ProfileHandler.Update
PUT  /profile/password       -> ProfileHandler.ChangePassword
```

## 6. Project Structure

```
golanglaundry/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Viper config loader
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── dashboard_handler.go
│   │   ├── customer_handler.go
│   │   ├── service_handler.go
│   │   ├── order_handler.go
│   │   ├── payment_handler.go
│   │   ├── expense_handler.go
│   │   ├── inventory_handler.go
│   │   ├── report_handler.go
│   │   ├── user_handler.go
│   │   ├── setting_handler.go
│   │   └── profile_handler.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── customer_service.go
│   │   ├── service_service.go
│   │   ├── order_service.go
│   │   ├── payment_service.go
│   │   ├── expense_service.go
│   │   ├── inventory_service.go
│   │   ├── report_service.go
│   │   └── user_service.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── customer_repository.go
│   │   ├── service_repository.go
│   │   ├── order_repository.go
│   │   ├── payment_repository.go
│   │   ├── expense_repository.go
│   │   ├── inventory_repository.go
│   │   └── setting_repository.go
│   ├── models/
│   │   ├── user.go
│   │   ├── customer.go
│   │   ├── service.go
│   │   ├── order.go
│   │   ├── order_tracking.go
│   │   ├── payment.go
│   │   ├── expense.go
│   │   ├── inventory_item.go
│   │   └── setting.go
│   ├── middleware/
│   │   ├── auth.go              # Session auth middleware
│   │   ├── permission.go        # Role-based access
│   │   ├── csrf.go              # CSRF protection
│   │   └── logger.go            # Request logging
│   ├── routes/
│   │   └── routes.go            # All route registration
│   ├── validation/
│   │   └── validator.go         # Validator instance + custom rules
│   └── helpers/
│       ├── response.go          # JSON/HTML response helpers
│       ├── pagination.go        # Pagination helper
│       ├── number.go            # Number formatting
│       └── template.go          # Template render helpers
├── templates/
│   ├── layouts/
│   │   └── main.html            # Main layout (sidebar, navbar)
│   ├── auth/
│   │   └── login.html
│   ├── dashboard/
│   │   └── index.html
│   ├── customers/
│   │   ├── index.html
│   │   ├── form.html
│   │   └── show.html
│   ├── services/
│   │   ├── index.html
│   │   └── form.html
│   ├── orders/
│   │   ├── index.html
│   │   ├── form.html
│   │   └── show.html
│   ├── payments/
│   │   ├── index.html
│   │   └── form.html
│   ├── expenses/
│   │   ├── index.html
│   │   └── form.html
│   ├── inventory/
│   │   ├── index.html
│   │   └── form.html
│   ├── reports/
│   │   └── index.html
│   ├── users/
│   │   ├── index.html
│   │   └── form.html
│   ├── settings/
│   │   ├── index.html
│   │   └── profile.html
│   └── partials/ -- shared components
│       ├── sidebar.html
│       ├── navbar.html
│       ├── pagination.html
│       └── alerts.html
├── static/
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_customers.up.sql
│   ├── 000002_create_customers.down.sql
│   ├── 000003_create_services.up.sql
│   ├── 000003_create_services.down.sql
│   ├── 000004_create_orders.up.sql
│   ├── 000004_create_orders.down.sql
│   ├── 000005_create_order_tracking.up.sql
│   ├── 000005_create_order_tracking.down.sql
│   ├── 000006_create_payments.up.sql
│   ├── 000006_create_payments.down.sql
│   ├── 000007_create_expenses.up.sql
│   ├── 000007_create_expenses.down.sql
│   ├── 000008_create_inventory.up.sql
│   ├── 000008_create_inventory.down.sql
│   ├── 000009_create_settings.up.sql
│   ├── 000009_create_settings.down.sql
│   └── 000010_seed_data.up.sql
├── docs/
│   └── api/                     # Swagger output
├── scripts/
│   └── dev.sh
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── .air.toml
├── go.mod
├── go.sum
└── README.md
```

## 7. UI Design System

- **Color:** Slate/indigo palette (sidebar: slate-900, accent: indigo-600)
- **Cards:** White bg, rounded-xl, shadow-sm, border border-slate-200
- **Tables:** Thead bg-slate-50 with sticky header
- **Buttons:** btn-primary (indigo-600), btn-danger (red-600), btn-ghost (slate)
- **Badges:** Status badges (IN=emerald, OUT=red, processing=amber)
- **Forms:** Label above input, rounded-lg, border-slate-300, focus:ring-indigo-500
- **Layout:** Fixed sidebar (64px or 240px), top navbar, main content area
- **Responsive:** Sidebar collapses on mobile

## 8. Security

- **Password:** bcrypt (cost 12)
- **Session:** Cookie-based, HttpOnly, Secure, SameSite=Lax, encrypted
- **CSRF:** Token in session + hidden field in all POST forms
- **Input:** Validator tags on models + sanitize output in templates
- **AuthZ:** Middleware per role (admin can do everything, kasir limited, pegawai read-only)
- **XSS:** Go html/template auto-escape
- **SQLi:** GORM parameterized queries

## 9. Development Phases

| Phase | Features | Est. Tasks |
|-------|----------|------------|
| 1 | Project scaffold + setup (Go, Docker, config, DB) | 5-7 |
| 2 | Auth (login, logout, session, middleware) + Dashboard | 8-10 |
| 3 | Master data (Customers, Services) | 8-10 |
| 4 | Orders (CRUD, status tracking, auto numbering) | 10-12 |
| 5 | Payments | 5-7 |
| 6 | Expenses + Inventory | 8-10 |
| 7 | Reports + Export | 6-8 |
| 8 | User management + Settings + Profile | 8-10 |
| 9 | Polish + Docker deployment + README | 5-7 |

## 10. Future Considerations (YAGNI — not building now)

- Multi-outlet support
- Mobile API (REST/JSON endpoints)
- Real-time notifications (WebSocket)
- Invoice printing to thermal printer
- WhatsApp integration for order notifications
- Multi-currency
- Dark mode toggle
