package album

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*usecaseModel.Album, error)
	LikeAlbum(ctx context.Context, request *usecaseModel.AlbumLikeRequest) error
}
