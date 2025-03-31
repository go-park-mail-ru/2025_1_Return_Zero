package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
)

func NewAuthUsecase(authRepo auth.Repository) auth.Usecase {
	return AuthUsecase{authRepo: authRepo}
}

type AuthUsecase struct {
	authRepo auth.Repository
}

func (u AuthUsecase) CreateSession(ID int64) string {
	return u.authRepo.CreateSession(ID)
}

func (u AuthUsecase) DeleteSession(SID string) {
	u.authRepo.DeleteSession(SID)
}

func (u AuthUsecase) GetSession(SID string) (int64, error) {
	return u.authRepo.GetSession(SID)
}