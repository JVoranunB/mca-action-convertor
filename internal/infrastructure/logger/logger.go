// internal/infrastructure/logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global logger instance
	Log *zap.Logger
)

// Initialize sets up the logger
func Initialize(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return nil, err
	}

	// Set the global logger
	zap.ReplaceGlobals(Log)

	return Log, nil
}

// GetLogger returns a named logger for a specific component
func GetLogger(component string) *zap.Logger {
	return Log.With(zap.String("component", component))
}

// Sugar returns a sugared logger
func Sugar() *zap.SugaredLogger {
	return Log.Sugar()
}
