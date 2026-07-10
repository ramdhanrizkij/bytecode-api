package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var M000005_seed_default_roles = gormigrate.Migration{
	ID: "000005",
	Migrate: func(tx *gorm.DB) error {
		type Role struct {
			Name        string
			Description string
		}
		roles := []Role{
			{Name: "superadmin", Description: "Super Administrator with full access"},
			{Name: "admin", Description: "Administrator"},
			{Name: "user", Description: "Regular user"},
		}

		return tx.Table("roles").
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoNothing: true,
			}).
			Create(&roles).Error
	},
	Rollback: func(tx *gorm.DB) error {
		type Role struct {
			Name string
		}
		return tx.Table("roles").
			Where("name IN ?", []string{"superadmin", "admin", "user"}).
			Delete(&Role{}).Error
	},
}

func init() {
	registerMigration(5, &M000005_seed_default_roles)
}
