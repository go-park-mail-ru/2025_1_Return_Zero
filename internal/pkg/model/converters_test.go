package model_test

import (
	"testing"
	"time"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock constants for testing
const (
	mockUserID      = int64(1)
	mockArtistID    = int64(1)
	mockAlbumID     = int64(1)
	mockTrackID     = int64(1)
	mockPlaylistID  = int64(1)
	mockStreamID    = int64(1)
	mockUserID2     = int64(2)
	mockArtistID2   = int64(2)
	mockAlbumID2    = int64(2)
	mockTrackID2    = int64(2)
	mockPlaylistID2 = int64(2)
	mockTrackID3    = int64(3)
	mockUserID777   = int64(777)
	mockUserID12345 = int64(12345)
	mockUsername1   = "user1"

	mockUsername  = "testuser"
	mockEmail     = "test@example.com"
	mockPassword  = "password123"
	mockAvatarURL = "avatar.jpg"
	mockSessionID = "test-session-id"

	mockAlbumTitle    = "Album Title"
	mockArtistTitle   = "Artist Title"
	mockTrackTitle    = "Track Title"
	mockPlaylistTitle = "Playlist Title"

	mockThumbnail   = "thumbnail.jpg"
	mockFileURL     = "http://example.com/track.mp3"
	mockNewUsername = "newuser"
	mockNewEmail    = "new@example.com"
	mockNewPassword = "newpass"
	mockOldPassword = "oldpass"
)

func TestPaginationFromDeliveryToUsecase(t *testing.T) {
	deliveryPagination := &delivery.Pagination{
		Offset: 0,
		Limit:  10,
	}
	ucPagination := model.PaginationFromDeliveryToUsecase(deliveryPagination)
	assert.Equal(t, deliveryPagination.Offset, ucPagination.Offset)
	assert.Equal(t, deliveryPagination.Limit, ucPagination.Limit)
}

func TestPaginationFromUsecaseToArtistProto(t *testing.T) {
	ucPagination := &usecase.Pagination{
		Offset: 0,
		Limit:  10,
	}
	protoPagination := model.PaginationFromUsecaseToArtistProto(ucPagination)
	assert.Equal(t, int64(ucPagination.Offset), protoPagination.Offset)
	assert.Equal(t, int64(ucPagination.Limit), protoPagination.Limit)
}

func TestPaginationFromUsecaseToAlbumProto(t *testing.T) {
	ucPagination := &usecase.Pagination{
		Offset: 0,
		Limit:  10,
	}
	protoPagination := model.PaginationFromUsecaseToAlbumProto(ucPagination)
	assert.Equal(t, int64(ucPagination.Offset), protoPagination.Offset)
	assert.Equal(t, int64(ucPagination.Limit), protoPagination.Limit)
}

func TestPaginationFromUsecaseToTrackProto(t *testing.T) {
	ucPagination := &usecase.Pagination{
		Offset: 0,
		Limit:  10,
	}
	protoPagination := model.PaginationFromUsecaseToTrackProto(ucPagination)
	assert.Equal(t, int64(ucPagination.Offset), protoPagination.Offset)
	assert.Equal(t, int64(ucPagination.Limit), protoPagination.Limit)
}

func TestAlbumConverters_UsecaseToDelivery(t *testing.T) {
	ucAlbumArtist := &usecase.AlbumArtist{
		ID:    1,
		Title: "Artist Title",
	}
	ucAlbum := &usecase.Album{
		ID:          1,
		Title:       "Album Title",
		Type:        usecase.AlbumTypeAlbum,
		Thumbnail:   "thumbnail.jpg",
		Artists:     []*usecase.AlbumArtist{ucAlbumArtist},
		ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		IsLiked:     true,
	}

	deliveryAlbum := model.AlbumFromUsecaseToDelivery(ucAlbum, ucAlbum.Artists)
	assert.Equal(t, ucAlbum.ID, deliveryAlbum.ID)
	assert.Equal(t, ucAlbum.Title, deliveryAlbum.Title)
	assert.Equal(t, delivery.AlbumType(ucAlbum.Type), deliveryAlbum.Type)
	assert.Equal(t, ucAlbum.Thumbnail, deliveryAlbum.Thumbnail)
	assert.Equal(t, ucAlbum.ReleaseDate, deliveryAlbum.ReleaseDate)
	assert.Equal(t, ucAlbum.IsLiked, deliveryAlbum.IsLiked)
	assert.Len(t, deliveryAlbum.Artists, 1)
	assert.Equal(t, ucAlbumArtist.ID, deliveryAlbum.Artists[0].ID)
	assert.Equal(t, ucAlbumArtist.Title, deliveryAlbum.Artists[0].Title)

	ucAlbums := []*usecase.Album{ucAlbum}
	deliveryAlbums := model.AlbumsFromUsecaseToDelivery(ucAlbums)
	assert.Len(t, deliveryAlbums, 1)
	assert.Equal(t, deliveryAlbum, deliveryAlbums[0])

	deliveryAlbumArtists := model.AlbumArtistsFromUsecaseToDelivery(ucAlbum.Artists)
	assert.Len(t, deliveryAlbumArtists, 1)
	assert.Equal(t, ucAlbumArtist.ID, deliveryAlbumArtists[0].ID)
	assert.Equal(t, ucAlbumArtist.Title, deliveryAlbumArtists[0].Title)
}

func TestAlbumFromProtoToUsecase(t *testing.T) {
	tests := []struct {
		name       string
		protoAlbum *albumProto.Album
		expected   *usecase.Album
	}{
		{
			name: "AlbumTypeAlbum",
			protoAlbum: &albumProto.Album{
				Id:          1,
				Title:       "Album Title",
				Type:        albumProto.AlbumType_AlbumTypeAlbum,
				Thumbnail:   "thumbnail.jpg",
				ReleaseDate: &timestamppb.Timestamp{Seconds: time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC).Unix()},
				IsFavorite:  true,
			},
			expected: &usecase.Album{
				ID:          1,
				Title:       "Album Title",
				Type:        usecase.AlbumTypeAlbum,
				Thumbnail:   "thumbnail.jpg",
				ReleaseDate: time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC),
				IsLiked:     true,
			},
		},
		{
			name: "AlbumTypeEP",
			protoAlbum: &albumProto.Album{
				Id:          2,
				Title:       "EP Title",
				Type:        albumProto.AlbumType_AlbumTypeEP,
				Thumbnail:   "ep_thumbnail.jpg",
				ReleaseDate: &timestamppb.Timestamp{Seconds: time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC).Unix()},
				IsFavorite:  false,
			},
			expected: &usecase.Album{
				ID:          2,
				Title:       "EP Title",
				Type:        usecase.AlbumTypeEP,
				Thumbnail:   "ep_thumbnail.jpg",
				ReleaseDate: time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC),
				IsLiked:     false,
			},
		},
		{
			name: "DefaultToAlbum",
			protoAlbum: &albumProto.Album{
				Id:          3,
				Title:       "Default Title",
				Type:        albumProto.AlbumType(99),
				Thumbnail:   "default_thumbnail.jpg",
				ReleaseDate: &timestamppb.Timestamp{Seconds: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()},
				IsFavorite:  true,
			},
			expected: &usecase.Album{
				ID:          3,
				Title:       "Default Title",
				Type:        usecase.AlbumTypeAlbum,
				Thumbnail:   "default_thumbnail.jpg",
				ReleaseDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				IsLiked:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucAlbum := model.AlbumFromProtoToUsecase(tt.protoAlbum)
			assert.Equal(t, tt.expected, ucAlbum)
		})
	}
}

func TestAlbumIdsFromUsecaseToAlbumProto(t *testing.T) {
	ucAlbumIds := []int64{1, 2, 3}
	protoAlbumIds := model.AlbumIdsFromUsecaseToAlbumProto(ucAlbumIds)
	assert.Len(t, protoAlbumIds, len(ucAlbumIds))
	for i, id := range ucAlbumIds {
		assert.Equal(t, id, protoAlbumIds[i].Id)
	}

	ucEmptyAlbumIds := []int64{}
	protoEmptyAlbumIds := model.AlbumIdsFromUsecaseToAlbumProto(ucEmptyAlbumIds)
	assert.Len(t, protoEmptyAlbumIds, 0)
}

func TestAlbumLikeRequestFromUsecaseToProto(t *testing.T) {
	ucLikeRequest := &usecase.AlbumLikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  true,
	}
	protoLikeRequest := model.AlbumLikeRequestFromUsecaseToProto(ucLikeRequest)
	assert.Equal(t, ucLikeRequest.AlbumID, protoLikeRequest.AlbumId.Id)
	assert.Equal(t, ucLikeRequest.UserID, protoLikeRequest.UserId.Id)
	assert.Equal(t, ucLikeRequest.IsLike, protoLikeRequest.IsLike)
}

func TestAlbumLikeRequestFromDeliveryToUsecase(t *testing.T) {
	isLike := true
	userID := mockUserID
	albumID := mockAlbumID

	ucLikeRequest := model.AlbumLikeRequestFromDeliveryToUsecase(isLike, userID, albumID)

	assert.Equal(t, albumID, ucLikeRequest.AlbumID)
	assert.Equal(t, isLike, ucLikeRequest.IsLike)
	assert.Equal(t, userID, ucLikeRequest.UserID)
}

func TestArtistWithTitle_ProtoToUsecase(t *testing.T) {
	protoArtistWithTitle1 := &artistProto.ArtistWithTitle{
		Id:    1,
		Title: "Artist One",
	}
	protoArtistWithTitle2 := &artistProto.ArtistWithTitle{
		Id:    2,
		Title: "Artist Two",
	}

	protoList := []*artistProto.ArtistWithTitle{protoArtistWithTitle1, protoArtistWithTitle2}
	ucList := model.ArtistWithTitleListFromProtoToUsecase(protoList)
	assert.Len(t, ucList, 2)
	assert.Equal(t, protoArtistWithTitle1.Id, ucList[0].ID)
	assert.Equal(t, protoArtistWithTitle1.Title, ucList[0].Title)
	assert.Equal(t, protoArtistWithTitle2.Id, ucList[1].ID)
	assert.Equal(t, protoArtistWithTitle2.Title, ucList[1].Title)

	protoMap := map[int64]*artistProto.ArtistWithTitleList{
		1: {Artists: []*artistProto.ArtistWithTitle{protoArtistWithTitle1}},
		2: {Artists: []*artistProto.ArtistWithTitle{protoArtistWithTitle2}},
	}
	ucMap := model.ArtistWithTitleMapFromProtoToUsecase(protoMap)
	assert.Len(t, ucMap, 2)
	assert.Len(t, ucMap[1], 1)
	assert.Equal(t, protoArtistWithTitle1.Id, ucMap[1][0].ID)
	assert.Equal(t, protoArtistWithTitle1.Title, ucMap[1][0].Title)
	assert.Len(t, ucMap[2], 1)
	assert.Equal(t, protoArtistWithTitle2.Id, ucMap[2][0].ID)
	assert.Equal(t, protoArtistWithTitle2.Title, ucMap[2][0].Title)
}

