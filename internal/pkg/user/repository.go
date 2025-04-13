package user

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	CreateUser(ctx context.Context, regData *repoModel.User) (*repoModel.User, error)
	GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error)
	LoginUser(ctx context.Context, logData *repoModel.User) (*repoModel.User, error)
	GetAvatar(ctx context.Context, username string) (string, error)
	UploadAvatar(ctx context.Context, avatarUrl string, username string) error
	ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.PrivacySettings) error
	DeleteUser(ctx context.Context, user *repoModel.User) error
	GetUserData(ctx context.Context, username string) (*repoModel.UserAndSettings, error)
	ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error 
	GetIdByUsername(ctx context.Context, username string) (int64, error)
	GetUserStats(ctx context.Context, username string) (*repoModel.UserStats, error)
}
