package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_label "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
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

func addAdminToContext(req *http.Request) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxExtractor.AdminContextKey{}, int64(1))
	return req.WithContext(ctx)
}

func addLabelToContext(req *http.Request, labelID int64) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxExtractor.LabelContextKey{}, labelID)
	return req.WithContext(ctx)
}

func setupTestLogger(req *http.Request) *http.Request {
	log := zap.NewNop()
	ctx := req.Context()
	ctx = loggerPkg.LoggerToContext(ctx, log.Sugar())
	return req.WithContext(ctx)
}

func verifyResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedStatus, rec.Code)

	if expectedStatus == http.StatusOK || expectedStatus == http.StatusCreated {
		var response APISuccessResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody["status"], response.Status)

		if responseBody, ok := response.Body.(map[string]interface{}); ok {
			expectedResponseBody, ok := expectedBody["body"].(map[string]interface{})
			if ok {
				// Удаляем поле release_date из ответа, если оно существует
				delete(responseBody, "release_date")

				// Теперь сравниваем ожидаемое и фактическое тело ответа
				assert.Equal(t, expectedResponseBody, responseBody)
				return
			}
		}
		assert.Equal(t, expectedBody["body"], response.Body)
	} else {
		var response APIErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody["status"], response.Status)
		assert.Equal(t, expectedBody["error"].(map[string]interface{})["message"], response.Error)
	}
}

