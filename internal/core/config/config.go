package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name string `env:"APP_NAME" envDefault:"bytecode-api"`
	Port int    `env:"APP_PORT" envDefault:"8080"`
	Env  string `env:"APP_ENV" envDefault:"development"`
}

// SwaggerConfig controls the generated API documentation endpoint.
type SwaggerConfig struct {
	Enabled  bool   `env:"SWAGGER_ENABLED" envDefault:"true"`
	Username string `env:"SWAGGER_USERNAME" envDefault:""`
	Password string `env:"SWAGGER_PASSWORD" envDefault:""`
}

// IsProduction reports whether the application is running in production.
func (a AppConfig) IsProduction() bool {
	return strings.EqualFold(a.Env, "production")
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
	Secret             string `env:"JWT_SECRET" envDefault:"your-super-secret-key"`
	ExpiryHours        int    `env:"JWT_EXPIRY_HOURS" envDefault:"24"`
	RefreshExpiryHours int    `env:"JWT_REFRESH_EXPIRY_HOURS" envDefault:"168"`
}

// RedisConfig holds optional Redis caching configuration.
type RedisConfig struct {
	Enabled         bool   `env:"REDIS_ENABLED" envDefault:"false"`
	Host            string `env:"REDIS_HOST" envDefault:"localhost"`
	Port            int    `env:"REDIS_PORT" envDefault:"6379"`
	Password        string `env:"REDIS_PASSWORD" envDefault:""`
	DB              int    `env:"REDIS_DB" envDefault:"0"`
	CacheTTLMinutes int    `env:"REDIS_CACHE_TTL_MINUTES" envDefault:"5"`
}

// StorageConfig holds pluggable object storage configuration.
type StorageConfig struct {
	Provider       string `env:"STORAGE_PROVIDER" envDefault:"local"`
	DefaultBucket  string `env:"STORAGE_DEFAULT_BUCKET" envDefault:"uploads"`
	BucketsRaw     string `env:"STORAGE_BUCKETS" envDefault:"uploads"`
	LocalPath      string `env:"STORAGE_LOCAL_PATH" envDefault:"storage"`
	BaseURL        string `env:"STORAGE_BASE_URL" envDefault:"/storage"`
	MinIOEndpoint  string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinIOPublicURL string `env:"MINIO_PUBLIC_URL" envDefault:"http://localhost:9000"`
	MinIOAccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	MinIOSecretKey string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	MinIOUseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinIORegion    string `env:"MINIO_REGION" envDefault:"us-east-1"`
}

// Buckets returns the configured storage buckets, ensuring the default bucket exists.
func (s StorageConfig) Buckets() []string {
	seen := map[string]struct{}{}
	buckets := make([]string, 0)

	add := func(bucket string) {
		bucket = strings.TrimSpace(bucket)
		if bucket == "" {
			return
		}
		if _, exists := seen[bucket]; exists {
			return
		}
		seen[bucket] = struct{}{}
		buckets = append(buckets, bucket)
	}

	for _, bucket := range strings.Split(s.BucketsRaw, ",") {
		add(bucket)
	}
	add(s.DefaultBucket)

	if len(buckets) == 0 {
		return []string{"uploads"}
	}

	return buckets
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

// Config is the top-level application configuration struct.
type Config struct {
	App     AppConfig
	Swagger SwaggerConfig
	DB      DBConfig
	JWT     JWTConfig
	Redis   RedisConfig
	Storage StorageConfig
	Log     LogConfig
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
