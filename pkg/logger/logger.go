package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger builds a zap logger with the provided service name and log level.
func NewLogger(serviceName, level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if level != "" {
		lvl := zapcore.InfoLevel
		if err := lvl.Set(strings.ToLower(level)); err != nil {
			return nil, fmt.Errorf("parse log level: %w", err)
		}
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}
	logger, err := cfg.Build(zap.AddCaller(), zap.Fields(zap.String("service", serviceName)))
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}
	return logger, nil
}
