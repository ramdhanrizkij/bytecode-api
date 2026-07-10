package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global structured logger instance.
var Log *zap.Logger

// Sugar is the global sugared (printf-style) logger instance.
var Sugar *zap.SugaredLogger

// NewLogger builds and returns a *zap.Logger configured for the given level.
//
// Behaviour by level:
//   - "debug" → development mode, console encoding, debug-and-above enabled.
//   - "info" / "warn" / "error" → production mode, JSON encoding.
//
// Caller information and stack traces at error level are always included.
func NewLogger(level string) (*zap.Logger, error) {
	zapLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	var cfg zap.Config

	if level == "debug" {
		// Development config: human-readable console output.
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// Production config: structured JSON output.
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
	}

	// Always output to stdout.
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	// Set the resolved level.
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	// Always include caller information.
	cfg.DisableCaller = false

	// Add stack traces for error-level and above.
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return logger, nil
}

// NewSugaredLogger wraps a *zap.Logger in a SugaredLogger for printf-style
// logging convenience. Use sparingly in hot paths — prefer the structured
// *zap.Logger for performance-sensitive code.
func NewSugaredLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

// InitGlobal initialises the package-level Log and Sugar globals.
// Call this once during application startup, before any logging occurs.
func InitGlobal(level string) error {
	l, err := NewLogger(level)
	if err != nil {
		return fmt.Errorf("failed to initialise global logger: %w", err)
	}

	Log = l
	Sugar = NewSugaredLogger(l)

	// Redirect any stray uses of the zap global logger to our instance.
	zap.ReplaceGlobals(l)

	return nil
}

// parseLevel converts a string level name to a zapcore.Level value.
func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level %q: must be one of debug|info|warn|error", level)
	}
}
