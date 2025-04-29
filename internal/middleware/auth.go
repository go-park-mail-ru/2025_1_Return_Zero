package middleware

import (
	"context"
	"net/http"

	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func Auth(authClient *authProto.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sessionID := sessionCookie.Value
			userIDProto, err := (*authClient).GetSession(r.Context(), model.SessionIDFromUsecaseToProto(sessionID))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			userID := model.UserIDFromProtoToUsecase(userIDProto)
			ctx := context.WithValue(r.Context(), ctxExtractor.UserContextKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
