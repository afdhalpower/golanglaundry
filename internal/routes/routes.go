package routes

import (
	"html/template"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/afdhalpower/golanglaundry/internal/handlers"
	"github.com/afdhalpower/golanglaundry/internal/middleware"
	"github.com/afdhalpower/golanglaundry/internal/repositories"
	"github.com/afdhalpower/golanglaundry/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Static files
	app.Use("/static", static.New("./static"))

	// Global middleware
	app.Use(middleware.Logger())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	dashboardRepo := repositories.NewDashboardRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	serviceRepo := repositories.NewServiceRepository(db)

	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	customerService := services.NewCustomerService(customerRepo)
	serviceService := services.NewServiceService(serviceRepo)
	orderService := services.NewOrderService(orderRepo, customerRepo, serviceRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	profileHandler := handlers.NewProfileHandler(authService, userRepo)
	customerHandler := handlers.NewCustomerHandler(customerService)
	serviceHandler := handlers.NewServiceHandler(serviceService)
	orderHandler := handlers.NewOrderHandler(orderService, customerService, serviceService)

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth routes (public)
	auth := app.Group("/auth")
	auth.Get("/login", authHandler.LoginPage)
	auth.Post("/login", authHandler.Login)
	auth.Get("/logout", authHandler.Logout)

	// Protected routes
	protected := app.Group("/", middleware.AuthRequired())

	// Dashboard
	protected.Get("/", dashboardHandler.Index)
	protected.Get("/dashboard", dashboardHandler.Index)

	// Profile
	protected.Get("/profile", profileHandler.Index)
	protected.Post("/profile", profileHandler.Update)
	protected.Post("/profile/password", profileHandler.ChangePassword)

	// Customers
	customers := protected.Group("/customers")
	customers.Get("/", customerHandler.Index)
	customers.Get("/new", customerHandler.New)
	customers.Post("/", customerHandler.Create)
	customers.Get("/:id", customerHandler.Show)
	customers.Get("/:id/edit", customerHandler.Edit)
	customers.Post("/:id", customerHandler.Update)
	customers.Post("/:id/delete", customerHandler.Delete)

	// Services
	services := protected.Group("/services")
	services.Get("/", serviceHandler.Index)
	services.Get("/new", serviceHandler.New)
	services.Post("/", serviceHandler.Create)
	services.Get("/:id/edit", serviceHandler.Edit)
	services.Post("/:id", serviceHandler.Update)
	services.Post("/:id/delete", serviceHandler.Delete)

	// Orders
	orders := protected.Group("/orders")
	orders.Get("/", orderHandler.Index)
	orders.Get("/new", orderHandler.New)
	orders.Post("/", orderHandler.Create)
	orders.Get("/:id", orderHandler.Show)
	orders.Post("/:id/status", orderHandler.UpdateStatus)
	orders.Post("/:id/delete", orderHandler.Delete)

	slog.Info("routes registered successfully")
}

func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"slice": func(s string, i, j int) string {
			if i > len(s) {
				return ""
			}
			if j > len(s) {
				j = len(s)
			}
			return s[i:j]
		},
		"upper": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return string(s[0]-32) + s[1:]
		},
		"loop": func(count int) []int {
			var items []int
			for i := 0; i < count; i++ {
				items = append(items, i)
			}
			return items
		},
		"nowDate": func() string {
			return time.Now().Format("2006-01-02")
		},
	}
}
