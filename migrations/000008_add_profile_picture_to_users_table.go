package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var M000008_add_profile_picture_to_users_table = gormigrate.Migration{
	ID: "000008",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			ProfilePicture *string `gorm:"type:text"`
		}
		return tx.Migrator().AddColumn(&User{}, "ProfilePicture")
	},
	Rollback: func(tx *gorm.DB) error {
		type User struct {
			ProfilePicture *string `gorm:"type:text"`
		}
		return tx.Migrator().DropColumn(&User{}, "ProfilePicture")
	},
}

func init() {
	registerMigration(8, &M000008_add_profile_picture_to_users_table)
}
