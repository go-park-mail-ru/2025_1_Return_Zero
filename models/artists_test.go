package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllArtists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		filters         Filters
		expectedArtists []Artist
	}{
		{
			name: "Valid filters with offset 0 and limit 3",
			filters: Filters{
				Offset: 0,
				Limit:  3,
			},
			expectedArtists: []Artist{
				{ID: 1, Title: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-000640888066-bwv7e8-t500x500.jpg"},
				{ID: 2, Title: "YOASOBI", Image: "https://i.scdn.co/image/ab67616100005174bfdd8a29d0c6bc6950055234"},
				{ID: 3, Title: "Kenshi Yonezu", Image: "https://i.scdn.co/image/ab6761610000e5ebd7ca899f6e53b54976a8594b"},
			},
		},
		{
			name: "Valid filters with offset 2 and limit 2",
			filters: Filters{
				Offset: 2,
				Limit:  2,
			},
			expectedArtists: []Artist{
				{ID: 3, Title: "Kenshi Yonezu", Image: "https://i.scdn.co/image/ab6761610000e5ebd7ca899f6e53b54976a8594b"},
				{ID: 4, Title: "RADWIMPS", Image: "https://i.scdn.co/image/ab6761610000e5ebc9d443fb5ced1dd32d106632"},
			},
		},
		{
			name: "Offset greater than artists length",
			filters: Filters{
				Offset: 10,
				Limit:  5,
			},
			expectedArtists: []Artist{},
		},
		{
			name: "Offset + limit greater than artists length",
			filters: Filters{
				Offset: 4,
				Limit:  5,
			},
			expectedArtists: []Artist{
				{ID: 5, Title: "Official HIGE DANdism", Image: "https://i.scdn.co/image/ab6761610000e5ebf9f7513528a90d1dde6d3aaa"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model := NewArtistsModel()
			artists := model.GetAll(testCase.filters)
			assert.Equal(t, testCase.expectedArtists, artists)
		})
	}
}

func TestSetArtistTestData(t *testing.T) {
	t.Parallel()

	testArtists := []Artist{
		{ID: 100, Title: "Test Artist 1", Image: "test-image-1.jpg"},
		{ID: 101, Title: "Test Artist 2", Image: "test-image-2.jpg"},
	}

	model := NewArtistsModel()

	model.SetTestData(testArtists)

	result := model.GetAll(Filters{Offset: 0, Limit: 10})
	assert.Equal(t, testArtists, result)
}
