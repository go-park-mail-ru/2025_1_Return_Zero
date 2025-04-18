package userAvatarFile

import (
	"context"
	"errors"
	"io"
)

var (
	ErrAvatarNotFound         = errors.New("avatar not found")
	ErrFailedToParseImage     = errors.New("failed to parse image")
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
	ErrFailedToEncodeWebp     = errors.New("failed to encode webp")
	ErrFailedToUploadAvatar   = errors.New("failed to upload avatar")
	ErrFileTooLarge           = errors.New("file too large")
)

type Repository interface {
	GetAvatarURL(ctx context.Context, fileKey string) (string, error)
	UploadUserAvatar(ctx context.Context, username string, fileContent io.Reader) (string, error)
}
