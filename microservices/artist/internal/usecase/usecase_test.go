package usecase

import (
	"context"
	"errors"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/mocks"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	mockArtistID    = int64(1)
	mockUserID      = int64(123)
	mockLabelID     = int64(456)
	mockArtistTitle = "Test Artist"
	mockDescription = "Test Description"
	mockThumbnail   = "test_thumbnail.jpg"
	mockAvatarURL   = "http://example.com/avatar.jpg"
	mockNewTitle    = "New Artist Title"
	mockTrackID     = int64(789)
	mockAlbumID     = int64(101)
	mockSearchQuery = "test query"
)

func setupTest(t *testing.T) (*mock_domain.MockRepository, *mock_domain.MockS3Repository, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockRepository(ctrl)
	mockS3Repo := mock_domain.NewMockS3Repository(ctrl)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, mockS3Repo, ctx
}

func TestGetArtistByID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedArtist := &repoModel.Artist{
		ID:          mockArtistID,
		Title:       mockArtistTitle,
		Description: mockDescription,
		Thumbnail:   mockThumbnail,
		IsFavorite:  true,
	}

	expectedStats := &repoModel.ArtistStats{
		ListenersCount: 100,
		FavoritesCount: 50,
	}

	mockRepo.EXPECT().GetArtistByID(ctx, mockArtistID, mockUserID).Return(expectedArtist, nil)
	mockRepo.EXPECT().GetArtistStats(ctx, mockArtistID).Return(expectedStats, nil)

	result, err := usecase.GetArtistByID(ctx, mockArtistID, mockUserID)

	require.NoError(t, err)
	assert.Equal(t, expectedArtist.ID, result.ID)
	assert.Equal(t, expectedArtist.Title, result.Title)
	assert.Equal(t, expectedStats.ListenersCount, result.ListenersCount)
	assert.Equal(t, expectedStats.FavoritesCount, result.FavoritesCount)
}

func TestGetArtistByIDRepositoryError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistByID(ctx, mockArtistID, mockUserID).Return(nil, expectedErr)

	result, err := usecase.GetArtistByID(ctx, mockArtistID, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistByIDStatsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedArtist := &repoModel.Artist{
		ID:          mockArtistID,
		Title:       mockArtistTitle,
		Description: mockDescription,
		Thumbnail:   mockThumbnail,
		IsFavorite:  true,
	}

	expectedErr := errors.New("stats error")
	mockRepo.EXPECT().GetArtistByID(ctx, mockArtistID, mockUserID).Return(expectedArtist, nil)
	mockRepo.EXPECT().GetArtistStats(ctx, mockArtistID).Return(nil, expectedErr)

	result, err := usecase.GetArtistByID(ctx, mockArtistID, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetAllArtists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedArtists := []*repoModel.Artist{
		{
			ID:          mockArtistID,
			Title:       mockArtistTitle,
			Description: mockDescription,
			Thumbnail:   mockThumbnail,
			IsFavorite:  true,
		},
	}

	mockRepo.EXPECT().GetAllArtists(ctx, gomock.Any(), mockUserID).Return(expectedArtists, nil)

	result, err := usecase.GetAllArtists(ctx, filters, mockUserID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, expectedArtists[0].Title, result.Artists[0].Title)
}

func TestGetAllArtistsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAllArtists(ctx, gomock.Any(), mockUserID).Return(nil, expectedErr)

	result, err := usecase.GetAllArtists(ctx, filters, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistTitleByID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	mockRepo.EXPECT().GetArtistTitleByID(ctx, mockArtistID).Return(mockArtistTitle, nil)

	result, err := usecase.GetArtistTitleByID(ctx, mockArtistID)

	require.NoError(t, err)
	assert.Equal(t, mockArtistTitle, result)
}

func TestGetArtistTitleByIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistTitleByID(ctx, mockArtistID).Return("", expectedErr)

	result, err := usecase.GetArtistTitleByID(ctx, mockArtistID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, result)
}

