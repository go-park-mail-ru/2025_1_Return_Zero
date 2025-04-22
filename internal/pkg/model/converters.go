package model

import (
	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

///////////////////////////////////// PAGINATION ////////////////////////////////////

func PaginationFromUsecaseToRepository(usecasePagination *usecase.Pagination) *repository.Pagination {
	return &repository.Pagination{
		Offset: usecasePagination.Offset,
		Limit:  usecasePagination.Limit,
	}
}

func PaginationFromDeliveryToUsecase(deliveryPagination *delivery.Pagination) *usecase.Pagination {
	return &usecase.Pagination{
		Offset: deliveryPagination.Offset,
		Limit:  deliveryPagination.Limit,
	}
}

func PaginationFromUsecaseToArtistProto(usecasePagination *usecase.Pagination) *artistProto.Pagination {
	return &artistProto.Pagination{
		Offset: int64(usecasePagination.Offset),
		Limit:  int64(usecasePagination.Limit),
	}
}

func PaginationFromUsecaseToAlbumProto(usecasePagination *usecase.Pagination) *albumProto.Pagination {
	return &albumProto.Pagination{
		Offset: int64(usecasePagination.Offset),
		Limit:  int64(usecasePagination.Limit),
	}
}

///////////////////////////////////// ALBUM ////////////////////////////////////

func AlbumsFromUsecaseToDelivery(usecaseAlbums []*usecase.Album) []*delivery.Album {
	albums := make([]*delivery.Album, 0, len(usecaseAlbums))
	for _, usecaseAlbum := range usecaseAlbums {
		albums = append(albums, AlbumFromUsecaseToDelivery(usecaseAlbum, usecaseAlbum.Artists))
	}
	return albums
}

func AlbumFromUsecaseToDelivery(usecaseAlbum *usecase.Album, usecaseAlbumArtists []*usecase.AlbumArtist) *delivery.Album {
	return &delivery.Album{
		ID:          usecaseAlbum.ID,
		Title:       usecaseAlbum.Title,
		Type:        delivery.AlbumType(usecaseAlbum.Type),
		Thumbnail:   usecaseAlbum.Thumbnail,
		Artists:     AlbumArtistsFromUsecaseToDelivery(usecaseAlbumArtists),
		ReleaseDate: usecaseAlbum.ReleaseDate,
	}
}

func AlbumArtistsFromUsecaseToDelivery(usecaseAlbumArtists []*usecase.AlbumArtist) []*delivery.AlbumArtist {
	albumArtists := make([]*delivery.AlbumArtist, 0, len(usecaseAlbumArtists))
	for _, usecaseAlbumArtist := range usecaseAlbumArtists {
		albumArtists = append(albumArtists, &delivery.AlbumArtist{
			ID:    usecaseAlbumArtist.ID,
			Title: usecaseAlbumArtist.Title,
		})
	}
	return albumArtists
}

func AlbumFromRepositoryToUsecase(repositoryAlbum *repository.Album, repositoryAlbumArtists []*repository.ArtistWithTitle) *usecase.Album {
	return &usecase.Album{
		ID:          repositoryAlbum.ID,
		Title:       repositoryAlbum.Title,
		Type:        usecase.AlbumType(repositoryAlbum.Type),
		Thumbnail:   repositoryAlbum.Thumbnail,
		Artists:     AlbumArtistsFromRepositoryToUsecase(repositoryAlbumArtists),
		ReleaseDate: repositoryAlbum.ReleaseDate,
	}
}

func AlbumArtistsFromRepositoryToUsecase(repositoryAlbumArtists []*repository.ArtistWithTitle) []*usecase.AlbumArtist {
	albumArtists := make([]*usecase.AlbumArtist, 0, len(repositoryAlbumArtists))
	for _, repoAlbumArtist := range repositoryAlbumArtists {
		albumArtists = append(albumArtists, &usecase.AlbumArtist{
			ID:    repoAlbumArtist.ID,
			Title: repoAlbumArtist.Title,
		})
	}
	return albumArtists
}

func AlbumFromProtoToUsecase(protoAlbum *albumProto.Album) *usecase.Album {
	var albumType usecase.AlbumType

	switch protoAlbum.Type {
	case albumProto.AlbumType_AlbumTypeAlbum:
		albumType = usecase.AlbumTypeAlbum
	case albumProto.AlbumType_AlbumTypeEP:
		albumType = usecase.AlbumTypeEP
	case albumProto.AlbumType_AlbumTypeSingle:
		albumType = usecase.AlbumTypeSingle
	case albumProto.AlbumType_AlbumTypeCompilation:
		albumType = usecase.AlbumTypeCompilation
	default:
		albumType = usecase.AlbumTypeAlbum
	}

	return &usecase.Album{
		ID:          protoAlbum.Id,
		Title:       protoAlbum.Title,
		Type:        albumType,
		Thumbnail:   protoAlbum.Thumbnail,
		ReleaseDate: protoAlbum.ReleaseDate.AsTime(),
	}
}

///////////////////////////////////// ARTIST ////////////////////////////////////

func ArtistWithTitleListFromProtoToUsecase(protoArtistWithTitleList []*artistProto.ArtistWithTitle) []*usecase.AlbumArtist {
	artistWithTitleList := make([]*usecase.AlbumArtist, 0, len(protoArtistWithTitleList))
	for _, protoArtistWithTitle := range protoArtistWithTitleList {
		artistWithTitleList = append(artistWithTitleList, &usecase.AlbumArtist{
			ID:    protoArtistWithTitle.Id,
			Title: protoArtistWithTitle.Title,
		})
	}
	return artistWithTitleList
}

func ArtistWithTitleMapFromProtoToUsecase(protoArtistWithTitleMap map[int64]*artistProto.ArtistWithTitleList) map[int64][]*usecase.AlbumArtist {
	artistWithTitleMap := make(map[int64][]*usecase.AlbumArtist, len(protoArtistWithTitleMap))
	for id, protoArtistWithTitleList := range protoArtistWithTitleMap {
		artistWithTitleMap[id] = ArtistWithTitleListFromProtoToUsecase(protoArtistWithTitleList.Artists)
	}
	return artistWithTitleMap
}

func ArtistsFromUsecaseToDelivery(usecaseArtists []*usecase.Artist) []*delivery.Artist {
	artists := make([]*delivery.Artist, 0, len(usecaseArtists))
	for _, usecaseArtist := range usecaseArtists {
		artists = append(artists, ArtistFromUsecaseToDelivery(usecaseArtist))
	}
	return artists
}

func ArtistFromUsecaseToDelivery(usecaseArtist *usecase.Artist) *delivery.Artist {
	return &delivery.Artist{
		ID:          usecaseArtist.ID,
		Title:       usecaseArtist.Title,
		Thumbnail:   usecaseArtist.Thumbnail,
		Description: usecaseArtist.Description,
	}
}

func ArtistsFromProtoToUsecase(protoArtists []*artistProto.Artist) []*usecase.Artist {
	artists := make([]*usecase.Artist, 0, len(protoArtists))
	for _, protoArtist := range protoArtists {
		artists = append(artists, ArtistFromProtoToUsecase(protoArtist))
	}
	return artists
}

func ArtistFromProtoToUsecase(protoArtist *artistProto.Artist) *usecase.Artist {
	return &usecase.Artist{
		ID:          protoArtist.Id,
		Title:       protoArtist.Title,
		Thumbnail:   protoArtist.Thumbnail,
		Description: protoArtist.Description,
	}
}

func ArtistDetailedFromProtoToUsecase(protoArtist *artistProto.ArtistDetailed) *usecase.ArtistDetailed {
	return &usecase.ArtistDetailed{
		Artist:    *ArtistFromProtoToUsecase(protoArtist.Artist),
		Favorites: protoArtist.FavoritesCount,
		Listeners: protoArtist.ListenersCount,
	}
}

func ArtistDetailedFromUsecaseToDelivery(usecaseArtistDetailed *usecase.ArtistDetailed) *delivery.ArtistDetailed {
	return &delivery.ArtistDetailed{
		Artist:    *ArtistFromUsecaseToDelivery(&usecaseArtistDetailed.Artist),
		Favorites: usecaseArtistDetailed.Favorites,
		Listeners: usecaseArtistDetailed.Listeners,
	}
}

///////////////////////////////////// TRACK ////////////////////////////////////

func TracksFromUsecaseToDelivery(usecaseTracks []*usecase.Track) []*delivery.Track {
	tracks := make([]*delivery.Track, 0, len(usecaseTracks))
	for _, usecaseTrack := range usecaseTracks {
		tracks = append(tracks, TrackFromUsecaseToDelivery(usecaseTrack))
	}
	return tracks
}

func TrackFromUsecaseToDelivery(usecaseTrack *usecase.Track) *delivery.Track {
	return &delivery.Track{
		ID:        usecaseTrack.ID,
		Title:     usecaseTrack.Title,
		Thumbnail: usecaseTrack.Thumbnail,
		Duration:  usecaseTrack.Duration,
		Album:     usecaseTrack.Album,
		AlbumID:   usecaseTrack.AlbumID,
		Artists:   TrackArtistsFromUsecaseToDelivery(usecaseTrack.Artists),
	}
}

func TracksDetailedFromUsecaseToDelivery(usecaseTracks []*usecase.TrackDetailed) []*delivery.TrackDetailed {
	tracks := make([]*delivery.TrackDetailed, 0, len(usecaseTracks))
	for _, usecaseTrack := range usecaseTracks {
		tracks = append(tracks, TrackDetailedFromUsecaseToDelivery(usecaseTrack))
	}
	return tracks
}

func TrackDetailedFromUsecaseToDelivery(usecaseTrackDetailed *usecase.TrackDetailed) *delivery.TrackDetailed {
	return &delivery.TrackDetailed{
		Track:   *TrackFromUsecaseToDelivery(&usecaseTrackDetailed.Track),
		FileUrl: usecaseTrackDetailed.FileUrl,
	}
}

func TrackArtistsFromUsecaseToDelivery(usecaseTrackArtists []*usecase.TrackArtist) []*delivery.TrackArtist {
	trackArtists := make([]*delivery.TrackArtist, 0, len(usecaseTrackArtists))
	for _, usecaseTrackArtist := range usecaseTrackArtists {
		trackArtists = append(trackArtists, &delivery.TrackArtist{
			ID:    usecaseTrackArtist.ID,
			Title: usecaseTrackArtist.Title,
			Role:  usecaseTrackArtist.Role,
		})
	}
	return trackArtists
}

func TrackFromRepositoryToUsecase(repositoryTrack *repository.Track, repositoryTrackArtists []*repository.ArtistWithRole, albumTitle string) *usecase.Track {
	return &usecase.Track{
		ID:        repositoryTrack.ID,
		Title:     repositoryTrack.Title,
		Thumbnail: repositoryTrack.Thumbnail,
		Duration:  repositoryTrack.Duration,
		AlbumID:   repositoryTrack.AlbumID,
		Album:     albumTitle,
		Artists:   TrackArtistsFromRepositoryToUsecase(repositoryTrackArtists),
	}
}

func TrackArtistsFromRepositoryToUsecase(repositoryTrackArtists []*repository.ArtistWithRole) []*usecase.TrackArtist {
	trackArtists := make([]*usecase.TrackArtist, 0, len(repositoryTrackArtists))
	for _, repositoryTrackArtist := range repositoryTrackArtists {
		trackArtists = append(trackArtists, &usecase.TrackArtist{
			ID:    repositoryTrackArtist.ID,
			Title: repositoryTrackArtist.Title,
			Role:  repositoryTrackArtist.Role,
		})
	}
	return trackArtists
}

func TrackWithFileKeyFromRepositoryToUsecase(repositoryTrack *repository.TrackWithFileKey, repositoryTrackArtists []*repository.ArtistWithRole, albumTitle string) *usecase.Track {
	return &usecase.Track{
		ID:        repositoryTrack.ID,
		Title:     repositoryTrack.Title,
		Thumbnail: repositoryTrack.Thumbnail,
		Duration:  repositoryTrack.Duration,
		AlbumID:   repositoryTrack.AlbumID,
		Album:     albumTitle,
		Artists:   TrackArtistsFromRepositoryToUsecase(repositoryTrackArtists),
	}
}

func TrackDetailedFromRepositoryToUsecase(repositoryTrack *repository.TrackWithFileKey, repositoryTrackArtists []*repository.ArtistWithRole, albumTitle string, fileUrl string) *usecase.TrackDetailed {
	return &usecase.TrackDetailed{
		Track:   *TrackWithFileKeyFromRepositoryToUsecase(repositoryTrack, repositoryTrackArtists, albumTitle),
		FileUrl: fileUrl,
	}
}

///////////////////////////////////// STREAM ////////////////////////////////////

func TrackStreamCreateDataFromDeliveryToUsecase(deliveryTrackStream *delivery.TrackStreamCreateData) *usecase.TrackStreamCreateData {
	return &usecase.TrackStreamCreateData{
		TrackID: deliveryTrackStream.TrackID,
		UserID:  deliveryTrackStream.UserID,
	}
}

func TrackStreamCreateDataFromUsecaseToRepository(usecaseTrackStream *usecase.TrackStreamCreateData) *repository.TrackStreamCreateData {
	return &repository.TrackStreamCreateData{
		TrackID: usecaseTrackStream.TrackID,
		UserID:  usecaseTrackStream.UserID,
	}
}

func TrackStreamUpdateDataFromUsecaseToRepository(usecaseTrackStream *usecase.TrackStreamUpdateData) *repository.TrackStreamUpdateData {
	return &repository.TrackStreamUpdateData{
		StreamID: usecaseTrackStream.StreamID,
		Duration: usecaseTrackStream.Duration,
	}
}

func TrackStreamUpdateDataFromDeliveryToUsecase(repositoryTrackStream *delivery.TrackStreamUpdateData, userID int64, streamID int64) *usecase.TrackStreamUpdateData {
	return &usecase.TrackStreamUpdateData{
		StreamID: streamID,
		Duration: repositoryTrackStream.Duration,
		UserID:   userID,
	}
}

// /////////////////////////////////// USER ////////////////////////////////////
func PrivacyRepositoryToUsecase(repositoryPrivacy *repository.UserPrivacySettings) *usecase.UserPrivacy {
	return &usecase.UserPrivacy{
		IsPublicPlaylists:       repositoryPrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: repositoryPrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: repositoryPrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  repositoryPrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  repositoryPrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: repositoryPrivacy.IsPublicArtistsListened,
	}
}

func StatisticsRepositoryToUsecase(repositoryStatistics *repository.UserStats) *usecase.UserStatistics {
	return &usecase.UserStatistics{
		MinutesListened: repositoryStatistics.MinutesListened,
		TracksListened:  repositoryStatistics.TracksListened,
		ArtistsListened: repositoryStatistics.ArtistsListened,
	}
}

func UserFullDataRepositoryToUsecase(repositoryUser *repository.UserFullData) *usecase.UserFullData {
	usecasePrivacy := PrivacyRepositoryToUsecase(repositoryUser.Privacy)
	usecaseStatistics := StatisticsRepositoryToUsecase(repositoryUser.Statistics)
	return &usecase.UserFullData{
		Username:   repositoryUser.Username,
		Email:      repositoryUser.Email,
		Privacy:    usecasePrivacy,
		Statistics: usecaseStatistics,
	}
}

func UserFullDataUsecaseToDelivery(usecaseUser *usecase.UserFullData) *delivery.UserFullData {
	return &delivery.UserFullData{
		Username:   usecaseUser.Username,
		AvatarUrl:  usecaseUser.AvatarUrl,
		Email:      usecaseUser.Email,
		Privacy:    PrivacyUsecaseToDelivery(usecaseUser.Privacy),
		Statistics: StatisticsUsecaseToDelivery(usecaseUser.Statistics),
	}
}

func PrivacyUsecaseToDelivery(usecasePrivacy *usecase.UserPrivacy) *delivery.Privacy {
	return &delivery.Privacy{
		IsPublicPlaylists:       usecasePrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: usecasePrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: usecasePrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  usecasePrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  usecasePrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: usecasePrivacy.IsPublicArtistsListened,
	}
}

func StatisticsUsecaseToDelivery(usecaseStatistics *usecase.UserStatistics) *delivery.Statistics {
	return &delivery.Statistics{
		MinutesListened: usecaseStatistics.MinutesListened,
		TracksListened:  usecaseStatistics.TracksListened,
		ArtistsListened: usecaseStatistics.ArtistsListened,
	}
}

func PrivacyFromUsecaseToRepository(usecasePrivacy *usecase.UserPrivacy) *repository.UserPrivacySettings {
	if usecasePrivacy == nil {
		return nil
	}
	return &repository.UserPrivacySettings{
		IsPublicPlaylists:       usecasePrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: usecasePrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: usecasePrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  usecasePrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  usecasePrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: usecasePrivacy.IsPublicArtistsListened,
	}
}

func ChangeDataFromUsecaseToRepository(usecaseUser *usecase.UserChangeSettings) *repository.ChangeUserData {

	return &repository.ChangeUserData{
		Password:    usecaseUser.Password,
		NewUsername: usecaseUser.NewUsername,
		NewEmail:    usecaseUser.NewEmail,
		NewPassword: usecaseUser.NewPassword,
	}
}

func PrivacyFromDeliveryToUsecase(deliveryPrivacy *delivery.Privacy) *usecase.UserPrivacy {
	if deliveryPrivacy == nil {
		return nil
	}
	return &usecase.UserPrivacy{
		IsPublicPlaylists:       deliveryPrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: deliveryPrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: deliveryPrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  deliveryPrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  deliveryPrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: deliveryPrivacy.IsPublicArtistsListened,
	}
}

func ChangeDataFromDeliveryToUsecase(deliveryUser *delivery.UserChangeSettings) *usecase.UserChangeSettings {
	privacyDelivery := PrivacyFromDeliveryToUsecase(deliveryUser.Privacy)
	return &usecase.UserChangeSettings{
		Password:    deliveryUser.Password,
		NewUsername: deliveryUser.NewUsername,
		NewEmail:    deliveryUser.NewEmail,
		NewPassword: deliveryUser.NewPassword,
		Privacy:     privacyDelivery,
	}
}
