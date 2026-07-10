package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var M000010_create_products_table = gormigrate.Migration{
	ID: "000010",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Exec(`
			DO $$
			BEGIN
				CREATE TYPE product_status AS ENUM ('DRAFT', 'ACTIVE', 'INACTIVE', 'ARCHIVED');
			EXCEPTION
				WHEN duplicate_object THEN NULL;
			END
			$$;
		`).Error; err != nil {
			return err
		}

		type ProductStatus string
		type Category struct {
			ID uuid.UUID `gorm:"type:uuid;primaryKey"`
		}
		type Product struct {
			ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
			CategoryID  uuid.UUID     `gorm:"type:uuid;not null"`
			Category    Category      `gorm:"foreignKey:CategoryID;constraint:OnDelete:NO ACTION"`
			SKU         string        `gorm:"type:varchar(100);uniqueIndex;not null"`
			Name        *string       `gorm:"type:varchar(100)"`
			Description *string       `gorm:"type:text"`
			Price       float64       `gorm:"type:numeric(18,2);not null"`
			Stock       int           `gorm:"not null;default:0"`
			Weight      *float64      `gorm:"type:numeric(10,2)"`
			Status      ProductStatus `gorm:"type:product_status;not null;default:'DRAFT'"`
			CreatedAt   time.Time     `gorm:"type:timestamptz;default:now()"`
			UpdatedAt   time.Time     `gorm:"type:timestamptz;default:now()"`
		}

		return tx.Migrator().AutoMigrate(&Product{})
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable("products"); err != nil {
			return err
		}
		return tx.Exec("DROP TYPE IF EXISTS product_status").Error
	},
}

func init() {
	registerMigration(10, &M000010_create_products_table)
}
