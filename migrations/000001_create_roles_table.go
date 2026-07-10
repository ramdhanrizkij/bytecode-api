package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000001_create_roles_table = gormigrate.Migration{
	ID: "000001",
	Migrate: func(tx *gorm.DB) error {
		type Role struct {
			ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			Name        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
			Description *string   `gorm:"type:text"`
			GuardName   string    `gorm:"type:varchar(50);default:'api'"`
			CreatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
			UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&Role{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("roles")
	},
}

func init() {
	registerMigration(1, &M000001_create_roles_table)
}
