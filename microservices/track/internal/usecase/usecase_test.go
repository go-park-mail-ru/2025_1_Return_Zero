package usecase

import (
	"context"
	"errors"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/mocks"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_domain.MockRepository, *mock_domain.MockS3Repository, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockRepository(ctrl)
	mockS3Repo := mock_domain.NewMockS3Repository(ctrl)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, mockS3Repo, ctx
}

func TestGetAllTracks(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	repoTracks := []*repository.Track{
		{
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockRepo.EXPECT().GetAllTracks(ctx, repoFilters, userID).Return(repoTracks, nil)

	tracks, err := u.GetAllTracks(ctx, filters, userID)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, "thumbnail1.jpg", tracks[0].Thumbnail)
	assert.Equal(t, int64(200), tracks[0].Duration)
	assert.Equal(t, int64(1), tracks[0].AlbumID)
	assert.True(t, tracks[0].IsFavorite)

	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
	assert.Equal(t, "thumbnail2.jpg", tracks[1].Thumbnail)
	assert.Equal(t, int64(300), tracks[1].Duration)
	assert.Equal(t, int64(1), tracks[1].AlbumID)
	assert.False(t, tracks[1].IsFavorite)
}

func TestGetAllTracksError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	expectedErr := errors.New("database error")
	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockRepo.EXPECT().GetAllTracks(ctx, repoFilters, userID).Return(nil, expectedErr)

	tracks, err := u.GetAllTracks(ctx, filters, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetTrackByID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	trackID := int64(1)
	userID := int64(1)

	track := repository.Track{
		ID:         1,
		Title:      "Track 1",
		Thumbnail:  "thumbnail1.jpg",
		Duration:   200,
		AlbumID:    1,
		IsFavorite: true,
	}

	repoTrack := &repository.TrackWithFileKey{
		Track:   track,
		FileKey: "file_key.mp3",
	}

	presignedURL := "https://example.com/tracks/file_key.mp3"

	mockRepo.EXPECT().GetTrackByID(ctx, trackID, userID).Return(repoTrack, nil)
	mockS3Repo.EXPECT().GetPresignedURL(repoTrack.FileKey).Return(presignedURL, nil)

	trackResult, err := u.GetTrackByID(ctx, trackID, userID)
	require.NoError(t, err)
	require.NotNil(t, trackResult)

	assert.Equal(t, int64(1), trackResult.ID)
	assert.Equal(t, "Track 1", trackResult.Title)
	assert.Equal(t, "thumbnail1.jpg", trackResult.Thumbnail)
	assert.Equal(t, int64(200), trackResult.Duration)
	assert.Equal(t, int64(1), trackResult.AlbumID)
	assert.Equal(t, presignedURL, trackResult.FileUrl)
	assert.True(t, trackResult.IsFavorite)
}

func TestGetTrackByIDRepositoryError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	trackID := int64(1)
	userID := int64(1)
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetTrackByID(ctx, trackID, userID).Return(nil, expectedErr)

	trackResult, err := u.GetTrackByID(ctx, trackID, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, trackResult)
}

func TestGetTrackByIDS3Error(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	trackID := int64(1)
	userID := int64(1)

	track := repository.Track{
		ID:         1,
		Title:      "Track 1",
		Thumbnail:  "thumbnail1.jpg",
		Duration:   200,
		AlbumID:    1,
		IsFavorite: true,
	}

	repoTrack := &repository.TrackWithFileKey{
		Track:   track,
		FileKey: "file_key.mp3",
	}

	expectedErr := errors.New("s3 error")

	mockRepo.EXPECT().GetTrackByID(ctx, trackID, userID).Return(repoTrack, nil)
	mockS3Repo.EXPECT().GetPresignedURL(repoTrack.FileKey).Return("", expectedErr)

	trackResult, err := u.GetTrackByID(ctx, trackID, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, trackResult)
}

func TestCreateStream(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	streamData := &usecase.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	repoStreamData := &repository.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	expectedStreamID := int64(42)

	mockRepo.EXPECT().CreateStream(ctx, repoStreamData).Return(expectedStreamID, nil)

	streamID, err := u.CreateStream(ctx, streamData)
	require.NoError(t, err)
	assert.Equal(t, expectedStreamID, streamID)
}

