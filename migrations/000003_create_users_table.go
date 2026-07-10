package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000003_create_users_table = gormigrate.Migration{
	ID: "000003",
	Migrate: func(tx *gorm.DB) error {
		type Role struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type User struct {
			ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			Name      string     `gorm:"type:varchar(100);not null"`
			Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
			Password  string     `gorm:"type:varchar(255);not null"`
			RoleID    *uuid.UUID `gorm:"type:uuid"`
			Role      *Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"`
			IsActive  bool       `gorm:"default:true"`
			CreatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
			UpdatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&User{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("users")
	},
}

func init() {
	registerMigration(3, &M000003_create_users_table)
}
