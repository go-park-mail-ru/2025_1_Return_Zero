package domain

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
)

type Usecase interface {
	GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*usecaseModel.Album, error)
	GetAlbumTitleByID(ctx context.Context, id int64) (string, error)
	GetAlbumTitleByIDs(ctx context.Context, ids []int64) (*usecaseModel.AlbumTitleMap, error)
	GetAlbumsByIDs(ctx context.Context, ids []int64) ([]*usecaseModel.Album, error)
}
