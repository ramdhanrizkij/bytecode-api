package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	migrationNamePattern     = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	existingMigrationPattern = regexp.MustCompile(`^(\d+)_.*\.go$`)
)

// Create generates the next registered Go migration stub.
func Create(directory, name string) (string, error) {
	name = strings.TrimSpace(name)
	if !migrationNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid migration name %q: use lowercase letters, numbers, and underscores, starting with a letter", name)
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf("create migration directory: %w", err)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("read migration directory: %w", err)
	}

	var maxVersion uint64
	for _, entry := range entries {
		matches := existingMigrationPattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		version, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return "", fmt.Errorf("parse migration version in %q: %w", entry.Name(), err)
		}
		if version > maxVersion {
			maxVersion = version
		}
	}

	version := maxVersion + 1
	baseName := fmt.Sprintf("%06d_%s", version, name)
	path := filepath.Join(directory, baseName+".go")
	variableName := fmt.Sprintf("M%06d_%s", version, name)
	content := fmt.Sprintf(`package migrations

import (
	"errors"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var %[1]s = gormigrate.Migration{
	ID: "%[2]06d",
	Migrate: func(tx *gorm.DB) error {
		// Define migration-local models and apply them with tx.Migrator().
		return errors.New("migration %[2]06d is not implemented")
	},
	Rollback: func(tx *gorm.DB) error {
		return errors.New("migration %[2]06d rollback is not implemented")
	},
}

func init() {
	registerMigration(%[2]d, &%[1]s)
}
`, variableName, version)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("create migration file %q: %w", path, err)
	}
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return "", fmt.Errorf("write migration file %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close migration file %q: %w", path, err)
	}

	return baseName, nil
}