func TestArtistConverters_UsecaseToDelivery(t *testing.T) {
	ucArtist := &usecase.Artist{
		ID:          1,
		Title:       "Artist Title",
		Thumbnail:   "artist.jpg",
		Description: "artist",
		IsLiked:     true,
	}

	deliveryArtist := model.ArtistFromUsecaseToDelivery(ucArtist)
	assert.Equal(t, ucArtist.ID, deliveryArtist.ID)
	assert.Equal(t, ucArtist.Title, deliveryArtist.Title)
	assert.Equal(t, ucArtist.Thumbnail, deliveryArtist.Thumbnail)
	assert.Equal(t, ucArtist.Description, deliveryArtist.Description)
	assert.Equal(t, ucArtist.IsLiked, deliveryArtist.IsLiked)

	ucArtists := []*usecase.Artist{ucArtist}
	deliveryArtists := model.ArtistsFromUsecaseToDelivery(ucArtists)
	assert.Len(t, deliveryArtists, 1)
	assert.Equal(t, deliveryArtist, deliveryArtists[0])
}

func TestArtistConverters_ProtoToUsecase(t *testing.T) {
	protoArtist := &artistProto.Artist{
		Id:          1,
		Title:       "Artist Title",
		Thumbnail:   "artist.png",
		Description: "artist",
		IsFavorite:  false,
	}

	ucArtist := model.ArtistFromProtoToUsecase(protoArtist)
	assert.Equal(t, protoArtist.Id, ucArtist.ID)
	assert.Equal(t, protoArtist.Title, ucArtist.Title)
	assert.Equal(t, protoArtist.Thumbnail, ucArtist.Thumbnail)
	assert.Equal(t, protoArtist.Description, ucArtist.Description)
	assert.Equal(t, protoArtist.IsFavorite, ucArtist.IsLiked)

	protoArtists := []*artistProto.Artist{protoArtist}
	ucArtists := model.ArtistsFromProtoToUsecase(protoArtists)
	assert.Len(t, ucArtists, 1)
	assert.Equal(t, ucArtist, ucArtists[0])
}

func TestArtistDetailedFromProtoToUsecase(t *testing.T) {
	protoArtist := &artistProto.Artist{
		Id:          1,
		Title:       "Artist Title",
		Thumbnail:   "artist.gif",
		Description: "artist",
		IsFavorite:  true,
	}
	protoArtistDetailed := &artistProto.ArtistDetailed{
		Artist:         protoArtist,
		FavoritesCount: 1,
		ListenersCount: 1,
	}

	ucArtistDetailed := model.ArtistDetailedFromProtoToUsecase(protoArtistDetailed)

	assert.Equal(t, protoArtist.Id, ucArtistDetailed.Artist.ID)
	assert.Equal(t, protoArtist.Title, ucArtistDetailed.Artist.Title)
	assert.Equal(t, protoArtist.Thumbnail, ucArtistDetailed.Artist.Thumbnail)
	assert.Equal(t, protoArtist.Description, ucArtistDetailed.Artist.Description)
	assert.Equal(t, protoArtist.IsFavorite, ucArtistDetailed.Artist.IsLiked)
	assert.Equal(t, protoArtistDetailed.FavoritesCount, ucArtistDetailed.Favorites)
	assert.Equal(t, protoArtistDetailed.ListenersCount, ucArtistDetailed.Listeners)
}

func TestArtistDetailedFromUsecaseToDelivery(t *testing.T) {
	ucArtist := usecase.Artist{
		ID:          1,
		Title:       "Artist Title",
		Thumbnail:   "artist.jpg",
		Description: "artist",
		IsLiked:     false,
	}
	ucArtistDetailed := &usecase.ArtistDetailed{
		Artist:    ucArtist,
		Favorites: 1,
		Listeners: 1,
	}

	deliveryArtistDetailed := model.ArtistDetailedFromUsecaseToDelivery(ucArtistDetailed)

	assert.Equal(t, ucArtist.ID, deliveryArtistDetailed.Artist.ID)
	assert.Equal(t, ucArtist.Title, deliveryArtistDetailed.Artist.Title)
	assert.Equal(t, ucArtist.Thumbnail, deliveryArtistDetailed.Artist.Thumbnail)
	assert.Equal(t, ucArtist.Description, deliveryArtistDetailed.Artist.Description)
	assert.Equal(t, ucArtist.IsLiked, deliveryArtistDetailed.Artist.IsLiked)
	assert.Equal(t, ucArtistDetailed.Favorites, deliveryArtistDetailed.Favorites)
	assert.Equal(t, ucArtistDetailed.Listeners, deliveryArtistDetailed.Listeners)
}

func TestTrackIdsFromUsecaseToArtistProto(t *testing.T) {
	ucTrackIds := []int64{1, 2, 3}
	protoTrackIds := model.TrackIdsFromUsecaseToArtistProto(ucTrackIds)
	assert.Len(t, protoTrackIds, len(ucTrackIds))
	for i, id := range ucTrackIds {
		assert.Equal(t, id, protoTrackIds[i].Id)
	}

	ucEmptyTrackIds := []int64{}
	protoEmptyTrackIds := model.TrackIdsFromUsecaseToArtistProto(ucEmptyTrackIds)
	assert.Len(t, protoEmptyTrackIds, 0)
}

func TestArtistWithRoleListFromProtoToUsecase(t *testing.T) {
	protoArtistWithRole1 := &artistProto.ArtistWithRole{
		Id:    1,
		Title: "Role Artist One",
		Role:  "Main",
	}
	protoArtistWithRole2 := &artistProto.ArtistWithRole{
		Id:    2,
		Title: "Role Artist Two",
		Role:  "Feature",
	}
	protoList := []*artistProto.ArtistWithRole{protoArtistWithRole1, protoArtistWithRole2}

	ucList := model.ArtistWithRoleListFromProtoToUsecase(protoList)

	assert.Len(t, ucList, 2)
	assert.Equal(t, protoArtistWithRole1.Id, ucList[0].ID)
	assert.Equal(t, protoArtistWithRole1.Title, ucList[0].Title)
	assert.Equal(t, protoArtistWithRole1.Role, ucList[0].Role)
	assert.Equal(t, protoArtistWithRole2.Id, ucList[1].ID)
	assert.Equal(t, protoArtistWithRole2.Title, ucList[1].Title)
	assert.Equal(t, protoArtistWithRole2.Role, ucList[1].Role)

	ucEmptyList := model.ArtistWithRoleListFromProtoToUsecase([]*artistProto.ArtistWithRole{})
	assert.Len(t, ucEmptyList, 0)
}

func TestUserIDFromUsecaseToProtoArtist(t *testing.T) {
	userID := mockUserID
	protoUserID := model.UserIDFromUsecaseToProtoArtist(userID)
	assert.Equal(t, userID, protoUserID.Id)
}

func TestArtistsListenedFromProtoToUsecase(t *testing.T) {
	protoArtistsListened := &artistProto.ArtistListened{
		ArtistsListened: 1,
	}
	artistsListened := model.ArtistsListenedFromProtoToUsecase(protoArtistsListened)
	assert.Equal(t, protoArtistsListened.ArtistsListened, artistsListened)
}

func TestArtistLikeRequestFromUsecaseToProto(t *testing.T) {
	ucLikeRequest := &usecase.ArtistLikeRequest{
		ArtistID: 1,
		UserID:   1,
		IsLike:   true,
	}
	protoLikeRequest := model.ArtistLikeRequestFromUsecaseToProto(ucLikeRequest)
	assert.Equal(t, ucLikeRequest.ArtistID, protoLikeRequest.ArtistId.Id)
	assert.Equal(t, ucLikeRequest.UserID, protoLikeRequest.UserId.Id)
	assert.Equal(t, ucLikeRequest.IsLike, protoLikeRequest.IsLike)
}

func TestArtistLikeRequestFromDeliveryToUsecase(t *testing.T) {
	isLike := false
	userID := mockUserID
	artistID := mockArtistID

	ucLikeRequest := model.ArtistLikeRequestFromDeliveryToUsecase(isLike, userID, artistID)

	assert.Equal(t, artistID, ucLikeRequest.ArtistID)
	assert.Equal(t, isLike, ucLikeRequest.IsLike)
	assert.Equal(t, userID, ucLikeRequest.UserID)
}

func TestTrackConverters_UsecaseToDelivery(t *testing.T) {
	ucTrackArtist := &usecase.TrackArtist{
		ID:    1,
		Title: "Track Artist",
		Role:  "Vocalist",
	}
	ucTrack := &usecase.Track{
		ID:        100,
		Title:     "Usecase Track",
		Thumbnail: "uc_track.webp",
		Duration:  200,
		Album:     "Usecase Album",
		AlbumID:   200,
		Artists:   []*usecase.TrackArtist{ucTrackArtist},
		IsLiked:   true,
	}

	deliveryTrackArtists := model.TrackArtistsFromUsecaseToDelivery([]*usecase.TrackArtist{ucTrackArtist})
	assert.Len(t, deliveryTrackArtists, 1)
	assert.Equal(t, ucTrackArtist.ID, deliveryTrackArtists[0].ID)
	assert.Equal(t, ucTrackArtist.Title, deliveryTrackArtists[0].Title)
	assert.Equal(t, ucTrackArtist.Role, deliveryTrackArtists[0].Role)

	deliveryTrack := model.TrackFromUsecaseToDelivery(ucTrack)
	assert.Equal(t, ucTrack.ID, deliveryTrack.ID)
	assert.Equal(t, ucTrack.Title, deliveryTrack.Title)
	assert.Equal(t, ucTrack.Thumbnail, deliveryTrack.Thumbnail)
	assert.Equal(t, ucTrack.Duration, deliveryTrack.Duration)
	assert.Equal(t, ucTrack.Album, deliveryTrack.Album)
	assert.Equal(t, ucTrack.AlbumID, deliveryTrack.AlbumID)
	assert.Equal(t, ucTrack.IsLiked, deliveryTrack.IsLiked)
	assert.Len(t, deliveryTrack.Artists, 1)
	assert.Equal(t, deliveryTrackArtists[0], deliveryTrack.Artists[0])

	ucTracks := []*usecase.Track{ucTrack}
	deliveryTracks := model.TracksFromUsecaseToDelivery(ucTracks)
	assert.Len(t, deliveryTracks, 1)
	assert.Equal(t, deliveryTrack, deliveryTracks[0])
}

func TestTrackDetailedConverters_UsecaseToDelivery(t *testing.T) {
	ucTrackArtist := &usecase.TrackArtist{
		ID:    2,
		Title: "Detailed Track Artist",
		Role:  "Guitarist",
	}
	ucBaseTrack := usecase.Track{
		ID:        101,
		Title:     "Detailed Usecase Track",
		Thumbnail: "detailed_uc_track.jpg",
		Duration:  240,
		Album:     "Detailed Usecase Album",
		AlbumID:   201,
		Artists:   []*usecase.TrackArtist{ucTrackArtist},
		IsLiked:   false,
	}
	ucTrackDetailed := &usecase.TrackDetailed{
		Track:   ucBaseTrack,
		FileUrl: "http://example.com/track.mp3",
	}

	deliveryTrackDetailed := model.TrackDetailedFromUsecaseToDelivery(ucTrackDetailed)
	assert.Equal(t, ucBaseTrack.ID, deliveryTrackDetailed.Track.ID)
	assert.Equal(t, ucBaseTrack.Title, deliveryTrackDetailed.Track.Title)
	assert.Equal(t, ucBaseTrack.Thumbnail, deliveryTrackDetailed.Track.Thumbnail)
	assert.Equal(t, ucBaseTrack.Duration, deliveryTrackDetailed.Track.Duration)
	assert.Equal(t, ucBaseTrack.Album, deliveryTrackDetailed.Track.Album)
	assert.Equal(t, ucBaseTrack.AlbumID, deliveryTrackDetailed.Track.AlbumID)
	assert.Equal(t, ucBaseTrack.IsLiked, deliveryTrackDetailed.Track.IsLiked)
	assert.Len(t, deliveryTrackDetailed.Track.Artists, 1)
	assert.Equal(t, ucTrackArtist.ID, deliveryTrackDetailed.Track.Artists[0].ID)
	assert.Equal(t, ucTrackArtist.Title, deliveryTrackDetailed.Track.Artists[0].Title)
	assert.Equal(t, ucTrackArtist.Role, deliveryTrackDetailed.Track.Artists[0].Role)
	assert.Equal(t, ucTrackDetailed.FileUrl, deliveryTrackDetailed.FileUrl)

	ucTracksDetailed := []*usecase.TrackDetailed{ucTrackDetailed}
	deliveryTracksDetailed := model.TracksDetailedFromUsecaseToDelivery(ucTracksDetailed)
	assert.Len(t, deliveryTracksDetailed, 1)
	assert.Equal(t, deliveryTrackDetailed, deliveryTracksDetailed[0])
}

