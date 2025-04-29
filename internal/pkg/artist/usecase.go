package artist

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetArtistByID(ctx context.Context, id int64) (*usecaseModel.ArtistDetailed, error)
	GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error)
	LikeArtist(ctx context.Context, request *usecaseModel.ArtistLikeRequest) error
}
