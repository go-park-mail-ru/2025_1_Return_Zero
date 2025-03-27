package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestIDKey struct{}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestId)
		logger := LoggerFromContext(ctx).With(zap.String("request_id", requestId))
		ctx = context.WithValue(ctx, LoggerKey{}, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
