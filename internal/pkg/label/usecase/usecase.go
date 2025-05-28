package usecase

import (
	"context"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/domain"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewLabelUsecase(labelRepo domain.Repository, userProto userProto.UserServiceClient, artistProto artistProto.ArtistServiceClient, albumProto albumProto.AlbumServiceClient, trackProto trackProto.TrackServiceClient) domain.Usecase {
	return &labelUsecase{
		labelRepository: labelRepo,
		userProto:       userProto,
		artistProto:     artistProto,
		albumProto:      albumProto,
		trackProto:      trackProto,
	}
}

type labelUsecase struct {
	labelRepository domain.Repository
	userProto       userProto.UserServiceClient
	artistProto     artistProto.ArtistServiceClient
	albumProto      albumProto.AlbumServiceClient
	trackProto      trackProto.TrackServiceClient
	S3Repository    domain.S3Repository
}

func (u *labelUsecase) CheckIsLabelUnique(ctx context.Context, labelName string) (bool, error) {
	exist, err := u.labelRepository.CheckIsLabelUnique(ctx, labelName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (u *labelUsecase) CreateLabel(ctx context.Context, label *usecaseModel.Label) (*usecaseModel.Label, error) {
	isExist, err := u.CheckIsLabelUnique(ctx, label.Name)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, customErrors.ErrLableExist
	}
	_, err = u.userProto.ChecksUsersByUsernames(ctx, &userProto.Usernames{
		Usernames: label.Members,
	})
	if err != nil {
		return nil, err
	}

	labelID, err := u.labelRepository.CreateLabel(ctx, label.Name)
	if err != nil {
		return nil, err
	}

	_, err = u.userProto.UpdateUsersLabelID(ctx, &userProto.RequestUpdateUserLabelID{
		LabelId:   labelID,
		Usernames: label.Members,
	})
	if err != nil {
		return nil, err
	}
	label.Id = labelID
	return label, nil
}

func (u *labelUsecase) GetLabel(ctx context.Context, id int64) (*usecaseModel.Label, error) {
	label, err := u.labelRepository.GetLabel(ctx, id)
	if err != nil {
		return nil, err
	}
	members, err := u.userProto.GetUsersByLabelID(ctx, &userProto.LabelID{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	label.Members = members.Usernames
	labelUsecase := usecaseModel.Label{
		Id:      label.ID,
		Name:    label.Name,
		Members: label.Members,
	}
	return &labelUsecase, nil
}

func (u *labelUsecase) CreateArtist(ctx context.Context, artist *usecaseModel.ArtistLoad) (*usecaseModel.Artist, error) {
	artistProto := &artistProto.ArtistLoad{
		Title:   artist.Title,
		Image:   artist.Image,
		LabelId: artist.LabelID,
	}
	protoCreatedArtist, err := u.artistProto.CreateArtist(ctx, artistProto)
	if err != nil {
		return nil, err
	}

	return model.ArtistFromProtoToUsecase(protoCreatedArtist), nil
}

func (u *labelUsecase) EditArtist(ctx context.Context, artist *usecaseModel.ArtistEdit) (*usecaseModel.Artist, error) {
	artistProto := &artistProto.ArtistEdit{
		ArtistId: artist.ArtistID,
		Image:    artist.Image,
		LabelId:  artist.LabelID,
		NewTitle: artist.NewTitle,
	}
	protoEditedArtist, err := u.artistProto.EditArtist(ctx, artistProto)
	if err != nil {
		return nil, err
	}

	return model.ArtistFromProtoToUsecase(protoEditedArtist), nil
}

func (u *labelUsecase) GetArtists(ctx context.Context, labelID int64, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	protoArtists, err := u.artistProto.GetArtistsLabelID(ctx, &artistProto.FiltersWithLabelID{
		Filters: &artistProto.Filters{
			Pagination: model.PaginationFromUsecaseToArtistProto(filters.Pagination),
		},
		LabelId: labelID,
	})
	if err != nil {
		return nil, err
	}

	artists := model.ArtistsFromProtoToUsecase(protoArtists.Artists)
	return artists, nil
}

func (u *labelUsecase) GetAlbumsByLabelID(ctx context.Context, labelID int64, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	protoAlbums, err := u.albumProto.GetAlbumsLabelID(ctx, &albumProto.FiltersWithLabelID{
		Filters: &albumProto.Filters{
			Pagination: model.PaginationFromUsecaseToAlbumProto(filters.Pagination),
		},
		LabelId: labelID,
	})
	if err != nil {
		return nil, err
	}

	albumIDs := make([]*artistProto.AlbumID, 0, len(protoAlbums.Albums))
	for _, protoAlbum := range protoAlbums.Albums {
		albumIDs = append(albumIDs, &artistProto.AlbumID{Id: protoAlbum.Id})
	}

	protoArtists, err := u.artistProto.GetArtistsByAlbumIDs(ctx, &artistProto.AlbumIDList{Ids: albumIDs})
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

func (u *labelUsecase) DeleteArtist(ctx context.Context, artist *usecaseModel.ArtistDelete) error {
	artistProto := &artistProto.ArtistDelete{
		ArtistId: artist.ArtistID,
		LabelId:  artist.LabelID,
	}
	_, err := u.artistProto.DeleteArtist(ctx, artistProto)
	if err != nil {
		return err
	}

	return nil
}

func (u *labelUsecase) CreateAlbum(ctx context.Context, album *usecaseModel.CreateAlbumRequest) (int64, string, error) {
	var albumType albumProto.AlbumType
	switch album.Type {
	case string(usecaseModel.AlbumTypeAlbum):
		albumType = albumProto.AlbumType_AlbumTypeAlbum
	case string(usecaseModel.AlbumTypeEP):
		albumType = albumProto.AlbumType_AlbumTypeEP
	case string(usecaseModel.AlbumTypeSingle):
		albumType = albumProto.AlbumType_AlbumTypeSingle
	case string(usecaseModel.AlbumTypeCompilation):
		albumType = albumProto.AlbumType_AlbumTypeCompilation
	default:
		albumType = albumProto.AlbumType_AlbumTypeAlbum
	}

	albumProto := &albumProto.CreateAlbumRequest{
		Type:    albumType,
		Title:   album.Title,
		LabelId: album.LabelID,
		Image:   album.Image,
	}
	protoCreatedAlbum, err := u.albumProto.CreateAlbum(ctx, albumProto)
	if err != nil {
		return -1, "", err
	}

	artistIDs := make([]*artistProto.ArtistID, 0, len(album.ArtistsIDs))
	for _, id := range album.ArtistsIDs {
		artistIDs = append(artistIDs, &artistProto.ArtistID{Id: id})
	}

	tracksListWithAlbumId := &trackProto.TracksListWithAlbumID{
		AlbumId: &trackProto.AlbumID{Id: protoCreatedAlbum.Id},
		Tracks:  model.TrackListLoadFromUsecaseToProto(album.Tracks),
		Cover:   album.Image,
	}

	tracksIds, err := u.trackProto.AddTracksToAlbum(ctx, tracksListWithAlbumId)
	if err != nil {
		return -1, "", err
	}

	tracksIdsUsecase := model.TracksIdsFromProtoToUsecase(tracksIds)

	_, err = u.artistProto.ConnectArtists(ctx, &artistProto.ArtistsIDWithAlbumID{
		ArtistIds: &artistProto.ArtistIDList{Ids: artistIDs},
		AlbumId:   &artistProto.AlbumID{Id: protoCreatedAlbum.Id},
		TrackIds:  model.TracksIdsFromUsecaseToProtoArtist(tracksIdsUsecase),
	})
	if err != nil {
		return -1, "", err
	}
	return protoCreatedAlbum.Id, protoCreatedAlbum.Url, nil
}

func (u *labelUsecase) UpdateLabel(ctx context.Context, labelID int64, toAdd, toRemove []string) error {
	_, err := u.userProto.RemoveUsersFromLabel(ctx, &userProto.RequestRemoveUserLabelID{
		LabelId:   labelID,
		Usernames: toRemove,
	})
	if err != nil {
		return err
	}
	_, err = u.userProto.UpdateUsersLabelID(ctx, &userProto.RequestUpdateUserLabelID{
		LabelId:   labelID,
		Usernames: toAdd,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *labelUsecase) DeleteAlbum(ctx context.Context, albumID, labelID int64) error {
	_, err := u.albumProto.DeleteAlbum(ctx, &albumProto.AlbumID{
		Id: albumID,
	})
	if err != nil {
		return err
	}
	_, err = u.trackProto.DeleteTracksByAlbumID(ctx, &trackProto.AlbumID{
		Id: albumID,
	})
	if err != nil {
		return err
	}
	return nil
}
