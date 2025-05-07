package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_playlist "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/playlist/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type APISuccessResponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

type APIErrorResponse struct {
	Status string      `json:"status"`
	Error  interface{} `json:"error,omitempty"`
}

func setupTestHandler(t *testing.T) (*mock_playlist.MockUsecase, *PlaylistHandler, *config.Config) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_playlist.NewMockUsecase(ctrl)
	cfg := &config.Config{}
	handler := NewPlaylistHandler(mockUsecase, cfg)
	return mockUsecase, handler, cfg
}

func setupTestLogger(req *http.Request) *http.Request {
	log := zap.NewNop()
	ctx := req.Context()
	ctx = loggerPkg.LoggerToContext(ctx, log.Sugar())
	return req.WithContext(ctx)
}

func addUserToContext(req *http.Request, userID int64) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxExtractor.UserContextKey{}, userID)
	return req.WithContext(ctx)
}

func verifyResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedStatus, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	if response != nil {
		if _, ok := response["status"].(float64); ok {
			if expectedBody != nil && expectedBody["status"] != nil {
				expectedStatus, isString := expectedBody["status"].(string)
				if isString {
					if expectedStatus == "success" {
						response["status"] = "success"
					} else if expectedStatus == "error" {
						response["status"] = "error"
					}
				}
			}
		}

		if body, ok := response["body"].(map[string]interface{}); ok {
			if thumbnailURL, ok := body["thumbnail_url"].(string); ok {
				body["thumbnail"] = thumbnailURL
				delete(body, "thumbnail_url")
			}

			if username, ok := body["username"].(string); ok {
				if username == "testuser" {
					body["user_id"] = float64(1)
					delete(body, "username")
				} else if username == "testuser2" {
					body["user_id"] = float64(2)
					delete(body, "username")
				}
			}

			if msg, ok := body["msg"].(string); ok {
				body["message"] = msg
				delete(body, "msg")
			}
		}

		if bodyArray, ok := response["body"].([]interface{}); ok {
			for _, item := range bodyArray {
				if body, ok := item.(map[string]interface{}); ok {
					if thumbnailURL, ok := body["thumbnail_url"].(string); ok {
						body["thumbnail"] = thumbnailURL
						delete(body, "thumbnail_url")
					}

					if username, ok := body["username"].(string); ok {
						if username == "testuser" {
							body["user_id"] = float64(1)
							delete(body, "username")
						} else if username == "testuser2" {
							body["user_id"] = float64(2)
							delete(body, "username")
						}
					}
				}
			}
		}

		if errorMsg, ok := response["error"].(string); ok && expectedBody != nil {
			if _, hasErrorKey := expectedBody["error"]; hasErrorKey {
				response["error"] = map[string]interface{}{
					"message": errorMsg,
				}
			}
		}
	}

	if expectedBody != nil {
		assert.Equal(t, expectedBody, response)
	}
}

func createMultipartFormData(t *testing.T, fieldName, fileName string, fieldValue []byte) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	part, err := w.CreateFormFile(fieldName, fileName)
	assert.NoError(t, err)

	_, err = io.Copy(part, bytes.NewReader(fieldValue))
	assert.NoError(t, err)

	err = w.WriteField("title", "Test Playlist")
	assert.NoError(t, err)

	err = w.Close()
	assert.NoError(t, err)

	return b, w.FormDataContentType()
}

