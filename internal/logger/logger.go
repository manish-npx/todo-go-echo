package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once     sync.Once
	instance *zap.Logger
)

// Init initializes a singleton zap logger.
func Init() error {
	var initErr error
	once.Do(func() {
		instance, initErr = zap.NewProduction()
	})
	return initErr
}

// L returns the logger instance.
func L() *zap.Logger {
	if instance == nil {
		// Fallback logger for safety in tests or early initialization.
		instance = zap.NewNop()
	}
	return instance
}

// Sync flushes buffered logs.
func Sync() {
	_ = L().Sync()
}
