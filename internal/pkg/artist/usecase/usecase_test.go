package usecase

import (
	"context"
	"errors"
	"testing"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetArtistByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	artistClient := artistProto.ArtistServiceClient(mockArtistClient)
	userClient := userProto.UserServiceClient(mockUserClient)

	artistUsecase := NewUsecase(artistClient, userClient)

	tests := []struct {
		name           string
		id             int64
		ctx            context.Context
		setupMocks     func()
		expectedArtist *usecase.ArtistDetailed
		expectedError  error
	}{
		{
			name: "Success with authenticated user",
			id:   1,
			ctx:  context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(2)),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetArtistByID(
					gomock.Any(),
					&artistProto.ArtistIDWithUserID{
						ArtistId: &artistProto.ArtistID{Id: 1},
						UserId:   &artistProto.UserID{Id: 2},
					},
				).Return(&artistProto.ArtistDetailed{
					Artist: &artistProto.Artist{
						Id:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test-thumbnail.jpg",
						IsFavorite:  true,
					},
					FavoritesCount: 100,
					ListenersCount: 500,
				}, nil)
			},
			expectedArtist: &usecase.ArtistDetailed{
				Artist: usecase.Artist{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test-thumbnail.jpg",
					IsLiked:     true,
				},
				Favorites: 100,
				Listeners: 500,
			},
			expectedError: nil,
		},
		{
			name: "Success with unauthenticated user",
			id:   1,
			ctx:  context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetArtistByID(
					gomock.Any(),
					&artistProto.ArtistIDWithUserID{
						ArtistId: &artistProto.ArtistID{Id: 1},
						UserId:   &artistProto.UserID{Id: -1},
					},
				).Return(&artistProto.ArtistDetailed{
					Artist: &artistProto.Artist{
						Id:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test-thumbnail.jpg",
						IsFavorite:  false,
					},
					FavoritesCount: 100,
					ListenersCount: 500,
				}, nil)
			},
			expectedArtist: &usecase.ArtistDetailed{
				Artist: usecase.Artist{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test-thumbnail.jpg",
					IsLiked:     false,
				},
				Favorites: 100,
				Listeners: 500,
			},
			expectedError: nil,
		},
		{
			name: "Error from service",
			id:   1,
			ctx:  context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetArtistByID(
					gomock.Any(),
					&artistProto.ArtistIDWithUserID{
						ArtistId: &artistProto.ArtistID{Id: 1},
						UserId:   &artistProto.UserID{Id: -1},
					},
				).Return(nil, status.Error(codes.NotFound, "artist not found"))
			},
			expectedArtist: nil,
			expectedError:  customErrors.ErrArtistNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			artist, err := artistUsecase.GetArtistByID(tt.ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtist, artist)
			}
		})
	}
}

func TestGetAllArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	// Convert mock to required type with type assertion
	artistClient := artistProto.ArtistServiceClient(mockArtistClient)
	userClient := userProto.UserServiceClient(mockUserClient)

	artistUsecase := NewUsecase(artistClient, userClient)

	tests := []struct {
		name            string
		filters         *usecase.ArtistFilters
		ctx             context.Context
		setupMocks      func()
		expectedArtists []*usecase.Artist
		expectedError   error
	}{
		{
			name: "Success with authenticated user",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			ctx: context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(1)),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetAllArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Artist 1",
							Description: "Description 1",
							Thumbnail:   "thumbnail1.jpg",
							IsFavorite:  true,
						},
						{
							Id:          2,
							Title:       "Artist 2",
							Description: "Description 2",
							Thumbnail:   "thumbnail2.jpg",
							IsFavorite:  false,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Description: "Description 1",
					Thumbnail:   "thumbnail1.jpg",
					IsLiked:     true,
				},
				{
					ID:          2,
					Title:       "Artist 2",
					Description: "Description 2",
					Thumbnail:   "thumbnail2.jpg",
					IsLiked:     false,
				},
			},
			expectedError: nil,
		},
		{
			name: "Success with unauthenticated user",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			ctx: context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetAllArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Artist 1",
							Description: "Description 1",
							Thumbnail:   "thumbnail1.jpg",
							IsFavorite:  false,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Description: "Description 1",
					Thumbnail:   "thumbnail1.jpg",
					IsLiked:     false,
				},
			},
			expectedError: nil,
		},
		{
			name: "Error from service",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			ctx: context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().GetAllArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedArtists: nil,
			expectedError:   errors.New("internal server error: internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			artists, err := artistUsecase.GetAllArtists(tt.ctx, tt.filters)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
		})
	}
}

func TestLikeArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	// Convert mock to required type with type assertion
	artistClient := artistProto.ArtistServiceClient(mockArtistClient)
	userClient := userProto.UserServiceClient(mockUserClient)

	artistUsecase := NewUsecase(artistClient, userClient)

	tests := []struct {
		name          string
		request       *usecase.ArtistLikeRequest
		ctx           context.Context
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success",
			request: &usecase.ArtistLikeRequest{
				ArtistID: 1,
				UserID:   2,
				IsLike:   true,
			},
			ctx: context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().LikeArtist(
					gomock.Any(),
					&artistProto.LikeRequest{
						ArtistId: &artistProto.ArtistID{Id: 1},
						UserId:   &artistProto.UserID{Id: 2},
						IsLike:   true,
					},
				).Return(&emptypb.Empty{}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Error from service",
			request: &usecase.ArtistLikeRequest{
				ArtistID: 1,
				UserID:   2,
				IsLike:   true,
			},
			ctx: context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().LikeArtist(
					gomock.Any(),
					&artistProto.LikeRequest{
						ArtistId: &artistProto.ArtistID{Id: 1},
						UserId:   &artistProto.UserID{Id: 2},
						IsLike:   true,
					},
				).Return(nil, status.Error(codes.NotFound, "artist not found"))
			},
			expectedError: customErrors.ErrArtistNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := artistUsecase.LikeArtist(tt.ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFavoriteArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	artistClient := artistProto.ArtistServiceClient(mockArtistClient)
	userClient := userProto.UserServiceClient(mockUserClient)

	artistUsecase := NewUsecase(artistClient, userClient)

	tests := []struct {
		name            string
		filters         *usecase.ArtistFilters
		username        string
		ctx             context.Context
		setupMocks      func()
		expectedArtists []*usecase.Artist
		expectedError   error
	}{
		{
			name: "Success - public favorites, different user",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(2)),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(&userProto.UserID{Id: 1}, nil)

				mockUserClient.EXPECT().GetUserPrivacyByID(
					gomock.Any(),
					&userProto.UserID{Id: 1},
				).Return(&userProto.PrivacySettings{
					IsPublicFavoriteArtists: true,
				}, nil)

				mockArtistClient.EXPECT().GetFavoriteArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Artist 1",
							Description: "Description 1",
							Thumbnail:   "thumbnail1.jpg",
							IsFavorite:  true,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Description: "Description 1",
					Thumbnail:   "thumbnail1.jpg",
					IsLiked:     true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - private favorites, same user",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(1)),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(&userProto.UserID{Id: 1}, nil)

				mockUserClient.EXPECT().GetUserPrivacyByID(
					gomock.Any(),
					&userProto.UserID{Id: 1},
				).Return(&userProto.PrivacySettings{
					IsPublicFavoriteArtists: false,
				}, nil)

				mockArtistClient.EXPECT().GetFavoriteArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Artist 1",
							Description: "Description 1",
							Thumbnail:   "thumbnail1.jpg",
							IsFavorite:  true,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Artist 1",
					Description: "Description 1",
					Thumbnail:   "thumbnail1.jpg",
					IsLiked:     true,
				},
			},
			expectedError: nil,
		},
		{
			name: "Private favorites, different user - empty result",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(2)),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(&userProto.UserID{Id: 1}, nil)

				mockUserClient.EXPECT().GetUserPrivacyByID(
					gomock.Any(),
					&userProto.UserID{Id: 1},
				).Return(&userProto.PrivacySettings{
					IsPublicFavoriteArtists: false,
				}, nil)
			},
			expectedArtists: []*usecase.Artist{},
			expectedError:   nil,
		},
		{
			name: "Error from GetIDByUsername",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.Background(),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedArtists: nil,
			expectedError:   customErrors.ErrUserNotFound,
		},
		{
			name: "Error from GetUserPrivacyByID",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.Background(),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(&userProto.UserID{Id: 1}, nil)

				mockUserClient.EXPECT().GetUserPrivacyByID(
					gomock.Any(),
					&userProto.UserID{Id: 1},
				).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedArtists: nil,
			expectedError:   errors.New("internal server error: internal error"),
		},
		{
			name: "Error from GetFavoriteArtists",
			filters: &usecase.ArtistFilters{
				Pagination: &usecase.Pagination{
					Offset: 0,
					Limit:  10,
				},
			},
			username: "testuser",
			ctx:      context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(1)),
			setupMocks: func() {
				mockUserClient.EXPECT().GetIDByUsername(
					gomock.Any(),
					&userProto.Username{Username: "testuser"},
				).Return(&userProto.UserID{Id: 1}, nil)

				mockUserClient.EXPECT().GetUserPrivacyByID(
					gomock.Any(),
					&userProto.UserID{Id: 1},
				).Return(&userProto.PrivacySettings{
					IsPublicFavoriteArtists: true,
				}, nil)

				mockArtistClient.EXPECT().GetFavoriteArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedArtists: nil,
			expectedError:   errors.New("internal server error: internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			artists, err := artistUsecase.GetFavoriteArtists(tt.ctx, tt.filters, tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
		})
	}
}

func TestSearchArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	artistClient := artistProto.ArtistServiceClient(mockArtistClient)
	userClient := userProto.UserServiceClient(mockUserClient)

	artistUsecase := NewUsecase(artistClient, userClient)

	tests := []struct {
		name            string
		query           string
		ctx             context.Context
		setupMocks      func()
		expectedArtists []*usecase.Artist
		expectedError   error
	}{
		{
			name:  "Success with authenticated user",
			query: "test",
			ctx:   context.WithValue(context.Background(), ctxExtractor.UserContextKey{}, int64(1)),
			setupMocks: func() {
				mockArtistClient.EXPECT().SearchArtists(
					gomock.Any(),
					&artistProto.Query{
						Query:  "test",
						UserId: &artistProto.UserID{Id: 1},
					},
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Test Artist",
							Description: "Test Description",
							Thumbnail:   "test-thumbnail.jpg",
							IsFavorite:  true,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test-thumbnail.jpg",
					IsLiked:     true,
				},
			},
			expectedError: nil,
		},
		{
			name:  "Success with unauthenticated user",
			query: "test",
			ctx:   context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().SearchArtists(
					gomock.Any(),
					&artistProto.Query{
						Query:  "test",
						UserId: &artistProto.UserID{Id: -1},
					},
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{
						{
							Id:          1,
							Title:       "Test Artist",
							Description: "Test Description",
							Thumbnail:   "test-thumbnail.jpg",
							IsFavorite:  false,
						},
					},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{
				{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test-thumbnail.jpg",
					IsLiked:     false,
				},
			},
			expectedError: nil,
		},
		{
			name:  "Error from service",
			query: "test",
			ctx:   context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().SearchArtists(
					gomock.Any(),
					&artistProto.Query{
						Query:  "test",
						UserId: &artistProto.UserID{Id: -1},
					},
				).Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedArtists: nil,
			expectedError:   errors.New("internal server error: internal error"),
		},
		{
			name:  "No results",
			query: "nonexistent",
			ctx:   context.Background(),
			setupMocks: func() {
				mockArtistClient.EXPECT().SearchArtists(
					gomock.Any(),
					&artistProto.Query{
						Query:  "nonexistent",
						UserId: &artistProto.UserID{Id: -1},
					},
				).Return(&artistProto.ArtistList{
					Artists: []*artistProto.Artist{},
				}, nil)
			},
			expectedArtists: []*usecase.Artist{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			artists, err := artistUsecase.SearchArtists(tt.ctx, tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArtists, artists)
			}
		})
	}
}
