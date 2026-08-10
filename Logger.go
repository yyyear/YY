package YY

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger      = buildLogger()
)

func buildLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = atomicLevel
	newLogger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		return zap.NewNop()
	}
	return newLogger
}

func SetLevel(level zapcore.Level) {
	atomicLevel.SetLevel(level)
}

// Fields is retained for compatibility with the original logging API.
type Fields map[string]interface{}

func Debug(msg string, info Fields) { loggerMessage(zapcore.DebugLevel, msg, info) }
func Info(msg string, info Fields)  { loggerMessage(zapcore.InfoLevel, msg, info) }
func Warn(msg string, info Fields)  { loggerMessage(zapcore.WarnLevel, msg, info) }
func Error(msg string, info Fields) { loggerMessage(zapcore.ErrorLevel, msg, info) }

func Sync() error {
	return logger.Sync()
}

func loggerMessage(level zapcore.Level, msg string, fields ...Fields) {
	checked := logger.Check(level, msg)
	if checked == nil {
		return
	}
	fieldCount := 0
	for _, field := range fields {
		fieldCount += len(field)
	}
	zapFields := make([]zap.Field, 0, fieldCount)
	for _, field := range fields {
		for key, value := range field {
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}
	checked.Write(zapFields...)
}
