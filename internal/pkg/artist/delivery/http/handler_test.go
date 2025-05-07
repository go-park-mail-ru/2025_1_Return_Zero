package artist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	mock_artist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTestHandler(t *testing.T) (*mock_artist.MockUsecase, *ArtistHandler, *config.Config) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_artist.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}
	handler := NewArtistHandler(mockUsecase, cfg)
	return mockUsecase, handler, cfg
}

func TestGetAllArtists(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "Success",
			query: "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAllArtists(
					gomock.Any(),
					gomock.Any(),
				).Return([]*usecaseModel.Artist{
					{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test.jpg",
						IsLiked:     true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Artist{
				{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test.jpg",
					IsLiked:     true,
				},
			},
		},
		{
			name:  "Invalid Pagination",
			query: "?offset=invalid&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:  "Usecase Error",
			query: "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAllArtists(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/artists"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			handler.GetAllArtists(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var artists []*deliveryModel.Artist
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &artists)
				assert.NoError(t, err)

				expectedArtists := tt.expectedBody.([]*deliveryModel.Artist)
				assert.Equal(t, len(expectedArtists), len(artists))
				assert.Equal(t, expectedArtists[0].ID, artists[0].ID)
				assert.Equal(t, expectedArtists[0].Title, artists[0].Title)
				assert.Equal(t, expectedArtists[0].Description, artists[0].Description)
				assert.Equal(t, expectedArtists[0].Thumbnail, artists[0].Thumbnail)
				assert.Equal(t, expectedArtists[0].IsLiked, artists[0].IsLiked)
			}
		})
	}
}

func TestGetArtistByID(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		artistID       string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "Success",
			artistID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetArtistByID(
					gomock.Any(),
					int64(1),
				).Return(&usecaseModel.ArtistDetailed{
					Artist: usecaseModel.Artist{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test.jpg",
						IsLiked:     true,
					},
					Favorites: 100,
					Listeners: 500,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &deliveryModel.ArtistDetailed{
				Artist: deliveryModel.Artist{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test.jpg",
					IsLiked:     true,
				},
				Favorites: 100,
				Listeners: 500,
			},
		},
		{
			name:     "Invalid ID",
			artistID: "invalid",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:     "Not Found",
			artistID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetArtistByID(
					gomock.Any(),
					int64(999),
				).Return(nil, customErrors.ErrArtistNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/artists/"+tt.artistID, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			vars := map[string]string{
				"id": tt.artistID,
			}
			req = mux.SetURLVars(req, vars)

			handler.GetArtistByID(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var artist deliveryModel.ArtistDetailed
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &artist)
				assert.NoError(t, err)

				expectedArtist := tt.expectedBody.(*deliveryModel.ArtistDetailed)
				assert.Equal(t, expectedArtist.Artist.ID, artist.Artist.ID)
				assert.Equal(t, expectedArtist.Artist.Title, artist.Artist.Title)
				assert.Equal(t, expectedArtist.Artist.Description, artist.Artist.Description)
				assert.Equal(t, expectedArtist.Artist.Thumbnail, artist.Artist.Thumbnail)
				assert.Equal(t, expectedArtist.Artist.IsLiked, artist.Artist.IsLiked)
				assert.Equal(t, expectedArtist.Favorites, artist.Favorites)
				assert.Equal(t, expectedArtist.Listeners, artist.Listeners)
			}
		})
	}
}

