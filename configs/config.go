package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	DB     DatabaseConfig
	JWT    JWTConfig
	SMTP   SMTPConfig
	Worker WorkerConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     string
	BaseURL  string
	Shutdown time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret           string
	Issuer           string
	AccessTTLMinutes int
	RefreshTTLHours  int
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

type WorkerConfig struct {
	Concurrency int
}

func Load() (Config, error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			panic(recovered)
		}
	}()

	return load()
}

func load() (cfg Config, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("config load failed: %v", recovered)
		}
	}()

	_ = godotenv.Load()

	cfg = Config{
		App: AppConfig{
			Name:     getRequired("APP_NAME"),
			Env:      getRequired("APP_ENV"),
			Port:     getRequired("APP_PORT"),
			BaseURL:  getRequired("APP_BASE_URL"),
			Shutdown: 10 * time.Second,
		},
		DB: DatabaseConfig{
			Host:     getRequired("DB_HOST"),
			Port:     getRequiredInt("DB_PORT"),
			User:     getRequired("DB_USER"),
			Password: getRequired("DB_PASSWORD"),
			Name:     getRequired("DB_NAME"),
			SSLMode:  getRequired("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:           getRequired("JWT_SECRET"),
			Issuer:           getRequired("JWT_ISSUER"),
			AccessTTLMinutes: getRequiredInt("JWT_ACCESS_TTL_MINUTES"),
			RefreshTTLHours:  getRequiredInt("JWT_REFRESH_TTL_HOURS"),
		},
		SMTP: SMTPConfig{
			Host:     getRequired("SMTP_HOST"),
			Port:     getRequiredInt("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     getRequired("SMTP_FROM"),
			FromName: getRequired("SMTP_FROM_NAME"),
		},
		Worker: WorkerConfig{
			Concurrency: getRequiredInt("WORKER_CONCURRENCY"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var issues []string

	if c.App.Port == "" {
		issues = append(issues, "APP_PORT is required")
	}
	if c.DB.Port <= 0 {
		issues = append(issues, "DB_PORT must be greater than zero")
	}
	if c.JWT.AccessTTLMinutes <= 0 {
		issues = append(issues, "JWT_ACCESS_TTL_MINUTES must be greater than zero")
	}
	if c.JWT.RefreshTTLHours <= 0 {
		issues = append(issues, "JWT_REFRESH_TTL_HOURS must be greater than zero")
	}
	if c.SMTP.Port <= 0 {
		issues = append(issues, "SMTP_PORT must be greater than zero")
	}
	if c.Worker.Concurrency <= 0 {
		issues = append(issues, "WORKER_CONCURRENCY must be greater than zero")
	}

	if len(issues) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(issues, "; "))
	}

	return nil
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}

func getRequired(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(fmt.Sprintf("missing required env var %s", key))
	}

	return value
}

func getRequiredInt(key string) int {
	value := getRequired(key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid integer env var %s: %v", key, err))
	}

	return parsed
}
