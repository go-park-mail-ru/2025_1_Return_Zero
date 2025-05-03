package album

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*repoModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error)
}
