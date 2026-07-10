package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
)

type localProvider struct {
	rootPath string
	baseURL  string
}

func newLocalProvider(cfg *config.StorageConfig) (Provider, error) {
	rootPath := strings.TrimSpace(cfg.LocalPath)
	if rootPath == "" {
		rootPath = "storage"
	}

	if err := os.MkdirAll(rootPath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create local storage root: %w", err)
	}

	for _, bucket := range cfg.Buckets() {
		if err := os.MkdirAll(filepath.Join(rootPath, bucket), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create local bucket %q: %w", bucket, err)
		}
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = "/storage"
	}

	return &localProvider{
		rootPath: rootPath,
		baseURL:  strings.TrimRight(baseURL, "/"),
	}, nil
}

func (p *localProvider) Put(ctx context.Context, req *PutObjectRequest) (*StoredObject, error) {
	if req == nil {
		return nil, fmt.Errorf("put request is required")
	}
	if req.Reader == nil {
		return nil, fmt.Errorf("put request reader is required")
	}

	objectPath, err := p.objectPath(req.Bucket, req.Key)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(objectPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create local object directory: %w", err)
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create local object: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, req.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write local object: %w", err)
	}

	return &StoredObject{
		Bucket:      p.bucketOrDefault(req.Bucket),
		Key:         normalizeKey(req.Key),
		Size:        size,
		ContentType: req.ContentType,
		URL:         p.URL(req.Bucket, req.Key),
	}, nil
}

func (p *localProvider) Delete(ctx context.Context, bucket, key string) error {
	objectPath, err := p.objectPath(bucket, key)
	if err != nil {
		return err
	}

	if err := os.Remove(objectPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete local object: %w", err)
	}

	return nil
}

func (p *localProvider) URL(bucket, key string) string {
	return p.baseURL + "/" + url.PathEscape(p.bucketOrDefault(bucket)) + "/" + escapeKey(normalizeKey(key))
}

func (p *localProvider) ProviderName() string {
	return ProviderLocal
}

func (p *localProvider) Close() error {
	return nil
}

func (p *localProvider) objectPath(bucket, key string) (string, error) {
	bucket = p.bucketOrDefault(bucket)
	key = normalizeKey(key)
	if key == "" {
		return "", fmt.Errorf("object key is required")
	}

	objectPath := filepath.Join(p.rootPath, bucket, filepath.FromSlash(key))
	cleanRoot, err := filepath.Abs(p.rootPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve storage root: %w", err)
	}
	cleanObject, err := filepath.Abs(objectPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve storage object path: %w", err)
	}

	rootPrefix := cleanRoot + string(filepath.Separator)
	if cleanObject != cleanRoot && !strings.HasPrefix(cleanObject, rootPrefix) {
		return "", fmt.Errorf("invalid object key path")
	}

	return cleanObject, nil
}

func (p *localProvider) bucketOrDefault(bucket string) string {
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return "uploads"
	}
	return bucket
}

func normalizeKey(key string) string {
	key = strings.TrimSpace(strings.ReplaceAll(key, "\\", "/"))
	key = strings.TrimLeft(key, "/")
	parts := make([]string, 0)
	for _, part := range strings.Split(key, "/") {
		part = strings.TrimSpace(part)
		if part == "" || part == "." || part == ".." {
			continue
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "/")
}

func escapeKey(key string) string {
	if key == "" {
		return ""
	}

	parts := strings.Split(key, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
