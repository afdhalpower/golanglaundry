# Phase 2: Authentication + Dashboard — Implementation Plan

> **For agentic workers:** Use `superpowers:subagent-driven-development` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Full session-based authentication (login, logout, middleware, profile, change password) + live dashboard with charts and recent activity.

**Architecture:** Session auth using Fiber's session middleware with PostgreSQL store. Middleware checks session before every protected route. Dashboard queries aggregate data from DB.

**Tech Stack:** gofiber/session/v2, bcrypt, HTML templates with Tailwind + Alpine.js, Chart.js for dashboard graphs.

## Global Constraints

- Go module path: `github.com/afdhalpower/golanglaundry`
- All source under `internal/` except `cmd/server/main.go`
- All templates in `templates/`
- All static assets under `static/`
- Session stored in PostgreSQL via GORM (or memory for dev)
- Password hashing with bcrypt (cost 12)
- Follow existing project structure patterns

---

### Task 1: Add Session Middleware + Auth Service

**Files:**
- Create `internal/services/auth_service.go`
- Modify `internal/middleware/auth.go`
- Modify `internal/config/database.go` (add session store)
- Modify `cmd/server/main.go` (wire up session)
- Modify `internal/routes/routes.go` (add auth routes)

- [ ] **Step 1: Install session + bcrypt dependency**

```bash
export GONOSUMCHECK=* GONOSUMDB=* GOINSECURE=* GOPROXY=direct
go get github.com/gofiber/fiber/v3/middleware/session
go get github.com/gofiber/storage/postgres/v3
go get golang.org/x/crypto
go mod tidy
```

- [ ] **Step 2: Create AuthService**

```go
// internal/services/auth_service.go
package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type LoginRequest struct {
	Email    string
	Password string
}

func (s *AuthService) Login(req LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("akun telah dinonaktifkan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	return user, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("password lama salah")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(hash))
}

func (s *AuthService) UpdateProfile(userID uint, name, email string) error {
	return s.userRepo.UpdateProfile(userID, name, email)
}
```

- [ ] **Step 3: Create UserRepository**

```go
// internal/repositories/user_repository.go
package repositories

import (
	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) UpdatePassword(id uint, hash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *UserRepository) UpdateProfile(id uint, name, email string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  name,
		"email": email,
	}).Error
}

func (r *UserRepository) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
```

- [ ] **Step 4: Update auth middleware to check session**

```go
// internal/middleware/auth.go
package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var Store *session.Store

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err != nil {
			return c.Redirect().To("/auth/login")
		}

		userID := sess.Get("user_id")
		if userID == nil {
			return c.Redirect().To("/auth/login")
		}

		c.Locals("user_id", userID)
		c.Locals("user_name", sess.Get("user_name"))
		c.Locals("user_role", sess.Get("user_role"))

		return c.Next()
	}
}

func RolePermission(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("Akses ditolak")
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).SendString("Akses ditolak")
	}
}
```

- [ ] **Step 5: Wire up session + auth in main.go and routes.go**

Update `cmd/server/main.go`:
- Initialize session store with config
- Pass store to middleware

Update `internal/routes/routes.go`:
- Add `/auth/login` GET/POST
- Add `/auth/logout` GET
- Add `/profile` GET/PUT
- Add `/profile/password` PUT

- [ ] **Step 6: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat: add authentication service, repository, and session middleware"
```

---

### Task 2: Create Auth Handler + Login Template

**Files:**
- Create `internal/handlers/auth_handler.go`
- Create `templates/layouts/main.html`
- Create `templates/auth/login.html`
- Create `templates/partials/alerts.html`

- [ ] **Step 1: Create auth handler**

```go
// internal/handlers/auth_handler.go
package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/afdhalpower/golanglaundry/internal/services"
	"github.com/afdhalpower/golanglaundry/internal/validation"
)

type AuthHandler struct {
	authService *services.AuthService
	sessionStore *session.Store
}

func NewAuthHandler(authService *services.AuthService, sessionStore *session.Store) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		sessionStore: sessionStore,
	}
}

