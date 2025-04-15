package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
)

type RequestIDKey struct{}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestId)
		logger := helpers.LoggerFromContext(ctx).With(zap.String("request_id", requestId))
		ctx = context.WithValue(ctx, helpers.LoggerKey{}, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
