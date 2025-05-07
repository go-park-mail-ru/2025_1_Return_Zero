package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customResponseWriter"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/metrics"
)

var pathPatterns = []*regexp.Regexp{
	// tracks handlers
	regexp.MustCompile(`^/api/v1/tracks$`),
	regexp.MustCompile(`^/api/v1/tracks/(\d+)$`),
	regexp.MustCompile(`^/api/v1/tracks/(\d+)/stream$`),
	regexp.MustCompile(`^/api/v1/streams/(\d+)$`),

	// albums handlers
	regexp.MustCompile(`^/api/v1/albums$`),
	regexp.MustCompile(`^/api/v1/albums/(\d+)$`),
	regexp.MustCompile(`^/api/v1/albums/(\d+)/tracks$`),

	// artists handlers
	regexp.MustCompile(`^/api/v1/artists$`),
	regexp.MustCompile(`^/api/v1/artists/(\d+)$`),
	regexp.MustCompile(`^/api/v1/artists/(\d+)/tracks$`),
	regexp.MustCompile(`^/api/v1/artists/(\d+)/albums$`),

	// auth handlers
	regexp.MustCompile(`^/api/v1/auth/signup$`),
	regexp.MustCompile(`^/api/v1/auth/login$`),
	regexp.MustCompile(`^/api/v1/auth/logout$`),
	regexp.MustCompile(`^/api/v1/auth/check$`),

	// user handlers
	regexp.MustCompile(`^/api/v1/user/me/avatar$`),
	regexp.MustCompile(`^/api/v1/user/me$`),
	regexp.MustCompile(`^/api/v1/user/([^/]+)$`),
	regexp.MustCompile(`^/api/v1/user/me/history$`),
}

var pathReplacements = []string{
	// tracks handlers
	"/api/v1/tracks",
	"/api/v1/tracks/{id}",
	"/api/v1/tracks/{id}/stream",
	"/api/v1/streams/{id}",

	// albums handlers
	"/api/v1/albums",
	"/api/v1/albums/{id}",
	"/api/v1/albums/{id}/tracks",

	// artists handlers
	"/api/v1/artists",
	"/api/v1/artists/{id}",
	"/api/v1/artists/{id}/tracks",
	"/api/v1/artists/{id}/albums",

	// auth handlers
	"/api/v1/auth/signup",
	"/api/v1/auth/login",
	"/api/v1/auth/logout",
	"/api/v1/auth/check",

	// user handlers
	"/api/v1/user/me/avatar",
	"/api/v1/user/me",
	"/api/v1/user/{username}",
	"/api/v1/user/me/history",
}

func ExtractPath(path string) string {
	for i, pattern := range pathPatterns {
		if pattern.MatchString(path) {
			return pathReplacements[i]
		}
	}

	return path
}

func MetricsMiddleware(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := customResponseWriter.NewResponseWriter(w)

			next.ServeHTTP(rec, r)

			duration := time.Since(start).Seconds()

			method := r.Method
			path := r.URL.Path
			status := strconv.Itoa(rec.StatusCode)

			normalizedPath := ExtractPath(path)

			m.HTTPTotalNumberOfRequests.WithLabelValues(method, normalizedPath, status).Inc()
			m.HTTPRequestDuration.WithLabelValues(method, normalizedPath).Observe(duration)

		})
	}
}