func (h *AuthHandler) LoginPage(c fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"title": "Login",
	}, "layouts/main")
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Render("auth/login", fiber.Map{
			"title": "Login",
			"error": "Email dan password wajib diisi",
		}, "layouts/main")
	}

	user, err := h.authService.Login(services.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return c.Render("auth/login", fiber.Map{
			"title": "Login",
			"error": err.Error(),
		}, "layouts/main")
	}

	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return c.Render("auth/login", fiber.Map{
			"title": "Login",
			"error": "Gagal memulai session",
		}, "layouts/main")
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_name", user.Name)
	sess.Set("user_role", user.Role)
	sess.Set("user_email", user.Email)

	if err := sess.Save(); err != nil {
		return c.Render("auth/login", fiber.Map{
			"title": "Login",
			"error": "Gagal menyimpan session",
		}, "layouts/main")
	}

	return c.Redirect().To("/dashboard")
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := h.sessionStore.Get(c)
	if err == nil {
		sess.Destroy()
	}
	return c.Redirect().To("/auth/login")
}
```

- [ ] **Step 2: Create main layout template**

```html
<!-- templates/layouts/main.html -->
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .title}}{{.title}} - {{end}}Laundry Management System</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body class="bg-slate-100 min-h-screen">
    {{if .hideLayout}}
        {{embed}}
    {{else}}
        {{if ne .title "Login"}}
            {{template "partials/navbar" .}}
            <div class="flex">
                {{template "partials/sidebar" .}}
                <main class="flex-1 p-6 ml-64">
                    {{template "partials/alerts" .}}
                    {{embed}}
                </main>
            </div>
        {{else}}
            <div class="min-h-screen flex items-center justify-center">
                {{embed}}
            </div>
        {{endif}}
    {{endif}}

    <script src="/static/js/app.js"></script>
</body>
</html>
```

- [ ] **Step 3: Create login template**

```html
<!-- templates/auth/login.html -->
<div class="w-full max-w-md">
    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
        <div class="text-center mb-8">
            <h1 class="text-2xl font-bold text-slate-900">Laundry Management</h1>
            <p class="text-sm text-slate-500 mt-1">Masuk ke akun Anda</p>
        </div>

        {{if .error}}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm mb-6">
            {{.error}}
        </div>
        {{end}}

        <form action="/auth/login" method="POST" class="space-y-5">
            <div>
                <label for="email" class="block text-sm font-medium text-slate-700 mb-1">Email</label>
                <input type="email" id="email" name="email" required
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition">
            </div>
            <div>
                <label for="password" class="block text-sm font-medium text-slate-700 mb-1">Password</label>
                <input type="password" id="password" name="password" required
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition">
            </div>
            <button type="submit"
                class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2.5 px-4 rounded-lg transition duration-150">
                Masuk
            </button>
        </form>
    </div>
</div>
```

- [ ] **Step 4: Create alerts partial**

```html
<!-- templates/partials/alerts.html -->
{{if .success}}
<div class="bg-emerald-50 border border-emerald-200 text-emerald-700 px-4 py-3 rounded-lg text-sm mb-4 flex items-center justify-between">
    <span>{{.success}}</span>
    <button onclick="this.parentElement.remove()" class="text-emerald-500 hover:text-emerald-700">&times;</button>
</div>
{{end}}
{{if .error}}
<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm mb-4 flex items-center justify-between">
    <span>{{.error}}</span>
    <button onclick="this.parentElement.remove()" class="text-red-500 hover:text-red-700">&times;</button>
</div>
{{end}}
```

- [ ] **Step 5: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add auth handler with login page and session management"
```

---

### Task 3: Create Dashboard Handler + Template

**Files:**
- Create `internal/handlers/dashboard_handler.go`
- Create `internal/services/dashboard_service.go`
- Create `internal/repositories/dashboard_repository.go`
- Create `templates/dashboard/index.html`
- Create `templates/partials/navbar.html`
- Create `templates/partials/sidebar.html`

- [ ] **Step 1: Create dashboard service**

```go
// internal/services/dashboard_service.go
package services

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/repositories"
)

type DashboardStats struct {
	TotalOrdersToday  int64
	OrdersInProgress  int64
	OrdersCompleted   int64
	RevenueToday      float64
	RevenueThisMonth  float64
	TotalCustomers    int64
	TotalExpensesThisMonth float64
}

type DashboardService struct {
	dashboardRepo *repositories.DashboardRepository
}

func NewDashboardService(dashboardRepo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

func (s *DashboardService) GetStats() (*DashboardStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	ordersToday, err := s.dashboardRepo.CountOrdersSince(todayStart)
	if err != nil {
		return nil, err
	}

	ordersInProgress, err := s.dashboardRepo.CountOrdersByStatus([]string{"dicuci", "dikeringkan", "disetrika", "siap_diambil"})
	if err != nil {
		return nil, err
	}

	ordersCompleted, err := s.dashboardRepo.CountOrdersByStatus([]string{"sudah_diambil"})
	if err != nil {
		return nil, err
	}

	revenueToday, err := s.dashboardRepo.SumPaymentsSince(todayStart)
	if err != nil {
		return nil, err
	}

	revenueMonth, err := s.dashboardRepo.SumPaymentsSince(monthStart)
	if err != nil {
		return nil, err
	}

	totalCustomers, err := s.dashboardRepo.CountAllCustomers()
	if err != nil {
		return nil, err
	}

	expensesMonth, err := s.dashboardRepo.SumExpensesSince(monthStart)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalOrdersToday:     ordersToday,
		OrdersInProgress:     ordersInProgress,
		OrdersCompleted:      ordersCompleted,
		RevenueToday:         revenueToday,
		RevenueThisMonth:     revenueMonth,
		TotalCustomers:       totalCustomers,
		TotalExpensesThisMonth: expensesMonth,
	}, nil
}
```

