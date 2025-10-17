package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisConnection manages Redis client connection
type RedisConnection struct {
	Client *redis.Client
	logger *zap.Logger
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(redisURL string, logger *zap.Logger) (*RedisConnection, error) {
	// Parse Redis URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Create Redis client
	client := redis.NewClient(opts)

	return &RedisConnection{
		Client: client,
		logger: logger,
	}, nil
}

// Ping validates the Redis connection
func (r *RedisConnection) Ping(ctx context.Context) error {
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (r *RedisConnection) Close() error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return err
		}
		r.logger.Info("Redis connection closed")
	}
	return nil
}
