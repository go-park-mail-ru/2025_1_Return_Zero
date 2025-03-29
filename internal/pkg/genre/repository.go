package genre

import (
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	GetGenresByAlbumID(albumID int64) ([]*repoModel.Genre, error)
}
