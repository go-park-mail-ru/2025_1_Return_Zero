package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
)

type Repository interface {
	GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters) ([]*repoModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*repoModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error)
	GetAlbumsByIDs(ctx context.Context, ids []int64) ([]*repoModel.Album, error)
	CreateStream(ctx context.Context, albumID int64, userID int64) error
}
