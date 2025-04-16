package helpers

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"go.uber.org/zap"
)

type LoggerKey struct{}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger := ctx.Value(LoggerKey{}).(*zap.SugaredLogger)
	return logger
}

func LoggerToContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, logger)
}

type UserContextKey struct{}

func UserFromContext(ctx context.Context) (*usecaseModel.User, bool) {
	user, ok := ctx.Value(UserContextKey{}).(*usecaseModel.User)
	if !ok {
		return nil, false
	}
	return user, true
}