func TestTrackIdsFromUsecaseToTrackProto(t *testing.T) {
	ucTracks := []*usecase.Track{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	protoTrackIds := model.TrackIdsFromUsecaseToTrackProto(ucTracks)
	assert.Len(t, protoTrackIds, len(ucTracks))
	for i, ucTrack := range ucTracks {
		assert.Equal(t, ucTrack.ID, protoTrackIds[i].Id)
	}

	ucEmptyTracks := []*usecase.Track{}
	protoEmptyTrackIds := model.TrackIdsFromUsecaseToTrackProto(ucEmptyTracks)
	assert.Len(t, protoEmptyTrackIds, 0)
}

func TestTrackFromProtoToUsecase(t *testing.T) {
	protoTrack := &trackProto.Track{
		Id:         1,
		Title:      "Track Title",
		Thumbnail:  "track.svg",
		Duration:   200,
		AlbumId:    2,
		IsFavorite: true,
	}
	protoAlbum := &albumProto.AlbumTitle{
		Title: "Album Title",
	}
	protoArtists := &artistProto.ArtistWithRoleList{
		Artists: []*artistProto.ArtistWithRole{
			{Id: 1, Title: "Artist A", Role: "main"},
			{Id: 2, Title: "Artist B", Role: "featured"},
		},
	}

	ucTrack := model.TrackFromProtoToUsecase(protoTrack, protoAlbum, protoArtists)

	assert.Equal(t, protoTrack.Id, ucTrack.ID)
	assert.Equal(t, protoTrack.Title, ucTrack.Title)
	assert.Equal(t, protoTrack.Thumbnail, ucTrack.Thumbnail)
	assert.Equal(t, protoTrack.Duration, ucTrack.Duration)
	assert.Equal(t, protoTrack.AlbumId, ucTrack.AlbumID)
	assert.Equal(t, protoAlbum.Title, ucTrack.Album)
	assert.Equal(t, protoTrack.IsFavorite, ucTrack.IsLiked)
	assert.Len(t, ucTrack.Artists, 2)
	assert.Equal(t, protoArtists.Artists[0].Id, ucTrack.Artists[0].ID)
	assert.Equal(t, protoArtists.Artists[0].Title, ucTrack.Artists[0].Title)
	assert.Equal(t, protoArtists.Artists[0].Role, ucTrack.Artists[0].Role)
	assert.Equal(t, protoArtists.Artists[1].Id, ucTrack.Artists[1].ID)
	assert.Equal(t, protoArtists.Artists[1].Title, ucTrack.Artists[1].Title)
	assert.Equal(t, protoArtists.Artists[1].Role, ucTrack.Artists[1].Role)
}

func TestTrackDetailedFromProtoToUsecase(t *testing.T) {
	baseProtoTrack := &trackProto.Track{
		Id:         1,
		Title:      "Detailed Track Title",
		Thumbnail:  "detailed_track.jpeg",
		Duration:   200,
		AlbumId:    2,
		IsFavorite: false,
	}
	protoTrackDetailed := &trackProto.TrackDetailed{
		Track:   baseProtoTrack,
		FileUrl: "https://example.com/track.mp3",
	}
	protoAlbum := &albumProto.AlbumTitle{
		Title: "Album Title",
	}
	protoArtists := &artistProto.ArtistWithRoleList{
		Artists: []*artistProto.ArtistWithRole{
			{Id: 1, Title: "Artist A", Role: "main"},
		},
	}

	ucTrackDetailed := model.TrackDetailedFromProtoToUsecase(protoTrackDetailed, protoAlbum, protoArtists)

	assert.Equal(t, baseProtoTrack.Id, ucTrackDetailed.Track.ID)
	assert.Equal(t, baseProtoTrack.Title, ucTrackDetailed.Track.Title)
	assert.Equal(t, baseProtoTrack.Thumbnail, ucTrackDetailed.Track.Thumbnail)
	assert.Equal(t, baseProtoTrack.Duration, ucTrackDetailed.Track.Duration)
	assert.Equal(t, baseProtoTrack.AlbumId, ucTrackDetailed.Track.AlbumID)
	assert.Equal(t, protoAlbum.Title, ucTrackDetailed.Track.Album)
	assert.Equal(t, baseProtoTrack.IsFavorite, ucTrackDetailed.Track.IsLiked)
	assert.Len(t, ucTrackDetailed.Track.Artists, 1)
	assert.Equal(t, protoArtists.Artists[0].Id, ucTrackDetailed.Track.Artists[0].ID)
	assert.Equal(t, protoArtists.Artists[0].Title, ucTrackDetailed.Track.Artists[0].Title)
	assert.Equal(t, protoArtists.Artists[0].Role, ucTrackDetailed.Track.Artists[0].Role)
	assert.Equal(t, protoTrackDetailed.FileUrl, ucTrackDetailed.FileUrl)
}

func TestTrackIDListFromArtistToTrackProto(t *testing.T) {
	protoArtistTrackIDList := &artistProto.TrackIDList{
		Ids: []*artistProto.TrackID{
			{Id: 1},
			{Id: 2},
		},
	}
	userID := int64(777)

	trackProtoTrackIDList := model.TrackIDListFromArtistToTrackProto(protoArtistTrackIDList, userID)

	assert.Len(t, trackProtoTrackIDList.Ids, len(protoArtistTrackIDList.Ids))
	for i, artistTrackID := range protoArtistTrackIDList.Ids {
		assert.Equal(t, artistTrackID.Id, trackProtoTrackIDList.Ids[i].Id)
	}
	assert.NotNil(t, trackProtoTrackIDList.UserId)
	assert.Equal(t, userID, trackProtoTrackIDList.UserId.Id)

	emptyProtoArtistTrackIDList := &artistProto.TrackIDList{Ids: []*artistProto.TrackID{}}
	emptyTrackProtoTrackIDList := model.TrackIDListFromArtistToTrackProto(emptyProtoArtistTrackIDList, userID)
	assert.Len(t, emptyTrackProtoTrackIDList.Ids, 0)
	assert.NotNil(t, emptyTrackProtoTrackIDList.UserId)
	assert.Equal(t, userID, emptyTrackProtoTrackIDList.UserId.Id)
}

func TestTrackLikeRequestFromUsecaseToProto(t *testing.T) {
	ucLikeRequest := &usecase.TrackLikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}
	protoLikeRequest := model.TrackLikeRequestFromUsecaseToProto(ucLikeRequest)
	assert.Equal(t, ucLikeRequest.TrackID, protoLikeRequest.TrackId.Id)
	assert.Equal(t, ucLikeRequest.UserID, protoLikeRequest.UserId.Id)
	assert.Equal(t, ucLikeRequest.IsLike, protoLikeRequest.IsLike)
}

func TestTrackLikeRequestFromDeliveryToUsecase(t *testing.T) {
	isLike := false
	userID := mockUserID
	trackID := mockTrackID

	ucLikeRequest := model.TrackLikeRequestFromDeliveryToUsecase(isLike, userID, trackID)

	assert.Equal(t, trackID, ucLikeRequest.TrackID)
	assert.Equal(t, isLike, ucLikeRequest.IsLike)
	assert.Equal(t, userID, ucLikeRequest.UserID)
}

func TestUserIDFromUsecaseToProtoTrack(t *testing.T) {
	userID := mockUserID
	protoUserID := model.UserIDFromUsecaseToProtoTrack(userID)
	assert.Equal(t, userID, protoUserID.Id)
}

func TestTracksListenedFromProtoToUsecase(t *testing.T) {
	protoTracksListened := &trackProto.TracksListened{
		Tracks: 1,
	}
	tracksListened := model.TracksListenedFromProtoToUsecase(protoTracksListened)
	assert.Equal(t, protoTracksListened.Tracks, tracksListened)
}

func TestMinutesListenedFromProtoToUsecase(t *testing.T) {
	protoMinutesListened := &trackProto.MinutesListened{
		Minutes: 1,
	}
	minutesListened := model.MinutesListenedFromProtoToUsecase(protoMinutesListened)
	assert.Equal(t, protoMinutesListened.Minutes, minutesListened)
}

func TestTrackStreamCreateDataFromDeliveryToUsecase(t *testing.T) {
	deliveryTrackStream := &delivery.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}
	ucTrackStream := model.TrackStreamCreateDataFromDeliveryToUsecase(deliveryTrackStream)
	assert.Equal(t, deliveryTrackStream.TrackID, ucTrackStream.TrackID)
	assert.Equal(t, deliveryTrackStream.UserID, ucTrackStream.UserID)
}

func TestTrackStreamUpdateDataFromDeliveryToUsecase(t *testing.T) {
	deliveryTrackStream := &delivery.TrackStreamUpdateData{
		Duration: 1,
	}
	userID := mockUserID
	streamID := mockStreamID

	ucTrackStream := model.TrackStreamUpdateDataFromDeliveryToUsecase(deliveryTrackStream, userID, streamID)

	assert.Equal(t, streamID, ucTrackStream.StreamID)
	assert.Equal(t, deliveryTrackStream.Duration, ucTrackStream.Duration)
	assert.Equal(t, userID, ucTrackStream.UserID)
}

func TestTrackStreamCreateDataFromUsecaseToProto(t *testing.T) {
	ucTrackStream := &usecase.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}
	protoTrackStream := model.TrackStreamCreateDataFromUsecaseToProto(ucTrackStream)
	assert.Equal(t, ucTrackStream.TrackID, protoTrackStream.TrackId.Id)
	assert.Equal(t, ucTrackStream.UserID, protoTrackStream.UserId.Id)
}

func TestTrackStreamUpdateDataFromUsecaseToProto(t *testing.T) {
	ucTrackStream := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 1,
		UserID:   1,
	}
	protoTrackStream := model.TrackStreamUpdateDataFromUsecaseToProto(ucTrackStream)
	assert.Equal(t, ucTrackStream.StreamID, protoTrackStream.StreamId.Id)
	assert.Equal(t, ucTrackStream.Duration, protoTrackStream.Duration)
	assert.Equal(t, ucTrackStream.UserID, protoTrackStream.UserId.Id)
}

func TestArtistIdsFromUsecaseToArtistProto(t *testing.T) {
	artistIDs := []int64{1, 2, 3}
	protoArtistIDList := model.ArtistIdsFromUsecaseToArtistProto(artistIDs)
	assert.Len(t, protoArtistIDList.Ids, len(artistIDs))
	for i, id := range artistIDs {
		assert.Equal(t, id, protoArtistIDList.Ids[i].Id)
	}

	emptyArtistIDs := []int64{}
	emptyProtoArtistIDList := model.ArtistIdsFromUsecaseToArtistProto(emptyArtistIDs)
	assert.Len(t, emptyProtoArtistIDList.Ids, 0)
}

