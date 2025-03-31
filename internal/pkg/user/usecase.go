package user

import (
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	CreateUser(user *usecaseModel.User) (*usecaseModel.User, string, error)
	GetUserBySID(SID string) (*usecaseModel.User, error)
	LoginUser(user *usecaseModel.User) (*usecaseModel.User, string, error)
	Logout(SID string)
}