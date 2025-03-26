package middleware

import (
	"net/http"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := LoggerFromContext(r.Context())
		defer logger.Sync()
		logger.Infow("access", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