func TestLabelHandler_CreateLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		body           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Create Label Success",
			body: `{"label_name": "label1", "usernames": ["user1", "user2"]}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateLabel(gomock.Any(), gomock.Any()).Return(&usecase.Label{
					Id:      1,
					Name:    "label1",
					Members: []string{"user1", "user2"},
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusCreated),
				"body": map[string]interface{}{
					"id":         float64(1),
					"usernames":  []interface{}{"user1", "user2"},
					"label_name": "label1",
				},
			},
		},
		{
			name: "Create Label - not admin",
			body: `{"label_name": "label1", "usernames": ["user1", "user2"]}`,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "Unauthorized",
				},
			},
		},
		{
			name: "Create Label - wrong json",
			body: `{"label_name": "label1", "usernames": ["user1", "user2"`,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateLabel(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "Failed to read JSON",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body := bytes.NewBufferString(tt.body)
			req, err := http.NewRequest("POST", "/label", body)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Create Label - not admin" {
				req = addAdminToContext(req)
			}

			rec := httptest.NewRecorder()
			handler.CreateLabel(rec, req)
			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_GetLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		labelID        string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Get Label - not admin",
			labelID: "1",
			mockBehavior: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "Unauthorized",
				},
			},
		},
		{
			name:    "Get Label - wrong query",
			labelID: "wrong",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetLabel(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "Invalid label ID",
				},
			},
		},
		{
			name:    "Get Label - success",
			labelID: "1",
			mockBehavior: func() {
				mockUsecase.EXPECT().GetLabel(gomock.Any(), int64(1)).Return(&usecase.Label{
					Id:      1,
					Name:    "label1",
					Members: []string{"user1", "user2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": map[string]interface{}{
					"id":         float64(1),
					"label_name": "label1",
					"usernames":  []interface{}{"user1", "user2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			req, err := http.NewRequest("GET", "/label/"+tt.labelID, nil)
			vars := map[string]string{
				"id": tt.labelID,
			}
			req = mux.SetURLVars(req, vars)

			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Get Label - not admin" {
				req = addAdminToContext(req)
			}

			rec := httptest.NewRecorder()
			handler.GetLabel(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_UpdateLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		body           string
		to_add         []string
		to_remove      []string
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "Update Label Success",
			body:      `{"label_id": 1, "to_add": ["user3"], "to_remove": ["user1"]}`,
			to_add:    []string{"user3"},
			to_remove: []string{"user1"},
			labelID:   1,
			mockBehavior: func() {
				mockUsecase.EXPECT().UpdateLabel(gomock.Any(), int64(1), []string{"user3"}, []string{"user1"}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body":   "Label edited succesfully",
			},
		},
		{
			name:      "Update Label - not admin",
			body:      `{"label_id": "1", "to_add": ["user3"], "to_remove": ["user1"]}`,
			to_add:    []string{"user3"},
			to_remove: []string{"user1"},
			labelID:   1,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "Unauthorized",
				},
			},
		},
		{
			name:      "Update Label - wrong json",
			body:      `{"label_id": "1", "to_add": ["user3"], "to_remove": ["user1"`,
			to_add:    []string{"user3"},
			to_remove: []string{"user1"},
			labelID:   1,
			mockBehavior: func() {
				mockUsecase.EXPECT().UpdateLabel(gomock.Any(), int64(1), []string{"user3"}, []string{"user1"}).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "Failed to read JSON",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body := bytes.NewBufferString(tt.body)
			req, err := http.NewRequest("PUT", "/label", body)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Update Label - not admin" {
				req = addAdminToContext(req)
			}

			rec := httptest.NewRecorder()
			handler.UpdateLabel(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func createArtistMultipartFormData(t *testing.T, title string, thumbnailName string, thumbnailData []byte) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if thumbnailData != nil {
		part, err := w.CreateFormFile("thumbnail", thumbnailName)
		assert.NoError(t, err)
		_, err = io.Copy(part, bytes.NewReader(thumbnailData))
		assert.NoError(t, err)
	}

	err := w.WriteField("title", title)
	assert.NoError(t, err)

	err = w.Close()
	assert.NoError(t, err)

	return b, w.FormDataContentType()
}

func TestLabelHandler_CreateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		title          string
		image          []byte
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Create Artist Success",
			title:   "Test Artist",
			image:   []byte("fake image data"),
			labelID: 1,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateArtist(
					gomock.Any(),
					gomock.Any(),
				).Return(&usecase.Artist{
					ID:          1,
					Title:       "Test Artist",
					Thumbnail:   "test_artist.jpg",
					Description: "some description",
					IsLiked:     false,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusCreated),
				"body": map[string]interface{}{
					"id":            float64(1),
					"title":         "Test Artist",
					"thumbnail_url": "test_artist.jpg",
					"description":   "some description",
					"is_liked":      false,
				},
			},
		},
		{
			name:           "Create Artist - not label",
			title:          "Test Artist",
			image:          []byte("fake image data"),
			labelID:        0,
			mockBehavior:   func() {},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusForbidden),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:    "Create Artist - no title",
			title:   "",
			image:   []byte("fake image data"),
			labelID: 1,
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateArtist(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "title is empty",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body, str := createArtistMultipartFormData(t, tt.title, "test_artist.jpg", tt.image)
			req, err := http.NewRequest("POST", "/label/artist", &body)
			req.Header.Set("Content-Type", str)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Create Artist - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.CreateArtist(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func createEditArtistMultipartFormData(t *testing.T, artistID int64, title string, thumbnailName string, thumbnailData []byte) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if thumbnailData != nil {
		part, err := w.CreateFormFile("thumbnail", thumbnailName)
		assert.NoError(t, err)
		_, err = io.Copy(part, bytes.NewReader(thumbnailData))
		assert.NoError(t, err)
	}

	err := w.WriteField("title", title)
	assert.NoError(t, err)

	err = w.WriteField("artist_id", strconv.FormatInt(artistID, 10))
	assert.NoError(t, err)

	err = w.Close()
	assert.NoError(t, err)

	return b, w.FormDataContentType()
}

func TestLabelHandler_EditArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		artistID       int64
		newTitle       string
		image          []byte
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "Update Artist Success",
			artistID: 1,
			newTitle: "Updated Artist",
			image:    []byte("fake image data"),
			labelID:  1,
			mockBehavior: func() {
				mockUsecase.EXPECT().EditArtist(
					gomock.Any(),
					gomock.Any(),
				).Return(&usecase.Artist{
					ID:          1,
					Title:       "Updated Artist",
					Thumbnail:   "test_artist.jpg",
					Description: "some description",
					IsLiked:     false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": map[string]interface{}{
					"id":            float64(1),
					"title":         "Updated Artist",
					"thumbnail_url": "test_artist.jpg",
					"description":   "some description",
					"is_liked":      false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body, str := createEditArtistMultipartFormData(t, tt.artistID, tt.newTitle, "test_artist.jpg", tt.image)
			req, err := http.NewRequest("PUT", "/label/artist", &body)
			req.Header.Set("Content-Type", str)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Update Artist - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.EditArtist(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_GetArtists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Get Artists - not label",
			labelID: 0,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusForbidden),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:    "Get Artists - success",
			labelID: 1,
			mockBehavior: func() {
				mockUsecase.EXPECT().GetArtists(gomock.Any(), int64(1), gomock.Any()).Return([]*usecase.Artist{
					{
						ID:          1,
						Title:       "Artist 1",
						Thumbnail:   "artist1.jpg",
						Description: "Description 1",
						IsLiked:     false,
					},
					{
						ID:          2,
						Title:       "Artist 2",
						Thumbnail:   "artist2.jpg",
						Description: "Description 2",
						IsLiked:     true,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": []interface{}{
					map[string]interface{}{
						"id":            float64(1),
						"title":         "Artist 1",
						"thumbnail_url": "artist1.jpg",
						"description":   "Description 1",
						"is_liked":      false,
					},
					map[string]interface{}{
						"id":            float64(2),
						"title":         "Artist 2",
						"thumbnail_url": "artist2.jpg",
						"description":   "Description 2",
						"is_liked":      true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, err := http.NewRequest("GET", "/label/artists", nil)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Get Artists - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.GetArtists(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_DeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		artistID       int64
		labelID        int64
		body           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:     "Delete Artist - success",
			artistID: 1,
			body:     `{"artist_id": 1}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DeleteArtist(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": map[string]interface{}{
					"artist_id": float64(1),
				},
			},
		},
		{
			name:           "Delete Artist - not label",
			artistID:       1,
			body:           `{"artist_id": 1}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:     "Delete Artist - wrong json",
			artistID: 1,
			body:     `{"artist_id": 1`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DeleteArtist(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "unexpected EOF",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body := bytes.NewBufferString(tt.body)
			req, err := http.NewRequest("DELETE", "/label/artist", body)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Delete Artist - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.DeleteArtist(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func concatenateArtistIDs(artistIDs []int64) string {
	var result string
	for i, id := range artistIDs {
		if i > 0 {
			result += ","
		}
		result += strconv.FormatInt(id, 10)
	}
	return result
}

func createAlbumMultipartFormData(t *testing.T, artistIDs []int64, albumType string, title string, image []byte, tracks []*delivery.CreateTrackRequest) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err := w.WriteField("artists_ids", concatenateArtistIDs(artistIDs))
	assert.NoError(t, err)

	err = w.WriteField("type", albumType)
	assert.NoError(t, err)

	err = w.WriteField("title", title)
	assert.NoError(t, err)

	if image != nil {
		part, err := w.CreateFormFile("thumbnail", "album_image.jpg")
		assert.NoError(t, err)
		_, err = io.Copy(part, bytes.NewReader(image))
		assert.NoError(t, err)
	}

	for i, track := range tracks {
		trackPart, err := w.CreateFormFile("tracks[]", fmt.Sprintf("track_%d.mp3", i))
		assert.NoError(t, err)
		_, err = io.Copy(trackPart, bytes.NewReader(track.Track))
		assert.NoError(t, err)

		err = w.WriteField("track_titles[]", track.Title)
		assert.NoError(t, err)
	}

	err = w.Close()
	assert.NoError(t, err)

	return b, w.FormDataContentType()
}

func TestLabelHandler_CreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		artistIDs      []int64
		Type           string
		title          string
		image          []byte
		tracks         []*delivery.CreateTrackRequest
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "Create Album Success",
			artistIDs: []int64{1, 2},
			Type:      "album",
			title:     "Test Album",
			image:     []byte("fake image data"),
			tracks: []*delivery.CreateTrackRequest{
				{
					Title: "Track 1",
					Track: []byte("fake track data 1"),
				},
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateAlbum(
					gomock.Any(),
					gomock.Any(),
				).Return(int64(1), "album-url.jpg", nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusCreated),
				"body": map[string]interface{}{
					"id":            float64(1),
					"title":         "Test Album",
					"type":          "album",
					"thumbnail_url": "album-url.jpg",
					"is_liked":      false,
					"artists": []interface{}{
						map[string]interface{}{
							"id":    float64(1),
							"title": "",
						},
						map[string]interface{}{
							"id":    float64(2),
							"title": "",
						},
					},
				},
			},
		},
		{
			name:      "Create Album - not label",
			artistIDs: []int64{1, 2},
			Type:      "album",
			title:     "Test Album",
			image:     []byte("fake image data"),
			tracks: []*delivery.CreateTrackRequest{
				{
					Title: "Track 1",
					Track: []byte("fake track data 1"),
				},
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:      "Create Album - no title",
			artistIDs: []int64{1, 2},
			Type:      "album",
			title:     "",
			image:     []byte("fake image data"),
			tracks: []*delivery.CreateTrackRequest{
				{
					Title: "Track 1",
					Track: []byte("fake track data 1"),
				},
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateAlbum(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "title is empty",
				},
			},
		},
		{
			name:      "Create Album - no type",
			artistIDs: []int64{1, 2},
			Type:      "",
			title:     "album",
			image:     []byte("fake image data"),
			tracks: []*delivery.CreateTrackRequest{
				{
					Title: "Track 1",
					Track: []byte("fake track data 1"),
				},
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateAlbum(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "type is empty",
				},
			},
		},
		{
			name:      "Create Album - wrong type",
			artistIDs: []int64{1, 2},
			Type:      "wrong_type",
			title:     "album",
			image:     []byte("fake image data"),
			tracks: []*delivery.CreateTrackRequest{
				{
					Title: "Track 1",
					Track: []byte("fake track data 1"),
				},
			},
			mockBehavior: func() {
				mockUsecase.EXPECT().CreateAlbum(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "type is invalid",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body, str := createAlbumMultipartFormData(t, tt.artistIDs, tt.Type, tt.title, tt.image, tt.tracks)
			req, err := http.NewRequest("POST", "/label/album", &body)
			req.Header.Set("Content-Type", str)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Create Album - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.CreateAlbum(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_DeleteAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		albumID        int64
		labelID        int64
		body           string
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Delete Album - success",
			albumID: 1,
			body:    `{"album_id": 1}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DeleteAlbum(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": "Album deleted successfully",
			},
		},
		{
			name:           "Delete Album - not label",
			albumID:        1,
			body:           `{"album_id": 1}`,
			mockBehavior:   func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:     "Delete Album - wrong json",
			albumID:  1,
			body:     `{"album_id": 1`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DeleteAlbum(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusBadRequest),
				"error": map[string]interface{}{
					"message": "Failed to read JSON",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body := bytes.NewBufferString(tt.body)
			req, err := http.NewRequest("DELETE", "/label/album", body)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Delete Album - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.DeleteAlbum(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}

func TestLabelHandler_GetAlbumsByLabelID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_label.NewMockUsecase(ctrl)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	handler := NewLabelHandler(mockUsecase, cfg)

	tests := []struct {
		name           string
		labelID        int64
		mockBehavior   func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:    "Get Albums - not label",
			labelID: 0,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusUnauthorized),
				"error": map[string]interface{}{
					"message": "user not in label",
				},
			},
		},
		{
			name:    "Get Albums - success",
			labelID: 1,
			mockBehavior: func() {
				mockUsecase.EXPECT().GetAlbumsByLabelID(gomock.Any(), int64(1), gomock.Any()).Return([]*usecase.Album{
					{
						ID:            1,
						Type:          "album",
						Title:         "Album 1",
						IsLiked:       false,
						Thumbnail:     "album1.jpg",
						Artists: []*usecase.AlbumArtist{
							{
								ID:    1,
								Title: "Artist 1",
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": int(http.StatusOK),
				"body": []interface{}{
					map[string]interface{}{
						"id":            float64(1),
						"title":         "Album 1",
						"thumbnail_url": "album1.jpg",
						"type":          "album",
						"is_liked":      false,
						"release_date":  "0001-01-01T00:00:00Z",
						"artists": []interface{}{
							map[string]interface{}{
								"id":    float64(1),
								"title": "Artist 1",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req, err := http.NewRequest("GET", "/label/albums", nil)
			assert.NoError(t, err)

			req = setupTestLogger(req)
			if tt.name != "Get Albums - not label" {
				req = addAdminToContext(req)
				req = addLabelToContext(req, tt.labelID)
			}

			rec := httptest.NewRecorder()
			handler.GetAlbumsByLabelID(rec, req)

			verifyResponse(t, rec, tt.expectedStatus, tt.expectedBody)
		})
	}
}