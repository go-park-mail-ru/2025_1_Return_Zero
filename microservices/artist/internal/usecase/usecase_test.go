package usecase

import (
	"context"
	"errors"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/mocks"
	artistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
	"github.com/stretchr/testify/assert"
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

func TestGetArtistByID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)
	userID := int64(2)

	mockArtist := &repoModel.Artist{
		ID:          artistID,
		Title:       "Test Artist",
		Description: "Test Description",
		Thumbnail:   "test.jpg",
		IsFavorite:  true,
	}

	mockStats := &repoModel.ArtistStats{
		ListenersCount: 100,
		FavoritesCount: 50,
	}

	mockRepo.EXPECT().GetArtistByID(ctx, artistID, userID).Return(mockArtist, nil)
	mockRepo.EXPECT().GetArtistStats(ctx, artistID).Return(mockStats, nil)

	result, err := usecase.GetArtistByID(ctx, artistID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, artistID, result.ID)
	assert.Equal(t, mockArtist.Title, result.Title)
	assert.Equal(t, mockArtist.Description, result.Description)
	assert.Equal(t, mockArtist.Thumbnail, result.Thumbnail)
	assert.Equal(t, mockArtist.IsFavorite, result.IsFavorite)
	assert.Equal(t, mockStats.ListenersCount, result.ListenersCount)
	assert.Equal(t, mockStats.FavoritesCount, result.FavoritesCount)
}

func TestGetArtistByIDErrorFromRepo(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)
	userID := int64(2)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetArtistByID(ctx, artistID, userID).Return(nil, expectedErr)

	result, err := usecase.GetArtistByID(ctx, artistID, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetArtistByIDErrorFromStats(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)
	userID := int64(2)

	mockArtist := &repoModel.Artist{
		ID:          artistID,
		Title:       "Test Artist",
		Description: "Test Description",
		Thumbnail:   "test.jpg",
		IsFavorite:  true,
	}

	expectedErr := errors.New("stats error")
	mockRepo.EXPECT().GetArtistByID(ctx, artistID, userID).Return(mockArtist, nil)
	mockRepo.EXPECT().GetArtistStats(ctx, artistID).Return(nil, expectedErr)

	result, err := usecase.GetArtistByID(ctx, artistID, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetAllArtists(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	userID := int64(1)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockArtists := []*repoModel.Artist{
		{
			ID:          1,
			Title:       "Artist 1",
			Description: "Description 1",
			Thumbnail:   "thumb1.jpg",
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Artist 2",
			Description: "Description 2",
			Thumbnail:   "thumb2.jpg",
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().GetAllArtists(ctx, repoFilters, userID).Return(mockArtists, nil)

	result, err := usecase.GetAllArtists(ctx, filters, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, mockArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, mockArtists[0].Title, result.Artists[0].Title)
	assert.Equal(t, mockArtists[1].ID, result.Artists[1].ID)
	assert.Equal(t, mockArtists[1].Title, result.Artists[1].Title)
}

func TestGetArtistTitleByID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)
	expectedTitle := "Test Artist"

	mockRepo.EXPECT().GetArtistTitleByID(ctx, artistID).Return(expectedTitle, nil)

	title, err := usecase.GetArtistTitleByID(ctx, artistID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTitle, title)
}

func TestGetArtistsByTrackID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	trackID := int64(1)

	mockArtists := []*repoModel.ArtistWithRole{
		{
			ID:    1,
			Title: "Artist 1",
			Role:  "main",
		},
		{
			ID:    2,
			Title: "Artist 2",
			Role:  "featured",
		},
	}

	mockRepo.EXPECT().GetArtistsByTrackID(ctx, trackID).Return(mockArtists, nil)

	result, err := usecase.GetArtistsByTrackID(ctx, trackID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, mockArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, mockArtists[0].Title, result.Artists[0].Title)
	assert.Equal(t, mockArtists[0].Role, result.Artists[0].Role)
}

func TestGetArtistsByTrackIDs(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	trackIDs := []int64{1, 2}

	mockArtistsMap := map[int64][]*repoModel.ArtistWithRole{
		1: {
			{
				ID:    1,
				Title: "Artist 1",
				Role:  "main",
			},
		},
		2: {
			{
				ID:    2,
				Title: "Artist 2",
				Role:  "featured",
			},
		},
	}

	mockRepo.EXPECT().GetArtistsByTrackIDs(ctx, trackIDs).Return(mockArtistsMap, nil)

	result, err := usecase.GetArtistsByTrackIDs(ctx, trackIDs)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Len(t, result.Artists[1].Artists, 1)
	assert.Len(t, result.Artists[2].Artists, 1)
	assert.Equal(t, mockArtistsMap[1][0].ID, result.Artists[1].Artists[0].ID)
	assert.Equal(t, mockArtistsMap[2][0].ID, result.Artists[2].Artists[0].ID)
}

func TestGetArtistsByAlbumID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	albumID := int64(1)

	mockArtists := []*repoModel.ArtistWithTitle{
		{
			ID:    1,
			Title: "Artist 1",
		},
		{
			ID:    2,
			Title: "Artist 2",
		},
	}

	mockRepo.EXPECT().GetArtistsByAlbumID(ctx, albumID).Return(mockArtists, nil)

	result, err := usecase.GetArtistsByAlbumID(ctx, albumID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, mockArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, mockArtists[0].Title, result.Artists[0].Title)
}

