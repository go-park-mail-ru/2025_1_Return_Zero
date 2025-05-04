package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
)

func NewUsecase(albumClient *albumProto.AlbumServiceClient, artistClient *artistProto.ArtistServiceClient) album.Usecase {
	return &albumUsecase{albumClient: albumClient, artistClient: artistClient}
}

type albumUsecase struct {
	albumClient  *albumProto.AlbumServiceClient
	artistClient *artistProto.ArtistServiceClient
}

func (u *albumUsecase) GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoFilters := &albumProto.FiltersWithUserID{
		Filters: &albumProto.Filters{
			Pagination: model.PaginationFromUsecaseToAlbumProto(filters.Pagination),
		},
		UserId: &albumProto.UserID{Id: userID},
	}

	protoAlbums, err := (*u.albumClient).GetAllAlbums(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	albumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		albumIDs = append(albumIDs, &artistProto.AlbumID{Id: protoAlbum.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: albumIDs})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	artistWithTitleMap := model.ArtistWithTitleMapFromProtoToUsecase(protoArtists.Artists)

	albums := make([]*usecaseModel.Album, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		usecaseAlbum := model.AlbumFromProtoToUsecase(protoAlbum)
		usecaseAlbum.Artists = artistWithTitleMap[protoAlbum.Id]
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}

func (u *albumUsecase) GetAlbumsByArtistID(ctx context.Context, artistID int64, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoAlbumIDs, err := (*u.artistClient).GetAlbumIDsByArtistID(ctx, &artistProto.ArtistID{Id: artistID})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	albumIDs := make([]*albumProto.AlbumID, 0, len(protoAlbumIDs.Ids))
	for _, protoAlbumID := range protoAlbumIDs.Ids {
		newAlbumID := &albumProto.AlbumID{Id: protoAlbumID.Id}
		albumIDs = append(albumIDs, newAlbumID)
	}

	protoAlbums, err := (*u.albumClient).GetAlbumsByIDs(ctx, &albumProto.AlbumIDListWithUserID{
		Ids:    &albumProto.AlbumIDList{Ids: albumIDs},
		UserId: &albumProto.UserID{Id: userID},
	})

	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	artistAlbumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbumIDs.Ids))
	for _, protoAlbumID := range protoAlbumIDs.Ids {
		artistAlbumIDs = append(artistAlbumIDs, &artistProto.AlbumID{Id: protoAlbumID.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: artistAlbumIDs})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	artistWithTitleMap := model.ArtistWithTitleMapFromProtoToUsecase(protoArtists.Artists)

	albums := make([]*usecaseModel.Album, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		usecaseAlbum := model.AlbumFromProtoToUsecase(protoAlbum)
		usecaseAlbum.Artists = artistWithTitleMap[protoAlbum.Id]
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}

func (u *albumUsecase) GetAlbumByID(ctx context.Context, id int64) (*usecaseModel.Album, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoAlbum, err := (*u.albumClient).GetAlbumByID(ctx, &albumProto.AlbumIDWithUserID{
		AlbumId: &albumProto.AlbumID{Id: id},
		UserId:  &albumProto.UserID{Id: userID},
	})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumID(ctx, &artistProto.AlbumID{Id: id})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	artistWithTitleList := model.ArtistWithTitleListFromProtoToUsecase(protoArtists.Artists)

	usecaseAlbum := model.AlbumFromProtoToUsecase(protoAlbum)
	usecaseAlbum.Artists = artistWithTitleList
	return usecaseAlbum, nil
}

func (u *albumUsecase) LikeAlbum(ctx context.Context, request *usecaseModel.AlbumLikeRequest) error {
	protoRequest := model.AlbumLikeRequestFromUsecaseToProto(request)
	_, err := (*u.albumClient).LikeAlbum(ctx, protoRequest)
	if err != nil {
		return customErrors.HandleAlbumGRPCError(err)
	}
	return nil
}

func (u *albumUsecase) GetFavoriteAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters, userID int64) ([]*usecaseModel.Album, error) {
	protoFilters := &albumProto.FiltersWithUserID{
		Filters: &albumProto.Filters{
			Pagination: model.PaginationFromUsecaseToAlbumProto(filters.Pagination),
		},
		UserId: &albumProto.UserID{Id: userID},
	}

	protoAlbums, err := (*u.albumClient).GetFavoriteAlbums(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	albumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		albumIDs = append(albumIDs, &artistProto.AlbumID{Id: protoAlbum.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: albumIDs})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	artistWithTitleMap := model.ArtistWithTitleMapFromProtoToUsecase(protoArtists.Artists)

	albums := make([]*usecaseModel.Album, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		usecaseAlbum := model.AlbumFromProtoToUsecase(protoAlbum)
		usecaseAlbum.Artists = artistWithTitleMap[protoAlbum.Id]
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}

func (u *albumUsecase) SearchAlbums(ctx context.Context, query string) ([]*usecaseModel.Album, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoRequest := &albumProto.Query{
		Query:  query,
		UserId: &albumProto.UserID{Id: userID},
	}

	protoAlbums, err := (*u.albumClient).SearchAlbums(ctx, protoRequest)
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	albumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		albumIDs = append(albumIDs, &artistProto.AlbumID{Id: protoAlbum.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: albumIDs})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	artistWithTitleMap := model.ArtistWithTitleMapFromProtoToUsecase(protoArtists.Artists)

	albums := make([]*usecaseModel.Album, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		usecaseAlbum := model.AlbumFromProtoToUsecase(protoAlbum)
		usecaseAlbum.Artists = artistWithTitleMap[protoAlbum.Id]
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}
