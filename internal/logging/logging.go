package logging

import "go.uber.org/zap"

// New returns a zap logger configured for production or development based on env.
func New(env string) *zap.Logger {
	if env == "production" {
		logger, _ := zap.NewProduction()
		return logger
	}
	logger, _ := zap.NewDevelopment()
	return logger
}
