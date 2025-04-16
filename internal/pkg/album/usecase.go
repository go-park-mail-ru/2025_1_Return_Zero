package album

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error)
	GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]*usecaseModel.Album, error)
}
