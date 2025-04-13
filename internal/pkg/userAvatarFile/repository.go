package userAvatarFile

import (
	"io"
)

type Repository interface {
	GetAvatarURL(fileKey string) (string, error)
	UploadUserAvatar(username string, fileContent io.Reader) (string, error)
}
