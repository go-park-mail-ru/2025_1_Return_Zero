package ctxExtractor

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type UserContextKey struct{}

func UserFromContext(ctx context.Context) (int64, bool) {
	user, ok := ctx.Value(UserContextKey{}).(int64)
	if !ok {
		return -1, false
	}
	return user, true
}

func UserToContext(ctx context.Context, user *usecaseModel.User) context.Context {
	return context.WithValue(ctx, UserContextKey{}, user)
}
