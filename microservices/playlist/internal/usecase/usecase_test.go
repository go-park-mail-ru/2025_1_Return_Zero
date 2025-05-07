package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/mocks"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/usecase"
	"github.com/stretchr/testify/assert"
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

func TestGetPlaylistByID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	playlistID := int64(1)
	userID := int64(2)
	isLiked := true

	mockPlaylist := &repoModel.Playlist{
		ID:        playlistID,
		Title:     "Test Playlist",
		Thumbnail: "test.jpg",
		UserID:    userID,
		IsPublic:  true,
	}

	mockPlaylistWithIsLiked := &repoModel.PlaylistWithIsLiked{
		Playlist: mockPlaylist,
		IsLiked:  isLiked,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)
	mockRepo.EXPECT().GetPlaylistWithIsLikedByID(ctx, playlistID, userID).Return(mockPlaylistWithIsLiked, nil)

	result, err := usecase.GetPlaylistByID(ctx, &usecaseModel.GetPlaylistByIDRequest{
		PlaylistID: playlistID,
		UserID:     userID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, playlistID, result.Playlist.ID)
	assert.Equal(t, mockPlaylist.Title, result.Playlist.Title)
	assert.Equal(t, mockPlaylist.Thumbnail, result.Playlist.Thumbnail)
	assert.Equal(t, mockPlaylist.UserID, result.Playlist.UserID)
	assert.Equal(t, isLiked, result.IsLiked)
}

func TestCreatePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	playlistID := int64(1)
	userID := int64(2)

	mockPlaylist := &repoModel.Playlist{
		ID:        playlistID,
		Title:     "Test Playlist",
		Thumbnail: "test.jpg",
		UserID:    userID,
		IsPublic:  true,
	}

	mockCreatePlaylistRequest := &repoModel.CreatePlaylistRequest{
		Title:     "Test Playlist",
		UserID:    userID,
		Thumbnail: "test.jpg",
		IsPublic:  true,
	}

	mockRepo.EXPECT().CreatePlaylist(ctx, mockCreatePlaylistRequest).Return(mockPlaylist, nil)

	result, err := usecase.CreatePlaylist(ctx, &usecaseModel.CreatePlaylistRequest{
		Title:     "Test Playlist",
		UserID:    userID,
		Thumbnail: "test.jpg",
		IsPublic:  true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockPlaylist.ID, result.ID)
	assert.Equal(t, mockPlaylist.Title, result.Title)
	assert.Equal(t, mockPlaylist.Thumbnail, result.Thumbnail)
	assert.Equal(t, mockPlaylist.UserID, result.UserID)
}

func TestUploadPlaylistThumbnail(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)

	randomImage := io.NopCloser(bytes.NewReader([]byte("test")))

	mockS3Repo.EXPECT().UploadThumbnail(ctx, randomImage, "Test Playlist").Return("test.jpg", nil)

	result, err := usecase.UploadPlaylistThumbnail(ctx, &usecaseModel.UploadPlaylistThumbnailRequest{
		Title:     "Test Playlist",
		Thumbnail: randomImage,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test.jpg", result)
}

func TestGetCombinedPlaylistsByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)

	mockPlaylistList := &repoModel.PlaylistList{
		Playlists: []*repoModel.Playlist{
			{
				ID:        1,
				Title:     "Test Playlist",
				Thumbnail: "test.jpg",
				UserID:    userID,
				IsPublic:  true,
			},
		},
	}

	mockRepo.EXPECT().GetCombinedPlaylistsByUserID(ctx, userID).Return(mockPlaylistList, nil)

	result, err := usecase.GetCombinedPlaylistsByUserID(ctx, &usecaseModel.GetCombinedPlaylistsByUserIDRequest{
		UserID: userID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAddTrackToPlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)
	trackID := int64(3)

	mockRepo.EXPECT().AddTrackToPlaylist(ctx, &repoModel.AddTrackToPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    trackID,
	}).Return(nil)

	err := usecase.AddTrackToPlaylist(ctx, &usecaseModel.AddTrackToPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    trackID,
	})

	assert.NoError(t, err)
}

func TestRemoveTrackFromPlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)
	trackID := int64(3)

	mockRepo.EXPECT().RemoveTrackFromPlaylist(ctx, &repoModel.RemoveTrackFromPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    trackID,
	}).Return(nil)

	err := usecase.RemoveTrackFromPlaylist(ctx, &usecaseModel.RemoveTrackFromPlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		TrackID:    trackID,
	})

	assert.NoError(t, err)
}

