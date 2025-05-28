package middleware

import (
	"context"
	"net/http"

	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func Auth(authClient *authProto.AuthServiceClient, userClient *userProto.UserServiceClient) func(http.Handler) http.Handler {
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
			userFront, err := (*userClient).GetUserByID(r.Context(), model.UserIDFromUsecaseToProtoUser(userID))
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if userFront.Username == "admin" && userFront.Email == "admin@admin.ru" {
				ctx = context.WithValue(ctx, ctxExtractor.AdminContextKey{}, userID)
			}
			labelIDProto, err := (*userClient).GetLabelIDByUserID(r.Context(), model.UserIDFromUsecaseToProtoUser(userID))
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			labelID := model.LabelIDFromProtoToUsecase(labelIDProto)
			ctx = context.WithValue(ctx, ctxExtractor.LabelContextKey{}, labelID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
