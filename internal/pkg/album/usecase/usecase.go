package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
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
	protoFilters := &albumProto.Filters{
		Pagination: model.PaginationFromUsecaseToAlbumProto(filters.Pagination),
	}

	protoAlbums, err := (*u.albumClient).GetAllAlbums(ctx, protoFilters)
	if err != nil {
		return nil, err
	}

	albumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		albumIDs = append(albumIDs, &artistProto.AlbumID{Id: protoAlbum.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: albumIDs})
	if err != nil {
		return nil, err
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
	protoAlbumIDs, err := (*u.artistClient).GetAlbumIDsByArtistID(ctx, &artistProto.ArtistID{Id: artistID})
	if err != nil {
		return nil, err
	}

	albumIDs := make([]*albumProto.AlbumID, 0, len(protoAlbumIDs.Ids))
	for _, protoAlbumID := range protoAlbumIDs.Ids {
		newAlbumID := &albumProto.AlbumID{Id: protoAlbumID.Id}
		albumIDs = append(albumIDs, newAlbumID)
	}

	protoAlbums, err := (*u.albumClient).GetAlbumsByIDs(ctx, &albumProto.AlbumIDList{Ids: albumIDs})
	if err != nil {
		return nil, err
	}

	artistAlbumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbumIDs.Ids))
	for _, protoAlbumID := range protoAlbumIDs.Ids {
		artistAlbumIDs = append(artistAlbumIDs, &artistProto.AlbumID{Id: protoAlbumID.Id})
	}

	protoArtists, err := (*u.artistClient).GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: artistAlbumIDs})
	if err != nil {
		return nil, err
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
