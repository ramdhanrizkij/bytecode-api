package bootstrap

import (
	"github.com/ramdhanrizki/bytecode-api/configs"
	platformQueue "github.com/ramdhanrizki/bytecode-api/internal/platform/queue"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
	"gorm.io/gorm"
)

func NewQueue(cfg configs.Config, logger sharedLogger.Logger, db *gorm.DB) (sharedQueue.Publisher, error) {
	return platformQueue.NewPostgresQueue(db, logger, cfg.Worker.Concurrency)
}

func NewQueueConsumer(cfg configs.Config, logger sharedLogger.Logger, db *gorm.DB) (sharedQueue.Consumer, error) {
	return platformQueue.NewPostgresQueue(db, logger, cfg.Worker.Concurrency)
}
