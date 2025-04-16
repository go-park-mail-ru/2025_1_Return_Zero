package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_album "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/mocks"
	mock_artist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	mock_track "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/mocks"
	mock_trackFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/trackFile/mocks"
	mock_user "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/mocks"
)

func TestTrackUsecase_GetAllTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name           string
		filters        *usecaseModel.TrackFilters
		mockSetup      func()
		expectedTracks []*usecaseModel.Track
		expectedError  error
	}{
		{
			name: "Success",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetAllTracks(ctx, gomock.Any()).
					Return([]*repository.Track{
						{
							ID:      1,
							Title:   "Track 1",
							AlbumID: 10,
						},
						{
							ID:      2,
							Title:   "Track 2",
							AlbumID: 20,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackIDs(ctx, []int64{1, 2}).
					Return(map[int64][]*repository.ArtistWithRole{
						1: {
							{
								ID:    100,
								Title: "Artist 100",
								Role:  "main",
							},
						},
						2: {
							{
								ID:    200,
								Title: "Artist 200",
								Role:  "main",
							},
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByIDs(ctx, []int64{10, 20}).
					Return(map[int64]string{
						10: "Album 10",
						20: "Album 20",
					}, nil)
			},
			expectedTracks: []*usecaseModel.Track{
				{
					ID:      1,
					Title:   "Track 1",
					AlbumID: 10,
					Album:   "Album 10",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    100,
							Title: "Artist 100",
							Role:  "main",
						},
					},
				},
				{
					ID:      2,
					Title:   "Track 2",
					AlbumID: 20,
					Album:   "Album 20",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    200,
							Title: "Artist 200",
							Role:  "main",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error_GetAllTracks",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetAllTracks(ctx, gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedTracks: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Error_GetArtistsByTrackIDs",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetAllTracks(ctx, gomock.Any()).
					Return([]*repository.Track{
						{
							ID:      1,
							Title:   "Track 1",
							AlbumID: 10,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackIDs(ctx, []int64{1}).
					Return(nil, errors.New("artist error"))
			},
			expectedTracks: nil,
			expectedError:  errors.New("artist error"),
		},
		{
			name: "Error_GetAlbumTitleByIDs",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetAllTracks(ctx, gomock.Any()).
					Return([]*repository.Track{
						{
							ID:      1,
							Title:   "Track 1",
							AlbumID: 10,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackIDs(ctx, []int64{1}).
					Return(map[int64][]*repository.ArtistWithRole{
						1: {
							{
								ID:    100,
								Title: "Artist 100",
								Role:  "main",
							},
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByIDs(ctx, []int64{10}).
					Return(nil, errors.New("album error"))
			},
			expectedTracks: nil,
			expectedError:  errors.New("album error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tracks, err := usecase.GetAllTracks(ctx, tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, tracks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
		})
	}
}

func TestTrackUsecase_GetTrackByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name          string
		trackID       int64
		mockSetup     func()
		expectedTrack *usecaseModel.TrackDetailed
		expectedError error
	}{
		{
			name:    "Success",
			trackID: 1,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTrackByID(ctx, int64(1)).
					Return(&repository.TrackWithFileKey{
						Track: repository.Track{
							ID:       1,
							Title:    "Test Track",
							Duration: 200,
							AlbumID:  10,
						},
						FileKey: "test-file-key",
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackID(ctx, int64(1)).
					Return([]*repository.ArtistWithRole{
						{
							ID:    100,
							Title: "Test Artist",
							Role:  "main",
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByID(ctx, int64(10)).
					Return("Test Album", nil)

				mockTrackFileRepo.EXPECT().
					GetPresignedURL("test-file-key").
					Return("https://example.com/track.mp3", nil)
			},
			expectedTrack: &usecaseModel.TrackDetailed{
				Track: usecaseModel.Track{
					ID:       1,
					Title:    "Test Track",
					Duration: 200,
					AlbumID:  10,
					Album:    "Test Album",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    100,
							Title: "Test Artist",
							Role:  "main",
						},
					},
				},
				FileUrl: "https://example.com/track.mp3",
			},
			expectedError: nil,
		},
		{
			name:    "Error_GetTrackByID",
			trackID: 2,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTrackByID(ctx, int64(2)).
					Return(nil, errors.New("database error"))
			},
			expectedTrack: nil,
			expectedError: errors.New("database error"),
		},
		{
			name:    "Error_GetArtistsByTrackID",
			trackID: 3,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTrackByID(ctx, int64(3)).
					Return(&repository.TrackWithFileKey{
						Track: repository.Track{
							ID:       3,
							Title:    "Test Track",
							Duration: 200,
							AlbumID:  10,
						},
						FileKey: "test-file-key",
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackID(ctx, int64(3)).
					Return(nil, errors.New("artist error"))
			},
			expectedTrack: nil,
			expectedError: errors.New("artist error"),
		},
		{
			name:    "Error_GetAlbumTitleByID",
			trackID: 4,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTrackByID(ctx, int64(4)).
					Return(&repository.TrackWithFileKey{
						Track: repository.Track{
							ID:       4,
							Title:    "Test Track",
							Duration: 200,
							AlbumID:  10,
						},
						FileKey: "test-file-key",
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackID(ctx, int64(4)).
					Return([]*repository.ArtistWithRole{
						{
							ID:    100,
							Title: "Test Artist",
							Role:  "main",
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByID(ctx, int64(10)).
					Return("", errors.New("album error"))
			},
			expectedTrack: nil,
			expectedError: errors.New("album error"),
		},
		{
			name:    "Error_GetPresignedURL",
			trackID: 5,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTrackByID(ctx, int64(5)).
					Return(&repository.TrackWithFileKey{
						Track: repository.Track{
							ID:       5,
							Title:    "Test Track",
							Duration: 200,
							AlbumID:  10,
						},
						FileKey: "test-file-key",
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackID(ctx, int64(5)).
					Return([]*repository.ArtistWithRole{
						{
							ID:    100,
							Title: "Test Artist",
							Role:  "main",
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByID(ctx, int64(10)).
					Return("Test Album", nil)

				mockTrackFileRepo.EXPECT().
					GetPresignedURL("test-file-key").
					Return("", errors.New("file error"))
			},
			expectedTrack: nil,
			expectedError: errors.New("file error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			track, err := usecase.GetTrackByID(ctx, tt.trackID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, track)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTrack, track)
			}
		})
	}
}

func TestTrackUsecase_GetTracksByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name           string
		artistID       int64
		mockSetup      func()
		expectedTracks []*usecaseModel.Track
		expectedError  error
	}{
		{
			name:     "Success",
			artistID: 100,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTracksByArtistID(ctx, int64(100)).
					Return([]*repository.Track{
						{
							ID:      1,
							Title:   "Track 1",
							AlbumID: 10,
						},
						{
							ID:      2,
							Title:   "Track 2",
							AlbumID: 20,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackIDs(ctx, []int64{1, 2}).
					Return(map[int64][]*repository.ArtistWithRole{
						1: {
							{
								ID:    100,
								Title: "Artist 100",
								Role:  "main",
							},
						},
						2: {
							{
								ID:    100,
								Title: "Artist 100",
								Role:  "main",
							},
							{
								ID:    200,
								Title: "Artist 200",
								Role:  "featured",
							},
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByIDs(ctx, []int64{10, 20}).
					Return(map[int64]string{
						10: "Album 10",
						20: "Album 20",
					}, nil)
			},
			expectedTracks: []*usecaseModel.Track{
				{
					ID:      1,
					Title:   "Track 1",
					AlbumID: 10,
					Album:   "Album 10",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    100,
							Title: "Artist 100",
							Role:  "main",
						},
					},
				},
				{
					ID:      2,
					Title:   "Track 2",
					AlbumID: 20,
					Album:   "Album 20",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    100,
							Title: "Artist 100",
							Role:  "main",
						},
						{
							ID:    200,
							Title: "Artist 200",
							Role:  "featured",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "Error_GetTracksByArtistID",
			artistID: 101,
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetTracksByArtistID(ctx, int64(101)).
					Return(nil, errors.New("database error"))
			},
			expectedTracks: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tracks, err := usecase.GetTracksByArtistID(ctx, tt.artistID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, tracks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
		})
	}
}

func TestTrackUsecase_CreateStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name          string
		stream        *usecaseModel.TrackStreamCreateData
		mockSetup     func()
		expectedID    int64
		expectedError error
	}{
		{
			name: "Success",
			stream: &usecaseModel.TrackStreamCreateData{
				TrackID: 1,
				UserID:  100,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					CreateStream(ctx, gomock.Any()).
					Return(int64(1), nil)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "Error",
			stream: &usecaseModel.TrackStreamCreateData{
				TrackID: 2,
				UserID:  200,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					CreateStream(ctx, gomock.Any()).
					Return(int64(0), errors.New("database error"))
			},
			expectedID:    0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := usecase.CreateStream(ctx, tt.stream)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Equal(t, tt.expectedID, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestTrackUsecase_UpdateStreamDuration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name          string
		streamUpdate  *usecaseModel.TrackStreamUpdateData
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			streamUpdate: &usecaseModel.TrackStreamUpdateData{
				StreamID: 1,
				UserID:   100,
				Duration: 120,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetStreamByID(ctx, int64(1)).
					Return(&repository.TrackStream{
						ID:     1,
						UserID: 100,
					}, nil)

				mockTrackRepo.EXPECT().
					UpdateStreamDuration(ctx, gomock.Any()).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error_GetStreamByID",
			streamUpdate: &usecaseModel.TrackStreamUpdateData{
				StreamID: 2,
				UserID:   200,
				Duration: 120,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetStreamByID(ctx, int64(2)).
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "Error_PermissionDenied",
			streamUpdate: &usecaseModel.TrackStreamUpdateData{
				StreamID: 3,
				UserID:   300,
				Duration: 120,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetStreamByID(ctx, int64(3)).
					Return(&repository.TrackStream{
						ID:     3,
						UserID: 301,
					}, nil)
			},
			expectedError: track.ErrStreamPermissionDenied,
		},
		{
			name: "Error_UpdateStreamDuration",
			streamUpdate: &usecaseModel.TrackStreamUpdateData{
				StreamID: 4,
				UserID:   400,
				Duration: 120,
			},
			mockSetup: func() {
				mockTrackRepo.EXPECT().
					GetStreamByID(ctx, int64(4)).
					Return(&repository.TrackStream{
						ID:     4,
						UserID: 400,
					}, nil)

				mockTrackRepo.EXPECT().
					UpdateStreamDuration(ctx, gomock.Any()).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := usecase.UpdateStreamDuration(ctx, tt.streamUpdate)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackUsecase_GetLastListenedTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrackRepo := mock_track.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockTrackFileRepo := mock_trackFile.NewMockRepository(ctrl)
	mockUserRepo := mock_user.NewMockRepository(ctrl)

	usecase := NewUsecase(mockTrackRepo, mockArtistRepo, mockAlbumRepo, mockTrackFileRepo, mockUserRepo)
	ctx := context.Background()

	tests := []struct {
		name           string
		username       string
		filters        *usecaseModel.TrackFilters
		mockSetup      func()
		expectedTracks []*usecaseModel.Track
		expectedError  error
	}{
		{
			name:     "Success",
			username: "testuser",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetIDByUsername(ctx, "testuser").
					Return(int64(100), nil)

				mockTrackRepo.EXPECT().
					GetStreamsByUserID(ctx, int64(100), gomock.Any()).
					Return([]*repository.TrackStream{
						{
							ID:      1,
							TrackID: 10,
							UserID:  100,
						},
						{
							ID:      2,
							TrackID: 20,
							UserID:  100,
						},
					}, nil)

				mockTrackRepo.EXPECT().
					GetTracksByIDs(ctx, []int64{10, 20}).
					Return(map[int64]*repository.Track{
						10: {
							ID:      10,
							Title:   "Track 10",
							AlbumID: 100,
						},
						20: {
							ID:      20,
							Title:   "Track 20",
							AlbumID: 200,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByTrackIDs(ctx, []int64{10, 20}).
					Return(map[int64][]*repository.ArtistWithRole{
						10: {
							{
								ID:    1000,
								Title: "Artist 1000",
								Role:  "main",
							},
						},
						20: {
							{
								ID:    2000,
								Title: "Artist 2000",
								Role:  "main",
							},
						},
					}, nil)

				mockAlbumRepo.EXPECT().
					GetAlbumTitleByIDs(ctx, []int64{100, 200}).
					Return(map[int64]string{
						100: "Album 100",
						200: "Album 200",
					}, nil)
			},
			expectedTracks: []*usecaseModel.Track{
				{
					ID:      10,
					Title:   "Track 10",
					AlbumID: 100,
					Album:   "Album 100",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    1000,
							Title: "Artist 1000",
							Role:  "main",
						},
					},
				},
				{
					ID:      20,
					Title:   "Track 20",
					AlbumID: 200,
					Album:   "Album 200",
					Artists: []*usecaseModel.TrackArtist{
						{
							ID:    2000,
							Title: "Artist 2000",
							Role:  "main",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "Error_GetIDByUsername",
			username: "nonexistent",
			filters: &usecaseModel.TrackFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().
					GetIDByUsername(ctx, "nonexistent").
					Return(int64(0), errors.New("user not found"))
			},
			expectedTracks: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tracks, err := usecase.GetLastListenedTracks(ctx, tt.username, tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, tracks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTracks, tracks)
			}
		})
	}
}
