package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	mock_track "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type APISuccessResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

type APIErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func setupTestLogger(req *http.Request) *http.Request {
	log := zap.NewNop()
	ctx := req.Context()
	ctx = loggerPkg.LoggerToContext(ctx, log.Sugar())
	return req.WithContext(ctx)
}

func verifyResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedStatus, rec.Code)

	expectedStatusStr := expectedBody["status"].(string)

	if expectedStatusStr == "success" {
		var response APISuccessResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, expectedStatus, response.Status)

		if expectedBody["body"] != nil {
			assert.NotNil(t, response.Body, "Response body should not be nil")

			switch expectedBodyValue := expectedBody["body"].(type) {
			case []map[string]interface{}:
				responseBodyArray, ok := response.Body.([]interface{})
				assert.True(t, ok, "Expected response body to be an array")

				if len(expectedBodyValue) > 0 && len(responseBodyArray) > 0 {
					expectedTrack := expectedBodyValue[0]
					responseTrack, ok := responseBodyArray[0].(map[string]interface{})
					assert.True(t, ok, "Expected first item to be a map")

					assert.Equal(t, expectedTrack["id"], responseTrack["id"])
					assert.Equal(t, expectedTrack["title"], responseTrack["title"])
				}
			case map[string]interface{}:
				responseBodyMap, ok := response.Body.(map[string]interface{})
				assert.True(t, ok, "Expected response body to be an object")

				for key, expectedValue := range expectedBodyValue {
					assert.Equal(t, expectedValue, responseBodyMap[key])
				}
			}
		}
	} else if expectedStatusStr == "error" {
		var response APIErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, expectedStatus, response.Status)

		if expectedBody["error"] != nil {
			expectedError := expectedBody["error"].(map[string]interface{})
			assert.Equal(t, expectedError["message"], response.Error)
		}
	}
}

