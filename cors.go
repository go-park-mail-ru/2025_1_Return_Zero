package main

import (
	"net/http"
	"strconv"
	"strings"
)

type Cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

func (c *Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		var allowedOrigin string

		for _, o := range c.AllowOrigins {
			if c.AllowCredentials && o == "*" {
				allowedOrigin = o
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
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.AllowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.AllowHeaders, ","))

			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))

			if c.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Add("Vary", "Origin")

			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		if c.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Add("Vary", "Origin")
		next.ServeHTTP(w, r)
		return
	})
}