func TestArtistStreamCreateDataListFromUsecaseToProto(t *testing.T) {
	userID := mockUserID
	artistIDs := []int64{1, 2, 3}

	protoArtistStreamCreateDataList := model.ArtistStreamCreateDataListFromUsecaseToProto(userID, artistIDs)

	assert.Equal(t, userID, protoArtistStreamCreateDataList.UserId.Id)
	assert.Len(t, protoArtistStreamCreateDataList.ArtistIds.Ids, len(artistIDs))
	for i, id := range artistIDs {
		assert.Equal(t, id, protoArtistStreamCreateDataList.ArtistIds.Ids[i].Id)
	}

	emptyArtistIDs := []int64{}
	emptyProtoArtistStreamCreateDataList := model.ArtistStreamCreateDataListFromUsecaseToProto(userID, emptyArtistIDs)
	assert.Equal(t, userID, emptyProtoArtistStreamCreateDataList.UserId.Id)
	assert.Len(t, emptyProtoArtistStreamCreateDataList.ArtistIds.Ids, 0)
}

// Tests for User-related converters
func TestUserFullDataUsecaseToDelivery(t *testing.T) {
	ucPrivacy := &usecase.UserPrivacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: true,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  true,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: true,
	}

	ucStatistics := &usecase.UserStatistics{
		MinutesListened: 120,
		TracksListened:  50,
		ArtistsListened: 10,
	}

	ucUserFullData := &usecase.UserFullData{
		Username:   mockUsername,
		Email:      mockEmail,
		AvatarUrl:  mockAvatarURL,
		Privacy:    ucPrivacy,
		Statistics: ucStatistics,
	}

	deliveryUserFullData := model.UserFullDataUsecaseToDelivery(ucUserFullData)

	assert.Equal(t, ucUserFullData.Username, deliveryUserFullData.Username)
	assert.Equal(t, ucUserFullData.Email, deliveryUserFullData.Email)
	assert.Equal(t, ucUserFullData.AvatarUrl, deliveryUserFullData.AvatarUrl)

	assert.Equal(t, ucPrivacy.IsPublicPlaylists, deliveryUserFullData.Privacy.IsPublicPlaylists)
	assert.Equal(t, ucPrivacy.IsPublicMinutesListened, deliveryUserFullData.Privacy.IsPublicMinutesListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteArtists, deliveryUserFullData.Privacy.IsPublicFavoriteArtists)
	assert.Equal(t, ucPrivacy.IsPublicTracksListened, deliveryUserFullData.Privacy.IsPublicTracksListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteTracks, deliveryUserFullData.Privacy.IsPublicFavoriteTracks)
	assert.Equal(t, ucPrivacy.IsPublicArtistsListened, deliveryUserFullData.Privacy.IsPublicArtistsListened)

	assert.Equal(t, ucStatistics.MinutesListened, deliveryUserFullData.Statistics.MinutesListened)
	assert.Equal(t, ucStatistics.TracksListened, deliveryUserFullData.Statistics.TracksListened)
	assert.Equal(t, ucStatistics.ArtistsListened, deliveryUserFullData.Statistics.ArtistsListened)
}

func TestPrivacyUsecaseToDelivery(t *testing.T) {
	ucPrivacy := &usecase.UserPrivacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	deliveryPrivacy := model.PrivacyUsecaseToDelivery(ucPrivacy)

	assert.Equal(t, ucPrivacy.IsPublicPlaylists, deliveryPrivacy.IsPublicPlaylists)
	assert.Equal(t, ucPrivacy.IsPublicMinutesListened, deliveryPrivacy.IsPublicMinutesListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteArtists, deliveryPrivacy.IsPublicFavoriteArtists)
	assert.Equal(t, ucPrivacy.IsPublicTracksListened, deliveryPrivacy.IsPublicTracksListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteTracks, deliveryPrivacy.IsPublicFavoriteTracks)
	assert.Equal(t, ucPrivacy.IsPublicArtistsListened, deliveryPrivacy.IsPublicArtistsListened)
}

func TestStatisticsUsecaseToDelivery(t *testing.T) {
	ucStatistics := &usecase.UserStatistics{
		MinutesListened: 120,
		TracksListened:  50,
		ArtistsListened: 10,
	}

	deliveryStatistics := model.StatisticsUsecaseToDelivery(ucStatistics)

	assert.Equal(t, ucStatistics.MinutesListened, deliveryStatistics.MinutesListened)
	assert.Equal(t, ucStatistics.TracksListened, deliveryStatistics.TracksListened)
	assert.Equal(t, ucStatistics.ArtistsListened, deliveryStatistics.ArtistsListened)
}

func TestPrivacyFromUsecaseToRepository(t *testing.T) {
	ucPrivacy := &usecase.UserPrivacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	repoPrivacy := model.PrivacyFromUsecaseToRepository(ucPrivacy)

	assert.Equal(t, ucPrivacy.IsPublicPlaylists, repoPrivacy.IsPublicPlaylists)
	assert.Equal(t, ucPrivacy.IsPublicMinutesListened, repoPrivacy.IsPublicMinutesListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteArtists, repoPrivacy.IsPublicFavoriteArtists)
	assert.Equal(t, ucPrivacy.IsPublicTracksListened, repoPrivacy.IsPublicTracksListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteTracks, repoPrivacy.IsPublicFavoriteTracks)
	assert.Equal(t, ucPrivacy.IsPublicArtistsListened, repoPrivacy.IsPublicArtistsListened)

	nilRepoPrivacy := model.PrivacyFromUsecaseToRepository(nil)
	assert.Nil(t, nilRepoPrivacy)
}

func TestPrivacyFromDeliveryToUsecase(t *testing.T) {
	deliveryPrivacy := &delivery.Privacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	ucPrivacy := model.PrivacyFromDeliveryToUsecase(deliveryPrivacy)

	assert.Equal(t, deliveryPrivacy.IsPublicPlaylists, ucPrivacy.IsPublicPlaylists)
	assert.Equal(t, deliveryPrivacy.IsPublicMinutesListened, ucPrivacy.IsPublicMinutesListened)
	assert.Equal(t, deliveryPrivacy.IsPublicFavoriteArtists, ucPrivacy.IsPublicFavoriteArtists)
	assert.Equal(t, deliveryPrivacy.IsPublicTracksListened, ucPrivacy.IsPublicTracksListened)
	assert.Equal(t, deliveryPrivacy.IsPublicFavoriteTracks, ucPrivacy.IsPublicFavoriteTracks)
	assert.Equal(t, deliveryPrivacy.IsPublicArtistsListened, ucPrivacy.IsPublicArtistsListened)

	nilUcPrivacy := model.PrivacyFromDeliveryToUsecase(nil)
	assert.Nil(t, nilUcPrivacy)
}

func TestChangeDataFromDeliveryToUsecase(t *testing.T) {
	deliveryPrivacy := &delivery.Privacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	deliveryUserChangeSettings := &delivery.UserChangeSettings{
		Privacy:     deliveryPrivacy,
		Password:    "oldpass",
		NewUsername: "newuser",
		NewEmail:    "new@example.com",
		NewPassword: "newpass",
	}

	ucUserChangeSettings := model.ChangeDataFromDeliveryToUsecase(deliveryUserChangeSettings)

	assert.Equal(t, deliveryUserChangeSettings.Password, ucUserChangeSettings.Password)
	assert.Equal(t, deliveryUserChangeSettings.NewUsername, ucUserChangeSettings.NewUsername)
	assert.Equal(t, deliveryUserChangeSettings.NewEmail, ucUserChangeSettings.NewEmail)
	assert.Equal(t, deliveryUserChangeSettings.NewPassword, ucUserChangeSettings.NewPassword)

	// Check Privacy
	assert.Equal(t, deliveryPrivacy.IsPublicPlaylists, ucUserChangeSettings.Privacy.IsPublicPlaylists)
	assert.Equal(t, deliveryPrivacy.IsPublicMinutesListened, ucUserChangeSettings.Privacy.IsPublicMinutesListened)
	assert.Equal(t, deliveryPrivacy.IsPublicFavoriteArtists, ucUserChangeSettings.Privacy.IsPublicFavoriteArtists)
	assert.Equal(t, deliveryPrivacy.IsPublicTracksListened, ucUserChangeSettings.Privacy.IsPublicTracksListened)
	assert.Equal(t, deliveryPrivacy.IsPublicFavoriteTracks, ucUserChangeSettings.Privacy.IsPublicFavoriteTracks)
	assert.Equal(t, deliveryPrivacy.IsPublicArtistsListened, ucUserChangeSettings.Privacy.IsPublicArtistsListened)
}

func TestPlaylistsFromUsecaseToDelivery(t *testing.T) {
	ucPlaylists := []*usecase.Playlist{
		{
			ID:        1,
			Title:     "Playlist 1",
			Thumbnail: "playlist1.jpg",
			Username:  "user1",
		},
		{
			ID:        2,
			Title:     "Playlist 2",
			Thumbnail: "playlist2.jpg",
			Username:  "user2",
		},
	}

	deliveryPlaylists := model.PlaylistsFromUsecaseToDelivery(ucPlaylists)

	assert.Len(t, deliveryPlaylists, 2)
	assert.Equal(t, ucPlaylists[0].ID, deliveryPlaylists[0].ID)
	assert.Equal(t, ucPlaylists[0].Title, deliveryPlaylists[0].Title)
	assert.Equal(t, ucPlaylists[0].Thumbnail, deliveryPlaylists[0].Thumbnail)
	assert.Equal(t, ucPlaylists[0].Username, deliveryPlaylists[0].Username)
	assert.Equal(t, ucPlaylists[1].ID, deliveryPlaylists[1].ID)
	assert.Equal(t, ucPlaylists[1].Title, deliveryPlaylists[1].Title)
	assert.Equal(t, ucPlaylists[1].Thumbnail, deliveryPlaylists[1].Thumbnail)
	assert.Equal(t, ucPlaylists[1].Username, deliveryPlaylists[1].Username)
}

func TestPlaylistFromUsecaseToDelivery(t *testing.T) {
	ucPlaylist := &usecase.Playlist{
		ID:        1,
		Title:     "Playlist 1",
		Thumbnail: "playlist1.jpg",
		Username:  "user1",
	}

	deliveryPlaylist := model.PlaylistFromUsecaseToDelivery(ucPlaylist)

	assert.Equal(t, ucPlaylist.ID, deliveryPlaylist.ID)
	assert.Equal(t, ucPlaylist.Title, deliveryPlaylist.Title)
	assert.Equal(t, ucPlaylist.Thumbnail, deliveryPlaylist.Thumbnail)
	assert.Equal(t, ucPlaylist.Username, deliveryPlaylist.Username)
}

func TestPlaylistFromProtoToUsecase(t *testing.T) {
	username := mockUsername1
	protoPlaylist := &playlistProto.Playlist{
		Id:        1,
		Title:     "Playlist 1",
		Thumbnail: "playlist1.jpg",
	}

	ucPlaylist := model.PlaylistFromProtoToUsecase(protoPlaylist, username)

	assert.Equal(t, protoPlaylist.Id, ucPlaylist.ID)
	assert.Equal(t, protoPlaylist.Title, ucPlaylist.Title)
	assert.Equal(t, protoPlaylist.Thumbnail, ucPlaylist.Thumbnail)
	assert.Equal(t, username, ucPlaylist.Username)
}

