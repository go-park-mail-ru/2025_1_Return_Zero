package user

import (
	"context"
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrUsernameExist = errors.New("user with this username already exists")
	ErrEmailExist    = errors.New("user with this email already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrCreateSalt    = errors.New("failed to create salt")
	ErrWrongPassword = errors.New("wrong password")
)

type Repository interface {
	GetFullUserData(ctx context.Context, username string) (*repoModel.UserFullData, error)
	GetUserPrivacy(ctx context.Context, id int64) (*repoModel.UserPrivacySettings, error)
	GetUserStats(ctx context.Context, id int64) (*repoModel.UserStats, error)
	GetUserData(ctx context.Context, id int64) (*repoModel.User, error)
	GetIDByUsername(ctx context.Context, username string) (int64, error)
	ChangeUserPrivacySettings(ctx context.Context, username string, privacySettings *repoModel.UserPrivacySettings) error
	DeleteUser(ctx context.Context, user *repoModel.User) error
	ChangeUserData(ctx context.Context, username string, changeData *repoModel.ChangeUserData) error
	UploadAvatar(ctx context.Context, avatarUrl string, username string) error
	GetAvatar(ctx context.Context, username string) (string, error)
	LoginUser(ctx context.Context, logData *repoModel.User) (*repoModel.User, error)
	GetUserByID(ctx context.Context, ID int64) (*repoModel.User, error)
	CreateUser(ctx context.Context, regData *repoModel.User) (*repoModel.User, error)
}
