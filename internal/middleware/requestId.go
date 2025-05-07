package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
)

type RequestIDKey struct{}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestId)
		logger := loggerPkg.LoggerFromContext(ctx).With(zap.String("request_id", requestId))
		ctx = loggerPkg.LoggerToContext(ctx, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
