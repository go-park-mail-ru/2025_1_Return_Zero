package album

import (
	"context"
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrAlbumNotFound = errors.New("album not found")
)

type Repository interface {
	GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*repoModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]*repoModel.Album, error)
}
