package album

import (
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrAlbumNotFound = errors.New("album not found")
)

type Repository interface {
	GetAllAlbums(filters *repoModel.AlbumFilters) ([]*repoModel.Album, error)
	GetAlbumByID(id int64) (*repoModel.Album, error)
	GetAlbumTitleByID(id int64) (string, error)
	GetAlbumsByArtistID(artistID int64) ([]*repoModel.Album, error)
}
