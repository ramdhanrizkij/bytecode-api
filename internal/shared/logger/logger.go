package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

type ZapLogger struct {
	base *zap.Logger
}

func NewZapLogger(base *zap.Logger) *ZapLogger {
	return &ZapLogger{base: base}
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.base.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.base.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.base.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.base.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.base.Fatal(msg, fields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{base: l.base.With(fields...)}
}

func (l *ZapLogger) Sync() error {
	return l.base.Sync()
}

func (l *ZapLogger) Base() *zap.Logger {
	return l.base
}