func TestCreateStreamError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	streamData := &usecase.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	repoStreamData := &repository.TrackStreamCreateData{
		TrackID: 1,
		UserID:  1,
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().CreateStream(ctx, repoStreamData).Return(int64(0), expectedErr)

	streamID, err := u.CreateStream(ctx, streamData)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(0), streamID)
}

func TestUpdateStreamDuration(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	updateData := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		UserID:   1,
		Duration: 200,
	}

	repoStream := &repository.TrackStream{
		ID:       1,
		UserID:   1,
		TrackID:  1,
		Duration: 0,
	}

	repoUpdateData := &repository.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 200,
	}

	mockRepo.EXPECT().GetStreamByID(ctx, updateData.StreamID).Return(repoStream, nil)
	mockRepo.EXPECT().UpdateStreamDuration(ctx, repoUpdateData).Return(nil)

	err := u.UpdateStreamDuration(ctx, updateData)
	require.NoError(t, err)
}

func TestUpdateStreamDurationNotFound(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	updateData := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		UserID:   1,
		Duration: 200,
	}

	expectedErr := trackErrors.ErrStreamNotFound

	mockRepo.EXPECT().GetStreamByID(ctx, updateData.StreamID).Return(nil, expectedErr)

	err := u.UpdateStreamDuration(ctx, updateData)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUpdateStreamDurationPermissionDenied(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	updateData := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		UserID:   1,
		Duration: 200,
	}

	repoStream := &repository.TrackStream{
		ID:       1,
		UserID:   2, // Different user ID
		TrackID:  1,
		Duration: 0,
	}

	mockRepo.EXPECT().GetStreamByID(ctx, updateData.StreamID).Return(repoStream, nil)

	err := u.UpdateStreamDuration(ctx, updateData)
	assert.Error(t, err)
	assert.Equal(t, trackErrors.ErrStreamPermissionDenied, err)
}