- [ ] **Step 2: Create dashboard repository**

```go
// internal/repositories/dashboard_repository.go
package repositories

import (
	"time"

	"github.com/afdhalpower/golanglaundry/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) CountOrdersSince(since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("created_at >= ?", since).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) CountOrdersByStatus(statuses []string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("status IN ?", statuses).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) SumPaymentsSince(since time.Time) (float64, error) {
	var sum struct {
		Total float64
	}
	err := r.db.Model(&models.Payment{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("created_at >= ? AND status = 'lunas'", since).
		Scan(&sum).Error
	return sum.Total, err
}

func (r *DashboardRepository) CountAllCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Customer{}).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) SumExpensesSince(since time.Time) (float64, error) {
	var sum struct {
		Total float64
	}
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("date >= ?", since).
		Scan(&sum).Error
	return sum.Total, err
}
```

- [ ] **Step 3: Create dashboard handler**

```go
// internal/handlers/dashboard_handler.go
package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) Index(c fiber.Ctx) error {
	stats, err := h.dashboardService.GetStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal memuat dashboard")
	}

	return c.Render("dashboard/index", fiber.Map{
		"title":  "Dashboard",
		"stats":  stats,
		"role":   c.Locals("user_role"),
		"name":   c.Locals("user_name"),
	}, "layouts/main")
}
```

- [ ] **Step 4: Create dashboard template**

```html
<!-- templates/dashboard/index.html -->
<div class="space-y-6">
    <div>
        <h1 class="text-2xl font-bold text-slate-900">Dashboard</h1>
        <p class="text-sm text-slate-500 mt-1">Selamat datang kembali, {{.name}}!</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Order Hari Ini</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">{{.stats.TotalOrdersToday}}</p>
                </div>
                <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Diproses</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">{{.stats.OrdersInProgress}}</p>
                </div>
                <div class="w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Pendapatan Hari Ini</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">Rp {{printf "%.0f" .stats.RevenueToday}}</p>
                </div>
                <div class="w-10 h-10 bg-emerald-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Pelanggan</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">{{.stats.TotalCustomers}}</p>
                </div>
                <div class="w-10 h-10 bg-sky-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
                </div>
            </div>
        </div>
    </div>

    <!-- Second row -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Pendapatan Bulan Ini</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">Rp {{printf "%.0f" .stats.RevenueThisMonth}}</p>
                </div>
                <div class="w-10 h-10 bg-emerald-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Selesai</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">{{.stats.OrdersCompleted}}</p>
                </div>
                <div class="w-10 h-10 bg-emerald-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm text-slate-500">Pengeluaran Bulan Ini</p>
                    <p class="text-2xl font-bold text-slate-900 mt-1">Rp {{printf "%.0f" .stats.TotalExpensesThisMonth}}</p>
                </div>
                <div class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"></path></svg>
                </div>
            </div>
        </div>
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <h3 class="text-sm font-semibold text-slate-900 mb-4">Pendapatan Harian (7 Hari)</h3>
            <canvas id="revenueChart" height="200"></canvas>
        </div>
        <div class="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
            <h3 class="text-sm font-semibold text-slate-900 mb-4">Order Mingguan</h3>
            <canvas id="orderChart" height="200"></canvas>
        </div>
    </div>
</div>

<script>
document.addEventListener('DOMContentLoaded', function() {
    // Revenue Chart
    new Chart(document.getElementById('revenueChart'), {
        type: 'line',
        data: {
            labels: ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'],
            datasets: [{
                label: 'Pendapatan',
                data: [0, 0, 0, 0, 0, 0, 0],
                borderColor: '#6366f1',
                backgroundColor: 'rgba(99, 102, 241, 0.1)',
                fill: true,
                tension: 0.3
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: { legend: { display: false } },
            scales: {
                y: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } },
                x: { grid: { display: false } }
            }
        }
    });

    // Order Chart
    new Chart(document.getElementById('orderChart'), {
        type: 'bar',
        data: {
            labels: ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min'],
            datasets: [{
                label: 'Order',
                data: [0, 0, 0, 0, 0, 0, 0],
                backgroundColor: '#818cf8',
                borderRadius: 6
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: { legend: { display: false } },
            scales: {
                y: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } },
                x: { grid: { display: false } }
            }
        }
    });
});
</script>
```

