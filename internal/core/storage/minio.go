package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
)

type minioProvider struct {
	client        *minio.Client
	defaultBucket string
	publicURL     string
}

func newMinIOProvider(cfg *config.StorageConfig, log *zap.Logger) (Provider, error) {
	endpoint := strings.TrimSpace(cfg.MinIOEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT is required when STORAGE_PROVIDER=minio")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
		Region: cfg.MinIORegion,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	provider := &minioProvider{
		client:        client,
		defaultBucket: defaultBucket(cfg.DefaultBucket),
		publicURL:     strings.TrimRight(strings.TrimSpace(cfg.MinIOPublicURL), "/"),
	}

	ctx := context.Background()
	for _, bucket := range cfg.Buckets() {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to check minio bucket %q: %w", bucket, err)
		}
		if exists {
			continue
		}
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: cfg.MinIORegion}); err != nil {
			return nil, fmt.Errorf("failed to create minio bucket %q: %w", bucket, err)
		}
	}

	log.Info("minio storage connected",
		zap.String("endpoint", endpoint),
		zap.String("default_bucket", provider.defaultBucket),
	)

	return provider, nil
}

func (p *minioProvider) Put(ctx context.Context, req *PutObjectRequest) (*StoredObject, error) {
	if req == nil {
		return nil, fmt.Errorf("put request is required")
	}
	if req.Reader == nil {
		return nil, fmt.Errorf("put request reader is required")
	}

	bucket := p.bucketOrDefault(req.Bucket)
	key := normalizeKey(req.Key)
	if key == "" {
		return nil, fmt.Errorf("object key is required")
	}

	info, err := p.client.PutObject(ctx, bucket, key, req.Reader, req.Size, minio.PutObjectOptions{
		ContentType: req.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload object to minio: %w", err)
	}

	return &StoredObject{
		Bucket:      bucket,
		Key:         key,
		Size:        info.Size,
		ContentType: req.ContentType,
		URL:         p.URL(bucket, key),
	}, nil
}

func (p *minioProvider) Delete(ctx context.Context, bucket, key string) error {
	if err := p.client.RemoveObject(ctx, p.bucketOrDefault(bucket), normalizeKey(key), minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete minio object: %w", err)
	}
	return nil
}

func (p *minioProvider) URL(bucket, key string) string {
	bucket = p.bucketOrDefault(bucket)
	key = normalizeKey(key)

	if p.publicURL != "" {
		return p.publicURL + "/" + url.PathEscape(bucket) + "/" + escapeKey(key)
	}

	return "/" + url.PathEscape(bucket) + "/" + escapeKey(key)
}

func (p *minioProvider) ProviderName() string {
	return ProviderMinIO
}

func (p *minioProvider) Close() error {
	return nil
}

func (p *minioProvider) bucketOrDefault(bucket string) string {
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return p.defaultBucket
	}
	return bucket
}

func defaultBucket(bucket string) string {
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return "uploads"
	}
	return bucket
}
