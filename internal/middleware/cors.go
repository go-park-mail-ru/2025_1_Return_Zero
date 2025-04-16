package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
)

func CorsMiddleware(cfg config.Cors) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			var allowedOrigin string

			for _, o := range cfg.AllowedOrigins {
				if cfg.AllowCredentials && o == "*" {
					allowedOrigin = origin
					break
				}
				if o == origin {
					allowedOrigin = o
					break
				}
			}

			if allowedOrigin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))

				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))

				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				w.Header().Add("Vary", "Origin")

				w.WriteHeader(http.StatusNoContent)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Add("Vary", "Origin")
			next.ServeHTTP(w, r)
		})
	}
}
