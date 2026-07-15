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
		Level: slog.LevelInfo,
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
