package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerKey    = "logger"
	RequestIDKey = "request_id"
)

func NewZapLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	return logger.Sugar()
}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger := ctx.Value(LoggerKey).(*zap.SugaredLogger)
	return logger
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := NewZapLogger()
		ctx := context.WithValue(r.Context(), LoggerKey, logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
