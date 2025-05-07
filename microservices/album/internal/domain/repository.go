package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
)

type Repository interface {
	GetAllAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64, userID int64) (*repoModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumTitleByIDs(ctx context.Context, ids []int64) (map[int64]string, error)
	GetAlbumsByIDs(ctx context.Context, ids []int64, userID int64) ([]*repoModel.Album, error)
	CreateStream(ctx context.Context, albumID int64, userID int64) error
	LikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error
	CheckAlbumExists(ctx context.Context, albumID int64) (bool, error)
	UnlikeAlbum(ctx context.Context, request *repoModel.LikeRequest) error
	GetFavoriteAlbums(ctx context.Context, filters *repoModel.AlbumFilters, userID int64) ([]*repoModel.Album, error)
	SearchAlbums(ctx context.Context, query string, userID int64) ([]*repoModel.Album, error)
}
