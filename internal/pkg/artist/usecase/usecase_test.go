package usecase

import (
	"context"
	"errors"
	"testing"

	mock_artist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestArtistUsecase_GetArtistByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockRepository(ctrl)
	usecase := NewUsecase(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name           string
		artistID       int64
		mockSetup      func()
		expectedArtist *usecaseModel.ArtistDetailed
		expectedError  error
	}{
		{
			name:     "Success",
			artistID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetArtistByID(ctx, int64(1)).
					Return(&repository.Artist{
						ID:          1,
						Title:       "Test Artist",
						Thumbnail:   "test-thumbnail.jpg",
						Description: "Test Description",
					}, nil)

				mockRepo.EXPECT().
					GetArtistStats(ctx, int64(1)).
					Return(&repository.ArtistStats{
						ListenersCount: 1000,
						FavoritesCount: 500,
					}, nil)
			},
			expectedArtist: &usecaseModel.ArtistDetailed{
				Artist: usecaseModel.Artist{
					ID:          1,
					Title:       "Test Artist",
					Thumbnail:   "test-thumbnail.jpg",
					Description: "Test Description",
				},
				Listeners: 1000,
				Favorites: 500,
			},
			expectedError: nil,
		},
		{
			name:     "Error_GetArtistByID",
			artistID: 2,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetArtistByID(ctx, int64(2)).
					Return(nil, errors.New("database error"))
			},
			expectedArtist: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name:     "Error_GetArtistStats",
			artistID: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetArtistByID(ctx, int64(3)).
					Return(&repository.Artist{
						ID:          3,
						Title:       "Test Artist",
						Thumbnail:   "test-thumbnail.jpg",
						Description: "Test Description",
					}, nil)

				mockRepo.EXPECT().
					GetArtistStats(ctx, int64(3)).
					Return(nil, errors.New("stats error"))
			},
			expectedArtist: nil,
			expectedError:  errors.New("stats error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			artist, err := usecase.GetArtistByID(ctx, tt.artistID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, artist)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtist, artist)
			}
		})
	}
}

func TestArtistUsecase_GetAllArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockRepository(ctrl)
	usecase := NewUsecase(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name            string
		filters         *usecaseModel.ArtistFilters
		mockSetup       func()
		expectedArtists []*usecaseModel.Artist
		expectedError   error
	}{
		{
			name: "Success",
			filters: &usecaseModel.ArtistFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllArtists(ctx, gomock.Any()).
					Return([]*repository.Artist{
						{
							ID:          1,
							Title:       "Artist 1",
							Thumbnail:   "thumbnail1.jpg",
							Description: "Description 1",
						},
						{
							ID:          2,
							Title:       "Artist 2",
							Thumbnail:   "thumbnail2.jpg",
							Description: "Description 2",
						},
					}, nil)
			},
			expectedArtists: []*usecaseModel.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Thumbnail:   "thumbnail1.jpg",
					Description: "Description 1",
				},
				{
					ID:          2,
					Title:       "Artist 2",
					Thumbnail:   "thumbnail2.jpg",
					Description: "Description 2",
				},
			},
			expectedError: nil,
		},
		{
			name: "Error_GetAllArtists",
			filters: &usecaseModel.ArtistFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllArtists(ctx, gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedArtists: nil,
			expectedError:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			artists, err := usecase.GetAllArtists(ctx, tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, artists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
		})
	}
}
