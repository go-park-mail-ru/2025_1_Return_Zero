package usecase

import (
	"context"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(artistClient *artistProto.ArtistServiceClient) artist.Usecase {
	return &artistUsecase{
		artistClient: artistClient,
	}
}

type artistUsecase struct {
	artistClient *artistProto.ArtistServiceClient
}

func (u *artistUsecase) GetArtistByID(ctx context.Context, id int64) (*usecaseModel.ArtistDetailed, error) {
	protoArtist, err := (*u.artistClient).GetArtistByID(ctx, &artistProto.ArtistID{Id: id})
	if err != nil {
		return nil, err
	}

	return model.ArtistDetailedFromProtoToUsecase(protoArtist), nil
}

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	protoFilters := &artistProto.Filters{
		Pagination: model.PaginationFromUsecaseToProto(filters.Pagination),
	}

	protoArtists, err := (*u.artistClient).GetAllArtists(ctx, protoFilters)
	if err != nil {
		return nil, err
	}

	return model.ArtistsFromProtoToUsecase(protoArtists.Artists), nil
}
