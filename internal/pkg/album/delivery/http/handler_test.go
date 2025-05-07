package album

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	mock_album "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album/mocks"
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

func setupTestHandler(t *testing.T) (*mock_album.MockUsecase, *AlbumHandler, *config.Config) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_album.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}
	handler := NewAlbumHandler(mockUsecase, cfg)
	return mockUsecase, handler, cfg
}

func TestGetAllAlbums(t *testing.T) {
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
				mockUsecase.EXPECT().GetAllAlbums(
					gomock.Any(),
					gomock.Any(),
				).Return([]*usecaseModel.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Type:        usecaseModel.AlbumTypeAlbum,
						Thumbnail:   "test.jpg",
						ReleaseDate: time.Now(),
						Artists: []*usecaseModel.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
						IsLiked: true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Album{
				{
					ID:          1,
					Title:       "Test Album",
					Type:        deliveryModel.AlbumTypeAlbum,
					Thumbnail:   "test.jpg",
					ReleaseDate: time.Now(),
					Artists: []*deliveryModel.AlbumArtist{
						{
							ID:    1,
							Title: "Test Artist",
						},
					},
					IsLiked: true,
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
				mockUsecase.EXPECT().GetAllAlbums(
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
			req, _ := http.NewRequest(http.MethodGet, "/albums"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			handler.GetAllAlbums(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var albums []*deliveryModel.Album
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &albums)
				assert.NoError(t, err)

				expectedAlbums := tt.expectedBody.([]*deliveryModel.Album)
				assert.Equal(t, len(expectedAlbums), len(albums))
				assert.Equal(t, expectedAlbums[0].ID, albums[0].ID)
				assert.Equal(t, expectedAlbums[0].Title, albums[0].Title)
				assert.Equal(t, expectedAlbums[0].Type, albums[0].Type)
				assert.Equal(t, expectedAlbums[0].Thumbnail, albums[0].Thumbnail)
				assert.Equal(t, expectedAlbums[0].IsLiked, albums[0].IsLiked)
				assert.Equal(t, len(expectedAlbums[0].Artists), len(albums[0].Artists))
				assert.Equal(t, expectedAlbums[0].Artists[0].ID, albums[0].Artists[0].ID)
				assert.Equal(t, expectedAlbums[0].Artists[0].Title, albums[0].Artists[0].Title)
			}
		})
	}
}

func TestGetAlbumsByArtistID(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		artistID       string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "Success",
			artistID: "1",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAlbumsByArtistID(
					gomock.Any(),
					int64(1),
					gomock.Any(),
				).Return([]*usecaseModel.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Type:        usecaseModel.AlbumTypeAlbum,
						Thumbnail:   "test.jpg",
						ReleaseDate: time.Now(),
						Artists: []*usecaseModel.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
						IsLiked: true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Album{
				{
					ID:          1,
					Title:       "Test Album",
					Type:        deliveryModel.AlbumTypeAlbum,
					Thumbnail:   "test.jpg",
					ReleaseDate: time.Now(),
					Artists: []*deliveryModel.AlbumArtist{
						{
							ID:    1,
							Title: "Test Artist",
						},
					},
					IsLiked: true,
				},
			},
		},
		{
			name:     "Invalid Artist ID",
			artistID: "invalid",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:     "Invalid Pagination",
			artistID: "1",
			query:    "?offset=invalid&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:     "Usecase Error",
			artistID: "1",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAlbumsByArtistID(
					gomock.Any(),
					int64(1),
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
			req, _ := http.NewRequest(http.MethodGet, "/artists/"+tt.artistID+"/albums"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			vars := map[string]string{
				"id": tt.artistID,
			}
			req = mux.SetURLVars(req, vars)

			handler.GetAlbumsByArtistID(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var albums []*deliveryModel.Album
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &albums)
				assert.NoError(t, err)

				expectedAlbums := tt.expectedBody.([]*deliveryModel.Album)
				assert.Equal(t, len(expectedAlbums), len(albums))
				assert.Equal(t, expectedAlbums[0].ID, albums[0].ID)
				assert.Equal(t, expectedAlbums[0].Title, albums[0].Title)
				assert.Equal(t, expectedAlbums[0].Type, albums[0].Type)
				assert.Equal(t, expectedAlbums[0].Thumbnail, albums[0].Thumbnail)
				assert.Equal(t, expectedAlbums[0].IsLiked, albums[0].IsLiked)
			}
		})
	}
}