func TestPlaylistsFromProtoToUsecase(t *testing.T) {
	username := mockUsername1
	protoPlaylists := []*playlistProto.Playlist{
		{
			Id:        1,
			Title:     "Playlist 1",
			Thumbnail: "playlist1.jpg",
		},
		{
			Id:        2,
			Title:     "Playlist 2",
			Thumbnail: "playlist2.jpg",
		},
	}

	ucPlaylists := model.PlaylistsFromProtoToUsecase(protoPlaylists, username)

	assert.Len(t, ucPlaylists, 2)
	assert.Equal(t, protoPlaylists[0].Id, ucPlaylists[0].ID)
	assert.Equal(t, protoPlaylists[0].Title, ucPlaylists[0].Title)
	assert.Equal(t, protoPlaylists[0].Thumbnail, ucPlaylists[0].Thumbnail)
	assert.Equal(t, username, ucPlaylists[0].Username)
	assert.Equal(t, protoPlaylists[1].Id, ucPlaylists[1].ID)
	assert.Equal(t, protoPlaylists[1].Title, ucPlaylists[1].Title)
	assert.Equal(t, protoPlaylists[1].Thumbnail, ucPlaylists[1].Thumbnail)
	assert.Equal(t, username, ucPlaylists[1].Username)
}

func TestPlaylistWithIsLikedFromProtoToUsecase(t *testing.T) {
	username := mockUsername1
	protoPlaylist := &playlistProto.Playlist{
		Id:        1,
		Title:     "Playlist 1",
		Thumbnail: "playlist1.jpg",
	}
	protoPlaylistWithIsLiked := &playlistProto.PlaylistWithIsLiked{
		Playlist: protoPlaylist,
		IsLiked:  true,
	}

	ucPlaylistWithIsLiked := model.PlaylistWithIsLikedFromProtoToUsecase(protoPlaylistWithIsLiked, username)

	assert.Equal(t, protoPlaylist.Id, ucPlaylistWithIsLiked.Playlist.ID)
	assert.Equal(t, protoPlaylist.Title, ucPlaylistWithIsLiked.Playlist.Title)
	assert.Equal(t, protoPlaylist.Thumbnail, ucPlaylistWithIsLiked.Playlist.Thumbnail)
	assert.Equal(t, username, ucPlaylistWithIsLiked.Playlist.Username)
	assert.Equal(t, protoPlaylistWithIsLiked.IsLiked, ucPlaylistWithIsLiked.IsLiked)
}

func TestPlaylistWithIsLikedFromUsecaseToDelivery(t *testing.T) {
	ucPlaylist := usecase.Playlist{
		ID:        1,
		Title:     "Playlist 1",
		Thumbnail: "playlist1.jpg",
		Username:  "user1",
	}
	ucPlaylistWithIsLiked := &usecase.PlaylistWithIsLiked{
		Playlist: ucPlaylist,
		IsLiked:  true,
	}

	deliveryPlaylistWithIsLiked := model.PlaylistWithIsLikedFromUsecaseToDelivery(ucPlaylistWithIsLiked)

	assert.Equal(t, ucPlaylist.ID, deliveryPlaylistWithIsLiked.Playlist.ID)
	assert.Equal(t, ucPlaylist.Title, deliveryPlaylistWithIsLiked.Playlist.Title)
	assert.Equal(t, ucPlaylist.Thumbnail, deliveryPlaylistWithIsLiked.Playlist.Thumbnail)
	assert.Equal(t, ucPlaylist.Username, deliveryPlaylistWithIsLiked.Playlist.Username)
	assert.Equal(t, ucPlaylistWithIsLiked.IsLiked, deliveryPlaylistWithIsLiked.IsLiked)
}

func TestLikePlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	userID := mockUserID
	playlistID := mockPlaylistID2
	isLike := true

	ucLikeRequest := model.LikePlaylistRequestFromDeliveryToUsecase(userID, playlistID, isLike)

	assert.Equal(t, userID, ucLikeRequest.UserID)
	assert.Equal(t, playlistID, ucLikeRequest.PlaylistID)
	assert.Equal(t, isLike, ucLikeRequest.IsLike)
}

func TestLikePlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucLikeRequest := &usecase.LikePlaylistRequest{
		UserID:     1,
		PlaylistID: 2,
		IsLike:     true,
	}

	protoLikeRequest := model.LikePlaylistRequestFromUsecaseToProto(ucLikeRequest)

	assert.Equal(t, ucLikeRequest.UserID, protoLikeRequest.UserId)
	assert.Equal(t, ucLikeRequest.PlaylistID, protoLikeRequest.PlaylistId)
	assert.Equal(t, ucLikeRequest.IsLike, protoLikeRequest.IsLike)
}

func TestGetPlaylistsToAddRequestFromDeliveryToUsecase(t *testing.T) {
	trackID := mockTrackID
	userID := mockUserID2

	ucGetRequest := model.GetPlaylistsToAddRequestFromDeliveryToUsecase(trackID, userID)

	assert.Equal(t, trackID, ucGetRequest.TrackID)
	assert.Equal(t, userID, ucGetRequest.UserID)
}

func TestGetPlaylistsToAddRequestFromUsecaseToProto(t *testing.T) {
	ucGetRequest := &usecase.GetPlaylistsToAddRequest{
		UserID:  1,
		TrackID: 2,
	}

	protoGetRequest := model.GetPlaylistsToAddRequestFromUsecaseToProto(ucGetRequest)

	assert.Equal(t, ucGetRequest.UserID, protoGetRequest.UserId)
	assert.Equal(t, ucGetRequest.TrackID, protoGetRequest.TrackId)
}

func TestSessionIDFromProtoToUsecase(t *testing.T) {
	sessionID := mockSessionID
	protoSessionID := &authProto.SessionID{
		SessionId: sessionID,
	}

	ucSessionID := model.SessionIDFromProtoToUsecase(protoSessionID)

	assert.Equal(t, sessionID, ucSessionID)
}

func TestUserIDFromProtoToUsecase(t *testing.T) {
	userID := mockUserID12345
	protoUserID := &authProto.UserID{
		Id: userID,
	}

	ucUserID := model.UserIDFromProtoToUsecase(protoUserID)

	assert.Equal(t, userID, ucUserID)
}

func TestSessionIDFromUsecaseToProto(t *testing.T) {
	sessionID := mockSessionID

	protoSessionID := model.SessionIDFromUsecaseToProto(sessionID)

	assert.Equal(t, sessionID, protoSessionID.SessionId)
}

func TestUserIDFromUsecaseToProto(t *testing.T) {
	userID := mockUserID12345

	protoUserID := model.UserIDFromUsecaseToProto(userID)

	assert.Equal(t, userID, protoUserID.Id)
}

// Tests for remaining User-related converters
func TestRegisterDataFromUsecaseToProto(t *testing.T) {
	ucUser := &usecase.User{
		Username: mockUsername,
		Email:    mockEmail,
		Password: mockPassword,
	}

	protoRegisterData := model.RegisterDataFromUsecaseToProto(ucUser)

	assert.Equal(t, ucUser.Username, protoRegisterData.Username)
	assert.Equal(t, ucUser.Email, protoRegisterData.Email)
	assert.Equal(t, ucUser.Password, protoRegisterData.Password)
}

func TestUserFromProtoToUsecase(t *testing.T) {
	protoUser := &userProto.UserFront{
		Id:       123,
		Username: mockUsername,
		Email:    mockEmail,
		Avatar:   mockAvatarURL,
	}

	ucUser := model.UserFromProtoToUsecase(protoUser)

	assert.Equal(t, protoUser.Id, ucUser.ID)
	assert.Equal(t, protoUser.Username, ucUser.Username)
	assert.Equal(t, protoUser.Email, ucUser.Email)
	assert.Equal(t, protoUser.Avatar, ucUser.AvatarUrl)
}

func TestUserIDFromUsecaseToProtoUser(t *testing.T) {
	userID := mockUserID12345

	protoUserID := model.UserIDFromUsecaseToProtoUser(userID)

	assert.Equal(t, userID, protoUserID.Id)
}

func TestUserIDFromProtoToUsecaseUser(t *testing.T) {
	userID := mockUserID12345
	protoUserID := &userProto.UserID{
		Id: userID,
	}

	ucUserID := model.UserIDFromProtoToUsecaseUser(protoUserID)

	assert.Equal(t, userID, ucUserID)
}

func TestLoginDataFromUsecaseToProto(t *testing.T) {
	ucUser := &usecase.User{
		Username: mockUsername,
		Email:    mockEmail,
		Password: mockPassword,
	}

	protoLoginData := model.LoginDataFromUsecaseToProto(ucUser)

	assert.Equal(t, ucUser.Username, protoLoginData.Username)
	assert.Equal(t, ucUser.Email, protoLoginData.Email)
	assert.Equal(t, ucUser.Password, protoLoginData.Password)
}

func TestAvatarDataFromUsecaseToProto(t *testing.T) {
	fileURL := mockAvatarURL
	userID := mockUserID12345

	protoAvatarData := model.AvatarDataFromUsecaseToProto(fileURL, userID)

	assert.Equal(t, fileURL, protoAvatarData.AvatarPath)
	assert.Equal(t, userID, protoAvatarData.Id)
}

func TestDeleteUserFromUsecaseToProto(t *testing.T) {
	ucUser := &usecase.User{
		Username: mockUsername,
		Email:    mockEmail,
		Password: mockPassword,
	}

	protoUserDelete := model.DeleteUserFromUsecaseToProto(ucUser)

	assert.Equal(t, ucUser.Username, protoUserDelete.Username)
	assert.Equal(t, ucUser.Email, protoUserDelete.Email)
	assert.Equal(t, ucUser.Password, protoUserDelete.Password)
}

func TestUsernameFromUsecaseToProto(t *testing.T) {
	username := mockUsername

	protoUsername := model.UsernameFromUsecaseToProto(username)

	assert.Equal(t, username, protoUsername.Username)
}

