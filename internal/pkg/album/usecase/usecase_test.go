package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_album "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/mocks"
	mock_artist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAlbumUsecase_GetAllAlbums(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	usecase := NewUsecase(mockAlbumRepo, mockArtistRepo)
	ctx := context.Background()

	date1, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	date2, _ := time.Parse(time.RFC3339, "2023-02-01T00:00:00Z")

	tests := []struct {
		name           string
		filters        *usecaseModel.AlbumFilters
		mockSetup      func()
		expectedAlbums []*usecaseModel.Album
		expectedError  error
	}{
		{
			name: "Success",
			filters: &usecaseModel.AlbumFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAllAlbums(ctx, gomock.Any()).
					Return([]*repository.Album{
						{
							ID:          1,
							Title:       "Album 1",
							Thumbnail:   "thumbnail1.jpg",
							ReleaseDate: date1,
						},
						{
							ID:          2,
							Title:       "Album 2",
							Thumbnail:   "thumbnail2.jpg",
							ReleaseDate: date2,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByAlbumIDs(ctx, []int64{1, 2}).
					Return(map[int64][]*repository.ArtistWithTitle{
						1: {
							{
								ID:    10,
								Title: "Artist 10",
							},
						},
						2: {
							{
								ID:    20,
								Title: "Artist 20",
							},
						},
					}, nil)
			},
			expectedAlbums: []*usecaseModel.Album{
				{
					ID:          1,
					Title:       "Album 1",
					Thumbnail:   "thumbnail1.jpg",
					ReleaseDate: date1,
					Artists: []*usecaseModel.AlbumArtist{
						{
							ID:    10,
							Title: "Artist 10",
						},
					},
				},
				{
					ID:          2,
					Title:       "Album 2",
					Thumbnail:   "thumbnail2.jpg",
					ReleaseDate: date2,
					Artists: []*usecaseModel.AlbumArtist{
						{
							ID:    20,
							Title: "Artist 20",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error_GetAllAlbums",
			filters: &usecaseModel.AlbumFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAllAlbums(ctx, gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			expectedAlbums: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Error_GetArtistsByAlbumIDs",
			filters: &usecaseModel.AlbumFilters{
				Pagination: &usecaseModel.Pagination{
					Limit:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAllAlbums(ctx, gomock.Any()).
					Return([]*repository.Album{
						{
							ID:          1,
							Title:       "Album 1",
							Thumbnail:   "thumbnail1.jpg",
							ReleaseDate: date1,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByAlbumIDs(ctx, []int64{1}).
					Return(nil, errors.New("artist repo error"))
			},
			expectedAlbums: nil,
			expectedError:  errors.New("artist repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			albums, err := usecase.GetAllAlbums(ctx, tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, albums)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlbums, albums)
			}
		})
	}
}

func TestAlbumUsecase_GetAlbumsByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumRepo := mock_album.NewMockRepository(ctrl)
	mockArtistRepo := mock_artist.NewMockRepository(ctrl)
	usecase := NewUsecase(mockAlbumRepo, mockArtistRepo)
	ctx := context.Background()

	date3, _ := time.Parse(time.RFC3339, "2023-03-00T00:00:00Z")
	date4, _ := time.Parse(time.RFC3339, "2023-04-00T00:00:00Z")
	date5, _ := time.Parse(time.RFC3339, "2023-05-00T00:00:00Z")

	tests := []struct {
		name           string
		artistID       int64
		mockSetup      func()
		expectedAlbums []*usecaseModel.Album
		expectedError  error
	}{
		{
			name:     "Success",
			artistID: 1,
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAlbumsByArtistID(ctx, int64(1)).
					Return([]*repository.Album{
						{
							ID:          10,
							Title:       "Artist 1 Album",
							Thumbnail:   "thumbnail10.jpg",
							ReleaseDate: date3,
						},
						{
							ID:          11,
							Title:       "Another Artist 1 Album",
							Thumbnail:   "thumbnail11.jpg",
							ReleaseDate: date4,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByAlbumIDs(ctx, []int64{10, 11}).
					Return(map[int64][]*repository.ArtistWithTitle{
						10: {
							{
								ID:    1,
								Title: "Artist 1",
							},
							{
								ID:    2,
								Title: "Artist 2",
							},
						},
						11: {
							{
								ID:    1,
								Title: "Artist 1",
							},
						},
					}, nil)
			},
			expectedAlbums: []*usecaseModel.Album{
				{
					ID:          10,
					Title:       "Artist 1 Album",
					Thumbnail:   "thumbnail10.jpg",
					ReleaseDate: date3,
					Artists: []*usecaseModel.AlbumArtist{
						{
							ID:    1,
							Title: "Artist 1",
						},
						{
							ID:    2,
							Title: "Artist 2",
						},
					},
				},
				{
					ID:          11,
					Title:       "Another Artist 1 Album",
					Thumbnail:   "thumbnail11.jpg",
					ReleaseDate: date4,
					Artists: []*usecaseModel.AlbumArtist{
						{
							ID:    1,
							Title: "Artist 1",
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "Error_GetAlbumsByArtistID",
			artistID: 2,
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAlbumsByArtistID(ctx, int64(2)).
					Return(nil, errors.New("database error"))
			},
			expectedAlbums: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name:     "Error_GetArtistsByAlbumIDs",
			artistID: 3,
			mockSetup: func() {
				mockAlbumRepo.EXPECT().
					GetAlbumsByArtistID(ctx, int64(3)).
					Return([]*repository.Album{
						{
							ID:          30,
							Title:       "Artist 3 Album",
							Thumbnail:   "thumbnail30.jpg",
							ReleaseDate: date5,
						},
					}, nil)

				mockArtistRepo.EXPECT().
					GetArtistsByAlbumIDs(ctx, []int64{30}).
					Return(nil, errors.New("artist repo error"))
			},
			expectedAlbums: nil,
			expectedError:  errors.New("artist repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			albums, err := usecase.GetAlbumsByArtistID(ctx, tt.artistID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, albums)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlbums, albums)
			}
		})
	}
}