func TestGetArtistsByTrackID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedArtists := []*repoModel.ArtistWithRole{
		{
			ID:    mockArtistID,
			Title: mockArtistTitle,
			Role:  "vocalist",
		},
	}

	mockRepo.EXPECT().GetArtistsByTrackID(ctx, mockTrackID).Return(expectedArtists, nil)

	result, err := usecase.GetArtistsByTrackID(ctx, mockTrackID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, expectedArtists[0].Title, result.Artists[0].Title)
	assert.Equal(t, expectedArtists[0].Role, result.Artists[0].Role)
}

func TestGetArtistsByTrackIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsByTrackID(ctx, mockTrackID).Return(nil, expectedErr)

	result, err := usecase.GetArtistsByTrackID(ctx, mockTrackID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistsByTrackIDs(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	trackIDs := []int64{mockTrackID, mockTrackID + 1}
	expectedArtists := map[int64][]*repoModel.ArtistWithRole{
		mockTrackID: {
			{
				ID:    mockArtistID,
				Title: mockArtistTitle,
				Role:  "vocalist",
			},
		},
	}

	mockRepo.EXPECT().GetArtistsByTrackIDs(ctx, trackIDs).Return(expectedArtists, nil)

	result, err := usecase.GetArtistsByTrackIDs(ctx, trackIDs)

	require.NoError(t, err)
	assert.Contains(t, result.Artists, mockTrackID)
	assert.Len(t, result.Artists[mockTrackID].Artists, 1)
}

func TestGetArtistsByTrackIDsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	trackIDs := []int64{mockTrackID}
	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsByTrackIDs(ctx, trackIDs).Return(nil, expectedErr)

	result, err := usecase.GetArtistsByTrackIDs(ctx, trackIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistsByAlbumID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedArtists := []*repoModel.ArtistWithTitle{
		{
			ID:    mockArtistID,
			Title: mockArtistTitle,
		},
	}

	mockRepo.EXPECT().GetArtistsByAlbumID(ctx, mockAlbumID).Return(expectedArtists, nil)

	result, err := usecase.GetArtistsByAlbumID(ctx, mockAlbumID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, expectedArtists[0].Title, result.Artists[0].Title)
}

func TestGetArtistsByAlbumIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsByAlbumID(ctx, mockAlbumID).Return(nil, expectedErr)

	result, err := usecase.GetArtistsByAlbumID(ctx, mockAlbumID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistsByAlbumIDs(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	albumIDs := []int64{mockAlbumID, mockAlbumID + 1}
	expectedArtists := map[int64][]*repoModel.ArtistWithTitle{
		mockAlbumID: {
			{
				ID:    mockArtistID,
				Title: mockArtistTitle,
			},
		},
	}

	mockRepo.EXPECT().GetArtistsByAlbumIDs(ctx, albumIDs).Return(expectedArtists, nil)

	result, err := usecase.GetArtistsByAlbumIDs(ctx, albumIDs)

	require.NoError(t, err)
	assert.Contains(t, result.Artists, mockAlbumID)
	assert.Len(t, result.Artists[mockAlbumID].Artists, 1)
}

func TestGetArtistsByAlbumIDsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	albumIDs := []int64{mockAlbumID}
	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsByAlbumIDs(ctx, albumIDs).Return(nil, expectedErr)

	result, err := usecase.GetArtistsByAlbumIDs(ctx, albumIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetAlbumIDsByArtistID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedAlbumIDs := []int64{mockAlbumID, mockAlbumID + 1}
	mockRepo.EXPECT().GetAlbumIDsByArtistID(ctx, mockArtistID).Return(expectedAlbumIDs, nil)

	result, err := usecase.GetAlbumIDsByArtistID(ctx, mockArtistID)

	require.NoError(t, err)
	assert.Equal(t, expectedAlbumIDs, result)
}

func TestGetAlbumIDsByArtistIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAlbumIDsByArtistID(ctx, mockArtistID).Return(nil, expectedErr)

	result, err := usecase.GetAlbumIDsByArtistID(ctx, mockArtistID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetTrackIDsByArtistID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedTrackIDs := []int64{mockTrackID, mockTrackID + 1}
	mockRepo.EXPECT().GetTrackIDsByArtistID(ctx, mockArtistID).Return(expectedTrackIDs, nil)

	result, err := usecase.GetTrackIDsByArtistID(ctx, mockArtistID)

	require.NoError(t, err)
	assert.Equal(t, expectedTrackIDs, result)
}

func TestGetTrackIDsByArtistIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetTrackIDsByArtistID(ctx, mockArtistID).Return(nil, expectedErr)

	result, err := usecase.GetTrackIDsByArtistID(ctx, mockArtistID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestCreateStreamsByArtistIDs(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	streamData := &usecaseModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{mockArtistID},
		UserID:    mockUserID,
	}

	mockRepo.EXPECT().CreateStreamsByArtistIDs(ctx, gomock.Any()).Return(nil)

	err := usecase.CreateStreamsByArtistIDs(ctx, streamData)

	require.NoError(t, err)
}

func TestCreateStreamsByArtistIDsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	streamData := &usecaseModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{mockArtistID},
		UserID:    mockUserID,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CreateStreamsByArtistIDs(ctx, gomock.Any()).Return(expectedErr)

	err := usecase.CreateStreamsByArtistIDs(ctx, streamData)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetArtistsListenedByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedCount := int64(42)
	mockRepo.EXPECT().GetArtistsListenedByUserID(ctx, mockUserID).Return(expectedCount, nil)

	result, err := usecase.GetArtistsListenedByUserID(ctx, mockUserID)

	require.NoError(t, err)
	assert.Equal(t, expectedCount, result)
}

func TestGetArtistsListenedByUserIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsListenedByUserID(ctx, mockUserID).Return(int64(0), expectedErr)

	result, err := usecase.GetArtistsListenedByUserID(ctx, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(0), result)
}

