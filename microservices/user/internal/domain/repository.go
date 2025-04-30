package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
)

type Repository interface {
	CreateUser(ctx context.Context, regData *repoModel.RegisterData) (*repoModel.User, error)
	LoginUser(ctx context.Context, logData *repoModel.LoginData) (*repoModel.User, error)
	GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error)
	UploadAvatar(ctx context.Context, avatarUrl string, ID int64) error
	GetIDByUsername(ctx context.Context, username string) (int64, error)
	DeleteUser(ctx context.Context, userRepo *repoModel.UserDelete) error
	ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error
	ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.PrivacySettings) error
	GetUserPrivacy(ctx context.Context, id int64) (*repoModel.PrivacySettings, error)
	GetFullUserData(ctx context.Context, username string) (*repoModel.UserFullData, error)
}

type S3Repository interface {
	GetPresignedURL(fileKey string) (string, error)
	UploadUserAvatar(ctx context.Context, username string, fileContent []byte) (string, error)
	DeleteUserAvatar(ctx context.Context, fileKey string) error
	GetAvatarURL(ctx context.Context, fileKey string) (string, error)
	GetAvatarKey(ctx context.Context, username string) (string, error)
}