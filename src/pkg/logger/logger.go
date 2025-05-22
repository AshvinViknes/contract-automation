package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Initialize initializes the logger instance once
func Initialize(level string, isDevelopment bool) *zap.Logger {
	once.Do(func() {
		var config zap.Config

		if !isDevelopment {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Set log level
		logLevel, parseErr := zapcore.ParseLevel(level)
		if parseErr != nil {
			logLevel = zapcore.DebugLevel
		}
		config.Level = zap.NewAtomicLevelAt(logLevel)

		var err error
		log, err = config.Build()
		if err != nil {
			panic(err)
		}
	})
	return log
}

// GetLogger returns the logger instance, initializing it if necessary
func GetLogger() *zap.Logger {
	if log == nil {
		panic("Logger not initialized. Call Initialize() first")
	}
	return log
}

// Info logs info level message
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs error level message
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs fatal level message and exits
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Debug logs debug level message
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}
