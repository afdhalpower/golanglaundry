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

// DSN returns the PostgreSQL connection string
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