func TestUpdateStreamDurationUpdateError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	updateData := &usecase.TrackStreamUpdateData{
		StreamID: 1,
		UserID:   1,
		Duration: 200,
	}

	repoStream := &repository.TrackStream{
		ID:       1,
		UserID:   1,
		TrackID:  1,
		Duration: 0,
	}

	repoUpdateData := &repository.TrackStreamUpdateData{
		StreamID: 1,
		Duration: 200,
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetStreamByID(ctx, updateData.StreamID).Return(repoStream, nil)
	mockRepo.EXPECT().UpdateStreamDuration(ctx, repoUpdateData).Return(expectedErr)

	err := u.UpdateStreamDuration(ctx, updateData)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetLastListenedTracks(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoStreams := []*repository.TrackStream{
		{
			ID:       1,
			UserID:   1,
			TrackID:  1,
			Duration: 200,
		},
		{
			ID:       2,
			UserID:   1,
			TrackID:  2,
			Duration: 300,
		},
	}

	repoTracks := map[int64]*repository.Track{
		1: {
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		2: {
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	mockRepo.EXPECT().GetStreamsByUserID(ctx, userID, repoFilters).Return(repoStreams, nil)

	trackIDs := []int64{1, 2}
	mockRepo.EXPECT().GetTracksByIDs(ctx, trackIDs, userID).Return(repoTracks, nil)

	tracks, err := u.GetLastListenedTracks(ctx, userID, filters)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
}

func TestGetLastListenedTracksNoStreams(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockRepo.EXPECT().GetStreamsByUserID(ctx, userID, repoFilters).Return([]*repository.TrackStream{}, nil)

	tracks, err := u.GetLastListenedTracks(ctx, userID, filters)
	require.NoError(t, err)
	assert.Empty(t, tracks)
}

func TestGetLastListenedTracksStreamError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetStreamsByUserID(ctx, userID, repoFilters).Return(nil, expectedErr)

	tracks, err := u.GetLastListenedTracks(ctx, userID, filters)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetLastListenedTracksGetTracksError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoStreams := []*repository.TrackStream{
		{
			ID:       1,
			UserID:   1,
			TrackID:  1,
			Duration: 200,
		},
		{
			ID:       2,
			UserID:   1,
			TrackID:  2,
			Duration: 300,
		},
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetStreamsByUserID(ctx, userID, repoFilters).Return(repoStreams, nil)
	mockRepo.EXPECT().GetTracksByIDs(ctx, []int64{1, 2}, userID).Return(nil, expectedErr)

	tracks, err := u.GetLastListenedTracks(ctx, userID, filters)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetTracksByIDs(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	trackIDs := []int64{1, 2}

	repoTracks := map[int64]*repository.Track{
		1: {
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		2: {
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	mockRepo.EXPECT().GetTracksByIDs(ctx, trackIDs, userID).Return(repoTracks, nil)

	tracks, err := u.GetTracksByIDs(ctx, trackIDs, userID)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
}

func TestGetTracksByIDsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	trackIDs := []int64{1, 2}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetTracksByIDs(ctx, trackIDs, userID).Return(nil, expectedErr)

	tracks, err := u.GetTracksByIDs(ctx, trackIDs, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetTracksByIDsFiltered(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	trackIDs := []int64{1, 2}
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoTracks := []*repository.Track{
		{
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	mockRepo.EXPECT().GetTracksByIDsFiltered(ctx, trackIDs, repoFilters, userID).Return(repoTracks, nil)

	tracks, err := u.GetTracksByIDsFiltered(ctx, trackIDs, filters, userID)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
}

func TestGetTracksByIDsFilteredError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	trackIDs := []int64{1, 2}
	filters := &usecase.TrackFilters{
		Pagination: &usecase.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repository.TrackFilters{
		Pagination: &repository.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetTracksByIDsFiltered(ctx, trackIDs, repoFilters, userID).Return(nil, expectedErr)

	tracks, err := u.GetTracksByIDsFiltered(ctx, trackIDs, filters, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetAlbumIDByTrackID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	trackID := int64(1)
	expectedAlbumID := int64(42)

	mockRepo.EXPECT().GetAlbumIDByTrackID(ctx, trackID).Return(expectedAlbumID, nil)

	albumID, err := u.GetAlbumIDByTrackID(ctx, trackID)
	require.NoError(t, err)
	assert.Equal(t, expectedAlbumID, albumID)
}

func TestGetAlbumIDByTrackIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	trackID := int64(1)
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetAlbumIDByTrackID(ctx, trackID).Return(int64(0), expectedErr)

	albumID, err := u.GetAlbumIDByTrackID(ctx, trackID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(0), albumID)
}

func TestGetTracksByAlbumID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	albumID := int64(1)
	userID := int64(1)

	repoTracks := []*repository.Track{
		{
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	mockRepo.EXPECT().GetTracksByAlbumID(ctx, albumID, userID).Return(repoTracks, nil)

	tracks, err := u.GetTracksByAlbumID(ctx, albumID, userID)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
}

func TestGetTracksByAlbumIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	albumID := int64(1)
	userID := int64(1)
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetTracksByAlbumID(ctx, albumID, userID).Return(nil, expectedErr)

	tracks, err := u.GetTracksByAlbumID(ctx, albumID, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestGetMinutesListenedByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	expectedMinutes := int64(120)

	mockRepo.EXPECT().GetMinutesListenedByUserID(ctx, userID).Return(expectedMinutes, nil)

	minutes, err := u.GetMinutesListenedByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedMinutes, minutes)
}

func TestGetMinutesListenedByUserIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetMinutesListenedByUserID(ctx, userID).Return(int64(0), expectedErr)

	minutes, err := u.GetMinutesListenedByUserID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(0), minutes)
}

func TestGetTracksListenedByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	expectedCount := int64(42)

	mockRepo.EXPECT().GetTracksListenedByUserID(ctx, userID).Return(expectedCount, nil)

	count, err := u.GetTracksListenedByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestGetTracksListenedByUserIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	userID := int64(1)
	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetTracksListenedByUserID(ctx, userID).Return(int64(0), expectedErr)

	count, err := u.GetTracksListenedByUserID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(0), count)
}

func TestLikeTrackSuccess(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}

	repoLikeRequest := &repository.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(true, nil)
	mockRepo.EXPECT().LikeTrack(ctx, repoLikeRequest).Return(nil)

	err := u.LikeTrack(ctx, likeRequest)
	require.NoError(t, err)
}

func TestLikeTrackCheckExistsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(false, expectedErr)

	err := u.LikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLikeTrackNotFound(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(false, nil)

	err := u.LikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.IsType(t, trackErrors.ErrTrackNotFound, err)
}

func TestLikeTrackLikeError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  true,
	}

	repoLikeRequest := &repository.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(true, nil)
	mockRepo.EXPECT().LikeTrack(ctx, repoLikeRequest).Return(expectedErr)

	err := u.LikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUnlikeTrackSuccess(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  false,
	}

	repoLikeRequest := &repository.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(true, nil)
	mockRepo.EXPECT().UnlikeTrack(ctx, repoLikeRequest).Return(nil)

	err := u.LikeTrack(ctx, likeRequest)
	require.NoError(t, err)
}

func TestUnlikeTrackError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecase.LikeRequest{
		TrackID: 1,
		UserID:  1,
		IsLike:  false,
	}

	repoLikeRequest := &repository.LikeRequest{
		TrackID: 1,
		UserID:  1,
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().CheckTrackExists(ctx, likeRequest.TrackID).Return(true, nil)
	mockRepo.EXPECT().UnlikeTrack(ctx, repoLikeRequest).Return(expectedErr)

	err := u.LikeTrack(ctx, likeRequest)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetFavoriteTracks(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	favoriteRequest := &usecase.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &usecase.TrackFilters{
			Pagination: &usecase.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	repoFavoriteRequest := &repository.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &repository.TrackFilters{
			Pagination: &repository.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	repoTracks := []*repository.Track{
		{
			ID:         1,
			Title:      "Track 1",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Track 2",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: true,
		},
	}

	mockRepo.EXPECT().GetFavoriteTracks(ctx, repoFavoriteRequest).Return(repoTracks, nil)

	tracks, err := u.GetFavoriteTracks(ctx, favoriteRequest)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Track 1", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track 2", tracks[1].Title)
}

func TestGetFavoriteTracksError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	favoriteRequest := &usecase.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &usecase.TrackFilters{
			Pagination: &usecase.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	repoFavoriteRequest := &repository.FavoriteRequest{
		RequestUserID: 1,
		ProfileUserID: 2,
		Filters: &repository.TrackFilters{
			Pagination: &repository.Pagination{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().GetFavoriteTracks(ctx, repoFavoriteRequest).Return(nil, expectedErr)

	tracks, err := u.GetFavoriteTracks(ctx, favoriteRequest)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}

func TestSearchTracks(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	query := "test track"
	userID := int64(1)

	repoTracks := []*repository.Track{
		{
			ID:         1,
			Title:      "Test Track",
			Thumbnail:  "thumbnail1.jpg",
			Duration:   200,
			AlbumID:    1,
			IsFavorite: true,
		},
		{
			ID:         2,
			Title:      "Track Test",
			Thumbnail:  "thumbnail2.jpg",
			Duration:   300,
			AlbumID:    1,
			IsFavorite: false,
		},
	}

	mockRepo.EXPECT().SearchTracks(ctx, query, userID).Return(repoTracks, nil)

	tracks, err := u.SearchTracks(ctx, query, userID)
	require.NoError(t, err)
	require.Len(t, tracks, 2)

	assert.Equal(t, int64(1), tracks[0].ID)
	assert.Equal(t, "Test Track", tracks[0].Title)
	assert.Equal(t, int64(2), tracks[1].ID)
	assert.Equal(t, "Track Test", tracks[1].Title)
}

func TestSearchTracksError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	u := NewTrackUsecase(mockRepo, mockS3Repo)

	query := "test track"
	userID := int64(1)

	expectedErr := errors.New("database error")

	mockRepo.EXPECT().SearchTracks(ctx, query, userID).Return(nil, expectedErr)

	tracks, err := u.SearchTracks(ctx, query, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tracks)
}
