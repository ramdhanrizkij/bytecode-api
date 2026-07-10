package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
)

const (
	ProviderLocal = "local"
	ProviderMinIO = "minio"
)

// PutObjectRequest carries an object upload request.
type PutObjectRequest struct {
	Bucket      string
	Key         string
	Reader      io.Reader
	Size        int64
	ContentType string
}

// StoredObject describes a stored object.
type StoredObject struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// Provider is the common interface for object storage backends.
type Provider interface {
	Put(ctx context.Context, req *PutObjectRequest) (*StoredObject, error)
	Delete(ctx context.Context, bucket, key string) error
	URL(bucket, key string) string
	ProviderName() string
	Close() error
}

// NewProvider creates a storage provider from configuration.
func NewProvider(cfg *config.StorageConfig, log *zap.Logger) (Provider, error) {
	if cfg == nil {
		cfg = &config.StorageConfig{}
	}

	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" {
		provider = ProviderLocal
	}

	switch provider {
	case ProviderLocal:
		return newLocalProvider(cfg)
	case ProviderMinIO:
		return newMinIOProvider(cfg, log)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Provider)
	}
}
