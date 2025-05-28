package usecase

import (
	"context"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/mocks"
	"errors"

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

func TestCreateAlbum(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
	usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

	createAlbumRequest := &usecaseModel.CreateAlbumRequest{
		Title:   "New Album",
		Type:    "album",
		LabelID: 1,
	}

	mockRepo.EXPECT().CreateAlbum(ctx, gomock.Any()).Return(int64(1), nil)

	mockS3Repo.EXPECT().UploadAlbumAvatar(ctx, "New Album", gomock.Any()).Return("thumbnail_url", nil)

	albumID, thumbnail_url, err := usecase.CreateAlbum(ctx, createAlbumRequest)
	require.NoError(t, err)
	assert.Equal(t, int64(1), albumID)
	assert.Equal(t, "thumbnail_url", thumbnail_url)
}


func TestCreateAlbumError(t *testing.T) {
    mockRepo, ctx := setupTest(t)
    mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
    usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

    createAlbumRequest := &usecaseModel.CreateAlbumRequest{
        Title:   "New Album",
        Type:    "album",
        LabelID: 1,
        Image:   []byte("test image"),
    }

    mockS3Repo.EXPECT().UploadAlbumAvatar(ctx, "New Album", gomock.Any()).Return("thumbnail_url", nil)
    
    mockRepo.EXPECT().CreateAlbum(ctx, gomock.Any()).Return(int64(0), errors.New("failed to create album"))

    albumID, thumbnail_url, err := usecase.CreateAlbum(ctx, createAlbumRequest)
    require.Error(t, err)
    assert.Equal(t, int64(0), albumID)
    assert.Empty(t, thumbnail_url)
}

func TestDeleteAlbum(t *testing.T) {
    mockRepo, ctx := setupTest(t)
    mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
    usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

    albumID := int64(1)

    mockRepo.EXPECT().CheckAlbumExists(ctx, albumID).Return(true, nil)
    
    mockRepo.EXPECT().DeleteAlbum(ctx, albumID).Return(nil)

    err := usecase.DeleteAlbum(ctx, albumID)
    require.NoError(t, err)
}

func TestDeleteAlbumNotFound(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
	usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

	albumID := int64(1)

	mockRepo.EXPECT().CheckAlbumExists(ctx, albumID).Return(false, nil)

	err := usecase.DeleteAlbum(ctx, albumID)
	require.Error(t, err)
	assert.Equal(t, "album not found", err.Error())
}

func TestDeleteAlbumError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
	usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

	albumID := int64(1)

	mockRepo.EXPECT().CheckAlbumExists(ctx, albumID).Return(true, nil)
	
	mockRepo.EXPECT().DeleteAlbum(ctx, albumID).Return(errors.New("failed to delete album"))

	err := usecase.DeleteAlbum(ctx, albumID)
	require.Error(t, err)
	assert.Equal(t, "failed to delete album", err.Error())
}

func TestGetAlbumsLabelID(t *testing.T) {
    mockRepo, ctx := setupTest(t)
    mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
    usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

    labelID := int64(1)
    expectedAlbums := []*repoModel.Album{
        {ID: 1, Title: "Album 1"},
        {ID: 2, Title: "Album 2"},
    }
    filters := &usecaseModel.AlbumFilters{
        Pagination: &usecaseModel.Pagination{
            Limit:  10,
            Offset: 0,
        },
    }

    mockRepo.EXPECT().GetAlbumsLabelID(ctx, gomock.Any(), labelID).Return(expectedAlbums, nil)

    albums, err := usecase.GetAlbumsLabelID(ctx, filters, labelID)
    require.NoError(t, err)
    assert.Len(t, albums, 2)
    assert.Equal(t, expectedAlbums[0].ID, albums[0].ID)
    assert.Equal(t, expectedAlbums[0].Title, albums[0].Title)
    assert.Equal(t, expectedAlbums[1].ID, albums[1].ID)
    assert.Equal(t, expectedAlbums[1].Title, albums[1].Title)
}

func TestGetAlbumsLabelIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))
	usecase := NewAlbumUsecase(mockRepo, mockS3Repo)

	labelID := int64(1)
	filters := &usecaseModel.AlbumFilters{
		Pagination: &usecaseModel.Pagination{
			Limit:  10,
			Offset: 0,
		},
	}

	mockRepo.EXPECT().GetAlbumsLabelID(ctx, gomock.Any(), labelID).Return(nil, assert.AnError)

	albums, err := usecase.GetAlbumsLabelID(ctx, filters, labelID)
	require.Error(t, err)
	assert.Nil(t, albums)
}