func TestGetAlbumByID(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		albumID        string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:    "Success",
			albumID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAlbumByID(
					gomock.Any(),
					int64(1),
				).Return(&usecaseModel.Album{
					ID:          1,
					Title:       "Test Album",
					Type:        usecaseModel.AlbumTypeAlbum,
					Thumbnail:   "test.jpg",
					ReleaseDate: time.Now(),
					Artists: []*usecaseModel.AlbumArtist{
						{
							ID:    1,
							Title: "Test Artist",
						},
					},
					IsLiked: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &deliveryModel.Album{
				ID:          1,
				Title:       "Test Album",
				Type:        deliveryModel.AlbumTypeAlbum,
				Thumbnail:   "test.jpg",
				ReleaseDate: time.Now(),
				Artists: []*deliveryModel.AlbumArtist{
					{
						ID:    1,
						Title: "Test Artist",
					},
				},
				IsLiked: true,
			},
		},
		{
			name:    "Invalid Album ID",
			albumID: "invalid",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:    "Not Found",
			albumID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAlbumByID(
					gomock.Any(),
					int64(999),
				).Return(nil, customErrors.ErrAlbumNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/albums/"+tt.albumID, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			vars := map[string]string{
				"id": tt.albumID,
			}
			req = mux.SetURLVars(req, vars)

			handler.GetAlbumByID(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var album deliveryModel.Album
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &album)
				assert.NoError(t, err)

				expectedAlbum := tt.expectedBody.(*deliveryModel.Album)
				assert.Equal(t, expectedAlbum.ID, album.ID)
				assert.Equal(t, expectedAlbum.Title, album.Title)
				assert.Equal(t, expectedAlbum.Type, album.Type)
				assert.Equal(t, expectedAlbum.Thumbnail, album.Thumbnail)
				assert.Equal(t, expectedAlbum.IsLiked, album.IsLiked)
				assert.Equal(t, len(expectedAlbum.Artists), len(album.Artists))
				assert.Equal(t, expectedAlbum.Artists[0].ID, album.Artists[0].ID)
				assert.Equal(t, expectedAlbum.Artists[0].Title, album.Artists[0].Title)
			}
		})
	}
}

