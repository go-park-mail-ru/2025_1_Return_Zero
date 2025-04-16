package middleware

import (
	"context"
	"net/http"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

func Auth(userUsecase user.Usecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sessionID := sessionCookie.Value
			user, err := userUsecase.GetUserBySID(r.Context(), sessionID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), helpers.UserContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