func TestGetPlaylistTrackIds(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID,
		IsPublic: true,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	mockRepo.EXPECT().GetPlaylistTrackIds(ctx, &repoModel.GetPlaylistTrackIdsRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return([]int64{1, 2, 3}, nil)

	result, err := usecase.GetPlaylistTrackIds(ctx, &usecaseModel.GetPlaylistTrackIdsRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFailGetPlaylistTrackIds(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	repoPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID + 1,
		IsPublic: false,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(repoPlaylist, nil)

	result, err := usecase.GetPlaylistTrackIds(ctx, &usecaseModel.GetPlaylistTrackIdsRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdatePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID,
		IsPublic: true,
	}
	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	mockRepo.EXPECT().UpdatePlaylist(ctx, &repoModel.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      "Test Playlist",
	}).Return(mockPlaylist, nil)

	result, err := usecase.UpdatePlaylist(ctx, &usecaseModel.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      "Test Playlist",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestFailUpdatePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID,
		IsPublic: true,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	mockRepo.EXPECT().UpdatePlaylist(ctx, &repoModel.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      "Test Playlist",
	}).Return(nil, errors.New("error"))

	result, err := usecase.UpdatePlaylist(ctx, &usecaseModel.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      "Test Playlist",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUnauthorizedUpdatePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID + 1,
		IsPublic: true,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	result, err := usecase.UpdatePlaylist(ctx, &usecaseModel.UpdatePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		Title:      "Test Playlist",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRemovePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID,
		IsPublic: true,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	mockRepo.EXPECT().RemovePlaylist(ctx, &repoModel.RemovePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return(nil)

	err := usecase.RemovePlaylist(ctx, &usecaseModel.RemovePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	})

	assert.NoError(t, err)
}

func TestFailRemovePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(nil, errors.New("error"))

	err := usecase.RemovePlaylist(ctx, &usecaseModel.RemovePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	})

	assert.Error(t, err)
}

func TestUnauthorizedRemovePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	mockPlaylist := &repoModel.Playlist{
		ID:       playlistID,
		UserID:   userID + 1,
		IsPublic: true,
	}

	mockRepo.EXPECT().GetPlaylistByID(ctx, playlistID).Return(mockPlaylist, nil)

	err := usecase.RemovePlaylist(ctx, &usecaseModel.RemovePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	})

	assert.Error(t, err)
}

func TestGetPlaylistsToAdd(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	trackID := int64(3)

	mockPlaylistWithIncluded := &repoModel.PlaylistWithIsIncludedTrack{
		Playlist: &repoModel.Playlist{
			ID:        1,
			Title:     "Test Playlist",
			Thumbnail: "test.jpg",
			UserID:    userID,
			IsPublic:  true,
		},
		IsIncluded: true,
	}

	mockResponse := &repoModel.GetPlaylistsToAddResponse{
		Playlists: []*repoModel.PlaylistWithIsIncludedTrack{
			mockPlaylistWithIncluded,
		},
	}

	mockRepo.EXPECT().GetPlaylistsToAdd(ctx, &repoModel.GetPlaylistsToAddRequest{
		UserID:  userID,
		TrackID: trackID,
	}).Return(mockResponse, nil)

	result, err := usecase.GetPlaylistsToAdd(ctx, &usecaseModel.GetPlaylistsToAddRequest{
		UserID:  userID,
		TrackID: trackID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Playlists))
	assert.Equal(t, int64(1), result.Playlists[0].Playlist.ID)
	assert.Equal(t, true, result.Playlists[0].IsIncluded)
}

func TestUpdatePlaylistsPublisityByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	isPublic := true

	mockRepo.EXPECT().UpdatePlaylistsPublisityByUserID(ctx, &repoModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   userID,
		IsPublic: isPublic,
	}).Return(nil)

	err := usecase.UpdatePlaylistsPublisityByUserID(ctx, &usecaseModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   userID,
		IsPublic: isPublic,
	})

	assert.NoError(t, err)
}

func TestLikePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	// Test liking a playlist
	mockRepo.EXPECT().LikePlaylist(ctx, &repoModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return(nil)

	err := usecase.LikePlaylist(ctx, &usecaseModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		IsLike:     true,
	})

	assert.NoError(t, err)

	// Test unliking a playlist
	mockRepo.EXPECT().UnlikePlaylist(ctx, &repoModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return(nil)

	err = usecase.LikePlaylist(ctx, &usecaseModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		IsLike:     false,
	})

	assert.NoError(t, err)
}

