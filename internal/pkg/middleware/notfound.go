package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := NewZapLogger()
		logger.Warnw("not found", "method", r.Method, "url", r.URL.String())

		helpers.WriteJSONError(w, http.StatusNotFound, "Resource not found")
	})
}
