package usecase

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/domain"
)

func NewAuthUsecase(authRepository domain.Repository) domain.Usecase {
	return &authUsecase{
		authRepo: authRepository,
	}
}

type authUsecase struct {
	authRepo domain.Repository
}

func (u *authUsecase) CreateSession(ctx context.Context, userID int64) (string, error) {
	sessionID, err := u.authRepo.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (u *authUsecase) DeleteSession(ctx context.Context, sessionID string) error {
	err := u.authRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) GetSession(ctx context.Context, sessionID string) (int64, error) {
	userID, err := u.authRepo.GetSession(ctx, sessionID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}
