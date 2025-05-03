package ctxExtractor

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type UserContextKey struct{}

func UserFromContext(ctx context.Context) (*usecaseModel.User, bool) {
	user, ok := ctx.Value(UserContextKey{}).(*usecaseModel.User)
	if !ok {
		return nil, false
	}
	return user, true
}

func UserToContext(ctx context.Context, user *usecaseModel.User) context.Context {
	return context.WithValue(ctx, UserContextKey{}, user)
}
