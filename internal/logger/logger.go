package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func InitializeLogger() error {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevel()
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}
