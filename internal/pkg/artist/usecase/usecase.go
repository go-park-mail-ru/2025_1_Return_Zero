package usecase

import (
	"context"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
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
	var userID int64
	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	} else {
		userID = user.ID
	}

	protoArtist, err := (*u.artistClient).GetArtistByID(ctx, &artistProto.ArtistIDWithUserID{
		ArtistId: &artistProto.ArtistID{Id: id},
		UserId:   &artistProto.UserID{Id: userID},
	})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistDetailedFromProtoToUsecase(protoArtist), nil
}

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	var userID int64
	user, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	} else {
		userID = user.ID
	}

	protoFilters := &artistProto.FiltersWithUserID{
		Filters: &artistProto.Filters{
			Pagination: model.PaginationFromUsecaseToArtistProto(filters.Pagination),
		},
		UserId: &artistProto.UserID{Id: userID},
	}

	protoArtists, err := (*u.artistClient).GetAllArtists(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistsFromProtoToUsecase(protoArtists.Artists), nil
}

func (u *artistUsecase) LikeArtist(ctx context.Context, request *usecaseModel.ArtistLikeRequest) error {
	protoRequest := model.ArtistLikeRequestFromUsecaseToProto(request)
	_, err := (*u.artistClient).LikeArtist(ctx, protoRequest)
	if err != nil {
		return customErrors.HandleArtistGRPCError(err)
	}
	return nil
}
