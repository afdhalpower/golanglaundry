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
		return c.SendString("Laundry Management System")
	})

	slog.Info("routes registered successfully")
}

func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
	}
}
