package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/json"
)

func generateCSRFToken(tokenLength int) (string, error) {
	b := make([]byte, tokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

func CSRFMiddleware(cfg config.CSRFConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				var token string
				cookie, err := r.Cookie(cfg.CSRFCookieName)
				if err != nil || cookie.Value == "" {
					newToken, err := generateCSRFToken(cfg.CSRFTokenLength)
					if err != nil {
						json.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate CSRF token", nil)
						return
					}
					token = newToken
					cookie = &http.Cookie{
						Name:     cfg.CSRFCookieName,
						Value:    token,
						Path:     "/",
						HttpOnly: true,
						Secure:   false,
					}
					http.SetCookie(w, cookie)
				} else {
					token = cookie.Value
				}

				w.Header().Set(cfg.CSRFHeaderName, token)
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(cfg.CSRFCookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					json.WriteErrorResponse(w, http.StatusForbidden, "CSRF token missing", nil)
					return
				}
				json.WriteErrorResponse(w, http.StatusInternalServerError, "Error reading CSRF cookie", nil)
				return
			}
			if cookie.Value == "" {
				json.WriteErrorResponse(w, http.StatusForbidden, "CSRF token missing", nil)
				return
			}

			token := r.Header.Get(cfg.CSRFHeaderName)
			if token == "" || token != cookie.Value {
				json.WriteErrorResponse(w, http.StatusForbidden, "Invalid CSRF token", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