func TestPrivacyFromProtoToUsecase(t *testing.T) {
	protoPrivacy := &userProto.PrivacySettings{
		Username:                mockUsername,
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	ucPrivacy := model.PrivacyFromProtoToUsecase(protoPrivacy)

	assert.Equal(t, protoPrivacy.IsPublicPlaylists, ucPrivacy.IsPublicPlaylists)
	assert.Equal(t, protoPrivacy.IsPublicMinutesListened, ucPrivacy.IsPublicMinutesListened)
	assert.Equal(t, protoPrivacy.IsPublicFavoriteArtists, ucPrivacy.IsPublicFavoriteArtists)
	assert.Equal(t, protoPrivacy.IsPublicTracksListened, ucPrivacy.IsPublicTracksListened)
	assert.Equal(t, protoPrivacy.IsPublicFavoriteTracks, ucPrivacy.IsPublicFavoriteTracks)
	assert.Equal(t, protoPrivacy.IsPublicArtistsListened, ucPrivacy.IsPublicArtistsListened)
}

func TestUserFullDataFromProtoToUsecase(t *testing.T) {
	protoPrivacy := &userProto.PrivacySettings{
		Username:                mockUsername,
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	protoUserFullData := &userProto.UserFullData{
		Username: mockUsername,
		Email:    mockEmail,
		Avatar:   mockAvatarURL,
		Privacy:  protoPrivacy,
	}

	ucUserFullData := model.UserFullDataFromProtoToUsecase(protoUserFullData)

	assert.Equal(t, protoUserFullData.Username, ucUserFullData.Username)
	assert.Equal(t, protoUserFullData.Email, ucUserFullData.Email)
	assert.Equal(t, protoUserFullData.Avatar, ucUserFullData.AvatarUrl)

	assert.Equal(t, protoPrivacy.IsPublicPlaylists, ucUserFullData.Privacy.IsPublicPlaylists)
	assert.Equal(t, protoPrivacy.IsPublicMinutesListened, ucUserFullData.Privacy.IsPublicMinutesListened)
	assert.Equal(t, protoPrivacy.IsPublicFavoriteArtists, ucUserFullData.Privacy.IsPublicFavoriteArtists)
	assert.Equal(t, protoPrivacy.IsPublicTracksListened, ucUserFullData.Privacy.IsPublicTracksListened)
	assert.Equal(t, protoPrivacy.IsPublicFavoriteTracks, ucUserFullData.Privacy.IsPublicFavoriteTracks)
	assert.Equal(t, protoPrivacy.IsPublicArtistsListened, ucUserFullData.Privacy.IsPublicArtistsListened)
}

func TestPrivacyFromUsecaseToProto(t *testing.T) {
	username := mockUsername
	ucPrivacy := &usecase.UserPrivacy{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	protoPrivacy := model.PrivacyFromUsecaseToProto(username, ucPrivacy)

	assert.Equal(t, username, protoPrivacy.Username)
	assert.Equal(t, ucPrivacy.IsPublicPlaylists, protoPrivacy.IsPublicPlaylists)
	assert.Equal(t, ucPrivacy.IsPublicMinutesListened, protoPrivacy.IsPublicMinutesListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteArtists, protoPrivacy.IsPublicFavoriteArtists)
	assert.Equal(t, ucPrivacy.IsPublicTracksListened, protoPrivacy.IsPublicTracksListened)
	assert.Equal(t, ucPrivacy.IsPublicFavoriteTracks, protoPrivacy.IsPublicFavoriteTracks)
	assert.Equal(t, ucPrivacy.IsPublicArtistsListened, protoPrivacy.IsPublicArtistsListened)
}

func TestChangeUserDataFromUsecaseToProto(t *testing.T) {
	username := mockUsername
	ucUserChangeSettings := &usecase.UserChangeSettings{
		Password:    mockOldPassword,
		NewUsername: mockNewUsername,
		NewEmail:    mockNewEmail,
		NewPassword: mockNewPassword,
	}

	protoChangeUserData := model.ChangeUserDataFromUsecaseToProto(username, ucUserChangeSettings)

	assert.Equal(t, username, protoChangeUserData.Username)
	assert.Equal(t, ucUserChangeSettings.NewUsername, protoChangeUserData.NewUsername)
	assert.Equal(t, ucUserChangeSettings.NewEmail, protoChangeUserData.NewEmail)
	assert.Equal(t, ucUserChangeSettings.NewPassword, protoChangeUserData.NewPassword)
	assert.Equal(t, ucUserChangeSettings.Password, protoChangeUserData.Password)
}

func TestFileKeyFromUsecaseToProto(t *testing.T) {
	avatarURL := mockAvatarURL

	protoFileKey := model.FileKeyFromUsecaseToProto(avatarURL)

	assert.Equal(t, avatarURL, protoFileKey.FileKey)
}

func TestAvatarUrlFromProtoToUsecase(t *testing.T) {
	avatarURL := mockAvatarURL
	protoAvatarURL := &userProto.AvatarUrl{
		Url: avatarURL,
	}

	ucAvatarURL := model.AvatarUrlFromProtoToUsecase(protoAvatarURL)

	assert.Equal(t, avatarURL, ucAvatarURL)
}

func TestAvatarImageFromUsecaseToProto(t *testing.T) {
	username := mockUsername
	image := []byte("image-data")

	protoAvatarImage := model.AvatarImageFromUsecaseToProto(username, image)

	assert.Equal(t, username, protoAvatarImage.Username)
	assert.Equal(t, image, protoAvatarImage.Image)
}

func TestFileKeyFromProtoToUsecase(t *testing.T) {
	fileKey := "file-key"
	protoFileKey := &userProto.FileKey{
		FileKey: fileKey,
	}

	ucFileKey := model.FileKeyFromProtoToUsecase(protoFileKey)

	assert.Equal(t, fileKey, ucFileKey)
}

func TestUpdatePlaylistsPublisityByUserIDRequestFromUsecaseToProto(t *testing.T) {
	isPublic := true
	userID := mockUserID

	protoRequest := model.UpdatePlaylistsPublisityByUserIDRequestFromUsecaseToProto(isPublic, userID)

	assert.Equal(t, isPublic, protoRequest.IsPublic)
	assert.Equal(t, userID, protoRequest.UserId)
}

func TestUploadPlaylistThumbnailRequestFromUsecaseToProto(t *testing.T) {
	title := mockPlaylistTitle
	thumbnail := []byte("thumbnail-data")

	protoRequest := model.UploadPlaylistThumbnailRequestFromUsecaseToProto(title, thumbnail)

	assert.Equal(t, title, protoRequest.Title)
	assert.Equal(t, thumbnail, protoRequest.Thumbnail)
}

func TestCreatePlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucRequest := &usecase.CreatePlaylistRequest{
		Title:  mockPlaylistTitle,
		UserID: mockUserID,
	}
	thumbnail := mockThumbnail
	isPublic := true

	protoRequest := model.CreatePlaylistRequestFromUsecaseToProto(ucRequest, thumbnail, isPublic)

	assert.Equal(t, ucRequest.Title, protoRequest.Title)
	assert.Equal(t, ucRequest.UserID, protoRequest.UserId)
	assert.Equal(t, thumbnail, protoRequest.Thumbnail)
	assert.Equal(t, isPublic, protoRequest.IsPublic)
}

func TestCreatePlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	deliveryRequest := &delivery.CreatePlaylistRequest{
		Title:     mockPlaylistTitle,
		Thumbnail: []byte(mockThumbnail),
	}
	userID := mockUserID

	ucRequest := model.CreatePlaylistRequestFromDeliveryToUsecase(deliveryRequest, userID)

	assert.Equal(t, deliveryRequest.Title, ucRequest.Title)
	assert.Equal(t, userID, ucRequest.UserID)
	assert.Equal(t, deliveryRequest.Thumbnail, ucRequest.Thumbnail)
}

func TestAddTrackToPlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	deliveryRequest := &delivery.AddTrackToPlaylistRequest{
		TrackID: mockTrackID,
	}
	userID := mockUserID2
	playlistID := mockTrackID3

	ucRequest := model.AddTrackToPlaylistRequestFromDeliveryToUsecase(deliveryRequest, userID, playlistID)

	assert.Equal(t, deliveryRequest.TrackID, ucRequest.TrackID)
	assert.Equal(t, userID, ucRequest.UserID)
	assert.Equal(t, playlistID, ucRequest.PlaylistID)
}

func TestRemoveTrackFromPlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	trackID := mockTrackID
	userID := mockUserID2
	playlistID := mockTrackID3

	ucRequest := model.RemoveTrackFromPlaylistRequestFromDeliveryToUsecase(trackID, userID, playlistID)

	assert.Equal(t, trackID, ucRequest.TrackID)
	assert.Equal(t, userID, ucRequest.UserID)
	assert.Equal(t, playlistID, ucRequest.PlaylistID)
}

func TestAddTrackToPlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucRequest := &usecase.AddTrackToPlaylistRequest{
		UserID:     mockUserID,
		PlaylistID: mockPlaylistID2,
		TrackID:    mockTrackID3,
	}

	protoRequest := model.AddTrackToPlaylistRequestFromUsecaseToProto(ucRequest)

	assert.Equal(t, ucRequest.PlaylistID, protoRequest.PlaylistId)
	assert.Equal(t, ucRequest.TrackID, protoRequest.TrackId)
	assert.Equal(t, ucRequest.UserID, protoRequest.UserId)
}

func TestRemoveTrackFromPlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucRequest := &usecase.RemoveTrackFromPlaylistRequest{
		UserID:     mockUserID,
		PlaylistID: mockPlaylistID2,
		TrackID:    mockTrackID3,
	}

	protoRequest := model.RemoveTrackFromPlaylistRequestFromUsecaseToProto(ucRequest)

	assert.Equal(t, ucRequest.PlaylistID, protoRequest.PlaylistId)
	assert.Equal(t, ucRequest.TrackID, protoRequest.TrackId)
	assert.Equal(t, ucRequest.UserID, protoRequest.UserId)
}

func TestUpdatePlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucRequest := &usecase.UpdatePlaylistRequest{
		UserID:     mockUserID,
		PlaylistID: mockPlaylistID2,
		Title:      "Updated Title",
	}
	thumbnail := "updated-thumbnail.jpg"

	protoRequest := model.UpdatePlaylistRequestFromUsecaseToProto(ucRequest, thumbnail)

	assert.Equal(t, ucRequest.PlaylistID, protoRequest.Id)
	assert.Equal(t, ucRequest.Title, protoRequest.Title)
	assert.Equal(t, thumbnail, protoRequest.Thumbnail)
	assert.Equal(t, ucRequest.UserID, protoRequest.UserId)
}

func TestUpdatePlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	deliveryRequest := &delivery.UpdatePlaylistRequest{
		Title:     "Updated Title",
		Thumbnail: []byte("updated-thumbnail.jpg"),
	}
	userID := mockUserID
	playlistID := mockPlaylistID2

	ucRequest := model.UpdatePlaylistRequestFromDeliveryToUsecase(deliveryRequest, userID, playlistID)

	assert.Equal(t, userID, ucRequest.UserID)
	assert.Equal(t, playlistID, ucRequest.PlaylistID)
	assert.Equal(t, deliveryRequest.Title, ucRequest.Title)
	assert.Equal(t, deliveryRequest.Thumbnail, ucRequest.Thumbnail)
}

func TestRemovePlaylistRequestFromUsecaseToProto(t *testing.T) {
	ucRequest := &usecase.RemovePlaylistRequest{
		UserID:     mockUserID,
		PlaylistID: mockPlaylistID2,
	}

	protoRequest := model.RemovePlaylistRequestFromUsecaseToProto(ucRequest)

	assert.Equal(t, ucRequest.UserID, protoRequest.UserId)
	assert.Equal(t, ucRequest.PlaylistID, protoRequest.PlaylistId)
}

func TestRemovePlaylistRequestFromDeliveryToUsecase(t *testing.T) {
	playlistID := mockPlaylistID
	userID := mockUserID2

	ucRequest := model.RemovePlaylistRequestFromDeliveryToUsecase(playlistID, userID)

	assert.Equal(t, playlistID, ucRequest.PlaylistID)
	assert.Equal(t, userID, ucRequest.UserID)
}

