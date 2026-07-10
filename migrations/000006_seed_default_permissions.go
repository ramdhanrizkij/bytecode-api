package migrations

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var migration000006PermissionNames = []string{
	"roles.view",
	"roles.create",
	"roles.edit",
	"roles.delete",
	"roles.assign-permission",
	"roles.remove-permission",
	"permissions.view",
	"permissions.create",
	"permissions.edit",
	"permissions.delete",
	"users.view",
	"users.create",
	"users.edit",
	"users.delete",
}

var M000006_seed_default_permissions = gormigrate.Migration{
	ID: "000006",
	Migrate: func(tx *gorm.DB) error {
		type Permission struct {
			Name        string
			Description string
		}
		permissions := []Permission{
			{Name: "roles.view", Description: "View roles list and details"},
			{Name: "roles.create", Description: "Create new roles"},
			{Name: "roles.edit", Description: "Update existing roles"},
			{Name: "roles.delete", Description: "Delete roles"},
			{Name: "roles.assign-permission", Description: "Assign permissions to roles"},
			{Name: "roles.remove-permission", Description: "Remove permissions from roles"},
			{Name: "permissions.view", Description: "View permissions list and details"},
			{Name: "permissions.create", Description: "Create new permissions"},
			{Name: "permissions.edit", Description: "Update existing permissions"},
			{Name: "permissions.delete", Description: "Delete permissions"},
			{Name: "users.view", Description: "View users list and details"},
			{Name: "users.create", Description: "Create new users"},
			{Name: "users.edit", Description: "Update existing users"},
			{Name: "users.delete", Description: "Delete users"},
		}
		if err := tx.Table("permissions").
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoNothing: true,
			}).
			Create(&permissions).Error; err != nil {
			return fmt.Errorf("seed permissions: %w", err)
		}

		var superadmin struct {
			ID uuid.UUID
		}
		if err := tx.Table("roles").Select("id").Where("name = ?", "superadmin").Take(&superadmin).Error; err != nil {
			return fmt.Errorf("find superadmin role: %w", err)
		}

		var seededPermissions []struct {
			ID uuid.UUID
		}
		if err := tx.Table("permissions").
			Select("id").
			Where("name IN ?", migration000006PermissionNames).
			Find(&seededPermissions).Error; err != nil {
			return fmt.Errorf("find seeded permissions: %w", err)
		}

		type RolePermission struct {
			RoleID       uuid.UUID
			PermissionID uuid.UUID
		}
		assignments := make([]RolePermission, 0, len(seededPermissions))
		for _, permission := range seededPermissions {
			assignments = append(assignments, RolePermission{
				RoleID:       superadmin.ID,
				PermissionID: permission.ID,
			})
		}

		return tx.Table("role_permissions").
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "permission_id"}},
				DoNothing: true,
			}).
			Create(&assignments).Error
	},
	Rollback: func(tx *gorm.DB) error {
		type Permission struct {
			Name string
		}
		return tx.Table("permissions").
			Where("name IN ?", migration000006PermissionNames).
			Delete(&Permission{}).Error
	},
}

func init() {
	registerMigration(6, &M000006_seed_default_permissions)
}
