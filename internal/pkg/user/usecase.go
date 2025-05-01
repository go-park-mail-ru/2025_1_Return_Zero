package user

import (
	"context"
	"io"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	CreateUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error)
	GetUserBySID(ctx context.Context, SID string) (*usecaseModel.User, error)
	LoginUser(ctx context.Context, user *usecaseModel.User) (*usecaseModel.User, string, error)
	Logout(ctx context.Context, SID string) error
	UploadAvatar(ctx context.Context, username string, fileAvatar io.Reader, ID int64) (string, error)
	ChangeUserData(ctx context.Context, username string, userChangeData *usecaseModel.UserChangeSettings) (*usecaseModel.UserFullData, error)
	DeleteUser(ctx context.Context, user *usecaseModel.User, SID string) error
	GetUserData(ctx context.Context, username string) (*usecaseModel.UserFullData, error)
	GetUserByID(ctx context.Context, id int64) (*usecaseModel.User, error)
}
