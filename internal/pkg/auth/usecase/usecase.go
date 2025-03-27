package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func NewAuthUsecase(authRepo auth.Repository) auth.Usecase {
	return AuthUsecase{authRepo: authRepo}
}

type AuthUsecase struct {
	authRepo auth.Repository
}

func (u AuthUsecase) CreateSession(ID uint) string {
	return u.authRepo.CreateSession(ID)
}

func (u AuthUsecase) DeleteSession(SID string) {
	u.authRepo.DeleteSession(SID)
}

func (u AuthUsecase) GetSession(SID string) (*model.Session, error) {
	return u.authRepo.GetSession(SID)
}