func TestLikeArtist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecaseModel.LikeRequest{
		ArtistID: mockArtistID,
		UserID:   mockUserID,
		IsLike:   true,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, mockArtistID).Return(true, nil)
	mockRepo.EXPECT().LikeArtist(ctx, gomock.Any()).Return(nil)

	err := usecase.LikeArtist(ctx, likeRequest)

	require.NoError(t, err)
}

func TestLikeArtistUnlike(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecaseModel.LikeRequest{
		ArtistID: mockArtistID,
		UserID:   mockUserID,
		IsLike:   false,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, mockArtistID).Return(true, nil)
	mockRepo.EXPECT().UnlikeArtist(ctx, gomock.Any()).Return(nil)

	err := usecase.LikeArtist(ctx, likeRequest)

	require.NoError(t, err)
}

func TestLikeArtistNotFound(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecaseModel.LikeRequest{
		ArtistID: mockArtistID,
		UserID:   mockUserID,
		IsLike:   true,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, mockArtistID).Return(false, nil)

	err := usecase.LikeArtist(ctx, likeRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "artist not found")
}

func TestLikeArtistCheckExistsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecaseModel.LikeRequest{
		ArtistID: mockArtistID,
		UserID:   mockUserID,
		IsLike:   true,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CheckArtistExists(ctx, mockArtistID).Return(false, expectedErr)

	err := usecase.LikeArtist(ctx, likeRequest)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLikeArtistLikeError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	likeRequest := &usecaseModel.LikeRequest{
		ArtistID: mockArtistID,
		UserID:   mockUserID,
		IsLike:   true,
	}

	expectedErr := errors.New("like error")
	mockRepo.EXPECT().CheckArtistExists(ctx, mockArtistID).Return(true, nil)
	mockRepo.EXPECT().LikeArtist(ctx, gomock.Any()).Return(expectedErr)

	err := usecase.LikeArtist(ctx, likeRequest)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetFavoriteArtists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedArtists := []*repoModel.Artist{
		{
			ID:          mockArtistID,
			Title:       mockArtistTitle,
			Description: mockDescription,
			Thumbnail:   mockThumbnail,
			IsFavorite:  true,
		},
	}

	mockRepo.EXPECT().GetFavoriteArtists(ctx, gomock.Any(), mockUserID).Return(expectedArtists, nil)

	result, err := usecase.GetFavoriteArtists(ctx, filters, mockUserID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
	assert.True(t, result.Artists[0].IsFavorite)
}

func TestGetFavoriteArtistsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetFavoriteArtists(ctx, gomock.Any(), mockUserID).Return(nil, expectedErr)

	result, err := usecase.GetFavoriteArtists(ctx, filters, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestSearchArtists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedArtists := []*repoModel.Artist{
		{
			ID:          mockArtistID,
			Title:       mockArtistTitle,
			Description: mockDescription,
			Thumbnail:   mockThumbnail,
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().SearchArtists(ctx, mockSearchQuery, mockUserID).Return(expectedArtists, nil)

	result, err := usecase.SearchArtists(ctx, mockSearchQuery, mockUserID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, expectedArtists[0].Title, result.Artists[0].Title)
}

func TestSearchArtistsError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	expectedErr := errors.New("search error")
	mockRepo.EXPECT().SearchArtists(ctx, mockSearchQuery, mockUserID).Return(nil, expectedErr)

	result, err := usecase.SearchArtists(ctx, mockSearchQuery, mockUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestCreateArtist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistLoad := &usecaseModel.ArtistLoad{
		Title:   mockArtistTitle,
		Image:   []byte("fake image data"),
		LabelID: mockLabelID,
	}

	expectedCreatedArtist := &repoModel.Artist{
		ID:          mockArtistID,
		Title:       mockArtistTitle,
		Description: "",
		Thumbnail:   mockAvatarURL,
		IsFavorite:  false,
	}

	mockS3Repo.EXPECT().UploadArtistAvatar(ctx, mockArtistTitle, artistLoad.Image).Return(mockAvatarURL, nil)
	mockRepo.EXPECT().CreateArtist(ctx, gomock.Any()).Return(expectedCreatedArtist, nil)

	result, err := usecase.CreateArtist(ctx, artistLoad)

	require.NoError(t, err)
	assert.Equal(t, expectedCreatedArtist.ID, result.ID)
	assert.Equal(t, expectedCreatedArtist.Title, result.Title)
	assert.Equal(t, expectedCreatedArtist.Thumbnail, result.Thumbnail)
}

func TestCreateArtistS3Error(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistLoad := &usecaseModel.ArtistLoad{
		Title:   mockArtistTitle,
		Image:   []byte("fake image data"),
		LabelID: mockLabelID,
	}

	expectedErr := errors.New("s3 upload error")
	mockS3Repo.EXPECT().UploadArtistAvatar(ctx, mockArtistTitle, artistLoad.Image).Return("", expectedErr)

	result, err := usecase.CreateArtist(ctx, artistLoad)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestCreateArtistRepositoryError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistLoad := &usecaseModel.ArtistLoad{
		Title:   mockArtistTitle,
		Image:   []byte("test image data"),
		LabelID: mockLabelID,
	}

	mockS3Repo.EXPECT().UploadArtistAvatar(ctx, mockArtistTitle, artistLoad.Image).Return(mockAvatarURL, nil)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CreateArtist(ctx, gomock.Any()).Return(nil, expectedErr)

	result, err := usecase.CreateArtist(ctx, artistLoad)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestEditArtist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistEdit := &usecaseModel.ArtistEdit{
		ArtistID: mockArtistID,
		NewTitle: mockNewTitle,
		Image:    []byte("new image data"),
		LabelID:  mockLabelID,
	}

	expectedArtist := &repoModel.Artist{
		ID:          mockArtistID,
		Title:       mockNewTitle,
		Description: mockDescription,
		Thumbnail:   mockAvatarURL,
		IsFavorite:  false,
	}

	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(mockLabelID, nil)
	mockRepo.EXPECT().GetArtistTitleByID(ctx, mockArtistID).Return(mockArtistTitle, nil)
	mockRepo.EXPECT().CheckArtistNameExist(ctx, mockArtistID).Return(true, nil)
	mockRepo.EXPECT().ChangeArtistTitle(ctx, mockNewTitle, mockArtistID).Return(nil)
	mockS3Repo.EXPECT().UploadArtistAvatar(ctx, mockNewTitle, artistEdit.Image).Return(mockAvatarURL, nil)
	mockRepo.EXPECT().UploadAvatar(ctx, mockArtistID, mockAvatarURL).Return(nil)
	mockRepo.EXPECT().GetArtistByIDWithoutUser(ctx, mockArtistID).Return(expectedArtist, nil)

	result, err := usecase.EditArtist(ctx, artistEdit)

	require.NoError(t, err)
	assert.Equal(t, expectedArtist.ID, result.ID)
	assert.Equal(t, expectedArtist.Title, result.Title)
	assert.Equal(t, mockLabelID, result.LabelID)
}

func TestEditArtistForbidden(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistEdit := &usecaseModel.ArtistEdit{
		ArtistID: mockArtistID,
		NewTitle: mockNewTitle,
		LabelID:  mockLabelID,
	}

	differentLabelID := mockLabelID + 1
	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(differentLabelID, nil)

	result, err := usecase.EditArtist(ctx, artistEdit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "you are not allowed to edit this artist")
	assert.Nil(t, result)
}

func TestEditArtistWithoutTitleChange(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistEdit := &usecaseModel.ArtistEdit{
		ArtistID: mockArtistID,
		NewTitle: "", // No title change
		Image:    []byte("new image data"),
		LabelID:  mockLabelID,
	}

	expectedArtist := &repoModel.Artist{
		ID:          mockArtistID,
		Title:       mockArtistTitle,
		Description: mockDescription,
		Thumbnail:   mockAvatarURL,
		IsFavorite:  false,
	}

	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(mockLabelID, nil)
	mockRepo.EXPECT().GetArtistTitleByID(ctx, mockArtistID).Return(mockArtistTitle, nil)
	mockS3Repo.EXPECT().UploadArtistAvatar(ctx, mockArtistTitle, artistEdit.Image).Return(mockAvatarURL, nil)
	mockRepo.EXPECT().UploadAvatar(ctx, mockArtistID, mockAvatarURL).Return(nil)
	mockRepo.EXPECT().GetArtistByIDWithoutUser(ctx, mockArtistID).Return(expectedArtist, nil)

	result, err := usecase.EditArtist(ctx, artistEdit)

	require.NoError(t, err)
	assert.Equal(t, expectedArtist.ID, result.ID)
	assert.Equal(t, expectedArtist.Title, result.Title)
}

func TestEditArtistNotFound(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistEdit := &usecaseModel.ArtistEdit{
		ArtistID: mockArtistID,
		NewTitle: mockNewTitle,
		LabelID:  mockLabelID,
	}

	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(mockLabelID, nil)
	mockRepo.EXPECT().GetArtistTitleByID(ctx, mockArtistID).Return(mockArtistTitle, nil)
	mockRepo.EXPECT().CheckArtistNameExist(ctx, mockArtistID).Return(false, nil)

	result, err := usecase.EditArtist(ctx, artistEdit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "artist not found")
	assert.Nil(t, result)
}

func TestGetArtistsLabelID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedArtists := []*repoModel.Artist{
		{
			ID:          mockArtistID,
			Title:       mockArtistTitle,
			Description: mockDescription,
			Thumbnail:   mockThumbnail,
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().GetArtistsLabelID(ctx, gomock.Any(), mockLabelID).Return(expectedArtists, nil)

	result, err := usecase.GetArtistsLabelID(ctx, filters, mockLabelID)

	require.NoError(t, err)
	assert.Len(t, result.Artists, 1)
	assert.Equal(t, expectedArtists[0].ID, result.Artists[0].ID)
}

func TestGetArtistsLabelIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Offset: 0,
			Limit:  10,
		},
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistsLabelID(ctx, gomock.Any(), mockLabelID).Return(nil, expectedErr)

	result, err := usecase.GetArtistsLabelID(ctx, filters, mockLabelID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestDeleteArtist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistDelete := &usecaseModel.ArtistDelete{
		ArtistID: mockArtistID,
		LabelID:  mockLabelID,
	}

	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(mockLabelID, nil)
	mockRepo.EXPECT().DeleteArtist(ctx, mockArtistID).Return(nil)

	err := usecase.DeleteArtist(ctx, artistDelete)

	require.NoError(t, err)
}

func TestDeleteArtistForbidden(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistDelete := &usecaseModel.ArtistDelete{
		ArtistID: mockArtistID,
		LabelID:  mockLabelID,
	}

	differentLabelID := mockLabelID + 1
	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(differentLabelID, nil)

	err := usecase.DeleteArtist(ctx, artistDelete)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "you are not allowed to edit this artist")
}

func TestDeleteArtistGetLabelIDError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistDelete := &usecaseModel.ArtistDelete{
		ArtistID: mockArtistID,
		LabelID:  mockLabelID,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(int64(0), expectedErr)

	err := usecase.DeleteArtist(ctx, artistDelete)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestDeleteArtistRepositoryError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistDelete := &usecaseModel.ArtistDelete{
		ArtistID: mockArtistID,
		LabelID:  mockLabelID,
	}

	expectedErr := errors.New("delete error")
	mockRepo.EXPECT().GetArtistLabelID(ctx, mockArtistID).Return(mockLabelID, nil)
	mockRepo.EXPECT().DeleteArtist(ctx, mockArtistID).Return(expectedErr)

	err := usecase.DeleteArtist(ctx, artistDelete)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestConnectArtists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistIDs := []int64{mockArtistID, mockArtistID + 1}
	trackIDs := []int64{mockTrackID, mockTrackID + 1}

	mockRepo.EXPECT().AddArtistsToAlbum(ctx, artistIDs, mockAlbumID).Return(nil)
	mockRepo.EXPECT().AddArtistsToTracks(ctx, artistIDs, trackIDs).Return(nil)

	err := usecase.ConnectArtists(ctx, artistIDs, mockAlbumID, trackIDs)

	require.NoError(t, err)
}

func TestConnectArtistsAlbumError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistIDs := []int64{mockArtistID}
	trackIDs := []int64{mockTrackID}

	expectedErr := errors.New("album connection error")
	mockRepo.EXPECT().AddArtistsToAlbum(ctx, artistIDs, mockAlbumID).Return(expectedErr)

	err := usecase.ConnectArtists(ctx, artistIDs, mockAlbumID, trackIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestConnectArtistsTracksError(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)
	usecase := NewArtistUsecase(mockRepo, mockS3Repo)

	artistIDs := []int64{mockArtistID}
	trackIDs := []int64{mockTrackID}

	expectedErr := errors.New("tracks connection error")
	mockRepo.EXPECT().AddArtistsToAlbum(ctx, artistIDs, mockAlbumID).Return(nil)
	mockRepo.EXPECT().AddArtistsToTracks(ctx, artistIDs, trackIDs).Return(expectedErr)

	err := usecase.ConnectArtists(ctx, artistIDs, mockAlbumID, trackIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
