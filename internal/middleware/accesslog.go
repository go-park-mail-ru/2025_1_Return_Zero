package middleware

import (
	"bufio"
	"net"
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

// Hijack implements the http.Hijacker interface to support WebSocket connections
func (w *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.LoggerFromContext(r.Context())
		defer logger.Sync()

		if websocketRequest := r.Header.Get("Upgrade"); websocketRequest == "websocket" {
			next.ServeHTTP(w, r)
			return
		}

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
