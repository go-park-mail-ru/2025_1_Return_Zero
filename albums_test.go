package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
	"github.com/stretchr/testify/assert"
)

func TestListAlbums(t *testing.T) {
	t.Parallel()
	testAlbums := []models.Album{
		{
			ID:     1,
			Title:  "Album 1",
			Artist: "Artist 1",
			Image:  "Image 1",
		},
		{
			ID:     2,
			Title:  "Album 2",
			Artist: "Artist 2",
			Image:  "Image 2",
		},
		{
			ID:     3,
			Title:  "Album 3",
			Artist: "Artist 3",
			Image:  "Image 3",
		},
		{
			ID:     4,
			Title:  "Album 4",
			Artist: "Artist 4",
			Image:  "Image 4",
		},
		{
			ID:     5,
			Title:  "Album 5",
			Artist: "Artist 5",
			Image:  "Image 5",
		},
	}

	model := models.NewAlbumsModel()
	model.SetTestData(testAlbums)

	albumsHandler := &AlbumsHandler{
		Model: model,
	}

	testCases := []struct {
		name           string
		offset         string
		limit          string
		expectedStatus int
		expectedAlbums []models.Album
	}{
		{
			name:           "Max limit",
			offset:         "0",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedAlbums: testAlbums,
		},
		{
			name:           "0 limit",
			offset:         "0",
			limit:          "0",
			expectedStatus: http.StatusOK,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Negative offset",
			offset:         "-1",
			limit:          "5",
			expectedStatus: http.StatusBadRequest,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Negative limit",
			offset:         "0",
			limit:          "-1",
			expectedStatus: http.StatusBadRequest,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Offset greater than total",
			offset:         "10",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Offset and limit greater than total",
			offset:         "10",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Normal offset and limit",
			offset:         "2",
			limit:          "3",
			expectedStatus: http.StatusOK,
			expectedAlbums: []models.Album{
				testAlbums[2],
				testAlbums[3],
				testAlbums[4],
			},
		},
		{
			name:           "Invalid offset",
			offset:         "invalid",
			limit:          "5",
			expectedStatus: http.StatusBadRequest,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Invalid limit",
			offset:         "0",
			limit:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedAlbums: []models.Album{},
		},
		{
			name:           "Empty offset and limit (default values)",
			offset:         "",
			limit:          "",
			expectedStatus: http.StatusOK,
			expectedAlbums: testAlbums,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/albums", nil)
			query := req.URL.Query()
			query.Add("offset", testCase.offset)
			query.Add("limit", testCase.limit)
			req.URL.RawQuery = query.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(albumsHandler.List)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCase.expectedStatus, rr.Code)
			if testCase.expectedStatus == http.StatusOK {
				var actualAlbums []models.Album
				err := json.Unmarshal(rr.Body.Bytes(), &actualAlbums)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedAlbums, actualAlbums)
			}
		})
	}
}
