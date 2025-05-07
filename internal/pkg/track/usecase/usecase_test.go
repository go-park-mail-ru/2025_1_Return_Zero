package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	trackUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetAllTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockTrackClient.EXPECT().GetAllTracks(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetAllTracks(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
}

func TestGetAllTracksError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	mockTrackClient.EXPECT().GetAllTracks(gomock.Any(), gomock.Any()).Return(nil, errors.New("test error"))

	tracks, err := trackUC.GetAllTracks(ctx, filters)

	assert.Error(t, err)
	assert.Nil(t, tracks)
}

func TestGetTrackByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	trackID := int64(1)

	trackDetailed := &track.TrackDetailed{
		Track: &track.Track{
			Id:      1,
			Title:   "Test Track",
			AlbumId: 1,
		},
		FileUrl: "test_url",
	}

	artistList := &artist.ArtistWithRoleList{
		Artists: []*artist.ArtistWithRole{
			{
				Id:   1,
				Role: "singer",
			},
		},
	}

	albumTitle := &album.AlbumTitle{
		Title: "Test Album",
	}

	mockTrackClient.EXPECT().GetTrackByID(gomock.Any(), gomock.Any()).Return(trackDetailed, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackID(gomock.Any(), gomock.Any()).Return(artistList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByID(gomock.Any(), gomock.Any()).Return(albumTitle, nil)

	track, err := trackUC.GetTrackByID(ctx, trackID)

	assert.NoError(t, err)
	assert.NotNil(t, track)
	assert.Equal(t, int64(1), track.ID)
	assert.Equal(t, "Test Track", track.Title)
	assert.Equal(t, "test_url", track.FileUrl)
}

func TestLikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	likeRequest := &usecase.TrackLikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}

	mockTrackClient.EXPECT().LikeTrack(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

	err := trackUC.LikeTrack(ctx, likeRequest)

	assert.NoError(t, err)
}

func TestCreateStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	stream := &usecase.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	albumID := &track.AlbumID{Id: 1}
	streamID := &track.StreamID{Id: 1}
	artists := &artist.ArtistWithRoleList{
		Artists: []*artist.ArtistWithRole{
			{
				Id:   1,
				Role: "singer",
			},
		},
	}

	mockTrackClient.EXPECT().CreateStream(gomock.Any(), gomock.Any()).Return(streamID, nil)
	mockTrackClient.EXPECT().GetAlbumIDByTrackID(gomock.Any(), gomock.Any()).Return(albumID, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackID(gomock.Any(), gomock.Any()).Return(artists, nil)
	mockArtistClient.EXPECT().CreateStreamsByArtistIDs(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)
	mockAlbumClient.EXPECT().CreateStream(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

	id, err := trackUC.CreateStream(ctx, stream)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestUpdateStreamDuration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	stream := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 120,
	}

	mockTrackClient.EXPECT().UpdateStreamDuration(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

	err := trackUC.UpdateStreamDuration(ctx, stream)

	assert.NoError(t, err)
}

func TestGetPlaylistTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	playlistID := int64(1)

	playlistTrackIDs := &playlist.GetPlaylistTrackIdsResponse{
		TrackIds: []int64{1},
	}

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockPlaylistClient.EXPECT().GetPlaylistTrackIds(gomock.Any(), gomock.Any()).Return(playlistTrackIDs, nil)
	mockTrackClient.EXPECT().GetTracksByIDs(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetPlaylistTracks(ctx, playlistID)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
}

func TestSearchTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	query := "test"

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockTrackClient.EXPECT().SearchTracks(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.SearchTracks(ctx, query)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
}

func TestGetTracksByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	artistID := int64(1)
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	artistTrackIDs := &artist.TrackIDList{
		Ids: []*artist.TrackID{
			{Id: 1},
		},
	}

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockArtistClient.EXPECT().GetTrackIDsByArtistID(gomock.Any(), gomock.Any()).Return(artistTrackIDs, nil)
	mockTrackClient.EXPECT().GetTracksByIDsFiltered(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetTracksByArtistID(ctx, artistID, filters)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
}

func TestGetLastListenedTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	userID := int64(1)
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockTrackClient.EXPECT().GetLastListenedTracks(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetLastListenedTracks(ctx, userID, filters)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
}

func TestGetTracksByAlbumID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	albumID := int64(1)

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitle := &album.AlbumTitle{
		Title: "Test Album",
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockTrackClient.EXPECT().GetTracksByAlbumID(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByID(gomock.Any(), gomock.Any()).Return(albumTitle, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetTracksByAlbumID(ctx, albumID)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
}

func TestGetFavoriteTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	username := "testuser"
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	userID := &user.UserID{
		Id: 1,
	}

	userPrivacy := &user.PrivacySettings{
		IsPublicFavoriteTracks: true,
	}

	trackList := &track.TrackList{
		Tracks: []*track.Track{
			{
				Id:      1,
				Title:   "Test Track",
				AlbumId: 1,
			},
		},
	}

	albumTitleMap := &album.AlbumTitleMap{
		Titles: map[int64]*album.AlbumTitle{
			1: {Title: "Test Album"},
		},
	}

	artistsMap := &artist.ArtistWithRoleMap{
		Artists: map[int64]*artist.ArtistWithRoleList{
			1: {
				Artists: []*artist.ArtistWithRole{
					{
						Id:   1,
						Role: "singer",
					},
				},
			},
		},
	}

	mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), gomock.Any()).Return(userID, nil)
	mockUserClient.EXPECT().GetUserPrivacyByID(gomock.Any(), gomock.Any()).Return(userPrivacy, nil)
	mockTrackClient.EXPECT().GetFavoriteTracks(gomock.Any(), gomock.Any()).Return(trackList, nil)
	mockAlbumClient.EXPECT().GetAlbumTitleByIDs(gomock.Any(), gomock.Any()).Return(albumTitleMap, nil)
	mockArtistClient.EXPECT().GetArtistsByTrackIDs(gomock.Any(), gomock.Any()).Return(artistsMap, nil)

	tracks, err := trackUC.GetFavoriteTracks(ctx, filters, username)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracks))
	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
}

// Test for the privacy restriction case
func TestGetFavoriteTracksPrivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockAlbumClient := mocks.NewMockAlbumServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	trackUC := trackUsecase.NewUsecase(mockTrackClient, mockArtistClient, mockAlbumClient, mockPlaylistClient, mockUserClient)

	ctx := context.Background()
	username := "testuser"
	pagination := &usecase.Pagination{
		Limit:  10,
		Offset: 0,
	}
	filters := &usecase.TrackFilters{
		Pagination: pagination,
	}

	userID := &user.UserID{
		Id: 1,
	}

	userPrivacy := &user.PrivacySettings{
		IsPublicFavoriteTracks: false,
	}

	mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), gomock.Any()).Return(userID, nil)
	mockUserClient.EXPECT().GetUserPrivacyByID(gomock.Any(), gomock.Any()).Return(userPrivacy, nil)

	tracks, err := trackUC.GetFavoriteTracks(ctx, filters, username)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(tracks))
}
