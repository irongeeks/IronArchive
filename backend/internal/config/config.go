package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DatabaseURL      string
	RedisURL         string
	MeilisearchURL   string
	MeiliMasterKey   string
	JWTSecret        string
	ServerPort       string
	ServerHost       string
	LogLevel         string
	LogFormat        string
	EmailStoragePath string

	// Database connection pool configuration
	DBMaxConns          int32
	DBMinConns          int32
	DBMaxConnLifetime   time.Duration
	DBMaxConnIdleTime   time.Duration
	DBHealthCheckPeriod time.Duration

	// Timeout configuration
	MeilisearchTimeout time.Duration
	ShutdownTimeout    time.Duration
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional for development)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		RedisURL:         getEnv("REDIS_URL", ""),
		MeilisearchURL:   getEnv("MEILISEARCH_URL", ""),
		MeiliMasterKey:   getEnv("MEILI_MASTER_KEY", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		ServerHost:       getEnv("SERVER_HOST", "localhost"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		LogFormat:        getEnv("LOG_FORMAT", "json"),
		EmailStoragePath: getEnv("EMAIL_STORAGE_PATH", "./data/emails"),

		// Database pool configuration with sensible defaults
		DBMaxConns:          getEnvAsInt32("DB_MAX_CONNS", 25),
		DBMinConns:          getEnvAsInt32("DB_MIN_CONNS", 5),
		DBMaxConnLifetime:   getEnvAsDuration("DB_MAX_CONN_LIFETIME", 1*time.Hour),
		DBMaxConnIdleTime:   getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
		DBHealthCheckPeriod: getEnvAsDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),

		// Timeout configuration
		MeilisearchTimeout: getEnvAsDuration("MEILISEARCH_TIMEOUT", 5*time.Second),
		ShutdownTimeout:    getEnvAsDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	// Validate required configuration
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}
	if cfg.MeilisearchURL == "" {
		return nil, fmt.Errorf("MEILISEARCH_URL is required")
	}
	if cfg.MeiliMasterKey == "" {
		return nil, fmt.Errorf("MEILI_MASTER_KEY is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt32 retrieves an environment variable as int32 or returns a default value
func getEnvAsInt32(key string, defaultValue int32) int32 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(value)
}

// getEnvAsDuration retrieves an environment variable as time.Duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
