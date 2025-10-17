package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ironarchive/internal/config"
	"ironarchive/internal/database"
	"ironarchive/internal/utils"

	"go.uber.org/zap"
)

func main() {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := utils.NewLogger(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("IronArchive server starting...",
		zap.String("version", "1.0.0"),
		zap.String("port", cfg.ServerPort),
	)

	// Initialize PostgreSQL connection
	logger.Info("Connecting to PostgreSQL...", zap.String("url", maskConnectionString(cfg.DatabaseURL)))
	pgConn, err := database.NewPostgresConnection(cfg, logger)
	if err != nil {
		logger.Error("Failed to create PostgreSQL connection", zap.Error(err))
		os.Exit(1)
	}
	defer pgConn.Close()

	// Validate PostgreSQL connection
	if err := pgConn.Ping(ctx); err != nil {
		logger.Error("Failed to ping PostgreSQL", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("PostgreSQL connection successful")

	// Initialize Redis connection
	logger.Info("Connecting to Redis...", zap.String("url", maskConnectionString(cfg.RedisURL)))
	redisConn, err := database.NewRedisConnection(cfg.RedisURL, logger)
	if err != nil {
		logger.Error("Failed to create Redis connection", zap.Error(err))
		os.Exit(1)
	}
	defer func() {
		if err := redisConn.Close(); err != nil {
			logger.Warn("Failed to close Redis connection", zap.Error(err))
		}
	}()

	// Validate Redis connection
	if err := redisConn.Ping(ctx); err != nil {
		logger.Error("Failed to ping Redis", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Redis connection successful")

	// Initialize Meilisearch connection
	logger.Info("Connecting to Meilisearch...", zap.String("url", cfg.MeilisearchURL))
	meiliConn, err := database.NewMeilisearchConnection(cfg.MeilisearchURL, cfg.MeiliMasterKey, logger)
	if err != nil {
		logger.Error("Failed to create Meilisearch connection", zap.Error(err))
		os.Exit(1)
	}
	defer meiliConn.Close()

	// Validate Meilisearch connection
	pingCtx, pingCancel := context.WithTimeout(ctx, cfg.MeilisearchTimeout)
	defer pingCancel()
	if err := meiliConn.Ping(pingCtx); err != nil {
		logger.Error("Failed to ping Meilisearch", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Meilisearch connection successful")

	logger.Info("All service connections validated successfully")
	logger.Info(fmt.Sprintf("Server is ready on %s:%s", cfg.ServerHost, cfg.ServerPort))

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	// Perform cleanup
	select {
	case <-shutdownCtx.Done():
		logger.Warn("Shutdown timeout exceeded")
	default:
		logger.Info("Server stopped gracefully")
	}
}

// maskConnectionString masks sensitive information in connection strings
func maskConnectionString(connStr string) string {
	// Mask password in connection string for security
	// Handles formats like postgres://user:password@host/db and redis://password@host

	// Find password delimiter patterns
	if idx := strings.Index(connStr, "://"); idx != -1 {
		prefix := connStr[:idx+3] // Keep protocol (e.g., "postgres://")
		rest := connStr[idx+3:]

		// Check for user:password@host pattern
		if atIdx := strings.Index(rest, "@"); atIdx != -1 {
			beforeAt := rest[:atIdx]
			afterAt := rest[atIdx:]

			// Check if there's a colon (user:password format)
			if colonIdx := strings.Index(beforeAt, ":"); colonIdx != -1 {
				user := beforeAt[:colonIdx]
				return prefix + user + ":***@" + afterAt
			}
			// Just password@ format (Redis)
			return prefix + "***@" + afterAt
		}

		// No auth in connection string
		return connStr
	}

	// Fallback for non-standard formats
	if len(connStr) > 20 {
		return connStr[:10] + "***" + connStr[len(connStr)-10:]
	}
	return "***"
}

