package migrations

import (
	"errors"
	"fmt"
	"sort"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	// HistoryTable stores migration IDs managed by Gormigrate.
	HistoryTable = "gormigrate_migrations"

	legacyHistoryTable = "schema_migrations"
)

type registeredMigration struct {
	version   int64
	migration *gormigrate.Migration
}

var registry []registeredMigration

// Runner applies registered Go migrations through Gormigrate.
type Runner struct {
	db          *gorm.DB
	definitions []registeredMigration
	migrations  []*gormigrate.Migration
}

// registerMigration adds a migration to the registry. Definitions are sorted
// and validated by New, so Go file initialization order does not matter.
func registerMigration(version int64, migration *gormigrate.Migration) {
	registry = append(registry, registeredMigration{
		version:   version,
		migration: migration,
	})
}

// New creates a runner from the registered Go migration definitions.
func New(db *gorm.DB) (*Runner, error) {
	if db == nil {
		return nil, errors.New("migrations: database is required")
	}

	definitions, err := migrationDefinitions()
	if err != nil {
		return nil, err
	}

	migrationList := make([]*gormigrate.Migration, 0, len(definitions))
	for _, definition := range definitions {
		migrationList = append(migrationList, definition.migration)
	}

	return &Runner{
		db:          db,
		definitions: definitions,
		migrations:  migrationList,
	}, nil
}

// Migrate applies every migration that has not run yet.
func (r *Runner) Migrate() error {
	if err := r.bootstrapLegacyHistory(); err != nil {
		return err
	}
	if err := r.gormigrate().Migrate(); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// RollbackLast rolls back the most recently applied migration.
func (r *Runner) RollbackLast() error {
	if err := r.bootstrapLegacyHistory(); err != nil {
		return err
	}
	if err := r.gormigrate().RollbackLast(); err != nil {
		return fmt.Errorf("rollback last migration: %w", err)
	}
	return nil
}

// Refresh rolls back every applied migration and then applies all migrations again.
func (r *Runner) Refresh() error {
	if err := r.bootstrapLegacyHistory(); err != nil {
		return err
	}

	for {
		err := r.gormigrate().RollbackLast()
		if errors.Is(err, gormigrate.ErrNoRunMigration) {
			break
		}
		if err != nil {
			return fmt.Errorf("refresh migrations: %w", err)
		}
	}

	if err := r.gormigrate().Migrate(); err != nil {
		return fmt.Errorf("reapply migrations: %w", err)
	}
	return nil
}

func (r *Runner) gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(r.db, migrationOptions(), r.migrations)
}

func migrationOptions() *gormigrate.Options {
	return &gormigrate.Options{
		TableName:                 HistoryTable,
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}
}

func migrationDefinitions() ([]registeredMigration, error) {
	definitions := append([]registeredMigration(nil), registry...)
	sort.Slice(definitions, func(i, j int) bool {
		return definitions[i].version < definitions[j].version
	})
	if len(definitions) == 0 {
		return nil, errors.New("migrations: no Go migrations registered")
	}

	versions := make(map[int64]struct{}, len(definitions))
	ids := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		if definition.version <= 0 {
			return nil, fmt.Errorf("migrations: invalid version %d", definition.version)
		}
		if _, exists := versions[definition.version]; exists {
			return nil, fmt.Errorf("migrations: duplicate version %d", definition.version)
		}
		versions[definition.version] = struct{}{}

		if definition.migration == nil {
			return nil, fmt.Errorf("migrations: version %d has a nil definition", definition.version)
		}
		if definition.migration.ID == "" {
			return nil, fmt.Errorf("migrations: version %d has an empty ID", definition.version)
		}
		if _, exists := ids[definition.migration.ID]; exists {
			return nil, fmt.Errorf("migrations: duplicate ID %q", definition.migration.ID)
		}
		ids[definition.migration.ID] = struct{}{}
		if definition.migration.Migrate == nil || definition.migration.Rollback == nil {
			return nil, fmt.Errorf("migrations: %s must define migrate and rollback functions", definition.migration.ID)
		}
	}

	return definitions, nil
}

// bootstrapLegacyHistory allows a database previously managed by golang-migrate
// to move to Gormigrate without replaying migrations that already ran.
func (r *Runner) bootstrapLegacyHistory() error {
	if r.db.Migrator().HasTable(HistoryTable) || !r.db.Migrator().HasTable(legacyHistoryTable) {
		return nil
	}

	var state struct {
		Version int64
		Dirty   bool
	}
	result := r.db.Table(legacyHistoryTable).Select("version", "dirty").Take(&state)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if result.Error != nil {
		return fmt.Errorf("read legacy migration state: %w", result.Error)
	}
	if state.Dirty {
		return fmt.Errorf("legacy migration version %d is dirty; repair it before switching to Gormigrate", state.Version)
	}
	if state.Version <= 0 {
		return nil
	}

	legacyVersionFound := false
	applied := make([]*gormigrate.Migration, 0, len(r.definitions))
	for _, definition := range r.definitions {
		if definition.version == state.Version {
			legacyVersionFound = true
		}
		if definition.version > state.Version {
			continue
		}
		applied = append(applied, &gormigrate.Migration{
			ID: definition.migration.ID,
			Migrate: func(*gorm.DB) error {
				return nil
			},
		})
	}
	if !legacyVersionFound {
		return fmt.Errorf("legacy migration version %d has no matching Go migration", state.Version)
	}

	if err := gormigrate.New(r.db, migrationOptions(), applied).Migrate(); err != nil {
		return fmt.Errorf("import legacy migration history at version %d: %w", state.Version, err)
	}
	return nil
}