func TestGetArtistsByAlbumIDs(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	albumIDs := []int64{1, 2}

	mockArtistsMap := map[int64][]*repoModel.ArtistWithTitle{
		1: {
			{
				ID:    1,
				Title: "Artist 1",
			},
		},
		2: {
			{
				ID:    2,
				Title: "Artist 2",
			},
		},
	}

	mockRepo.EXPECT().GetArtistsByAlbumIDs(ctx, albumIDs).Return(mockArtistsMap, nil)

	result, err := usecase.GetArtistsByAlbumIDs(ctx, albumIDs)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Len(t, result.Artists[1].Artists, 1)
	assert.Len(t, result.Artists[2].Artists, 1)
	assert.Equal(t, mockArtistsMap[1][0].ID, result.Artists[1].Artists[0].ID)
	assert.Equal(t, mockArtistsMap[2][0].ID, result.Artists[2].Artists[0].ID)
}

func TestGetAlbumIDsByArtistID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)

	mockAlbumIDs := []int64{1, 2, 3}

	mockRepo.EXPECT().GetAlbumIDsByArtistID(ctx, artistID).Return(mockAlbumIDs, nil)

	albumIDs, err := usecase.GetAlbumIDsByArtistID(ctx, artistID)

	assert.NoError(t, err)
	assert.Equal(t, mockAlbumIDs, albumIDs)
}

func TestGetTrackIDsByArtistID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	artistID := int64(1)

	mockTrackIDs := []int64{1, 2, 3}

	mockRepo.EXPECT().GetTrackIDsByArtistID(ctx, artistID).Return(mockTrackIDs, nil)

	trackIDs, err := usecase.GetTrackIDsByArtistID(ctx, artistID)

	assert.NoError(t, err)
	assert.Equal(t, mockTrackIDs, trackIDs)
}

func TestCreateStreamsByArtistIDs(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)

	data := &usecaseModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{1, 2},
		UserID:    1,
	}

	repoData := &repoModel.ArtistStreamCreateDataList{
		ArtistIDs: []int64{1, 2},
		UserID:    1,
	}

	mockRepo.EXPECT().CreateStreamsByArtistIDs(ctx, repoData).Return(nil)

	err := usecase.CreateStreamsByArtistIDs(ctx, data)

	assert.NoError(t, err)
}

func TestGetArtistsListenedByUserID(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	userID := int64(1)
	expected := int64(5)

	mockRepo.EXPECT().GetArtistsListenedByUserID(ctx, userID).Return(expected, nil)

	count, err := usecase.GetArtistsListenedByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func TestLikeArtistSuccess(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
		IsLike:   true,
	}

	repoRequest := &repoModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, request.ArtistID).Return(true, nil)
	mockRepo.EXPECT().LikeArtist(ctx, repoRequest).Return(nil)

	err := usecase.LikeArtist(ctx, request)

	assert.NoError(t, err)
}

func TestLikeArtistUnlike(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
		IsLike:   false,
	}

	repoRequest := &repoModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, request.ArtistID).Return(true, nil)
	mockRepo.EXPECT().UnlikeArtist(ctx, repoRequest).Return(nil)

	err := usecase.LikeArtist(ctx, request)

	assert.NoError(t, err)
}

func TestLikeArtistNotFound(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		ArtistID: 1,
		UserID:   1,
		IsLike:   true,
	}

	mockRepo.EXPECT().CheckArtistExists(ctx, request.ArtistID).Return(false, nil)

	err := usecase.LikeArtist(ctx, request)

	assert.Error(t, err)
	assert.IsType(t, artistErrors.ErrArtistNotFound, err)
}

func TestGetFavoriteArtists(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	userID := int64(1)

	filters := &usecaseModel.Filters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoFilters := &repoModel.Filters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockArtists := []*repoModel.Artist{
		{
			ID:          1,
			Title:       "Artist 1",
			Description: "Description 1",
			Thumbnail:   "thumb1.jpg",
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Artist 2",
			Description: "Description 2",
			Thumbnail:   "thumb2.jpg",
			IsFavorite:  true,
		},
	}

	mockRepo.EXPECT().GetFavoriteArtists(ctx, repoFilters, userID).Return(mockArtists, nil)

	result, err := usecase.GetFavoriteArtists(ctx, filters, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, mockArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, mockArtists[0].Title, result.Artists[0].Title)
	assert.Equal(t, mockArtists[1].ID, result.Artists[1].ID)
	assert.Equal(t, mockArtists[1].Title, result.Artists[1].Title)
}

func TestSearchArtists(t *testing.T) {
	mockRepo, ctx := setupTest(t)

	usecase := NewArtistUsecase(mockRepo)
	query := "test query"
	userID := int64(1)

	mockArtists := []*repoModel.Artist{
		{
			ID:          1,
			Title:       "Test Artist",
			Description: "Description 1",
			Thumbnail:   "thumb1.jpg",
			IsFavorite:  false,
		},
		{
			ID:          2,
			Title:       "Another Test",
			Description: "Description 2",
			Thumbnail:   "thumb2.jpg",
			IsFavorite:  true,
		},
	}

	mockRepo.EXPECT().SearchArtists(ctx, query, userID).Return(mockArtists, nil)

	result, err := usecase.SearchArtists(ctx, query, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Artists, 2)
	assert.Equal(t, mockArtists[0].ID, result.Artists[0].ID)
	assert.Equal(t, mockArtists[0].Title, result.Artists[0].Title)
	assert.Equal(t, mockArtists[1].ID, result.Artists[1].ID)
	assert.Equal(t, mockArtists[1].Title, result.Artists[1].Title)
}
