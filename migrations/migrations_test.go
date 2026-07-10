package migrations

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationDefinitions(t *testing.T) {
	definitions, err := migrationDefinitions()

	require.NoError(t, err)
	require.NotEmpty(t, definitions)
	for index, definition := range definitions {
		expectedVersion := int64(index + 1)
		require.Equal(t, expectedVersion, definition.version)
		require.Equal(t, fmt.Sprintf("%06d", expectedVersion), definition.migration.ID)
	}
}

func TestCreate(t *testing.T) {
	directory := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(directory, "000004_existing.go"), []byte("package migrations"), 0o644))

	baseName, err := Create(directory, "create_orders_table")

	require.NoError(t, err)
	require.Equal(t, "000005_create_orders_table", baseName)
	createdPath := filepath.Join(directory, baseName+".go")
	require.FileExists(t, createdPath)
	content, err := os.ReadFile(createdPath)
	require.NoError(t, err)
	_, err = format.Source(content)
	require.NoError(t, err)
	require.Contains(t, string(content), "var M000005_create_orders_table = gormigrate.Migration{")
	require.Contains(t, string(content), "registerMigration(5, &M000005_create_orders_table)")
}