func TestGetPlaylistsToAddResponseFromProtoToUsecase(t *testing.T) {
	username := mockUsername
	protoPlaylist1 := &playlistProto.Playlist{
		Id:        mockPlaylistID,
		Title:     "Playlist 1",
		Thumbnail: "thumbnail1.jpg",
	}
	protoPlaylist2 := &playlistProto.Playlist{
		Id:        mockPlaylistID2,
		Title:     "Playlist 2",
		Thumbnail: "thumbnail2.jpg",
	}

	protoPlaylistWithFlag1 := &playlistProto.PlaylistWithIsIncludedTrack{
		Playlist:        protoPlaylist1,
		IsIncludedTrack: true,
	}
	protoPlaylistWithFlag2 := &playlistProto.PlaylistWithIsIncludedTrack{
		Playlist:        protoPlaylist2,
		IsIncludedTrack: false,
	}

	protoResponse := &playlistProto.GetPlaylistsToAddResponse{
		Playlists: []*playlistProto.PlaylistWithIsIncludedTrack{
			protoPlaylistWithFlag1,
			protoPlaylistWithFlag2,
		},
	}

	ucPlaylists := model.GetPlaylistsToAddResponseFromProtoToUsecase(protoResponse, username)

	assert.Len(t, ucPlaylists, 2)
	assert.Equal(t, protoPlaylist1.Id, ucPlaylists[0].Playlist.ID)
	assert.Equal(t, protoPlaylist1.Title, ucPlaylists[0].Playlist.Title)
	assert.Equal(t, protoPlaylist1.Thumbnail, ucPlaylists[0].Playlist.Thumbnail)
	assert.Equal(t, username, ucPlaylists[0].Playlist.Username)
	assert.Equal(t, protoPlaylistWithFlag1.IsIncludedTrack, ucPlaylists[0].IsIncluded)

	assert.Equal(t, protoPlaylist2.Id, ucPlaylists[1].Playlist.ID)
	assert.Equal(t, protoPlaylist2.Title, ucPlaylists[1].Playlist.Title)
	assert.Equal(t, protoPlaylist2.Thumbnail, ucPlaylists[1].Playlist.Thumbnail)
	assert.Equal(t, username, ucPlaylists[1].Playlist.Username)
	assert.Equal(t, protoPlaylistWithFlag2.IsIncludedTrack, ucPlaylists[1].IsIncluded)
}

func TestPlaylistsWithIsIncludedTrackFromUsecaseToDelivery(t *testing.T) {
	ucPlaylist1 := usecase.Playlist{
		ID:        mockPlaylistID,
		Title:     "Playlist 1",
		Thumbnail: "thumbnail1.jpg",
		Username:  "user1",
	}
	ucPlaylist2 := usecase.Playlist{
		ID:        mockPlaylistID2,
		Title:     "Playlist 2",
		Thumbnail: "thumbnail2.jpg",
		Username:  "user2",
	}

	ucPlaylistWithFlag1 := &usecase.PlaylistWithIsIncludedTrack{
		Playlist:   ucPlaylist1,
		IsIncluded: true,
	}
	ucPlaylistWithFlag2 := &usecase.PlaylistWithIsIncludedTrack{
		Playlist:   ucPlaylist2,
		IsIncluded: false,
	}

	ucPlaylists := []*usecase.PlaylistWithIsIncludedTrack{
		ucPlaylistWithFlag1,
		ucPlaylistWithFlag2,
	}

	deliveryPlaylists := model.PlaylistsWithIsIncludedTrackFromUsecaseToDelivery(ucPlaylists)

	assert.Len(t, deliveryPlaylists, 2)
	assert.Equal(t, ucPlaylist1.ID, deliveryPlaylists[0].Playlist.ID)
	assert.Equal(t, ucPlaylist1.Title, deliveryPlaylists[0].Playlist.Title)
	assert.Equal(t, ucPlaylist1.Thumbnail, deliveryPlaylists[0].Playlist.Thumbnail)
	assert.Equal(t, ucPlaylist1.Username, deliveryPlaylists[0].Playlist.Username)
	assert.Equal(t, ucPlaylistWithFlag1.IsIncluded, deliveryPlaylists[0].IsIncluded)

	assert.Equal(t, ucPlaylist2.ID, deliveryPlaylists[1].Playlist.ID)
	assert.Equal(t, ucPlaylist2.Title, deliveryPlaylists[1].Playlist.Title)
	assert.Equal(t, ucPlaylist2.Thumbnail, deliveryPlaylists[1].Playlist.Thumbnail)
	assert.Equal(t, ucPlaylist2.Username, deliveryPlaylists[1].Playlist.Username)
	assert.Equal(t, ucPlaylistWithFlag2.IsIncluded, deliveryPlaylists[1].IsIncluded)
}

func TestAlbumFromProtoToUsecaseEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		protoAlbumType albumProto.AlbumType
		expected       usecase.AlbumType
	}{
		{
			name:           "AlbumTypeSingle",
			protoAlbumType: albumProto.AlbumType_AlbumTypeSingle,
			expected:       usecase.AlbumTypeSingle,
		},
		{
			name:           "AlbumTypeCompilation",
			protoAlbumType: albumProto.AlbumType_AlbumTypeCompilation,
			expected:       usecase.AlbumTypeCompilation,
		},
		{
			name:           "Invalid type defaults to AlbumTypeAlbum",
			protoAlbumType: albumProto.AlbumType(999),
			expected:       usecase.AlbumTypeAlbum,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoAlbum := &albumProto.Album{
				Id:          mockAlbumID,
				Title:       "Test Album",
				Type:        tt.protoAlbumType,
				Thumbnail:   mockThumbnail,
				ReleaseDate: timestamppb.New(time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC)),
				IsFavorite:  true,
			}

			ucAlbum := model.AlbumFromProtoToUsecase(protoAlbum)

			assert.Equal(t, protoAlbum.Id, ucAlbum.ID)
			assert.Equal(t, protoAlbum.Title, ucAlbum.Title)
			assert.Equal(t, tt.expected, ucAlbum.Type)
			assert.Equal(t, protoAlbum.Thumbnail, ucAlbum.Thumbnail)
			assert.Equal(t, protoAlbum.ReleaseDate.AsTime(), ucAlbum.ReleaseDate)
			assert.Equal(t, protoAlbum.IsFavorite, ucAlbum.IsLiked)
		})
	}
}

func TestTrackRequestFromDeliveryToUsecase(t *testing.T) {
	deliveryTrack := &delivery.CreateTrackRequest{
		Title: mockTrackTitle,
		Track: []byte("track-data"),
	}

	ucTrack := model.TrackRequestFromDeliveryToUsecase(deliveryTrack)

	assert.Equal(t, deliveryTrack.Title, ucTrack.Title)
	assert.Equal(t, deliveryTrack.Track, ucTrack.Track)
}

func TestNewAlbumFromDeliveryToUsecase(t *testing.T) {
	deliveryTracks := []*delivery.CreateTrackRequest{
		{
			Title: "Track 1",
			Track: []byte("track1-data"),
		},
		{
			Title: "Track 2",
			Track: []byte("track2-data"),
		},
	}

	deliveryAlbum := &delivery.CreateAlbumRequest{
		Title:      mockAlbumTitle,
		Image:      []byte("album-image"),
		Type:       "album",
		LabelID:    1,
		Tracks:     deliveryTracks,
		ArtistsIDs: []int64{1, 2},
	}

	ucAlbum := model.NewAlbumFromDeliveryToUsecase(deliveryAlbum)

	assert.Equal(t, deliveryAlbum.Title, ucAlbum.Title)
	assert.Equal(t, deliveryAlbum.Image, ucAlbum.Image)
	assert.Equal(t, deliveryAlbum.Type, ucAlbum.Type)
	assert.Equal(t, deliveryAlbum.LabelID, ucAlbum.LabelID)
	assert.Equal(t, deliveryAlbum.ArtistsIDs, ucAlbum.ArtistsIDs)
	assert.Len(t, ucAlbum.Tracks, 2)
	assert.Equal(t, deliveryTracks[0].Title, ucAlbum.Tracks[0].Title)
	assert.Equal(t, deliveryTracks[0].Track, ucAlbum.Tracks[0].Track)
	assert.Equal(t, deliveryTracks[1].Title, ucAlbum.Tracks[1].Title)
	assert.Equal(t, deliveryTracks[1].Track, ucAlbum.Tracks[1].Track)
}

