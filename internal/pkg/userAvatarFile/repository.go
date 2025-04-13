package userAvatarFile

import (
	"context"
	"io"
)

type Repository interface {
	GetAvatarURL(fileKey string) (string, error)
	UploadUserAvatar(ctx context.Context, username string, fileContent io.Reader) (string, error)
	DeleteUserAvatar(ctx context.Context, username string) error
}
