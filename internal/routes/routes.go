package routes

import (
	"html/template"
	"log/slog"

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

	// Initialize services
	authService := services.NewAuthService(userRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	profileHandler := handlers.NewProfileHandler(authService, userRepo)

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
	}
}
