package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000011_create_addresses_table = gormigrate.Migration{
	ID: "000011",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type Address struct {
			ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			UserID        uuid.UUID `gorm:"type:uuid;not null"`
			User          User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
			Label         *string   `gorm:"type:varchar(50)"`
			RecipientName *string   `gorm:"type:varchar(100)"`
			Phone         *string   `gorm:"type:varchar(20)"`
			Address       *string   `gorm:"type:text"`
			City          *string   `gorm:"type:varchar(100)"`
			Province      *string   `gorm:"type:varchar(100)"`
			PostalCode    *string   `gorm:"type:varchar(100)"`
			CreatedAt     time.Time `gorm:"type:timestamptz;default:now()"`
			UpdatedAt     time.Time `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&Address{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("addresses")
	},
}

func init() {
	registerMigration(11, &M000011_create_addresses_table)
}
