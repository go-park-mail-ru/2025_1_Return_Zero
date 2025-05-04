package model

import (
	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
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

func PaginationFromUsecaseToTrackProto(usecasePagination *usecase.Pagination) *trackProto.Pagination {
	return &trackProto.Pagination{
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
		IsLiked:     usecaseAlbum.IsLiked,
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
		IsLiked:     protoAlbum.IsFavorite,
	}
}

func AlbumIdsFromUsecaseToAlbumProto(usecaseAlbumIds []int64) []*albumProto.AlbumID {
	albumIds := make([]*albumProto.AlbumID, 0, len(usecaseAlbumIds))
	for _, id := range usecaseAlbumIds {
		albumIds = append(albumIds, &albumProto.AlbumID{Id: id})
	}
	return albumIds
}

func AlbumLikeRequestFromUsecaseToProto(usecaseLikeRequest *usecase.AlbumLikeRequest) *albumProto.LikeRequest {
	return &albumProto.LikeRequest{
		AlbumId: &albumProto.AlbumID{Id: usecaseLikeRequest.AlbumID},
		UserId:  &albumProto.UserID{Id: usecaseLikeRequest.UserID},
		IsLike:  usecaseLikeRequest.IsLike,
	}
}

func AlbumLikeRequestFromDeliveryToUsecase(isLike bool, userID int64, albumID int64) *usecase.AlbumLikeRequest {
	return &usecase.AlbumLikeRequest{
		AlbumID: albumID,
		IsLike:  isLike,
		UserID:  userID,
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
		IsLiked:     usecaseArtist.IsLiked,
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
		IsLiked:     protoArtist.IsFavorite,
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

func TrackIdsFromUsecaseToArtistProto(usecaseTrackIds []int64) []*artistProto.TrackID {
	trackIds := make([]*artistProto.TrackID, 0, len(usecaseTrackIds))
	for _, id := range usecaseTrackIds {
		trackIds = append(trackIds, &artistProto.TrackID{Id: id})
	}
	return trackIds
}

func ArtistWithRoleListFromProtoToUsecase(protoArtistWithRoleList []*artistProto.ArtistWithRole) []*usecase.TrackArtist {
	artistWithRoleList := make([]*usecase.TrackArtist, 0, len(protoArtistWithRoleList))
	for _, protoArtistWithRole := range protoArtistWithRoleList {
		artistWithRoleList = append(artistWithRoleList, &usecase.TrackArtist{
			ID:    protoArtistWithRole.Id,
			Title: protoArtistWithRole.Title,
			Role:  protoArtistWithRole.Role,
		})
	}
	return artistWithRoleList
}

func UserIDFromUsecaseToProtoArtist(userID int64) *artistProto.UserID {
	return &artistProto.UserID{
		Id: userID,
	}
}

func ArtistsListenedFromProtoToUsecase(protoArtistsNum *artistProto.ArtistListened) int64 {
	return protoArtistsNum.ArtistsListened
}

func ArtistLikeRequestFromUsecaseToProto(usecaseLikeRequest *usecase.ArtistLikeRequest) *artistProto.LikeRequest {
	return &artistProto.LikeRequest{
		ArtistId: &artistProto.ArtistID{Id: usecaseLikeRequest.ArtistID},
		UserId:   &artistProto.UserID{Id: usecaseLikeRequest.UserID},
		IsLike:   usecaseLikeRequest.IsLike,
	}
}

func ArtistLikeRequestFromDeliveryToUsecase(isLike bool, userID int64, artistID int64) *usecase.ArtistLikeRequest {
	return &usecase.ArtistLikeRequest{
		ArtistID: artistID,
		IsLike:   isLike,
		UserID:   userID,
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
		IsLiked:   usecaseTrack.IsLiked,
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

func TrackIdsFromUsecaseToTrackProto(usecaseTrackIds []*usecase.Track) []*trackProto.TrackID {
	trackIds := make([]*trackProto.TrackID, 0, len(usecaseTrackIds))
	for _, usecaseTrack := range usecaseTrackIds {
		trackIds = append(trackIds, &trackProto.TrackID{Id: usecaseTrack.ID})
	}
	return trackIds
}

func TrackFromProtoToUsecase(protoTrack *trackProto.Track, protoAlbum *albumProto.AlbumTitle, protoArtists *artistProto.ArtistWithRoleList) *usecase.Track {
	return &usecase.Track{
		ID:        protoTrack.Id,
		Title:     protoTrack.Title,
		Thumbnail: protoTrack.Thumbnail,
		Duration:  protoTrack.Duration,
		AlbumID:   protoTrack.AlbumId,
		Album:     protoAlbum.Title,
		Artists:   ArtistWithRoleListFromProtoToUsecase(protoArtists.Artists),
		IsLiked:   protoTrack.IsFavorite,
	}
}

func TrackDetailedFromProtoToUsecase(protoTrack *trackProto.TrackDetailed, protoAlbum *albumProto.AlbumTitle, protoArtists *artistProto.ArtistWithRoleList) *usecase.TrackDetailed {
	return &usecase.TrackDetailed{
		Track:   *TrackFromProtoToUsecase(protoTrack.Track, protoAlbum, protoArtists),
		FileUrl: protoTrack.FileUrl,
	}
}

func TrackIDListFromArtistToTrackProto(protoArtist *artistProto.TrackIDList, userID int64) *trackProto.TrackIDList {
	trackIds := make([]*trackProto.TrackID, 0, len(protoArtist.Ids))
	for _, id := range protoArtist.Ids {
		trackIds = append(trackIds, &trackProto.TrackID{Id: id.Id})
	}
	return &trackProto.TrackIDList{Ids: trackIds, UserId: &trackProto.UserID{Id: userID}}
}

func TrackLikeRequestFromUsecaseToProto(usecaseLikeRequest *usecase.TrackLikeRequest) *trackProto.LikeRequest {
	return &trackProto.LikeRequest{
		TrackId: &trackProto.TrackID{Id: usecaseLikeRequest.TrackID},
		UserId:  &trackProto.UserID{Id: usecaseLikeRequest.UserID},
		IsLike:  usecaseLikeRequest.IsLike,
	}
}

func TrackLikeRequestFromDeliveryToUsecase(isLike bool, userID int64, trackID int64) *usecase.TrackLikeRequest {
	return &usecase.TrackLikeRequest{
		TrackID: trackID,
		IsLike:  isLike,
		UserID:  userID,
	}
}

func UserIDFromUsecaseToProtoTrack(userID int64) *trackProto.UserID {
	return &trackProto.UserID{
		Id: userID,
	}
}

func TracksListenedFromProtoToUsecase(protoTracksNum *trackProto.TracksListened) int64 {
	return protoTracksNum.Tracks
}

func MinutesListenedFromProtoToUsecase(protoMinutesNum *trackProto.MinutesListened) int64 {
	return protoMinutesNum.Minutes
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

func TrackStreamCreateDataFromUsecaseToProto(usecaseTrackStream *usecase.TrackStreamCreateData) *trackProto.TrackStreamCreateData {
	return &trackProto.TrackStreamCreateData{
		TrackId: &trackProto.TrackID{Id: usecaseTrackStream.TrackID},
		UserId:  &trackProto.UserID{Id: usecaseTrackStream.UserID},
	}
}

func TrackStreamUpdateDataFromUsecaseToProto(usecaseTrackStream *usecase.TrackStreamUpdateData) *trackProto.TrackStreamUpdateData {
	return &trackProto.TrackStreamUpdateData{
		StreamId: &trackProto.StreamID{Id: usecaseTrackStream.StreamID},
		Duration: usecaseTrackStream.Duration,
		UserId:   &trackProto.UserID{Id: usecaseTrackStream.UserID},
	}
}

func ArtistIdsFromUsecaseToArtistProto(artistIDs []int64) *artistProto.ArtistIDList {
	artistIds := make([]*artistProto.ArtistID, 0, len(artistIDs))
	for _, id := range artistIDs {
		artistIds = append(artistIds, &artistProto.ArtistID{Id: id})
	}
	return &artistProto.ArtistIDList{Ids: artistIds}
}

func ArtistStreamCreateDataListFromUsecaseToProto(userID int64, artistIDs []int64) *artistProto.ArtistStreamCreateDataList {
	return &artistProto.ArtistStreamCreateDataList{
		ArtistIds: ArtistIdsFromUsecaseToArtistProto(artistIDs),
		UserId:    &artistProto.UserID{Id: userID},
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

///////////////////////////////////// PLAYLIST ////////////////////////////////////

func PlaylistsFromProtoToUsecase(protoPlaylists []*playlistProto.Playlist, username string) []*usecase.Playlist {
	usecasePlaylists := make([]*usecase.Playlist, 0, len(protoPlaylists))
	for _, protoPlaylist := range protoPlaylists {
		usecasePlaylists = append(usecasePlaylists, PlaylistFromProtoToUsecase(protoPlaylist, username))
	}
	return usecasePlaylists
}

func PlaylistWithIsLikedFromProtoToUsecase(protoPlaylist *playlistProto.PlaylistWithIsLiked, username string) *usecase.PlaylistWithIsLiked {
	return &usecase.PlaylistWithIsLiked{
		Playlist: *PlaylistFromProtoToUsecase(protoPlaylist.Playlist, username),
		IsLiked:  protoPlaylist.IsLiked,
	}
}

func PlaylistWithIsLikedFromUsecaseToDelivery(usecasePlaylist *usecase.PlaylistWithIsLiked) *delivery.PlaylistWithIsLiked {
	return &delivery.PlaylistWithIsLiked{
		Playlist: *PlaylistFromUsecaseToDelivery(&usecasePlaylist.Playlist),
		IsLiked:  usecasePlaylist.IsLiked,
	}
}

func LikePlaylistRequestFromDeliveryToUsecase(userID int64, playlistID int64, isLike bool) *usecase.LikePlaylistRequest {
	return &usecase.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		IsLike:     isLike,
	}
}

func LikePlaylistRequestFromUsecaseToProto(usecaseLikePlaylist *usecase.LikePlaylistRequest) *playlistProto.LikePlaylistRequest {
	return &playlistProto.LikePlaylistRequest{
		UserId:     usecaseLikePlaylist.UserID,
		PlaylistId: usecaseLikePlaylist.PlaylistID,
		IsLike:     usecaseLikePlaylist.IsLike,
	}
}

func UpdatePlaylistsPublisityByUserIDRequestFromUsecaseToProto(isPublic bool, userID int64) *playlistProto.UpdatePlaylistsPublisityByUserIDRequest {
	return &playlistProto.UpdatePlaylistsPublisityByUserIDRequest{
		IsPublic: isPublic,
		UserId:   userID,
	}
}

func UploadPlaylistThumbnailRequestFromUsecaseToProto(title string, thumbnail []byte) *playlistProto.UploadPlaylistThumbnailRequest {
	return &playlistProto.UploadPlaylistThumbnailRequest{
		Title:     title,
		Thumbnail: thumbnail,
	}
}

func CreatePlaylistRequestFromUsecaseToProto(usecasePlaylist *usecase.CreatePlaylistRequest, thumbnail string, isPublic bool) *playlistProto.CreatePlaylistRequest {
	return &playlistProto.CreatePlaylistRequest{
		Title:     usecasePlaylist.Title,
		UserId:    usecasePlaylist.UserID,
		Thumbnail: thumbnail,
		IsPublic:  isPublic,
	}
}

func PlaylistFromProtoToUsecase(protoPlaylist *playlistProto.Playlist, username string) *usecase.Playlist {
	return &usecase.Playlist{
		ID:        protoPlaylist.Id,
		Title:     protoPlaylist.Title,
		Thumbnail: protoPlaylist.Thumbnail,
		Username:  username,
	}
}

func CreatePlaylistRequestFromDeliveryToUsecase(deliveryPlaylist *delivery.CreatePlaylistRequest, userID int64) *usecase.CreatePlaylistRequest {
	return &usecase.CreatePlaylistRequest{
		Title:     deliveryPlaylist.Title,
		UserID:    userID,
		Thumbnail: deliveryPlaylist.Thumbnail,
	}
}

func PlaylistFromUsecaseToDelivery(usecasePlaylist *usecase.Playlist) *delivery.Playlist {
	return &delivery.Playlist{
		ID:        usecasePlaylist.ID,
		Title:     usecasePlaylist.Title,
		Thumbnail: usecasePlaylist.Thumbnail,
		Username:  usecasePlaylist.Username,
	}
}

func PlaylistsFromUsecaseToDelivery(usecasePlaylists []*usecase.Playlist) []*delivery.Playlist {
	deliveryPlaylists := make([]*delivery.Playlist, 0, len(usecasePlaylists))
	for _, usecasePlaylist := range usecasePlaylists {
		deliveryPlaylists = append(deliveryPlaylists, PlaylistFromUsecaseToDelivery(usecasePlaylist))
	}
	return deliveryPlaylists
}

func AddTrackToPlaylistRequestFromDeliveryToUsecase(deliveryAddTrackToPlaylist *delivery.AddTrackToPlaylistRequest, userID int64, playlistID int64) *usecase.AddTrackToPlaylistRequest {
	return &usecase.AddTrackToPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    deliveryAddTrackToPlaylist.TrackID,
	}
}

func RemoveTrackFromPlaylistRequestFromDeliveryToUsecase(trackID int64, userID int64, playlistID int64) *usecase.RemoveTrackFromPlaylistRequest {
	return &usecase.RemoveTrackFromPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    trackID,
	}
}

func AddTrackToPlaylistRequestFromUsecaseToProto(usecaseAddTrackToPlaylist *usecase.AddTrackToPlaylistRequest) *playlistProto.AddTrackToPlaylistRequest {
	return &playlistProto.AddTrackToPlaylistRequest{
		PlaylistId: usecaseAddTrackToPlaylist.PlaylistID,
		TrackId:    usecaseAddTrackToPlaylist.TrackID,
		UserId:     usecaseAddTrackToPlaylist.UserID,
	}
}

func RemoveTrackFromPlaylistRequestFromUsecaseToProto(usecaseRemoveTrackFromPlaylist *usecase.RemoveTrackFromPlaylistRequest) *playlistProto.RemoveTrackFromPlaylistRequest {
	return &playlistProto.RemoveTrackFromPlaylistRequest{
		PlaylistId: usecaseRemoveTrackFromPlaylist.PlaylistID,
		TrackId:    usecaseRemoveTrackFromPlaylist.TrackID,
		UserId:     usecaseRemoveTrackFromPlaylist.UserID,
	}
}

func UpdatePlaylistRequestFromUsecaseToProto(usecaseUpdatePlaylist *usecase.UpdatePlaylistRequest, thumbnail string) *playlistProto.UpdatePlaylistRequest {
	return &playlistProto.UpdatePlaylistRequest{
		Id:        usecaseUpdatePlaylist.PlaylistID,
		Title:     usecaseUpdatePlaylist.Title,
		Thumbnail: thumbnail,
		UserId:    usecaseUpdatePlaylist.UserID,
	}
}

func UpdatePlaylistRequestFromDeliveryToUsecase(deliveryUpdatePlaylist *delivery.UpdatePlaylistRequest, userID int64, playlistID int64) *usecase.UpdatePlaylistRequest {
	return &usecase.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      deliveryUpdatePlaylist.Title,
		Thumbnail:  deliveryUpdatePlaylist.Thumbnail,
	}
}

func RemovePlaylistRequestFromUsecaseToProto(usecaseRemovePlaylist *usecase.RemovePlaylistRequest) *playlistProto.RemovePlaylistRequest {
	return &playlistProto.RemovePlaylistRequest{
		UserId:     usecaseRemovePlaylist.UserID,
		PlaylistId: usecaseRemovePlaylist.PlaylistID,
	}
}

func RemovePlaylistRequestFromDeliveryToUsecase(playlistID int64, userID int64) *usecase.RemovePlaylistRequest {
	return &usecase.RemovePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}
}

func GetPlaylistsToAddRequestFromDeliveryToUsecase(trackID int64, userID int64) *usecase.GetPlaylistsToAddRequest {
	return &usecase.GetPlaylistsToAddRequest{
		UserID:  userID,
		TrackID: trackID,
	}
}

func GetPlaylistsToAddRequestFromUsecaseToProto(usecaseGetPlaylistsToAdd *usecase.GetPlaylistsToAddRequest) *playlistProto.GetPlaylistsToAddRequest {
	return &playlistProto.GetPlaylistsToAddRequest{
		UserId:  usecaseGetPlaylistsToAdd.UserID,
		TrackId: usecaseGetPlaylistsToAdd.TrackID,
	}
}

func GetPlaylistsToAddResponseFromProtoToUsecase(proto *playlistProto.GetPlaylistsToAddResponse, username string) []*usecase.PlaylistWithIsIncludedTrack {
	usecasePlaylists := make([]*usecase.PlaylistWithIsIncludedTrack, 0, len(proto.Playlists))
	for _, protoPlaylist := range proto.Playlists {
		usecasePlaylists = append(usecasePlaylists, &usecase.PlaylistWithIsIncludedTrack{
			Playlist:   *PlaylistFromProtoToUsecase(protoPlaylist.Playlist, username),
			IsIncluded: protoPlaylist.IsIncludedTrack,
		})
	}
	return usecasePlaylists
}

func PlaylistsWithIsIncludedTrackFromUsecaseToDelivery(usecasePlaylists []*usecase.PlaylistWithIsIncludedTrack) []*delivery.PlaylistWithIsIncludedTrack {
	deliveryPlaylists := make([]*delivery.PlaylistWithIsIncludedTrack, 0, len(usecasePlaylists))
	for _, usecasePlaylist := range usecasePlaylists {
		deliveryPlaylists = append(deliveryPlaylists, &delivery.PlaylistWithIsIncludedTrack{
			Playlist:   *PlaylistFromUsecaseToDelivery(&usecasePlaylist.Playlist),
			IsIncluded: usecasePlaylist.IsIncluded,
		})
	}
	return deliveryPlaylists
}

func RegisterDataFromUsecaseToProto(regData *usecase.User) *userProto.RegisterData {
	return &userProto.RegisterData{
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
	}
}

func UserFromProtoToUsecase(protoUser *userProto.UserFront) *usecase.User {
	return &usecase.User{
		ID:        protoUser.Id,
		Username:  protoUser.Username,
		Email:     protoUser.Email,
		AvatarUrl: protoUser.Avatar,
	}
}

func UserIDFromUsecaseToProtoUser(userID int64) *userProto.UserID {
	return &userProto.UserID{
		Id: userID,
	}
}

func UserIDFromProtoToUsecaseUser(protoUserID *userProto.UserID) int64 {
	return protoUserID.Id
}

func LoginDataFromUsecaseToProto(loginData *usecase.User) *userProto.LoginData {
	return &userProto.LoginData{
		Username: loginData.Username,
		Email:    loginData.Email,
		Password: loginData.Password,
	}
}

func AvatarDataFromUsecaseToProto(fileURL string, id int64) *userProto.AvatarData {
	return &userProto.AvatarData{
		AvatarPath: fileURL,
		Id:         id,
	}
}

func DeleteUserFromUsecaseToProto(user *usecase.User) *userProto.UserDelete {
	return &userProto.UserDelete{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func UsernameFromUsecaseToProto(username string) *userProto.Username {
	return &userProto.Username{
		Username: username,
	}
}

func PrivacyFromProtoToUsecase(protoPrivacy *userProto.PrivacySettings) *usecase.UserPrivacy {
	return &usecase.UserPrivacy{
		IsPublicPlaylists:       protoPrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: protoPrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: protoPrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  protoPrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  protoPrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: protoPrivacy.IsPublicArtistsListened,
	}
}

func UserFullDataFromProtoToUsecase(protoUser *userProto.UserFullData) *usecase.UserFullData {
	privacyUsecase := PrivacyFromProtoToUsecase(protoUser.Privacy)
	return &usecase.UserFullData{
		Username:  protoUser.Username,
		Email:     protoUser.Email,
		AvatarUrl: protoUser.Avatar,
		Privacy:   privacyUsecase,
	}
}

func PrivacyFromUsecaseToProto(username string, usecasePrivacy *usecase.UserPrivacy) *userProto.PrivacySettings {
	return &userProto.PrivacySettings{
		Username:                username,
		IsPublicPlaylists:       usecasePrivacy.IsPublicPlaylists,
		IsPublicMinutesListened: usecasePrivacy.IsPublicMinutesListened,
		IsPublicFavoriteArtists: usecasePrivacy.IsPublicFavoriteArtists,
		IsPublicTracksListened:  usecasePrivacy.IsPublicTracksListened,
		IsPublicFavoriteTracks:  usecasePrivacy.IsPublicFavoriteTracks,
		IsPublicArtistsListened: usecasePrivacy.IsPublicArtistsListened,
	}
}

func ChangeUserDataFromUsecaseToProto(username string, usecaseUser *usecase.UserChangeSettings) *userProto.ChangeUserDataMessage {
	return &userProto.ChangeUserDataMessage{
		Username:    username,
		NewUsername: usecaseUser.NewUsername,
		NewEmail:    usecaseUser.NewEmail,
		NewPassword: usecaseUser.NewPassword,
		Password:    usecaseUser.Password,
	}
}

func FileKeyFromUsecaseToProto(avatarURL string) *userProto.FileKey {
	return &userProto.FileKey{
		FileKey: avatarURL,
	}
}

func AvatarUrlFromProtoToUsecase(protoAvatarURL *userProto.AvatarUrl) string {
	return protoAvatarURL.Url
}

func AvatarImageFromUsecaseToProto(username string, image []byte) *userProto.AvatarImage {
	return &userProto.AvatarImage{
		Username: username,
		Image:    image,
	}
}

func FileKeyFromProtoToUsecase(protoFileKey *userProto.FileKey) string {
	return protoFileKey.FileKey
}

// ///////////////////////////////////// AUTH ////////////////////////////////////
func SessionIDFromProtoToUsecase(protoSessionID *authProto.SessionID) string {
	return protoSessionID.SessionId
}

func UserIDFromProtoToUsecase(protoUserID *authProto.UserID) int64 {
	return protoUserID.Id
}

func SessionIDFromUsecaseToProto(sessionID string) *authProto.SessionID {
	return &authProto.SessionID{
		SessionId: sessionID,
	}
}

func UserIDFromUsecaseToProto(userID int64) *authProto.UserID {
	return &authProto.UserID{
		Id: userID,
	}
}
