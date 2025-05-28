package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/mocks"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_domain.MockRepository, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockRepository(ctrl)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, ctx
}

func TestCreateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().CreateArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(&artistProto.Artist{
		Id:        1,
		Title:     "new_artist",
		Thumbnail: "thumbnail_url",
	}, nil)

	artist, err := usecase.CreateArtist(ctx, &usecaseModel.ArtistLoad{
		Title: "new_artist",
		Image: []byte("test image"),
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), artist.ID)
	assert.Equal(t, "new_artist", artist.Title)
}

func TestCreateArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().CreateArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	artist, err := usecase.CreateArtist(ctx, &usecaseModel.ArtistLoad{
		Title: "error_artist",
		Image: []byte("test image"),
	})

	require.Error(t, err)
	assert.Nil(t, artist)
}

func TestEditArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().EditArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(&artistProto.Artist{
		Id:        1,
		Title:     "edited_artist",
		Thumbnail: "edited_thumbnail_url",
	}, nil)

	artist, err := usecase.EditArtist(ctx, &usecaseModel.ArtistEdit{
		ArtistID: 1,
		NewTitle: "edited_artist",
		Image:    []byte("edited image"),
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), artist.ID)
	assert.Equal(t, "edited_artist", artist.Title)
}

func TestEditArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().EditArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	artist, err := usecase.EditArtist(ctx, &usecaseModel.ArtistEdit{
		ArtistID: 1,
		NewTitle: "error_artist",
		Image:    []byte("edited image"),
	})

	require.Error(t, err)
	assert.Nil(t, artist)
}

func TestCreateLabel(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockRepo.EXPECT().CheckIsLabelUnique(
		gomock.Any(),
		gomock.Eq("new_label"),
	).Return(false, nil)

	mockRepo.EXPECT().CreateLabel(
		gomock.Any(),
		gomock.Eq("new_label"),
	).Return(int64(1), nil)

	mockUserProto.EXPECT().ChecksUsersByUsernames(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil)

	mockUserProto.EXPECT().UpdateUsersLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil)

	label, err := usecase.CreateLabel(ctx, &usecaseModel.Label{
		Name:    "new_label",
		Members: []string{"1", "2"},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), label.Id)
	assert.Equal(t, "new_label", label.Name)
}

func TestCreateLabelError(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockRepo.EXPECT().CheckIsLabelUnique(
		gomock.Any(),
		gomock.Eq("error_label"),
	).Return(false, assert.AnError)

	label, err := usecase.CreateLabel(ctx, &usecaseModel.Label{
		Name:    "error_label",
		Members: []string{"1", "2"},
	})

	require.Error(t, err)
	assert.Nil(t, label)
}

func TestGetLabel(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockRepo.EXPECT().GetLabel(
		gomock.Any(),
		gomock.Eq(int64(1)),
	).Return(&repoModel.Label{
		ID:   1,
		Name: "test_label",
	}, nil)

	mockUserProto.EXPECT().GetUsersByLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(&userProto.Usernames{
		Usernames: []string{"user1", "user2"},
	}, nil)

	label, err := usecase.GetLabel(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), label.Id)
	assert.Equal(t, "test_label", label.Name)
	assert.Equal(t, []string{"user1", "user2"}, label.Members)
}

func TestGetLabelError(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockRepo.EXPECT().GetLabel(
		gomock.Any(),
		gomock.Eq(int64(999)),
	).Return(nil, assert.AnError)

	label, err := usecase.GetLabel(ctx, 999)

	require.Error(t, err)
	assert.Nil(t, label)
}

func TestGetArtists(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	filters := &usecaseModel.ArtistFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockArtistProto.EXPECT().GetArtistsLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(&artistProto.ArtistList{
		Artists: []*artistProto.Artist{
			{Id: 1, Title: "artist1"},
			{Id: 2, Title: "artist2"},
		},
	}, nil)

	artists, err := usecase.GetArtists(ctx, 1, filters)

	require.NoError(t, err)
	assert.Len(t, artists, 2)
	assert.Equal(t, int64(1), artists[0].ID)
	assert.Equal(t, "artist1", artists[0].Title)
}

func TestGetArtistsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	filters := &usecaseModel.ArtistFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockArtistProto.EXPECT().GetArtistsLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	artists, err := usecase.GetArtists(ctx, 999, filters)

	require.Error(t, err)
	assert.Nil(t, artists)
}

func TestGetAlbumsByLabelID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockAlbumRepo.EXPECT().GetAlbumsLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(&album.AlbumList{
		Albums: []*album.Album{
			{Id: 1, Title: "album1"},
			{Id: 2, Title: "album2"},
		},
	}, nil)

	mockArtistProto.EXPECT().GetArtistsByAlbumIDs(
		gomock.Any(),
		gomock.Any(),
	).Return(&artistProto.ArtistWithTitleMap{
		Artists: map[int64]*artistProto.ArtistWithTitleList{
			1: {
				Artists: []*artistProto.ArtistWithTitle{
					{Id: 1, Title: "Artist 1"},
				},
			},
			2: {
				Artists: []*artistProto.ArtistWithTitle{
					{Id: 2, Title: "Artist 2"},
				},
			},
		},
	}, nil)

	albums, err := usecase.GetAlbumsByLabelID(ctx, 1, &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	})

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "album1", albums[0].Title)
}

func TestGetAlbumsByLabelIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockAlbumRepo.EXPECT().GetAlbumsLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	albums, err := usecase.GetAlbumsByLabelID(ctx, 1, &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	})

	require.Error(t, err)
	assert.Nil(t, albums)
}

func TestDeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().DeleteArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil)

	err := usecase.DeleteArtist(ctx, &usecaseModel.ArtistDelete{
		ArtistID: 1,
		LabelID:  1,
	})

	require.NoError(t, err)
}

func TestDeleteArtistError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockArtistProto := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistProto, nil, nil)

	ctx := context.Background()

	mockArtistProto.EXPECT().DeleteArtist(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	err := usecase.DeleteArtist(ctx, &usecaseModel.ArtistDelete{
		ArtistID: 1,
		LabelID:  1,
	})

	require.Error(t, err)
}

func TestCreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(ctrl)
	mockTrackRepo := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistRepo := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistRepo, mockAlbumRepo, mockTrackRepo)

	ctx := context.Background()

	mockAlbumRepo.EXPECT().CreateAlbum(
		gomock.Any(),
		gomock.Any(),
	).Return(&album.AlbumIDAndURL{
		Id:  1,
		Url: "thumbnail_url",
	}, nil)

	mockTrackRepo.EXPECT().AddTracksToAlbum(
		gomock.Any(),
		gomock.Any(),
	).Return(&track.TrackIdsList{
		Ids: []*track.TrackID{},
	}, nil).AnyTimes()

	mockArtistRepo.EXPECT().ConnectArtists(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil).AnyTimes()

	albumID, thumbnailURL, err := usecase.CreateAlbum(ctx, &usecaseModel.CreateAlbumRequest{
		ArtistsIDs: []int64{1},
		Title:      "new_album",
		Type:       "album",
		Image:      []byte("test image"),
		LabelID:    1,
		Tracks:     []*usecaseModel.CreateTrackRequest{},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), albumID)
	assert.Equal(t, "thumbnail_url", thumbnailURL)
}

func TestCreateAlbumError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(ctrl)
	mockTrackRepo := mocks.NewMockTrackServiceClient(ctrl)
	mockArtistRepo := mocks.NewMockArtistServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, mockArtistRepo, mockAlbumRepo, mockTrackRepo)

	ctx := context.Background()

	mockAlbumRepo.EXPECT().CreateAlbum(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	albumID, thumbnailURL, err := usecase.CreateAlbum(ctx, &usecaseModel.CreateAlbumRequest{
		ArtistsIDs: []int64{1},
		Title:      "error_album",
		Type:       "album",
		Image:      []byte("test image"),
		LabelID:    1,
		Tracks:     []*usecaseModel.CreateTrackRequest{},
	})

	require.Error(t, err)
	assert.Equal(t, int64(-1), albumID)
	assert.Empty(t, thumbnailURL)
}

func TestUpdateLabel(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockUserProto.EXPECT().RemoveUsersFromLabel(
		gomock.Any(),
		gomock.Any(),
	).Return(&userProto.Nothing{}, nil)

	mockUserProto.EXPECT().UpdateUsersLabelID(
		gomock.Any(),
		gomock.Any(),
	).Return(&userProto.Nothing{}, nil)

	err := usecase.UpdateLabel(ctx, 1, []string{"user3"}, []string{"user1"})

	require.NoError(t, err)
}

func TestUpdateLabelError(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	mockUserProto := mocks.NewMockUserServiceClient(gomock.NewController(t))
	mockArtistProto := mocks.NewMockArtistServiceClient(gomock.NewController(t))
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(gomock.NewController(t))
	mockTrackRepo := mocks.NewMockTrackServiceClient(gomock.NewController(t))

	usecase := NewLabelUsecase(mockRepo, mockUserProto, mockArtistProto, mockAlbumRepo, mockTrackRepo)

	mockUserProto.EXPECT().RemoveUsersFromLabel(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	err := usecase.UpdateLabel(ctx, 1, []string{"user3"}, []string{"user1"})

	require.Error(t, err)
}

func TestDeleteAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(ctrl)
	mockTrackRepo := mocks.NewMockTrackServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, nil, mockAlbumRepo, mockTrackRepo)

	ctx := context.Background()

	mockAlbumRepo.EXPECT().DeleteAlbum(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil)

	mockTrackRepo.EXPECT().DeleteTracksByAlbumID(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, nil)

	err := usecase.DeleteAlbum(ctx, 1, 1)

	require.NoError(t, err)
}

func TestDeleteAlbumError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAlbumRepo := mocks.NewMockAlbumServiceClient(ctrl)
	mockTrackRepo := mocks.NewMockTrackServiceClient(ctrl)

	usecase := NewLabelUsecase(nil, nil, nil, mockAlbumRepo, mockTrackRepo)

	ctx := context.Background()

	mockAlbumRepo.EXPECT().DeleteAlbum(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, assert.AnError)

	err := usecase.DeleteAlbum(ctx, 1, 1)

	require.Error(t, err)
}