func TestTrackHandler_GetAllTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "OK",
			query: "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAllTracks(gomock.Any(), gomock.Any()).Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:  "Invalid Pagination",
			query: "?offset=-1",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "invalid offset: should be greater than 0",
				},
			},
		},
		{
			name:  "Internal Error",
			query: "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAllTracks(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "internal error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, err := http.NewRequest("GET", "/tracks"+tt.query, nil)
			assert.NoError(t, err)

			// Add logger to context
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetAllTracks(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_GetTrackByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		trackID        string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "OK",
			trackID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetTrackByID(gomock.Any(), int64(1)).Return(&usecaseModel.TrackDetailed{
					Track: usecaseModel.Track{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
					FileUrl: "test.jpg",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"id":       float64(1),
					"title":    "Test Track",
					"duration": float64(180),
				},
			},
		},
		{
			name:    "Invalid ID",
			trackID: "invalid",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:    "Not Found",
			trackID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetTrackByID(gomock.Any(), int64(999)).Return(nil, customErrors.ErrTrackNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrTrackNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/tracks/"+tt.trackID, nil)
			vars := map[string]string{
				"id": tt.trackID,
			}
			req = mux.SetURLVars(req, vars)

			// Add logger to context
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetTrackByID(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_GetTracksByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		artistID       string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "OK",
			artistID: "1",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetTracksByArtistID(gomock.Any(), int64(1), gomock.Any()).Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
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
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:     "Invalid Pagination",
			artistID: "1",
			query:    "?offset=-1",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "invalid offset: should be greater than 0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/artists/"+tt.artistID+"/tracks"+tt.query, nil)
			vars := map[string]string{
				"id": tt.artistID,
			}
			req = mux.SetURLVars(req, vars)

			// Add logger to context
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetTracksByArtistID(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_CreateStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		trackID        string
		userID         int64
		authenticated  bool
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			trackID:       "1",
			userID:        1,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateStream(gomock.Any(), gomock.Any()).Return(int64(123), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"id": float64(123),
				},
			},
		},
		{
			name:           "Invalid Track ID",
			trackID:        "invalid",
			userID:         1,
			authenticated:  true,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:           "Unauthorized",
			trackID:        "1",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUnauthorized.Error(),
				},
			},
		},
		{
			name:          "Internal Error",
			trackID:       "1",
			userID:        1,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateStream(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "internal error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("POST", "/tracks/"+tt.trackID+"/stream", nil)
			vars := map[string]string{
				"id": tt.trackID,
			}
			req = mux.SetURLVars(req, vars)

			ctx := req.Context()
			if tt.authenticated {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}

			ctx = loggerPkg.LoggerToContext(ctx, zap.NewNop().Sugar())
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.CreateStream(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_UpdateStreamDuration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		streamID       string
		userID         int64
		authenticated  bool
		body           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			streamID:      "1",
			userID:        1,
			authenticated: true,
			body:          `{"duration": 180}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().UpdateStreamDuration(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"msg": "stream duration was successfully updated",
				},
			},
		},
		{
			name:           "Invalid Stream ID",
			streamID:       "invalid",
			userID:         1,
			authenticated:  true,
			body:           `{"duration": 180}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:           "Unauthorized",
			streamID:       "1",
			authenticated:  false,
			body:           `{"duration": 180}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Request Body",
			streamID:       "1",
			userID:         1,
			authenticated:  true,
			body:           `{"duration": -1}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "duration: -1 does not validate as range(0|999999999)",
				},
			},
		},
		{
			name:          "Internal Error",
			streamID:      "1",
			userID:        1,
			authenticated: true,
			body:          `{"duration": 180}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().UpdateStreamDuration(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "internal error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("PUT", "/streams/"+tt.streamID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			vars := map[string]string{
				"id": tt.streamID,
			}
			req = mux.SetURLVars(req, vars)

			ctx := req.Context()
			if tt.authenticated {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}

			ctx = loggerPkg.LoggerToContext(ctx, zap.NewNop().Sugar())
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.UpdateStreamDuration(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_GetLastListenedTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		query          string
		userID         int64
		authenticated  bool
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			query:         "?offset=0&limit=10",
			userID:        1,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().GetLastListenedTracks(gomock.Any(), int64(1), gomock.Any()).Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:           "Unauthorized",
			query:          "?offset=0&limit=10",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Pagination",
			query:          "?offset=-1",
			userID:         1,
			authenticated:  true,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "invalid offset: should be greater than 0",
				},
			},
		},
		{
			name:          "Not Found",
			query:         "?offset=0&limit=10",
			userID:        999,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().GetLastListenedTracks(gomock.Any(), int64(999), gomock.Any()).Return(nil, customErrors.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUserNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/users/me/history"+tt.query, nil)

			ctx := req.Context()
			if tt.authenticated {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}

			ctx = loggerPkg.LoggerToContext(ctx, zap.NewNop().Sugar())
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.GetLastListenedTracks(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestTrackHandler_LikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		trackID        string
		userID         int64
		authenticated  bool
		body           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:          "OK - Like",
			trackID:       "1",
			userID:        1,
			authenticated: true,
			body:          `{"is_like": true}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeTrack(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"msg": "track liked/unliked",
				},
			},
		},
		{
			name:          "OK - Unlike",
			trackID:       "1",
			userID:        1,
			authenticated: true,
			body:          `{"is_like": false}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeTrack(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"msg": "track liked/unliked",
				},
			},
		},
		{
			name:           "Invalid Track ID",
			trackID:        "invalid",
			userID:         1,
			authenticated:  true,
			body:           `{"is_like": true}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:           "Unauthorized",
			trackID:        "1",
			authenticated:  false,
			body:           `{"is_like": true}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Request Body",
			trackID:        "1",
			userID:         1,
			authenticated:  true,
			body:           `{invalid json}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "failed to unmarshal JSON: parse error: syntax error near offset 1 of '{invalid json}'",
				},
			},
		},
		{
			name:          "Track Not Found",
			trackID:       "999",
			userID:        1,
			authenticated: true,
			body:          `{"is_like": true}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().LikeTrack(gomock.Any(), gomock.Any()).Return(customErrors.ErrTrackNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrTrackNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("POST", "/tracks/"+tt.trackID+"/like", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			vars := map[string]string{
				"id": tt.trackID,
			}
			req = mux.SetURLVars(req, vars)

			ctx := req.Context()
			if tt.authenticated {
				ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, tt.userID)
			}

			ctx = loggerPkg.LoggerToContext(ctx, zap.NewNop().Sugar())
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.LikeTrack(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody.(map[string]interface{}))
		})
	}
}

func TestTrackHandler_GetPlaylistTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		playlistID     string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:       "OK",
			playlistID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetPlaylistTracks(gomock.Any(), int64(1)).Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:           "Invalid Playlist ID",
			playlistID:     "invalid",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:       "Playlist Not Found",
			playlistID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetPlaylistTracks(gomock.Any(), int64(999)).Return(nil, customErrors.ErrPlaylistNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/playlists/"+tt.playlistID+"/tracks", nil)
			vars := map[string]string{
				"id": tt.playlistID,
			}
			req = mux.SetURLVars(req, vars)

			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetPlaylistTracks(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody.(map[string]interface{}))
		})
	}
}

func TestTrackHandler_GetFavoriteTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		username       string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "OK",
			username: "testuser",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteTracks(gomock.Any(), gomock.Any(), "testuser").Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:           "Empty Username",
			username:       "",
			query:          "?offset=0&limit=10",
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Pagination",
			username:       "testuser",
			query:          "?offset=-1",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "invalid offset: should be greater than 0",
				},
			},
		},
		{
			name:     "User Not Found",
			username: "nonexistent",
			query:    "?offset=0&limit=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetFavoriteTracks(gomock.Any(), gomock.Any(), "nonexistent").Return(nil, customErrors.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrUserNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/users/"+tt.username+"/tracks"+tt.query, nil)
			vars := map[string]string{
				"username": tt.username,
			}
			req = mux.SetURLVars(req, vars)

			// Add logger to context
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetFavoriteTracks(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody.(map[string]interface{}))
		})
	}
}

func TestTrackHandler_SearchTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "OK",
			query: "?query=test",
			mockBehavior: func() {
				mockUsecase.EXPECT().SearchTracks(gomock.Any(), "test").Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:           "Empty Query",
			query:          "",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "query is empty",
				},
			},
		},
		{
			name:  "No Results",
			query: "?query=nonexistent",
			mockBehavior: func() {
				mockUsecase.EXPECT().SearchTracks(gomock.Any(), "nonexistent").Return([]*usecaseModel.Track{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body":   []interface{}{},
			},
		},
		{
			name:  "Internal Error",
			query: "?query=test",
			mockBehavior: func() {
				mockUsecase.EXPECT().SearchTracks(gomock.Any(), "test").Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "internal error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/tracks/search"+tt.query, nil)

			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.SearchTracks(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody.(map[string]interface{}))
		})
	}
}

func TestTrackHandler_GetTracksByAlbumID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_track.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewTrackHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		albumID        string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "OK",
			albumID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetTracksByAlbumID(gomock.Any(), int64(1)).Return([]*usecaseModel.Track{
					{
						ID:       1,
						Title:    "Test Track",
						Duration: 180,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []map[string]interface{}{
					{
						"id":       float64(1),
						"title":    "Test Track",
						"duration": float64(180),
					},
				},
			},
		},
		{
			name:           "Invalid Album ID",
			albumID:        "invalid",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
				},
			},
		},
		{
			name:    "Album Not Found",
			albumID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetTracksByAlbumID(gomock.Any(), int64(999)).Return(nil, customErrors.ErrAlbumNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrAlbumNotFound.Error(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/albums/"+tt.albumID+"/tracks", nil)
			vars := map[string]string{
				"id": tt.albumID,
			}
			req = mux.SetURLVars(req, vars)

			// Add logger to context
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetTracksByAlbumID(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}
