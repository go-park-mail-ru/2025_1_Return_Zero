package domain

import (
	"context"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
)

type Usecase interface {
	CreateUser(ctx context.Context, registerData *usecaseModel.RegisterData) (*usecaseModel.UserFront, error)
	LoginUser(ctx context.Context, loginData *usecaseModel.LoginData) (*usecaseModel.UserFront, error)
	UploadAvatar(ctx context.Context, avatarUrl string, id int64) error
	DeleteUser(ctx context.Context, deleteData *usecaseModel.UserDelete) error
	ChangeUserData(ctx context.Context, username string, changeData *usecaseModel.ChangeUserData) error
	ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *usecaseModel.PrivacySettings) error
	GetFullUserData(ctx context.Context, username string) (*usecaseModel.UserFullData, error)
	GetUserByID(ctx context.Context, id int64) (*usecaseModel.UserFront, error)
}
