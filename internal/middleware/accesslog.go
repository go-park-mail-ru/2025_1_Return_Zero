package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.LoggerFromContext(r.Context())
		defer logger.Sync()

		start := time.Now()

		sw := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)

		logger.Infow(
			"access",
			"method", r.Method,
			"url", r.URL.String(),
			"ip", r.RemoteAddr,
			"user-agent", r.UserAgent(),
			"status", sw.status,
			"duration", time.Since(start),
		)
	})
}
