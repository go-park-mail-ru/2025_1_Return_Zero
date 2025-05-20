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
	regexp.MustCompile(`^/api/v1/tracks/([0-9]+)$`),
	regexp.MustCompile(`^/api/v1/tracks/([0-9]+)/stream$`),
	regexp.MustCompile(`^/api/v1/tracks/([0-9]+)/like$`),
	regexp.MustCompile(`^/api/v1/tracks/search$`),
	regexp.MustCompile(`^/api/v1/streams/([0-9]+)$`),

	// albums handlers
	regexp.MustCompile(`^/api/v1/albums$`),
	regexp.MustCompile(`^/api/v1/albums/([0-9]+)$`),
	regexp.MustCompile(`^/api/v1/albums/search$`),
	regexp.MustCompile(`^/api/v1/albums/([0-9]+)/tracks$`),
	regexp.MustCompile(`^/api/v1/albums/([0-9]+)/like$`),

	// artists handlers
	regexp.MustCompile(`^/api/v1/artists$`),
	regexp.MustCompile(`^/api/v1/artists/([0-9]+)$`),
	regexp.MustCompile(`^/api/v1/artists/search$`),
	regexp.MustCompile(`^/api/v1/artists/([0-9]+)/tracks$`),
	regexp.MustCompile(`^/api/v1/artists/([0-9]+)/albums$`),
	regexp.MustCompile(`^/api/v1/artists/([0-9]+)/like$`),

	// playlists handlers
	regexp.MustCompile(`^/api/v1/playlists$`),
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)$`),
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)$`), // duplicate for PUT/DELETE, order preserved
	regexp.MustCompile(`^/api/v1/playlists/to-add$`),
	regexp.MustCompile(`^/api/v1/playlists/me$`),
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)/tracks$`),
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)/tracks/([0-9]+)$`),
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)/tracks$`), // GET playlist tracks
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)$`),        // GET playlist by ID
	regexp.MustCompile(`^/api/v1/playlists/([0-9]+)/like$`),
	regexp.MustCompile(`^/api/v1/playlists/search$`),

	// auth handlers
	regexp.MustCompile(`^/api/v1/auth/signup$`),
	regexp.MustCompile(`^/api/v1/auth/login$`),
	regexp.MustCompile(`^/api/v1/auth/logout$`),
	regexp.MustCompile(`^/api/v1/auth/check$`),

	// user handlers
	regexp.MustCompile(`^/api/v1/user/me/avatar$`),
	regexp.MustCompile(`^/api/v1/user/me$`),
	regexp.MustCompile(`^/api/v1/user/me$`), // for DELETE
	regexp.MustCompile(`^/api/v1/user/([^/]+)$`),
	regexp.MustCompile(`^/api/v1/user/me/history$`),
	regexp.MustCompile(`^/api/v1/user/([^/]+)/artists$`),
	regexp.MustCompile(`^/api/v1/user/([^/]+)/tracks$`),
	regexp.MustCompile(`^/api/v1/user/([^/]+)/playlists$`),
	regexp.MustCompile(`^/api/v1/user/me/albums$`),
}

var pathReplacements = []string{
	// tracks handlers
	"/api/v1/tracks",
	"/api/v1/tracks/{id}",
	"/api/v1/tracks/{id}/stream",
	"/api/v1/tracks/{id}/like",
	"/api/v1/tracks/search",
	"/api/v1/streams/{id}",

	// albums handlers
	"/api/v1/albums",
	"/api/v1/albums/{id}",
	"/api/v1/albums/search",
	"/api/v1/albums/{id}/tracks",
	"/api/v1/albums/{id}/like",

	// artists handlers
	"/api/v1/artists",
	"/api/v1/artists/{id}",
	"/api/v1/artists/search",
	"/api/v1/artists/{id}/tracks",
	"/api/v1/artists/{id}/albums",
	"/api/v1/artists/{id}/like",

	// playlists handlers
	"/api/v1/playlists",
	"/api/v1/playlists/{id}",
	"/api/v1/playlists/{id}", // PUT/DELETE
	"/api/v1/playlists/to-add",
	"/api/v1/playlists/me",
	"/api/v1/playlists/{id}/tracks",
	"/api/v1/playlists/{id}/tracks/{trackId}",
	"/api/v1/playlists/{id}/tracks", // GET
	"/api/v1/playlists/{id}",        // GET by ID
	"/api/v1/playlists/{id}/like",
	"/api/v1/playlists/search",

	// auth handlers
	"/api/v1/auth/signup",
	"/api/v1/auth/login",
	"/api/v1/auth/logout",
	"/api/v1/auth/check",

	// user handlers
	"/api/v1/user/me/avatar",
	"/api/v1/user/me", // PUT
	"/api/v1/user/me", // DELETE
	"/api/v1/user/{username}",
	"/api/v1/user/me/history",
	"/api/v1/user/{username}/artists",
	"/api/v1/user/{username}/tracks",
	"/api/v1/user/{username}/playlists",
	"/api/v1/user/me/albums",
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