func TestLikeAlbum(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		albumID        string
		userID         int64
		requestBody    map[string]interface{}
		mockBehavior   func()
		expectedStatus int
	}{
		{
			name:    "Success - Like",
			albumID: "1",
			userID:  42,
			requestBody: map[string]interface{}{
				"value": true,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeAlbum(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Success - Unlike",
			albumID: "1",
			userID:  42,
			requestBody: map[string]interface{}{
				"value": false,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeAlbum(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Unauthorized",
			albumID: "1",
			userID:  0,
			requestBody: map[string]interface{}{
				"value": true,
			},
			mockBehavior: func() {
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:    "Invalid Album ID",
			albumID: "invalid",
			userID:  42,
			requestBody: map[string]interface{}{
				"value": true,
			},
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request Body",
			albumID:     "1",
			userID:      42,
			requestBody: nil,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Usecase Error",
			albumID: "1",
			userID:  42,
			requestBody: map[string]interface{}{
				"value": true,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeAlbum(
					gomock.Any(),
					gomock.Any(),
				).Return(customErrors.ErrAlbumNotFound)
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
			req, _ := http.NewRequest(http.MethodPost, "/albums/"+tt.albumID+"/like", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}
			req = req.WithContext(ctx)

			vars := map[string]string{
				"id": tt.albumID,
			}
			req = mux.SetURLVars(req, vars)

			handler.LikeAlbum(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetFavoriteAlbums(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	tests := []struct {
		name           string
		userID         int64
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			userID: 42,
			query:  "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteAlbums(
					gomock.Any(),
					gomock.Any(),
					int64(42),
				).Return([]*usecaseModel.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Type:        usecaseModel.AlbumTypeAlbum,
						Thumbnail:   "test.jpg",
						ReleaseDate: time.Now(),
						Artists: []*usecaseModel.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
						IsLiked: true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Album{
				{
					ID:          1,
					Title:       "Test Album",
					Type:        deliveryModel.AlbumTypeAlbum,
					Thumbnail:   "test.jpg",
					ReleaseDate: time.Now(),
					Artists: []*deliveryModel.AlbumArtist{
						{
							ID:    1,
							Title: "Test Artist",
						},
					},
					IsLiked: true,
				},
			},
		},
		{
			name:   "Unauthorized",
			userID: 0,
			query:  "?offset=0&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   nil,
		},
		{
			name:   "Invalid Pagination",
			userID: 42,
			query:  "?offset=invalid&limit=10",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:   "Usecase Error",
			userID: 42,
			query:  "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteAlbums(
					gomock.Any(),
					gomock.Any(),
					int64(42),
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
			req, _ := http.NewRequest(http.MethodGet, "/user/me/albums"+tt.query, nil)

			ctx := context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar())
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}
			req = req.WithContext(ctx)

			handler.GetFavoriteAlbums(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var albums []*deliveryModel.Album
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &albums)
				assert.NoError(t, err)

				expectedAlbums := tt.expectedBody.([]*deliveryModel.Album)
				assert.Equal(t, len(expectedAlbums), len(albums))
				assert.Equal(t, expectedAlbums[0].ID, albums[0].ID)
				assert.Equal(t, expectedAlbums[0].Title, albums[0].Title)
				assert.Equal(t, expectedAlbums[0].Type, albums[0].Type)
				assert.Equal(t, expectedAlbums[0].Thumbnail, albums[0].Thumbnail)
				assert.Equal(t, expectedAlbums[0].IsLiked, albums[0].IsLiked)
			}
		})
	}
}

func TestSearchAlbums(t *testing.T) {
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
				mockUsecase.EXPECT().SearchAlbums(
					gomock.Any(),
					"test",
				).Return([]*usecaseModel.Album{
					{
						ID:          1,
						Title:       "Test Album",
						Type:        usecaseModel.AlbumTypeAlbum,
						Thumbnail:   "test.jpg",
						ReleaseDate: time.Now(),
						Artists: []*usecaseModel.AlbumArtist{
							{
								ID:    1,
								Title: "Test Artist",
							},
						},
						IsLiked: true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*deliveryModel.Album{
				{
					ID:          1,
					Title:       "Test Album",
					Type:        deliveryModel.AlbumTypeAlbum,
					Thumbnail:   "test.jpg",
					ReleaseDate: time.Now(),
					Artists: []*deliveryModel.AlbumArtist{
						{
							ID:    1,
							Title: "Test Artist",
						},
					},
					IsLiked: true,
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
				mockUsecase.EXPECT().SearchAlbums(
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
			req, _ := http.NewRequest(http.MethodGet, "/albums/search"+tt.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), loggerPkg.LoggerKey{}, zap.NewNop().Sugar()))

			handler.SearchAlbums(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && tt.expectedBody != nil {
				var response deliveryModel.APIResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)

				var albums []*deliveryModel.Album
				respBodyBytes, err := json.Marshal(response.Body)
				assert.NoError(t, err)

				err = json.Unmarshal(respBodyBytes, &albums)
				assert.NoError(t, err)

				expectedAlbums := tt.expectedBody.([]*deliveryModel.Album)
				assert.Equal(t, len(expectedAlbums), len(albums))
				assert.Equal(t, expectedAlbums[0].ID, albums[0].ID)
				assert.Equal(t, expectedAlbums[0].Title, albums[0].Title)
				assert.Equal(t, expectedAlbums[0].Type, albums[0].Type)
				assert.Equal(t, expectedAlbums[0].Thumbnail, albums[0].Thumbnail)
				assert.Equal(t, expectedAlbums[0].IsLiked, albums[0].IsLiked)
			}
		})
	}
}
