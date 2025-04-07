package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

type UserContextKey struct{}

func Auth(userUsecase user.Usecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sessionID := sessionCookie.Value
			user, err := userUsecase.GetUserBySID(sessionID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey{}, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(ctx context.Context) (*user.Usecase, bool) {
	user, ok := ctx.Value(UserContextKey{}).(*user.Usecase)
	if !ok {
		return nil, false
	}
	return user, true
}