package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global logger
	Log  *zap.Logger
	once sync.Once
)

// InitLogger initializes the logger with the given log level
func InitLogger(level string) *zap.Logger {
	once.Do(func() {
		var zapLevel zapcore.Level
		switch level {
		case "debug":
			zapLevel = zapcore.DebugLevel
		case "info":
			zapLevel = zapcore.InfoLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		default:
			zapLevel = zapcore.InfoLevel
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Create console encoder
		consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)

		// Create console core
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)

		// Create logger
		Log = zap.New(consoleCore,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})

	return Log
}

// GetLogger returns a singleton instance of zap.Logger
func GetLogger() *zap.Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error
		Log, err = config.Build()
		if err != nil {
			panic(err)
		}
	})

	return Log
}

// With adds structured context to the logger
func With(fields ...zapcore.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Debug logs a message at debug level
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs a message at info level
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a message at warn level
func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs a message at error level
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// WithError adds an error field to the logger
func WithError(err error) *zap.Logger {
	return GetLogger().With(zap.Error(err))
}

// WithField adds a field to the logger
func WithField(key string, value interface{}) *zap.Logger {
	return GetLogger().With(zap.Any(key, value))
}

// WithFields adds multiple fields to the logger
func WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return GetLogger().With(zapFields...)
}
