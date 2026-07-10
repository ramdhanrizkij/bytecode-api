package testhelper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	testcontainers_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	appmigrations "github.com/ramdhanrizkij/bytecode-api/migrations"
)

// TestDB holds the GORM instance and the container reference for integration tests.
type TestDB struct {
	DB        *gorm.DB
	Container *testcontainers_postgres.PostgresContainer
}

// SetupTestDB spawns a PostgreSQL container, runs migrations, and returns a GORM connection.
func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	container, err := testcontainers_postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		testcontainers_postgres.WithDatabase(dbName),
		testcontainers_postgres.WithUsername(dbUser),
		testcontainers_postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=Asia/Jakarta")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	db, err := gorm.Open(gorm_postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	}

	// Run migrations
	runMigrations(t, db)

	return &TestDB{
		DB:        db,
		Container: container,
	}
}

// Teardown closes the database connection and terminates the container.
func (tdb *TestDB) Teardown(t *testing.T) {
	sqlDB, err := tdb.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
	if err := tdb.Container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

// TruncateAll cleans up all tables.
func (tdb *TestDB) TruncateAll(t *testing.T) {
	tables := []string{"role_permissions", "refresh_tokens", "users", "permissions", "roles"}
	for _, table := range tables {
		if err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatalf("failed to truncate table %s: %s", table, err)
		}
	}
}

func runMigrations(t *testing.T, db *gorm.DB) {
	runner, err := appmigrations.New(db)
	if err != nil {
		t.Fatalf("failed to create migration runner: %s", err)
	}
	if err := runner.Migrate(); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}
}