func TestAlbumTypeConverter(t *testing.T) {
	tests := []struct {
		name      string
		albumType string
		expected  delivery.AlbumType
	}{
		{
			name:      "AlbumTypeAlbum",
			albumType: string(usecase.AlbumTypeAlbum),
			expected:  delivery.AlbumTypeAlbum,
		},
		{
			name:      "AlbumTypeEP",
			albumType: string(usecase.AlbumTypeEP),
			expected:  delivery.AlbumTypeEP,
		},
		{
			name:      "AlbumTypeSingle",
			albumType: string(usecase.AlbumTypeSingle),
			expected:  delivery.AlbumTypeSingle,
		},
		{
			name:      "AlbumTypeCompilation",
			albumType: string(usecase.AlbumTypeCompilation),
			expected:  delivery.AlbumTypeCompilation,
		},
		{
			name:      "Unknown type defaults to AlbumTypeAlbum",
			albumType: "unknown",
			expected:  delivery.AlbumTypeAlbum,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.AlbumTypeConverter(tt.albumType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArtistLoadFromUsecaseToProto(t *testing.T) {
	ucArtist := &usecase.ArtistLoad{
		Title:   mockArtistTitle,
		Image:   []byte("artist-image"),
		LabelID: 123,
	}

	protoArtist := model.ArtistLoadFromUsecaseToProto(ucArtist)

	assert.Equal(t, ucArtist.Title, protoArtist.Title)
	assert.Equal(t, ucArtist.Image, protoArtist.Image)
	assert.Equal(t, ucArtist.LabelID, protoArtist.LabelId)
}

func TestArtistLoadFromDeliveryToUsecase(t *testing.T) {
	deliveryArtist := &delivery.CreateArtistRequest{
		Title:   mockArtistTitle,
		Image:   []byte("artist-image"),
		LabelID: 123,
	}

	ucArtist := model.ArtistLoadFromDeliveryToUsecase(deliveryArtist)

	assert.Equal(t, deliveryArtist.Title, ucArtist.Title)
	assert.Equal(t, deliveryArtist.Image, ucArtist.Image)
	assert.Equal(t, deliveryArtist.LabelID, ucArtist.LabelID)
}

func TestArtistDeleteFromDeliveryToUsecase(t *testing.T) {
	deleteArtist := &delivery.DeleteArtistRequest{
		ArtistID: mockArtistID,
	}

	ucArtistDelete := model.ArtistDeleteFromDeliveryToUsecase(deleteArtist)

	assert.Equal(t, deleteArtist.ArtistID, ucArtistDelete.ArtistID)
}

func TestLabelIDFromUsecaseToProto(t *testing.T) {
	labelID := int64(123)

	protoLabelID := model.LabelIDFromUsecaseToProto(labelID)

	assert.Equal(t, labelID, protoLabelID.Id)
}

func TestLabelIDFromProtoToUsecase(t *testing.T) {
	protoLabelID := &userProto.LabelID{
		Id: 123,
	}

	labelID := model.LabelIDFromProtoToUsecase(protoLabelID)

	assert.Equal(t, protoLabelID.Id, labelID)
}

func TestLabelFromUsecaseToProto(t *testing.T) {
	ucLabel := &usecase.Label{
		Id:      123,
		Name:    "Test Label",
		Members: []string{"user1", "user2"},
	}

	protoLabel := model.LabelFromUsecaseToProto(ucLabel)

	assert.Equal(t, ucLabel.Id, protoLabel.Id)
	assert.Equal(t, ucLabel.Name, protoLabel.Name)
	assert.Equal(t, ucLabel.Members, protoLabel.Usernames)
}

func TestLabelFromProtoToUsecase(t *testing.T) {
	protoLabel := &userProto.Label{
		Id:        123,
		Name:      "Test Label",
		Usernames: []string{"user1", "user2"},
	}

	ucLabel := model.LabelFromProtoToUsecase(protoLabel)

	assert.Equal(t, protoLabel.Id, ucLabel.Id)
	assert.Equal(t, protoLabel.Name, ucLabel.Name)
	assert.Equal(t, protoLabel.Usernames, ucLabel.Members)
}

func TestLabelFromDeliveryToUsecase(t *testing.T) {
	deliveryLabel := &delivery.Label{
		Id:        123,
		LabelName: "Test Label",
		Usernames: []string{"user1", "user2"},
	}

	ucLabel := model.LabelFromDeliveryToUsecase(deliveryLabel)

	assert.Equal(t, deliveryLabel.Id, ucLabel.Id)
	assert.Equal(t, deliveryLabel.LabelName, ucLabel.Name)
	assert.Equal(t, deliveryLabel.Usernames, ucLabel.Members)
}

func TestLabelFromUsecaseToDelivery(t *testing.T) {
	ucLabel := &usecase.Label{
		Id:      123,
		Name:    "Test Label",
		Members: []string{"user1", "user2"},
	}

	deliveryLabel := model.LabelFromUsecaseToDelivery(ucLabel)

	assert.Equal(t, ucLabel.Id, deliveryLabel.Id)
	assert.Equal(t, ucLabel.Name, deliveryLabel.LabelName)
	assert.Equal(t, ucLabel.Members, deliveryLabel.Usernames)
}

func TestMembersFromProtoToUsecase(t *testing.T) {
	protoUsers := &userProto.UsersToFront{
		Users: []*userProto.UserFront{
			{
				Id:       1,
				Username: "user1",
				Email:    "user1@example.com",
				Avatar:   "avatar1.jpg",
			},
			{
				Id:       2,
				Username: "user2",
				Email:    "user2@example.com",
				Avatar:   "avatar2.jpg",
			},
		},
	}

	ucUsers := model.MembersFromProtoToUsecase(protoUsers)

	assert.Len(t, ucUsers, 2)
	assert.Equal(t, protoUsers.Users[0].Id, ucUsers[0].ID)
	assert.Equal(t, protoUsers.Users[0].Username, ucUsers[0].Username)
	assert.Equal(t, protoUsers.Users[0].Email, ucUsers[0].Email)
	assert.Equal(t, protoUsers.Users[0].Avatar, ucUsers[0].AvatarUrl)
	assert.Equal(t, protoUsers.Users[1].Id, ucUsers[1].ID)
	assert.Equal(t, protoUsers.Users[1].Username, ucUsers[1].Username)
	assert.Equal(t, protoUsers.Users[1].Email, ucUsers[1].Email)
	assert.Equal(t, protoUsers.Users[1].Avatar, ucUsers[1].AvatarUrl)
}

func TestArtistEditFromDeliveryToUsecase(t *testing.T) {
	editRequest := &delivery.EditArtistRequest{
		ArtistID: mockArtistID,
		NewTitle: "New Artist Title",
		Image:    []byte("new-image"),
		LabelID:  123,
	}

	ucArtistEdit := model.ArtistEditFromDeliveryToUsecase(editRequest)

	assert.Equal(t, editRequest.ArtistID, ucArtistEdit.ArtistID)
	assert.Equal(t, editRequest.NewTitle, ucArtistEdit.NewTitle)
	assert.Equal(t, editRequest.Image, ucArtistEdit.Image)
	assert.Equal(t, editRequest.LabelID, ucArtistEdit.LabelID)
}

func TestTrackLoadFromUsecaseToProto(t *testing.T) {
	ucTrack := &usecase.CreateTrackRequest{
		Title: mockTrackTitle,
		Track: []byte("track-data"),
	}

	protoTrack := model.TrackLoadFromUsecaseToProto(ucTrack)

	assert.Equal(t, ucTrack.Title, protoTrack.Title)
	assert.Equal(t, ucTrack.Track, protoTrack.File)
}

func TestTrackListLoadFromUsecaseToProto(t *testing.T) {
	ucTracks := []*usecase.CreateTrackRequest{
		{
			Title: "Track 1",
			Track: []byte("track1-data"),
		},
		{
			Title: "Track 2",
			Track: []byte("track2-data"),
		},
	}

	protoTracks := model.TrackListLoadFromUsecaseToProto(ucTracks)

	assert.Len(t, protoTracks, 2)
	assert.Equal(t, ucTracks[0].Title, protoTracks[0].Title)
	assert.Equal(t, ucTracks[0].Track, protoTracks[0].File)
	assert.Equal(t, ucTracks[1].Title, protoTracks[1].Title)
	assert.Equal(t, ucTracks[1].Track, protoTracks[1].File)
}

func TestTracksIdsFromProtoToUsecase(t *testing.T) {
	protoTracks := &trackProto.TrackIdsList{
		Ids: []*trackProto.TrackID{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	trackIds := model.TracksIdsFromProtoToUsecase(protoTracks)

	assert.Len(t, trackIds, 3)
	assert.Equal(t, protoTracks.Ids[0].Id, trackIds[0])
	assert.Equal(t, protoTracks.Ids[1].Id, trackIds[1])
	assert.Equal(t, protoTracks.Ids[2].Id, trackIds[2])
}

func TestTracksIdsFromUsecaseToProtoArtist(t *testing.T) {
	trackIDs := []int64{1, 2, 3}

	protoTrackIDList := model.TracksIdsFromUsecaseToProtoArtist(trackIDs)

	assert.Len(t, protoTrackIDList.Ids, 3)
	assert.Equal(t, trackIDs[0], protoTrackIDList.Ids[0].Id)
	assert.Equal(t, trackIDs[1], protoTrackIDList.Ids[1].Id)
	assert.Equal(t, trackIDs[2], protoTrackIDList.Ids[2].Id)
}

func TestCreateJamResponseFromUsecaseToDelivery(t *testing.T) {
	ucCreateJam := &usecase.CreateJamResponse{
		RoomID: "room123",
		HostID: "host456",
	}

	deliveryCreateJam := model.CreateJamResponseFromUsecaseToDelivery(ucCreateJam)

	assert.Equal(t, ucCreateJam.RoomID, deliveryCreateJam.RoomID)
	assert.Equal(t, ucCreateJam.HostID, deliveryCreateJam.HostID)
}

func TestCreateJamRequestFromDeliveryToUsecase(t *testing.T) {
	deliveryRequest := &delivery.CreateJamRequest{
		TrackID:  "track123",
		Position: 100,
	}

	ucRequest := model.CreateJamRequestFromDeliveryToUsecase(deliveryRequest, "user123")

	assert.Equal(t, "user123", ucRequest.UserID)
	assert.Equal(t, deliveryRequest.TrackID, ucRequest.TrackID)
	assert.Equal(t, deliveryRequest.Position, ucRequest.Position)
}

func TestJamMessageFromUsecaseToDelivery(t *testing.T) {
	ucJamMessage := &usecase.JamMessage{
		Type:       "play",
		TrackID:    "track123",
		Position:   100,
		Paused:     false,
		UserID:     "user123",
		HostID:     "host123",
		Users:      []string{"user1", "user2"},
		Loaded:     map[string]bool{"user1": true, "user2": false},
		UserImages: map[string]string{"user1": "image1", "user2": "image2"},
		UserNames:  map[string]string{"user1": "name1", "user2": "name2"},
	}

	deliveryJamMessage := model.JamMessageFromUsecaseToDelivery(ucJamMessage)

	assert.Equal(t, ucJamMessage.Type, deliveryJamMessage.Type)
	assert.Equal(t, ucJamMessage.TrackID, deliveryJamMessage.TrackID)
	assert.Equal(t, ucJamMessage.Position, deliveryJamMessage.Position)
	assert.Equal(t, ucJamMessage.Paused, deliveryJamMessage.Paused)
	assert.Equal(t, ucJamMessage.UserID, deliveryJamMessage.UserID)
	assert.Equal(t, ucJamMessage.HostID, deliveryJamMessage.HostID)
	assert.Equal(t, ucJamMessage.Users, deliveryJamMessage.Users)
	assert.Equal(t, ucJamMessage.Loaded, deliveryJamMessage.Loaded)
	assert.Equal(t, ucJamMessage.UserImages, deliveryJamMessage.UserImages)
	assert.Equal(t, ucJamMessage.UserNames, deliveryJamMessage.UserNames)
}

func TestJamMessageFromRepositoryToUsecase(t *testing.T) {
	repoJamMessage := &repository.JamMessage{
		Type:       "play",
		TrackID:    "track123",
		Position:   100,
		Paused:     false,
		UserID:     "user123",
		HostID:     "host123",
		Users:      []string{"user1", "user2"},
		Loaded:     map[string]bool{"user1": true, "user2": false},
		UserImages: map[string]string{"user1": "image1", "user2": "image2"},
		UserNames:  map[string]string{"user1": "name1", "user2": "name2"},
	}

	ucJamMessage := model.JamMessageFromRepositoryToUsecase(repoJamMessage)

	assert.Equal(t, repoJamMessage.Type, ucJamMessage.Type)
	assert.Equal(t, repoJamMessage.TrackID, ucJamMessage.TrackID)
	assert.Equal(t, repoJamMessage.Position, ucJamMessage.Position)
	assert.Equal(t, repoJamMessage.Paused, ucJamMessage.Paused)
	assert.Equal(t, repoJamMessage.UserID, ucJamMessage.UserID)
	assert.Equal(t, repoJamMessage.HostID, ucJamMessage.HostID)
	assert.Equal(t, repoJamMessage.Users, ucJamMessage.Users)
	assert.Equal(t, repoJamMessage.Loaded, ucJamMessage.Loaded)
	assert.Equal(t, repoJamMessage.UserImages, ucJamMessage.UserImages)
	assert.Equal(t, repoJamMessage.UserNames, ucJamMessage.UserNames)
}

func TestJamMessageFromDeliveryToUsecase(t *testing.T) {
	deliveryJamMessage := &delivery.JamMessage{
		Type:       "play",
		TrackID:    "track123",
		Position:   100,
		Paused:     false,
		UserID:     "user123",
		HostID:     "host123",
		Users:      []string{"user1", "user2"},
		Loaded:     map[string]bool{"user1": true, "user2": false},
		UserImages: map[string]string{"user1": "image1", "user2": "image2"},
		UserNames:  map[string]string{"user1": "name1", "user2": "name2"},
	}

	ucJamMessage := model.JamMessageFromDeliveryToUsecase(deliveryJamMessage)

	assert.Equal(t, deliveryJamMessage.Type, ucJamMessage.Type)
	assert.Equal(t, deliveryJamMessage.TrackID, ucJamMessage.TrackID)
	assert.Equal(t, deliveryJamMessage.Position, ucJamMessage.Position)
	assert.Equal(t, deliveryJamMessage.Paused, ucJamMessage.Paused)
	assert.Equal(t, "", ucJamMessage.UserID)
	assert.Equal(t, deliveryJamMessage.HostID, ucJamMessage.HostID)
	assert.Equal(t, deliveryJamMessage.Users, ucJamMessage.Users)
	assert.Equal(t, deliveryJamMessage.Loaded, ucJamMessage.Loaded)
	assert.Equal(t, deliveryJamMessage.UserImages, ucJamMessage.UserImages)
	assert.Equal(t, deliveryJamMessage.UserNames, ucJamMessage.UserNames)
}