func TestGetProfilePlaylists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)

	mockPlaylists := []*repoModel.Playlist{
		{
			ID:        1,
			Title:     "Test Playlist 1",
			Thumbnail: "test1.jpg",
			UserID:    userID,
			IsPublic:  true,
		},
		{
			ID:        2,
			Title:     "Test Playlist 2",
			Thumbnail: "test2.jpg",
			UserID:    userID,
			IsPublic:  false,
		},
	}

	mockResponse := &repoModel.GetProfilePlaylistsResponse{
		Playlists: mockPlaylists,
	}

	mockRepo.EXPECT().GetProfilePlaylists(ctx, &repoModel.GetProfilePlaylistsRequest{
		UserID: userID,
	}).Return(mockResponse, nil)

	result, err := usecase.GetProfilePlaylists(ctx, &usecaseModel.GetProfilePlaylistsRequest{
		UserID: userID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Playlists))
	assert.Equal(t, "Test Playlist 1", result.Playlists[0].Title)
	assert.Equal(t, "Test Playlist 2", result.Playlists[1].Title)
}

func TestSearchPlaylists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	query := "rock"

	mockPlaylists := []*repoModel.Playlist{
		{
			ID:        1,
			Title:     "Rock Playlist",
			Thumbnail: "rock.jpg",
			UserID:    userID,
			IsPublic:  true,
		},
		{
			ID:        2,
			Title:     "Classic Rock",
			Thumbnail: "classic_rock.jpg",
			UserID:    userID + 1,
			IsPublic:  true,
		},
	}

	mockPlaylistList := &repoModel.PlaylistList{
		Playlists: mockPlaylists,
	}

	mockRepo.EXPECT().SearchPlaylists(ctx, &repoModel.SearchPlaylistsRequest{
		UserID: userID,
		Query:  query,
	}).Return(mockPlaylistList, nil)

	result, err := usecase.SearchPlaylists(ctx, &usecaseModel.SearchPlaylistsRequest{
		UserID: userID,
		Query:  query,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Playlists))
	assert.Equal(t, "Rock Playlist", result.Playlists[0].Title)
	assert.Equal(t, "Classic Rock", result.Playlists[1].Title)
}

func TestFailLikePlaylist(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	playlistID := int64(1)

	// Test error when liking a playlist
	mockRepo.EXPECT().LikePlaylist(ctx, &repoModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return(errors.New("error"))

	err := usecase.LikePlaylist(ctx, &usecaseModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		IsLike:     true,
	})

	assert.Error(t, err)

	// Test error when unliking a playlist
	mockRepo.EXPECT().UnlikePlaylist(ctx, &repoModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
	}).Return(errors.New("error"))

	err = usecase.LikePlaylist(ctx, &usecaseModel.LikePlaylistRequest{
		UserID:     userID,
		PlaylistID: playlistID,
		IsLike:     false,
	})

	assert.Error(t, err)
}

func TestFailGetProfilePlaylists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)

	mockRepo.EXPECT().GetProfilePlaylists(ctx, &repoModel.GetProfilePlaylistsRequest{
		UserID: userID,
	}).Return(nil, errors.New("error"))

	result, err := usecase.GetProfilePlaylists(ctx, &usecaseModel.GetProfilePlaylistsRequest{
		UserID: userID,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFailSearchPlaylists(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	query := "rock"

	mockRepo.EXPECT().SearchPlaylists(ctx, &repoModel.SearchPlaylistsRequest{
		UserID: userID,
		Query:  query,
	}).Return(nil, errors.New("error"))

	result, err := usecase.SearchPlaylists(ctx, &usecaseModel.SearchPlaylistsRequest{
		UserID: userID,
		Query:  query,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFailGetPlaylistsToAdd(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	trackID := int64(3)

	mockRepo.EXPECT().GetPlaylistsToAdd(ctx, &repoModel.GetPlaylistsToAddRequest{
		UserID:  userID,
		TrackID: trackID,
	}).Return(nil, errors.New("error"))

	result, err := usecase.GetPlaylistsToAdd(ctx, &usecaseModel.GetPlaylistsToAddRequest{
		UserID:  userID,
		TrackID: trackID,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFailUpdatePlaylistsPublisityByUserID(t *testing.T) {
	mockRepo, mockS3Repo, ctx := setupTest(t)

	usecase := NewPlaylistUsecase(mockRepo, mockS3Repo)
	userID := int64(2)
	isPublic := true

	mockRepo.EXPECT().UpdatePlaylistsPublisityByUserID(ctx, &repoModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   userID,
		IsPublic: isPublic,
	}).Return(errors.New("error"))

	err := usecase.UpdatePlaylistsPublisityByUserID(ctx, &usecaseModel.UpdatePlaylistsPublisityByUserIDRequest{
		UserID:   userID,
		IsPublic: isPublic,
	})

	assert.Error(t, err)
}
