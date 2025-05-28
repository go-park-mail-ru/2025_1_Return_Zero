package ctxExtractor

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type UserContextKey struct{}
type LabelContextKey struct{}
type AdminContextKey struct{}

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

func LabelFromContext(ctx context.Context) (int64, bool) {
	label, ok := ctx.Value(LabelContextKey{}).(int64)
	if !ok {
		return -1, false
	}
	return label, true
}

func AdminFromContext(ctx context.Context) bool {
	_, ok := ctx.Value(AdminContextKey{}).(int64)
	return ok
}