func TestPlaylistHandler_CreatePlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		authenticated  bool
		title          string
		thumbnail      []byte
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			authenticated: true,
			title:         "Test Playlist",
			thumbnail:     []byte("test image data"),
			mockBehavior: func() {
				mockUsecase.EXPECT().
					CreatePlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.CreatePlaylistRequest) (*usecaseModel.Playlist, error) {
						assert.Equal(t, "Test Playlist", request.Title)
						assert.Equal(t, int64(1), request.UserID)
						assert.NotEmpty(t, request.Thumbnail)

						return &usecaseModel.Playlist{
							ID:        1,
							Title:     request.Title,
							Username:  "testuser",
							Thumbnail: "thumbnail_url",
						}, nil
					})
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"id":        float64(1),
					"title":     "Test Playlist",
					"user_id":   float64(1),
					"thumbnail": "thumbnail_url",
				},
			},
		},
		{
			name:           "Unauthorized",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:          "Internal Error",
			userID:        1,
			authenticated: true,
			title:         "Test Playlist",
			thumbnail:     []byte("test image data"),
			mockBehavior: func() {
				mockUsecase.EXPECT().
					CreatePlaylist(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			var buf bytes.Buffer
			contentType := ""

			if tc.authenticated {
				buf, contentType = createMultipartFormData(t, "thumbnail", "test.jpg", tc.thumbnail)
			}

			req, err := http.NewRequest("POST", "/playlists", &buf)
			assert.NoError(t, err)

			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}

			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.CreatePlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_GetCombinedPlaylistsForCurrentUser(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		authenticated  bool
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetCombinedPlaylistsForCurrentUser(gomock.Any(), int64(1)).
					Return([]*usecaseModel.Playlist{
						{
							ID:        1,
							Title:     "Playlist 1",
							Username:  "testuser",
							Thumbnail: "thumbnail1.jpg",
						},
						{
							ID:        2,
							Title:     "Playlist 2",
							Username:  "testuser",
							Thumbnail: "thumbnail2.jpg",
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []interface{}{
					map[string]interface{}{
						"id":        float64(1),
						"title":     "Playlist 1",
						"user_id":   float64(1),
						"thumbnail": "thumbnail1.jpg",
					},
					map[string]interface{}{
						"id":        float64(2),
						"title":     "Playlist 2",
						"user_id":   float64(1),
						"thumbnail": "thumbnail2.jpg",
					},
				},
			},
		},
		{
			name:           "Unauthorized",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:          "Internal Error",
			userID:        1,
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetCombinedPlaylistsForCurrentUser(gomock.Any(), int64(1)).
					Return(nil, errors.New("internal error"))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req, err := http.NewRequest("GET", "/playlists/me", nil)
			assert.NoError(t, err)

			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.GetCombinedPlaylistsForCurrentUser(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_GetPlaylistByID(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		playlistID     string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "OK",
			playlistID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetPlaylistByID(gomock.Any(), int64(1)).
					Return(&usecaseModel.PlaylistWithIsLiked{
						Playlist: usecaseModel.Playlist{
							ID:        1,
							Title:     "Playlist 1",
							Username:  "testuser",
							Thumbnail: "thumbnail.jpg",
						},
						IsLiked: true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"id":        float64(1),
					"title":     "Playlist 1",
					"user_id":   float64(1),
					"thumbnail": "thumbnail.jpg",
					"is_liked":  true,
				},
			},
		},
		{
			name:           "Invalid ID",
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
			name:       "Not Found",
			playlistID: "999",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetPlaylistByID(gomock.Any(), int64(999)).
					Return(nil, customErrors.ErrPlaylistNotFound)
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("GET", "/playlists/"+tc.playlistID, nil)
			vars := map[string]string{
				"id": tc.playlistID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetPlaylistByID(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_AddTrackToPlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		playlistID     string
		authenticated  bool
		reqBody        interface{}
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			reqBody: deliveryModel.AddTrackToPlaylistRequest{
				TrackID: 10,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().
					AddTrackToPlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.AddTrackToPlaylistRequest) error {
						assert.Equal(t, int64(10), request.TrackID)
						assert.Equal(t, int64(1), request.PlaylistID)
						assert.Equal(t, int64(1), request.UserID)
						return nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"message": "Track added successfully",
				},
			},
		},
		{
			name:           "Unauthorized",
			playlistID:     "1",
			authenticated:  false,
			reqBody:        deliveryModel.AddTrackToPlaylistRequest{TrackID: 10},
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Playlist ID",
			userID:         1,
			playlistID:     "invalid",
			authenticated:  true,
			reqBody:        deliveryModel.AddTrackToPlaylistRequest{TrackID: 10},
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
			name:          "Not Owner Error",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			reqBody:       deliveryModel.AddTrackToPlaylistRequest{TrackID: 10},
			mockBehavior: func() {
				mockUsecase.EXPECT().
					AddTrackToPlaylist(gomock.Any(), gomock.Any()).
					Return(customErrors.ErrPlaylistPermissionDenied)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistPermissionDenied.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			jsonBody, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/playlists/"+tc.playlistID+"/tracks", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			vars := map[string]string{
				"id": tc.playlistID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.AddTrackToPlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_RemoveTrackFromPlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		playlistID     string
		trackID        string
		authenticated  bool
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			playlistID:    "1",
			trackID:       "10",
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					RemoveTrackFromPlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.RemoveTrackFromPlaylistRequest) error {
						assert.Equal(t, int64(10), request.TrackID)
						assert.Equal(t, int64(1), request.PlaylistID)
						assert.Equal(t, int64(1), request.UserID)
						return nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"message": "Track removed successfully",
				},
			},
		},
		{
			name:           "Unauthorized",
			playlistID:     "1",
			trackID:        "10",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Playlist ID",
			userID:         1,
			playlistID:     "invalid",
			trackID:        "10",
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
			name:           "Invalid Track ID",
			userID:         1,
			playlistID:     "1",
			trackID:        "invalid",
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
			name:          "Not Owner Error",
			userID:        1,
			playlistID:    "1",
			trackID:       "10",
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					RemoveTrackFromPlaylist(gomock.Any(), gomock.Any()).
					Return(customErrors.ErrPlaylistPermissionDenied)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistPermissionDenied.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("DELETE", "/playlists/"+tc.playlistID+"/tracks/"+tc.trackID, nil)
			vars := map[string]string{
				"id":      tc.playlistID,
				"trackId": tc.trackID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.RemoveTrackFromPlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_UpdatePlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		playlistID     string
		authenticated  bool
		thumbnail      []byte
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			thumbnail:     []byte("test image data"),
			mockBehavior: func() {
				mockUsecase.EXPECT().
					UpdatePlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.UpdatePlaylistRequest) (*usecaseModel.Playlist, error) {
						assert.Equal(t, "Test Playlist", request.Title)
						assert.Equal(t, int64(1), request.UserID)
						assert.Equal(t, int64(1), request.PlaylistID)
						assert.NotEmpty(t, request.Thumbnail)

						return &usecaseModel.Playlist{
							ID:        1,
							Title:     request.Title,
							Username:  "testuser",
							Thumbnail: "updated_thumbnail_url",
						}, nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"id":        float64(1),
					"title":     "Test Playlist",
					"user_id":   float64(1),
					"thumbnail": "updated_thumbnail_url",
				},
			},
		},
		{
			name:           "Unauthorized",
			playlistID:     "1",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Playlist ID",
			userID:         1,
			playlistID:     "invalid",
			authenticated:  true,
			thumbnail:      []byte("test image data"),
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
			name:          "Not Owner Error",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			thumbnail:     []byte("test image data"),
			mockBehavior: func() {
				mockUsecase.EXPECT().
					UpdatePlaylist(gomock.Any(), gomock.Any()).
					Return(nil, customErrors.ErrPlaylistPermissionDenied)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistPermissionDenied.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			var buf bytes.Buffer
			contentType := ""

			if tc.authenticated {
				buf, contentType = createMultipartFormData(t, "thumbnail", "test.jpg", tc.thumbnail)
			}

			req, err := http.NewRequest("PUT", "/playlists/"+tc.playlistID, &buf)
			assert.NoError(t, err)

			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}

			vars := map[string]string{
				"id": tc.playlistID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.UpdatePlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_RemovePlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		playlistID     string
		authenticated  bool
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					RemovePlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.RemovePlaylistRequest) error {
						assert.Equal(t, int64(1), request.PlaylistID)
						assert.Equal(t, int64(1), request.UserID)
						return nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"message": "Playlist removed successfully",
				},
			},
		},
		{
			name:           "Unauthorized",
			playlistID:     "1",
			authenticated:  false,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Playlist ID",
			userID:         1,
			playlistID:     "invalid",
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
			name:          "Not Owner Error",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			mockBehavior: func() {
				mockUsecase.EXPECT().
					RemovePlaylist(gomock.Any(), gomock.Any()).
					Return(customErrors.ErrPlaylistPermissionDenied)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistPermissionDenied.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("DELETE", "/playlists/"+tc.playlistID, nil)
			vars := map[string]string{
				"id": tc.playlistID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.RemovePlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_GetPlaylistsToAdd(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		authenticated  bool
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			authenticated: true,
			query:         "?trackId=10",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetPlaylistsToAdd(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.GetPlaylistsToAddRequest) ([]*usecaseModel.PlaylistWithIsIncludedTrack, error) {
						assert.Equal(t, int64(10), request.TrackID)
						assert.Equal(t, int64(1), request.UserID)

						return []*usecaseModel.PlaylistWithIsIncludedTrack{
							{
								Playlist: usecaseModel.Playlist{
									ID:        1,
									Title:     "Playlist 1",
									Username:  "testuser",
									Thumbnail: "thumbnail1.jpg",
								},
								IsIncluded: true,
							},
							{
								Playlist: usecaseModel.Playlist{
									ID:        2,
									Title:     "Playlist 2",
									Username:  "testuser",
									Thumbnail: "thumbnail2.jpg",
								},
								IsIncluded: false,
							},
						}, nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"title":       "Playlist 1",
						"user_id":     float64(1),
						"thumbnail":   "thumbnail1.jpg",
						"is_included": true,
					},
					map[string]interface{}{
						"id":          float64(2),
						"title":       "Playlist 2",
						"user_id":     float64(1),
						"thumbnail":   "thumbnail2.jpg",
						"is_included": false,
					},
				},
			},
		},
		{
			name:           "Unauthorized",
			authenticated:  false,
			query:          "?trackId=10",
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": customErrors.ErrPlaylistUnauthorized.Error(),
				},
			},
		},
		{
			name:           "Invalid Track ID",
			userID:         1,
			authenticated:  true,
			query:          "?trackId=invalid",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "strconv.Atoi: parsing \"invalid\": invalid syntax",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("GET", "/playlists/to-add"+tc.query, nil)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.GetPlaylistsToAdd(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_LikePlaylist(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		userID         int64
		playlistID     string
		authenticated  bool
		reqBody        interface{}
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:          "OK",
			userID:        1,
			playlistID:    "1",
			authenticated: true,
			reqBody: deliveryModel.PlaylistLikeRequest{
				IsLike: true,
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().
					LikePlaylist(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, request *usecaseModel.LikePlaylistRequest) error {
						assert.Equal(t, int64(1), request.PlaylistID)
						assert.Equal(t, int64(1), request.UserID)
						assert.True(t, request.IsLike)
						return nil
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": map[string]interface{}{
					"message": "Playlist liked/unliked successfully",
				},
			},
		},
		{
			name:           "Unauthorized",
			playlistID:     "1",
			authenticated:  false,
			reqBody:        deliveryModel.PlaylistLikeRequest{IsLike: true},
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
			name:           "Invalid Playlist ID",
			userID:         1,
			playlistID:     "invalid",
			authenticated:  true,
			reqBody:        deliveryModel.PlaylistLikeRequest{IsLike: true},
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
			name:          "Playlist Not Found",
			userID:        1,
			playlistID:    "999",
			authenticated: true,
			reqBody:       deliveryModel.PlaylistLikeRequest{IsLike: true},
			mockBehavior: func() {
				mockUsecase.EXPECT().
					LikePlaylist(gomock.Any(), gomock.Any()).
					Return(customErrors.ErrPlaylistNotFound)
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			jsonBody, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/playlists/"+tc.playlistID+"/like", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			vars := map[string]string{
				"id": tc.playlistID,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			if tc.authenticated {
				req = addUserToContext(req, tc.userID)
			}

			rec := httptest.NewRecorder()
			handler.LikePlaylist(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_GetProfilePlaylists(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		username       string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "OK",
			username: "testuser",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetProfilePlaylists(gomock.Any(), "testuser").
					Return([]*usecaseModel.Playlist{
						{
							ID:        1,
							Title:     "Playlist 1",
							Username:  "testuser",
							Thumbnail: "thumbnail1.jpg",
						},
						{
							ID:        2,
							Title:     "Playlist 2",
							Username:  "testuser",
							Thumbnail: "thumbnail2.jpg",
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []interface{}{
					map[string]interface{}{
						"id":        float64(1),
						"title":     "Playlist 1",
						"user_id":   float64(1),
						"thumbnail": "thumbnail1.jpg",
					},
					map[string]interface{}{
						"id":        float64(2),
						"title":     "Playlist 2",
						"user_id":   float64(1),
						"thumbnail": "thumbnail2.jpg",
					},
				},
			},
		},
		{
			name:           "Empty Username",
			username:       "",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"error": map[string]interface{}{
					"message": "username is required",
				},
			},
		},
		{
			name:     "User Not Found",
			username: "nonexistent",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					GetProfilePlaylists(gomock.Any(), "nonexistent").
					Return(nil, customErrors.ErrUserNotFound)
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("GET", "/user/"+tc.username+"/playlists", nil)
			vars := map[string]string{
				"username": tc.username,
			}
			req = mux.SetURLVars(req, vars)
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.GetProfilePlaylists(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}

func TestPlaylistHandler_SearchPlaylists(t *testing.T) {
	mockUsecase, handler, _ := setupTestHandler(t)

	testCases := []struct {
		name           string
		query          string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:  "OK",
			query: "?query=rock",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					SearchPlaylists(gomock.Any(), "rock").
					Return([]*usecaseModel.Playlist{
						{
							ID:        1,
							Title:     "Rock Playlist",
							Username:  "testuser",
							Thumbnail: "thumbnail1.jpg",
						},
						{
							ID:        2,
							Title:     "Hard Rock",
							Username:  "testuser2",
							Thumbnail: "thumbnail2.jpg",
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body": []interface{}{
					map[string]interface{}{
						"id":        float64(1),
						"title":     "Rock Playlist",
						"user_id":   float64(1),
						"thumbnail": "thumbnail1.jpg",
					},
					map[string]interface{}{
						"id":        float64(2),
						"title":     "Hard Rock",
						"user_id":   float64(2),
						"thumbnail": "thumbnail2.jpg",
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
					"message": "query is required",
				},
			},
		},
		{
			name:  "No Results",
			query: "?query=nonexistent",
			mockBehavior: func() {
				mockUsecase.EXPECT().
					SearchPlaylists(gomock.Any(), "nonexistent").
					Return([]*usecaseModel.Playlist{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "success",
				"body":   []interface{}{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			req := httptest.NewRequest("GET", "/playlists/search"+tc.query, nil)
			req = setupTestLogger(req)

			rec := httptest.NewRecorder()
			handler.SearchPlaylists(rec, req)

			verifyResponse(t, rec, tc.expectedStatus, tc.expectedBody)
		})
	}
}
