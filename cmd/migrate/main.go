package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	appmigrations "github.com/ramdhanrizkij/bytecode-api/migrations"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	action := flags.String("action", "up", "migration action: up, down, refresh, or create")
	databaseURL := flags.String("database", os.Getenv("DATABASE_URL"), "PostgreSQL connection URL")
	directory := flags.String("dir", "migrations", "migration directory used by create")
	name := flags.String("name", "", "migration name used by create")
	if err := flags.Parse(args); err != nil {
		return err
	}

	switch *action {
	case "create":
		baseName, err := appmigrations.Create(*directory, *name)
		if err != nil {
			return err
		}
		fmt.Printf("created %s.go\n", baseName)
		return nil
	case "up", "down", "refresh":
		if strings.TrimSpace(*databaseURL) == "" {
			return errors.New("database URL is required; pass -database or set DATABASE_URL")
		}
	default:
		return fmt.Errorf("unsupported migration action %q", *action)
	}

	db, err := gorm.Open(postgres.Open(*databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	defer sqlDB.Close()

	runner, err := appmigrations.New(db)
	if err != nil {
		return err
	}

	switch *action {
	case "up":
		err = runner.Migrate()
	case "down":
		err = runner.RollbackLast()
	case "refresh":
		err = runner.Refresh()
	}
	if err != nil {
		return err
	}

	fmt.Printf("migration action %s completed successfully\n", *action)
	return nil
}
