package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000004_create_role_permissions_table = gormigrate.Migration{
	ID: "000004",
	Migrate: func(tx *gorm.DB) error {
		type Role struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type Permission struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type RolePermission struct {
			ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			RoleID       uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_role_permissions_role_id_permission_id,priority:1"`
			Role         Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
			PermissionID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_role_permissions_role_id_permission_id,priority:2"`
			Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
			CreatedAt    time.Time  `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&RolePermission{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("role_permissions")
	},
}

func init() {
	registerMigration(4, &M000004_create_role_permissions_table)
}
