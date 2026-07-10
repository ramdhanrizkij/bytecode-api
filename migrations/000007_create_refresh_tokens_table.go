package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000007_create_refresh_tokens_table = gormigrate.Migration{
	ID: "000007",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type RefreshToken struct {
			ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			UserID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_refresh_tokens_user_id"`
			User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
			TokenHash string     `gorm:"type:varchar(64);uniqueIndex;not null"`
			ExpiresAt time.Time  `gorm:"type:timestamptz;not null;index:idx_refresh_tokens_expires_at"`
			RevokedAt *time.Time `gorm:"type:timestamptz"`
			CreatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
			UpdatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&RefreshToken{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("refresh_tokens")
	},
}

func init() {
	registerMigration(7, &M000007_create_refresh_tokens_table)
}
