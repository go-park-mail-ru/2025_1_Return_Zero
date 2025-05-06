package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/mocks"
	albumErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
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

func TestGetAllAlbums(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	repoFilters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	repoAlbums := []*repoModel.Album{
		{
			ID:          1,
			Title:       "Album 1",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url1",
			ReleaseDate: time.Now(),
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Album 2",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url2",
			ReleaseDate: time.Now(),
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().GetAllAlbums(ctx, repoFilters, userID).Return(repoAlbums, nil)

	albums, err := usecase.GetAllAlbums(ctx, filters, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestGetAllAlbumsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	repoFilters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAllAlbums(ctx, repoFilters, userID).Return(nil, expectedErr)

	albums, err := usecase.GetAllAlbums(ctx, filters, userID)

	require.Error(t, err)
	assert.Nil(t, albums)
	assert.Equal(t, expectedErr, err)
}

func TestGetAlbumByID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)
	userID := int64(1)
	releaseDate := time.Now()

	repoAlbum := &repoModel.Album{
		ID:          albumID,
		Title:       "Test Album",
		Type:        repoModel.AlbumTypeAlbum,
		Thumbnail:   "thumbnail_url",
		ReleaseDate: releaseDate,
		IsFavorite:  true,
	}

	mockRepo.EXPECT().GetAlbumByID(ctx, albumID, userID).Return(repoAlbum, nil)

	album, err := usecase.GetAlbumByID(ctx, albumID, userID)

	require.NoError(t, err)
	assert.Equal(t, albumID, album.ID)
	assert.Equal(t, "Test Album", album.Title)
	assert.Equal(t, usecaseModel.AlbumTypeAlbum, album.Type)
	assert.Equal(t, "thumbnail_url", album.Thumbnail)
	assert.True(t, album.IsFavorite)
}

func TestGetAlbumByIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)
	userID := int64(1)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAlbumByID(ctx, albumID, userID).Return(nil, expectedErr)

	album, err := usecase.GetAlbumByID(ctx, albumID, userID)

	require.Error(t, err)
	assert.Nil(t, album)
	assert.Equal(t, expectedErr, err)
}

func TestGetAlbumTitleByID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)
	expectedTitle := "Test Album"

	mockRepo.EXPECT().GetAlbumTitleByID(ctx, albumID).Return(expectedTitle, nil)

	title, err := usecase.GetAlbumTitleByID(ctx, albumID)

	require.NoError(t, err)
	assert.Equal(t, expectedTitle, title)
}

func TestGetAlbumTitleByIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAlbumTitleByID(ctx, albumID).Return("", expectedErr)

	title, err := usecase.GetAlbumTitleByID(ctx, albumID)

	require.Error(t, err)
	assert.Equal(t, "", title)
	assert.Equal(t, expectedErr, err)
}

func TestGetAlbumTitleByIDs(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumIDs := []int64{1, 2}
	repoTitles := map[int64]string{
		1: "Album 1",
		2: "Album 2",
	}

	mockRepo.EXPECT().GetAlbumTitleByIDs(ctx, albumIDs).Return(repoTitles, nil)

	titles, err := usecase.GetAlbumTitleByIDs(ctx, albumIDs)

	require.NoError(t, err)
	assert.NotNil(t, titles)
	assert.Equal(t, "Album 1", titles.Titles[1].Title)
	assert.Equal(t, "Album 2", titles.Titles[2].Title)
}

func TestGetAlbumTitleByIDsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumIDs := []int64{1, 2}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAlbumTitleByIDs(ctx, albumIDs).Return(nil, expectedErr)

	titles, err := usecase.GetAlbumTitleByIDs(ctx, albumIDs)

	require.Error(t, err)
	assert.Nil(t, titles)
	assert.Equal(t, expectedErr, err)
}

