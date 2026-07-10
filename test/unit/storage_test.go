package unit

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/storage"
)

func TestLocalStoragePutAndDelete(t *testing.T) {
	tempDir := t.TempDir()
	provider, err := storage.NewProvider(&config.StorageConfig{
		Provider:      storage.ProviderLocal,
		LocalPath:     tempDir,
		BaseURL:       "/files",
		DefaultBucket: "uploads",
		BucketsRaw:    "uploads,images",
	}, zap.NewNop())
	assert.NoError(t, err)

	object, err := provider.Put(context.Background(), &storage.PutObjectRequest{
		Bucket:      "images",
		Key:         "avatars/user-1.txt",
		Reader:      strings.NewReader("hello storage"),
		Size:        int64(len("hello storage")),
		ContentType: "text/plain",
	})
	assert.NoError(t, err)
	assert.Equal(t, "images", object.Bucket)
	assert.Equal(t, "avatars/user-1.txt", object.Key)
	assert.Equal(t, "/files/images/avatars/user-1.txt", object.URL)

	savedPath := filepath.Join(tempDir, "images", "avatars", "user-1.txt")
	data, err := os.ReadFile(savedPath)
	assert.NoError(t, err)
	assert.Equal(t, "hello storage", string(data))

	err = provider.Delete(context.Background(), "images", "avatars/user-1.txt")
	assert.NoError(t, err)

	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err))
}

func TestStorageConfigBucketsIncludesDefaultBucket(t *testing.T) {
	cfg := config.StorageConfig{
		DefaultBucket: "uploads",
		BucketsRaw:    "images,documents",
	}

	buckets := cfg.Buckets()

	assert.Equal(t, []string{"images", "documents", "uploads"}, buckets)
}
