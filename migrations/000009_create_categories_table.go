package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000009_create_categories_table = gormigrate.Migration{
	ID: "000009",
	Migrate: func(tx *gorm.DB) error {
		type Category struct {
			ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			Name        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
			Description *string   `gorm:"type:text"`
			CreatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
			UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&Category{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("categories")
	},
}

func init() {
	registerMigration(9, &M000009_create_categories_table)
}
