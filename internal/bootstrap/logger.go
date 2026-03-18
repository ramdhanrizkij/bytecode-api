package bootstrap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ramdhanrizki/bytecode-api/configs"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

func NewLogger(cfg configs.Config) (sharedLogger.Logger, error) {
	level := zapcore.InfoLevel
	if cfg.App.Env == "development" {
		level = zapcore.DebugLevel
	}

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.App.Env == "development",
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	base, err := zapCfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}

	return sharedLogger.NewZapLogger(base).With(
		zap.String("app_name", cfg.App.Name),
		zap.String("app_env", cfg.App.Env),
	), nil
}
