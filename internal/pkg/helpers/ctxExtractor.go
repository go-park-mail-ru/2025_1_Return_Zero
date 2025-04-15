package helpers

import (
	"context"

	"go.uber.org/zap"
)

type LoggerKey struct{}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger := ctx.Value(LoggerKey{}).(*zap.SugaredLogger)
	return logger
}
