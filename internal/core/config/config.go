package config

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name string `env:"APP_NAME" envDefault:"bytecode-api"`
	Port int    `env:"APP_PORT" envDefault:"8080"`
	Env  string `env:"APP_ENV" envDefault:"development"`
}

// DBConfig holds database connection configuration.
type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"secret"`
	Name     string `env:"DB_NAME" envDefault:"bytecode_api"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

// DSN returns the PostgreSQL connection string for GORM.
func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode,
	)
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret      string `env:"JWT_SECRET" envDefault:"your-super-secret-key"`
	ExpiryHours int    `env:"JWT_EXPIRY_HOURS" envDefault:"24"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

// Config is the top-level application configuration struct.
type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
	Log LogConfig
}

// LoadConfig loads configuration from the .env file (if present) and environment
// variables. The .env file is optional — in production, env vars are set directly.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists; silently ignore if not found (production use case).
	if err := godotenv.Load(); err != nil {
		// Only non-ErrNotFound errors should be surfaced.
		// godotenv.Load returns an error if the file doesn't exist,
		// but we intentionally allow that.
		_ = err // Best-effort: ignore missing .env file
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Security validation: prevent running production with the default JWT secret.
	if cfg.App.Env == "production" && cfg.JWT.Secret == "your-super-secret-key" {
		return nil, errors.New("security error: JWT_SECRET must be changed from the default value in production")
	}

	return cfg, nil
}
