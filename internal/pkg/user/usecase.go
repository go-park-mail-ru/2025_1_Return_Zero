package user

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Usecase interface {
	CreateUser(regData *model.RegisterData) (*model.UserToFront, string, error)
	GetUserBySID(SID string) (*model.UserToFront, error)
	LoginUser(logData *model.LoginData) (*model.UserToFront, string, error)
	Logout(SID string)
}