func TestLikeArtist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		artistID       string
		userID         int64
		requestBody    map[string]interface{}
		mockBehavior   func()
		expectedStatus int
	}{
		{
			name:     "Success - Like",
			artistID: "1",
			userID:   42,
			requestBody: map[string]interface{}{
				"is_like": true,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeArtist(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Success - Unlike",
			artistID: "1",
			userID:   42,
			requestBody: map[string]interface{}{
				"is_like": false,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeArtist(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Unauthorized",
			artistID: "1",
			userID:   0,
			requestBody: map[string]interface{}{
				"is_like": true,
			},
			mockBehavior: func() {
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:     "Invalid Artist ID",
			artistID: "invalid",
			userID:   42,
			requestBody: map[string]interface{}{
				"is_like": true,
			},
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request Body",
			artistID:    "1",
			userID:      42,
			requestBody: nil,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Usecase Error",
			artistID: "1",
			userID:   42,
			requestBody: map[string]interface{}{
				"is_like": true,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeArtist(
					gomock.Any(),
					gomock.Any(),
				).Return(customErrors.ErrArtistNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/artists/"+tt.artistID+"/like", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}
			req = req.WithContext(ctx)

			vars := map[string]string{
				"id": tt.artistID,
			}
			req = mux.SetURLVars(req, vars)

			handler.LikeArtist(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetFavoriteArtists(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		username       string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "Success",
			username: "testuser",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteArtists(
					gomock.Any(),
					gomock.Any(),
					"testuser",
				).Return([]*usecaseModel.Artist{
					{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test.jpg",
						IsLiked:     true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Artist{
				{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test.jpg",
					IsLiked:     true,
				},
			},
		},
		{
			name:     "Empty Username",
			username: "",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   nil,
		},
		{
			name:     "Invalid Pagination",
			username: "testuser",
			query:    "?offset=invalid&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:     "User Not Found",
			username: "nonexistent",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteArtists(
					gomock.Any(),
					gomock.Any(),
					"nonexistent",
				).Return(nil, customErrors.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/user/"+tt.username+"/artists"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			vars := map[string]string{
				"username": tt.username,
			}
			req = mux.SetURLVars(req, vars)

			handler.GetFavoriteArtists(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var artists []*deliveryModel.Artist
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &artists)
				assert.NoError(t, err)

				expectedArtists := tt.expectedBody.([]*deliveryModel.Artist)
				assert.Equal(t, len(expectedArtists), len(artists))
				assert.Equal(t, expectedArtists[0].ID, artists[0].ID)
				assert.Equal(t, expectedArtists[0].Title, artists[0].Title)
				assert.Equal(t, expectedArtists[0].Description, artists[0].Description)
				assert.Equal(t, expectedArtists[0].Thumbnail, artists[0].Thumbnail)
				assert.Equal(t, expectedArtists[0].IsLiked, artists[0].IsLiked)
			}
		})
	}
}

func TestSearchArtists(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "Success",
			query: "?query=test",
			mockBehavior: func() {
				mockUsecase.EXPECT().SearchArtists(
					gomock.Any(),
					"test",
				).Return([]*usecaseModel.Artist{
					{
						ID:          1,
						Title:       "Test Artist",
						Description: "Test Description",
						Thumbnail:   "test.jpg",
						IsLiked:     true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Artist{
				{
					ID:          1,
					Title:       "Test Artist",
					Description: "Test Description",
					Thumbnail:   "test.jpg",
					IsLiked:     true,
				},
			},
		},
		{
			name:  "Empty Query",
			query: "",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:  "Usecase Error",
			query: "?query=test",
			mockBehavior: func() {
				mockUsecase.EXPECT().SearchArtists(
					gomock.Any(),
					"test",
				).Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/artists/search"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			handler.SearchArtists(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var artists []*deliveryModel.Artist
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &artists)
				assert.NoError(t, err)

				expectedArtists := tt.expectedBody.([]*deliveryModel.Artist)
				assert.Equal(t, len(expectedArtists), len(artists))
				assert.Equal(t, expectedArtists[0].ID, artists[0].ID)
				assert.Equal(t, expectedArtists[0].Title, artists[0].Title)
				assert.Equal(t, expectedArtists[0].Description, artists[0].Description)
				assert.Equal(t, expectedArtists[0].Thumbnail, artists[0].Thumbnail)
				assert.Equal(t, expectedArtists[0].IsLiked, artists[0].IsLiked)
			}
		})
	}
}
