package userAvatarFile

import (
	"io"
)

type Repository interface {
	GetPresignedURL(userKey string) (string, error)
	UploadUserAvatar(username string, fileContent io.Reader) (string, error)
}
