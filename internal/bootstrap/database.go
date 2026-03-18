package bootstrap

import (
	"gorm.io/gorm"

	"github.com/ramdhanrizki/bytecode-api/configs"
	"github.com/ramdhanrizki/bytecode-api/internal/platform/persistence/postgres"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

func NewDatabase(cfg configs.Config, logger sharedLogger.Logger) (*gorm.DB, error) {
	return postgres.NewGormDB(cfg.DB, logger)
}
