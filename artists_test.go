package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/models"
	"github.com/stretchr/testify/assert"
)

func TestListArtists(t *testing.T) {
	t.Parallel()
	testArtists := []models.Artist{
		{
			ID:    1,
			Title: "Artist 1",
			Image: "Image 1",
		},
		{
			ID:    2,
			Title: "Artist 2",
			Image: "Image 2",
		},
		{
			ID:    3,
			Title: "Artist 3",
			Image: "Image 3",
		},
		{
			ID:    4,
			Title: "Artist 4",
			Image: "Image 4",
		},
		{
			ID:    5,
			Title: "Artist 5",
			Image: "Image 5",
		},
	}

	model := models.NewArtistsModel()
	model.SetTestData(testArtists)

	artistsHandler := &ArtistsHandler{
		Model: model,
	}

	testCases := []struct {
		name            string
		offset          string
		limit           string
		expectedStatus  int
		expectedArtists []models.Artist
	}{
		{
			name:            "Max limit",
			offset:          "0",
			limit:           "5",
			expectedStatus:  http.StatusOK,
			expectedArtists: testArtists,
		},
		{
			name:            "0 limit",
			offset:          "0",
			limit:           "0",
			expectedStatus:  http.StatusOK,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Negative offset",
			offset:          "-1",
			limit:           "5",
			expectedStatus:  http.StatusBadRequest,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Negative limit",
			offset:          "0",
			limit:           "-1",
			expectedStatus:  http.StatusBadRequest,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Offset greater than total",
			offset:          "10",
			limit:           "5",
			expectedStatus:  http.StatusOK,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Offset and limit greater than total",
			offset:          "10",
			limit:           "5",
			expectedStatus:  http.StatusOK,
			expectedArtists: []models.Artist{},
		},
		{
			name:           "Normal offset and limit",
			offset:         "2",
			limit:          "3",
			expectedStatus: http.StatusOK,
			expectedArtists: []models.Artist{
				testArtists[2],
				testArtists[3],
				testArtists[4],
			},
		},
		{
			name:            "Invalid offset",
			offset:          "invalid",
			limit:           "5",
			expectedStatus:  http.StatusBadRequest,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Invalid limit",
			offset:          "0",
			limit:           "invalid",
			expectedStatus:  http.StatusBadRequest,
			expectedArtists: []models.Artist{},
		},
		{
			name:            "Empty offset and limit (default values)",
			offset:          "",
			limit:           "",
			expectedStatus:  http.StatusOK,
			expectedArtists: testArtists,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/artists", nil)
			query := req.URL.Query()
			query.Add("offset", testCase.offset)
			query.Add("limit", testCase.limit)
			req.URL.RawQuery = query.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(artistsHandler.List)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCase.expectedStatus, rr.Code)
			if testCase.expectedStatus == http.StatusOK {
				var actualArtists []models.Artist
				err := json.Unmarshal(rr.Body.Bytes(), &actualArtists)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedArtists, actualArtists)
			}
		})
	}
}
