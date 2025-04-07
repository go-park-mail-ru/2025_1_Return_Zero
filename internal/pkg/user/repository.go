package user

import (
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	CreateUser(user *repoModel.User) (*repoModel.User, error)
	GetUserByID(ID int64) (*repoModel.User, error)
	LoginUser(user *repoModel.User) (*repoModel.User, error)
	UploadAvatar(avatarUrl string, username string) error
	GetAvatar(username string) (string, error)
}
