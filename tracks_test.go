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

func TestGetTracks(t *testing.T) {
	t.Parallel()
	testTracks := []models.Track{
		{
			ID:     1,
			Title:  "Track 1",
			Artist: "Artist 1",
			Image:  "Image 1",
		},
		{
			ID:     2,
			Title:  "Track 2",
			Artist: "Artist 2",
			Image:  "Image 2",
		},
		{
			ID:     3,
			Title:  "Track 3",
			Artist: "Artist 3",
			Image:  "Image 3",
		},
		{
			ID:     4,
			Title:  "Track 4",
			Artist: "Artist 4",
			Image:  "Image 4",
		},
		{
			ID:     5,
			Title:  "Track 5",
			Artist: "Artist 5",
			Image:  "Image 5",
		},
	}

	model := models.NewTracksModel()
	model.SetTestData(testTracks)

	tracksHandler := &TracksHandler{
		Model: model,
	}

	testCases := []struct {
		name           string
		offset         string
		limit          string
		expectedStatus int
		expectedTracks []models.Track
	}{
		{
			name:           "Max limit",
			offset:         "0",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedTracks: testTracks,
		},
		{
			name:           "0 limit",
			offset:         "0",
			limit:          "0",
			expectedStatus: http.StatusOK,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Negative offset",
			offset:         "-1",
			limit:          "5",
			expectedStatus: http.StatusBadRequest,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Negative limit",
			offset:         "0",
			limit:          "-1",
			expectedStatus: http.StatusBadRequest,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Offset greater than total",
			offset:         "10",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Offset and limit greater than total",
			offset:         "10",
			limit:          "5",
			expectedStatus: http.StatusOK,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Normal offset and limit",
			offset:         "2",
			limit:          "3",
			expectedStatus: http.StatusOK,
			expectedTracks: []models.Track{
				testTracks[2],
				testTracks[3],
				testTracks[4],
			},
		},
		{
			name:           "Invalid offset",
			offset:         "invalid",
			limit:          "5",
			expectedStatus: http.StatusBadRequest,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Invalid limit",
			offset:         "0",
			limit:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedTracks: []models.Track{},
		},
		{
			name:           "Empty offset and limit (default values)",
			offset:         "",
			limit:          "",
			expectedStatus: http.StatusOK,
			expectedTracks: testTracks,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/tracks", nil)
			query := req.URL.Query()
			query.Add("offset", testCase.offset)
			query.Add("limit", testCase.limit)
			req.URL.RawQuery = query.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tracksHandler.List)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, testCase.expectedStatus, rr.Code)
			if testCase.expectedStatus == http.StatusOK {
				var actualTracks []models.Track
				err := json.Unmarshal(rr.Body.Bytes(), &actualTracks)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedTracks, actualTracks)
			}
		})
	}
}