- [ ] **Step 5: Create navbar and sidebar partials**

```html
<!-- templates/partials/navbar.html -->
<header class="bg-white border-b border-slate-200 sticky top-0 z-30">
    <div class="flex items-center justify-between h-16 px-6">
        <div class="flex items-center gap-3">
            <h2 class="text-lg font-semibold text-slate-900">{{.title}}</h2>
        </div>
        <div class="flex items-center gap-4">
            <span class="text-sm text-slate-600">{{.name}}</span>
            <span class="text-xs bg-indigo-100 text-indigo-700 px-2 py-1 rounded-full font-medium capitalize">{{.role}}</span>
            <a href="/auth/logout" class="text-sm text-slate-500 hover:text-red-600 transition">Logout</a>
        </div>
    </div>
</header>
```

```html
<!-- templates/partials/sidebar.html -->
<aside class="fixed left-0 top-0 h-full w-64 bg-slate-900 text-white z-40">
    <div class="flex items-center gap-2 h-16 px-6 border-b border-slate-700">
        <svg class="w-7 h-7 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path></svg>
        <span class="font-semibold text-sm">Laundry Manager</span>
    </div>
    <nav class="p-4 space-y-1">
        <a href="/dashboard" class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm {{if eq .title `Dashboard`}}bg-indigo-600 text-white{{else}}text-slate-300 hover:bg-slate-800{{end}} transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
            Dashboard
        </a>
        <a href="/customers" class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
            Pelanggan
        </a>
        <a href="/services" class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path></svg>
            Layanan
        </a>
        <a href="/orders" class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path></svg>
            Pesanan
        </a>
        <hr class="border-slate-700 my-3">
        <a href="/settings" class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-slate-300 hover:bg-slate-800 transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
            Pengaturan
        </a>
    </nav>
</aside>
```

- [ ] **Step 6: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat: add dashboard handler with stats cards and charts"
```

---

### Task 4: Create Profile Handler + Template

**Files:**
- Create `internal/handlers/profile_handler.go`
- Create `templates/settings/profile.html`

- [ ] **Step 1: Create profile handler**

```go
// internal/handlers/profile_handler.go
package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/afdhalpower/golanglaundry/internal/helpers"
	"github.com/afdhalpower/golanglaundry/internal/services"
)

type ProfileHandler struct {
	authService *services.AuthService
}

func NewProfileHandler(authService *services.AuthService) *ProfileHandler {
	return &ProfileHandler{authService: authService}
}

func (h *ProfileHandler) Index(c fiber.Ctx) error {
	return c.Render("settings/profile", fiber.Map{
		"title": "Profil",
	}, "layouts/main")
}

func (h *ProfileHandler) Update(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	name := c.FormValue("name")
	email := c.FormValue("email")

	if err := h.authService.UpdateProfile(userID, name, email); err != nil {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Gagal memperbarui profil",
		}, "layouts/main")
	}

	return c.Render("settings/profile", fiber.Map{
		"title":   "Profil",
		"success": "Profil berhasil diperbarui",
	}, "layouts/main")
}

