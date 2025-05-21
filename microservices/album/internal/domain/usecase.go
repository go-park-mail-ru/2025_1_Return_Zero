package domain

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
)

type Usecase interface {
	GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters, userID int64) ([]*usecaseModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64, userID int64) (*usecaseModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumTitleByIDs(ctx context.Context, ids []int64) (*usecaseModel.AlbumTitleMap, error)
	GetAlbumsByIDs(ctx context.Context, ids []int64, userID int64) ([]*usecaseModel.Album, error)
	CreateStream(ctx context.Context, albumID int64, userID int64) error
	LikeAlbum(ctx context.Context, request *usecaseModel.LikeRequest) error
	GetFavoriteAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters, userID int64) ([]*usecaseModel.Album, error)
	SearchAlbums(ctx context.Context, query string, userID int64) ([]*usecaseModel.Album, error)
	CreateAlbum(ctx context.Context, album *usecaseModel.CreateAlbumRequest) (int64, error) 
	DeleteAlbum(ctx context.Context, albumID int64) error
}
