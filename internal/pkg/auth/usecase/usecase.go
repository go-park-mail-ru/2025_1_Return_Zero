package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
)

func NewAuthUsecase(authRepo auth.Repository) auth.Usecase {
	return AuthUsecase{authRepo: authRepo}
}

type AuthUsecase struct {
	authRepo auth.Repository
}

func (u AuthUsecase) CreateSession(ctx context.Context, ID int64) (string, error) {
	SID, err := u.authRepo.CreateSession(ctx, ID)
	if err != nil {
		return "", err
	}
	return SID, nil
}

func (u AuthUsecase) DeleteSession(ctx context.Context, SID string) error {
	err := u.authRepo.DeleteSession(ctx, SID)
	if err != nil {
		return err
	}
	return nil
}

func (u AuthUsecase) GetSession(ctx context.Context, SID string) (int64, error) {
	id, err := u.authRepo.GetSession(ctx, SID)
	if err != nil {
		return -1, err
	}
	return id, nil
}
