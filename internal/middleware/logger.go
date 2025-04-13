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

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := NewZapLogger()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), LoggerKey{}, logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
