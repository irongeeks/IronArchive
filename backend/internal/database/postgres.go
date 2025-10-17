package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"ironarchive/internal/config"
)

// PostgresConnection manages PostgreSQL database connection
type PostgresConnection struct {
	Pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewPostgresConnection creates a new PostgreSQL connection with pool configuration
func NewPostgresConnection(cfg *config.Config, logger *zap.Logger) (*PostgresConnection, error) {
	ctx := context.Background()

	// Parse connection URL
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Apply connection pool configuration
	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.DBHealthCheckPeriod

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	logger.Info("PostgreSQL connection pool configured",
		zap.Int32("max_conns", cfg.DBMaxConns),
		zap.Int32("min_conns", cfg.DBMinConns),
		zap.Duration("max_conn_lifetime", cfg.DBMaxConnLifetime),
		zap.Duration("max_conn_idle_time", cfg.DBMaxConnIdleTime),
	)

	return &PostgresConnection{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Ping validates the database connection
func (p *PostgresConnection) Ping(ctx context.Context) error {
	if err := p.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

// Close closes the database connection pool
func (p *PostgresConnection) Close() {
	if p.Pool != nil {
		p.Pool.Close()
		p.logger.Info("PostgreSQL connection pool closed")
	}
}
