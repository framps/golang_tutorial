package game

import (
	"go.uber.org/zap"
)

var (
	// Log -
	Log *zap.SugaredLogger
)

// NewLog -
func NewLog() {
	logger, _ := zap.NewProduction()
	Log = logger.Sugar()
}
