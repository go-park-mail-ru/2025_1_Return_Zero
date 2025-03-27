package user

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Repository interface {
	CreateUser(regData *model.RegisterData) (*model.UserToFront, error)
	GetUserByID(ID uint) (*model.UserToFront, error)
	LoginUser(logData *model.LoginData) (*model.UserToFront, error)
}
