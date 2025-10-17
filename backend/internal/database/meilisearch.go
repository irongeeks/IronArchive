package database

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// MeilisearchConnection manages Meilisearch client connection
type MeilisearchConnection struct {
	Client meilisearch.ServiceManager
	logger *zap.Logger
}

// NewMeilisearchConnection creates a new Meilisearch connection
func NewMeilisearchConnection(meilisearchURL, masterKey string, logger *zap.Logger) (*MeilisearchConnection, error) {
	// Create Meilisearch client
	client := meilisearch.New(meilisearchURL, meilisearch.WithAPIKey(masterKey))

	return &MeilisearchConnection{
		Client: client,
		logger: logger,
	}, nil
}

// Ping validates the Meilisearch connection
func (m *MeilisearchConnection) Ping(ctx context.Context) error {
	// Check health endpoint
	health, err := m.Client.Health()
	if err != nil {
		return fmt.Errorf("failed to check Meilisearch health: %w", err)
	}

	if health.Status != "available" {
		return fmt.Errorf("Meilisearch is not available, status: %s", health.Status)
	}

	return nil
}

// Close performs cleanup for Meilisearch connection
func (m *MeilisearchConnection) Close() {
	// Meilisearch Go client doesn't require explicit connection closing
	m.logger.Info("Meilisearch connection closed")
}
