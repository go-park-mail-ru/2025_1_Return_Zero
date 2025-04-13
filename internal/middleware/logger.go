package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerKey struct{}

func NewZapLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger := ctx.Value(LoggerKey{}).(*zap.SugaredLogger)
	return logger
}

func Logger(next http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), LoggerKey{}, logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Logger(next, logger)
	}
}
