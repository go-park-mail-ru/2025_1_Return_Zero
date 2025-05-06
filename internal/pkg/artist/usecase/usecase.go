package usecase

import (
	"context"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(artistClient artistProto.ArtistServiceClient, userClient userProto.UserServiceClient) artist.Usecase {
	return &artistUsecase{
		artistClient: artistClient,
		userClient:   userClient,
	}
}

type artistUsecase struct {
	artistClient artistProto.ArtistServiceClient
	userClient   userProto.UserServiceClient
}

func (u *artistUsecase) GetArtistByID(ctx context.Context, id int64) (*usecaseModel.ArtistDetailed, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoArtist, err := u.artistClient.GetArtistByID(ctx, &artistProto.ArtistIDWithUserID{
		ArtistId: &artistProto.ArtistID{Id: id},
		UserId:   &artistProto.UserID{Id: userID},
	})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistDetailedFromProtoToUsecase(protoArtist), nil
}

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoFilters := &artistProto.FiltersWithUserID{
		Filters: &artistProto.Filters{
			Pagination: model.PaginationFromUsecaseToArtistProto(filters.Pagination),
		},
		UserId: &artistProto.UserID{Id: userID},
	}

	protoArtists, err := u.artistClient.GetAllArtists(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistsFromProtoToUsecase(protoArtists.Artists), nil
}

func (u *artistUsecase) LikeArtist(ctx context.Context, request *usecaseModel.ArtistLikeRequest) error {
	protoRequest := model.ArtistLikeRequestFromUsecaseToProto(request)
	_, err := u.artistClient.LikeArtist(ctx, protoRequest)
	if err != nil {
		return customErrors.HandleArtistGRPCError(err)
	}
	return nil
}

func (u *artistUsecase) GetFavoriteArtists(ctx context.Context, filters *usecaseModel.ArtistFilters, username string) ([]*usecaseModel.Artist, error) {
	profileUserID, err := u.userClient.GetIDByUsername(ctx, &userProto.Username{Username: username})
	if err != nil {
		return nil, customErrors.HandleUserGRPCError(err)
	}

	profilePrivacy, err := u.userClient.GetUserPrivacyByID(ctx, &userProto.UserID{Id: profileUserID.Id})
	if err != nil {
		return nil, customErrors.HandleUserGRPCError(err)
	}

	currentUserID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		currentUserID = -1
	}

	if !profilePrivacy.IsPublicFavoriteArtists && profileUserID.Id != currentUserID {
		return []*usecaseModel.Artist{}, nil
	}

	protoFilters := &artistProto.FiltersWithUserID{
		Filters: &artistProto.Filters{
			Pagination: model.PaginationFromUsecaseToArtistProto(filters.Pagination),
		},
		UserId: &artistProto.UserID{Id: profileUserID.Id},
	}

	protoArtists, err := u.artistClient.GetFavoriteArtists(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistsFromProtoToUsecase(protoArtists.Artists), nil
}

func (u *artistUsecase) SearchArtists(ctx context.Context, query string) ([]*usecaseModel.Artist, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoRequest := &artistProto.Query{
		Query:  query,
		UserId: &artistProto.UserID{Id: userID},
	}

	protoArtists, err := u.artistClient.SearchArtists(ctx, protoRequest)
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	return model.ArtistsFromProtoToUsecase(protoArtists.Artists), nil
}