func TestGetAlbumsByIDs(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumIDs := []int64{1, 2}
	userID := int64(1)
	releaseDate := time.Now()

	repoAlbums := []*repoModel.Album{
		{
			ID:          1,
			Title:       "Album 1",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url1",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Album 2",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url2",
			ReleaseDate: releaseDate,
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().GetAlbumsByIDs(ctx, albumIDs, userID).Return(repoAlbums, nil)

	albums, err := usecase.GetAlbumsByIDs(ctx, albumIDs, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestGetAlbumsByIDsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumIDs := []int64{1, 2}
	userID := int64(1)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetAlbumsByIDs(ctx, albumIDs, userID).Return(nil, expectedErr)

	albums, err := usecase.GetAlbumsByIDs(ctx, albumIDs, userID)

	require.Error(t, err)
	assert.Nil(t, albums)
	assert.Equal(t, expectedErr, err)
}

func TestCreateStream(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)
	userID := int64(1)

	mockRepo.EXPECT().CreateStream(ctx, albumID, userID).Return(nil)

	err := usecase.CreateStream(ctx, albumID, userID)

	require.NoError(t, err)
}

func TestCreateStreamError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	albumID := int64(1)
	userID := int64(1)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CreateStream(ctx, albumID, userID).Return(expectedErr)

	err := usecase.CreateStream(ctx, albumID, userID)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLikeAlbum(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  true,
	}

	repoRequest := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(true, nil)
	mockRepo.EXPECT().LikeAlbum(ctx, repoRequest).Return(nil)

	err := usecase.LikeAlbum(ctx, request)

	require.NoError(t, err)
}

func TestLikeAlbumNotFound(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  true,
	}

	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(false, nil)

	err := usecase.LikeAlbum(ctx, request)

	require.Error(t, err)
	assert.IsType(t, albumErrors.ErrAlbumNotFound, err)
}

func TestUnlikeAlbum(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  false,
	}

	repoRequest := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(true, nil)
	mockRepo.EXPECT().UnlikeAlbum(ctx, repoRequest).Return(nil)

	err := usecase.LikeAlbum(ctx, request)

	require.NoError(t, err)
}

func TestCheckAlbumExistsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  true,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(false, expectedErr)

	err := usecase.LikeAlbum(ctx, request)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLikeAlbumError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  true,
	}

	repoRequest := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(true, nil)
	mockRepo.EXPECT().LikeAlbum(ctx, repoRequest).Return(expectedErr)

	err := usecase.LikeAlbum(ctx, request)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUnlikeAlbumError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	request := &usecaseModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
		IsLike:  false,
	}

	repoRequest := &repoModel.LikeRequest{
		AlbumID: 1,
		UserID:  1,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CheckAlbumExists(ctx, request.AlbumID).Return(true, nil)
	mockRepo.EXPECT().UnlikeAlbum(ctx, repoRequest).Return(expectedErr)

	err := usecase.LikeAlbum(ctx, request)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestGetFavoriteAlbums(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	repoFilters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	releaseDate := time.Now()
	repoAlbums := []*repoModel.Album{
		{
			ID:          1,
			Title:       "Album 1",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url1",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Album 2",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url2",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
	}

	mockRepo.EXPECT().GetFavoriteAlbums(ctx, repoFilters, userID).Return(repoAlbums, nil)

	albums, err := usecase.GetFavoriteAlbums(ctx, filters, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Album 1", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Album 2", albums[1].Title)
	assert.True(t, albums[1].IsFavorite)
}

func TestGetFavoriteAlbumsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}
	userID := int64(1)

	repoFilters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().GetFavoriteAlbums(ctx, repoFilters, userID).Return(nil, expectedErr)

	albums, err := usecase.GetFavoriteAlbums(ctx, filters, userID)

	require.Error(t, err)
	assert.Nil(t, albums)
	assert.Equal(t, expectedErr, err)
}

func TestSearchAlbums(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	query := "test"
	userID := int64(1)
	releaseDate := time.Now()

	repoAlbums := []*repoModel.Album{
		{
			ID:          1,
			Title:       "Test Album",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url1",
			ReleaseDate: releaseDate,
			IsFavorite:  true,
		},
		{
			ID:          2,
			Title:       "Another Test",
			Type:        repoModel.AlbumTypeAlbum,
			Thumbnail:   "url2",
			ReleaseDate: releaseDate,
			IsFavorite:  false,
		},
	}

	mockRepo.EXPECT().SearchAlbums(ctx, query, userID).Return(repoAlbums, nil)

	albums, err := usecase.SearchAlbums(ctx, query, userID)

	require.NoError(t, err)
	assert.Len(t, albums, 2)
	assert.Equal(t, int64(1), albums[0].ID)
	assert.Equal(t, "Test Album", albums[0].Title)
	assert.True(t, albums[0].IsFavorite)
	assert.Equal(t, int64(2), albums[1].ID)
	assert.Equal(t, "Another Test", albums[1].Title)
	assert.False(t, albums[1].IsFavorite)
}

func TestSearchAlbumsError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	usecase := NewAlbumUsecase(mockRepo)

	query := "test"
	userID := int64(1)

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().SearchAlbums(ctx, query, userID).Return(nil, expectedErr)

	albums, err := usecase.SearchAlbums(ctx, query, userID)

	require.Error(t, err)
	assert.Nil(t, albums)
	assert.Equal(t, expectedErr, err)
}