func (h *ProfileHandler) ChangePassword(c fiber.Ctx) error {
	userID := helpers.LogAndGetUserID(c)
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")

	if newPassword == "" || len(newPassword) < 6 {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": "Password baru minimal 6 karakter",
		}, "layouts/main")
	}

	if err := h.authService.ChangePassword(userID, oldPassword, newPassword); err != nil {
		return c.Render("settings/profile", fiber.Map{
			"title": "Profil",
			"error": err.Error(),
		}, "layouts/main")
	}

	return c.Render("settings/profile", fiber.Map{
		"title":   "Profil",
		"success": "Password berhasil diubah",
	}, "layouts/main")
}
```

- [ ] **Step 2: Create profile template**

```html
<!-- templates/settings/profile.html -->
<div class="max-w-2xl mx-auto space-y-6">
    <div>
        <h1 class="text-2xl font-bold text-slate-900">Profil</h1>
        <p class="text-sm text-slate-500 mt-1">Kelola informasi profil dan password Anda</p>
    </div>

    <!-- Edit Profile -->
    <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
        <h3 class="text-lg font-semibold text-slate-900 mb-4">Informasi Profil</h3>
        <form action="/profile" method="POST" class="space-y-4">
            <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Nama</label>
                <input type="text" name="name" value="{{.user.Name}}" required
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 outline-none transition">
            </div>
            <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
                <input type="email" name="email" value="{{.user.Email}}" required
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 outline-none transition">
            </div>
            <button type="submit" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2.5 rounded-lg text-sm font-medium transition">
                Simpan Perubahan
            </button>
        </form>
    </div>

    <!-- Change Password -->
    <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
        <h3 class="text-lg font-semibold text-slate-900 mb-4">Ganti Password</h3>
        <form action="/profile/password" method="POST" class="space-y-4">
            <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Password Lama</label>
                <input type="password" name="old_password" required
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 outline-none transition">
            </div>
            <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Password Baru</label>
                <input type="password" name="new_password" required minlength="6"
                    class="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-500 outline-none transition">
            </div>
            <button type="submit" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2.5 rounded-lg text-sm font-medium transition">
                Ganti Password
            </button>
        </form>
    </div>
</div>
```

- [ ] **Step 3: Verify compile**

```bash
go build ./... && echo "BUILD OK"
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: add profile handler with edit profile and change password"
```

---

### Task 5: Wire Everything Together + Test

**Files:**
- Modify `cmd/server/main.go`
- Modify `internal/routes/routes.go`

- [ ] **Step 1: Update main.go with session store and handlers**

```go
// cmd/server/main.go (key additions)
import (
    "github.com/gofiber/fiber/v3/middleware/session"
    "github.com/gofiber/storage/postgres/v3"
    "github.com/afdhalpower/golanglaundry/internal/middleware"
    "github.com/afdhalpower/golanglaundry/internal/handlers"
    "github.com/afdhalpower/golanglaundry/internal/services"
    "github.com/afdhalpower/golanglaundry/internal/repositories"
)

// In main():
// Initialize session store
storage := postgres.New(postgres.Config{
    Table: "sessions",
})
sessionStore := session.New(session.Config{
    Storage:    storage,
    KeyLookup:  "cookie:laundry_session",
})
middleware.Store = sessionStore

// Initialize repositories
userRepo := repositories.NewUserRepository(db)

// Initialize services
authService := services.NewAuthService(userRepo)

// Initialize handlers
authHandler := handlers.NewAuthHandler(authService, sessionStore)
```

- [ ] **Step 2: Update routes.go**

```go
// Auth routes (public)
auth.Get("/login", authHandler.LoginPage)
auth.Post("/login", authHandler.Login)
auth.Get("/logout", authHandler.Logout)

// Dashboard
protected.Get("/dashboard", dashboardHandler.Index)
protected.Get("/", dashboardHandler.Index)  // redirect / to /dashboard

// Profile
protected.Get("/profile", profileHandler.Index)
protected.Post("/profile", profileHandler.Update)
protected.Post("/profile/password", profileHandler.ChangePassword)
```

- [ ] **Step 3: Full build verification**

```bash
go build ./... && echo "BUILD OK"
go vet ./... && echo "VET OK"
```

- [ ] **Step 4: Run server and test login flow**

```bash
# Kill any existing server
pkill -f "go run.*cmd/server" 2>/dev/null || true
sleep 1

# Start server
go run ./cmd/server/ &
sleep 3

# Test health
curl -s http://localhost:3000/health
# Expected: {"status":"ok"}

# Test login page loads
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/auth/login
# Expected: 200

# Test login
curl -s -c /tmp/cookies.txt -X POST http://localhost:3000/auth/login \
  -d "email=admin@laundry.com&password=admin123" \
  -o /dev/null -w "%{http_code} - %{redirect_url}"
# Expected: 200 or redirect

# Test dashboard (authenticated)
curl -s -b /tmp/cookies.txt http://localhost:3000/dashboard -o /dev/null -w "%{http_code}"
# Expected: 200

# Kill server
kill %1 2>/dev/null || true
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: wire auth, dashboard, and profile handlers into main server"
```

---

## Verification Checklist (Final)

- [ ] `go build ./...` passes with no errors
- [ ] `go vet ./...` passes with no errors
- [ ] Login page renders at `/auth/login`
- [ ] Login with admin@laundry.com / admin123 works
- [ ] After login, redirects to `/dashboard`
- [ ] Dashboard shows stats cards
- [ ] Logout clears session and redirects to login
- [ ] Profile page allows editing name/email and changing password
- [ ] AuthRequired middleware blocks unauthenticated access
- [ ] All templates render properly
