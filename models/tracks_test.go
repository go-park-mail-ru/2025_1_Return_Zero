package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTracks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		filters        Filters
		expectedTracks []Track
	}{
		{
			name: "Valid filters with offset 0 and limit 5",
			filters: Filters{
				Offset: 0,
				Limit:  5,
			},
			expectedTracks: []Track{
				{ID: 1, Title: "Lagtrain", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", Album: "Anticyclone"},
				{ID: 2, Title: "Lost Umbrella", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg", Album: "Anticyclone"},
				{ID: 3, Title: "Racing Into The Night", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", Album: "THE BOOK"},
				{ID: 4, Title: "Idol", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg", Album: "THE BOOK"},
				{ID: 5, Title: "Monster", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-ztyGtBiqtACBb5zy-WtrLrg-t500x500.jpg", Album: "THE BOOK"},
			},
		},
		{
			name: "Valid filters with offset 5 and limit 3",
			filters: Filters{
				Offset: 5,
				Limit:  3,
			},
			expectedTracks: []Track{
				{ID: 6, Title: "KICK BACK", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", Album: "BOOTLEG"},
				{ID: 7, Title: "Lemon", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg", Album: "BOOTLEG"},
				{ID: 8, Title: "Peace Sign", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000482219301-jrnq0h-t500x500.jpg", Album: "BOOTLEG"},
			},
		},
		{
			name: "Offset greater than tracks length",
			filters: Filters{
				Offset: 100,
				Limit:  5,
			},
			expectedTracks: []Track{},
		},
		{
			name: "Offset + limit greater than tracks length",
			filters: Filters{
				Offset: 18,
				Limit:  5,
			},
			expectedTracks: []Track{
				{ID: 19, Title: "Tabun", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-dumxejUZ4jURPErm-xUFVFw-t500x500.jpg", Album: "BOOTLEG"},
				{ID: 20, Title: "Ghost City Tokyo", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-ssoxHlQypZXAQKap-tEfJ6A-t500x500.jpg", Album: "BOOTLEG"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model := NewTracksModel()
			tracks := model.GetAll(testCase.filters)
			assert.Equal(t, testCase.expectedTracks, tracks)
		})
	}
}

func TestSetTestData(t *testing.T) {
	t.Parallel()

	testTracks := []Track{
		{ID: 100, Title: "Test Track 1", Artist: "Test Artist", Image: "test-image-1.jpg", Album: "Test Album"},
		{ID: 101, Title: "Test Track 2", Artist: "Test Artist", Image: "test-image-2.jpg", Album: "Test Album"},
	}

	model := NewTracksModel()

	model.SetTestData(testTracks)

	result := model.GetAll(Filters{Offset: 0, Limit: 10})
	assert.Equal(t, testTracks, result)
}
