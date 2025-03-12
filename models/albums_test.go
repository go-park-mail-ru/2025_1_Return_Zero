package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllAlbums(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		filters        Filters
		expectedAlbums []Album
	}{
		{
			name: "Valid filters with offset 0 and limit 3",
			filters: Filters{
				Offset: 0,
				Limit:  3,
			},
			expectedAlbums: []Album{
				{ID: 1, Title: "Anticyclone", Artist: "Inabakumori", Description: "Single", Image: "https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e"},
				{ID: 2, Title: "THE BOOK", Artist: "YOASOBI", Description: "Single", Image: "https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1"},
				{ID: 3, Title: "BOOTLEG", Artist: "Kenshi Yonezu", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
			},
		},
		{
			name: "Valid filters with offset 2 and limit 2",
			filters: Filters{
				Offset: 2,
				Limit:  2,
			},
			expectedAlbums: []Album{
				{ID: 3, Title: "BOOTLEG", Artist: "Kenshi Yonezu", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
				{ID: 4, Title: "Your Name.", Artist: "RADWIMPS", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
			},
		},
		{
			name: "Offset greater than albums length",
			filters: Filters{
				Offset: 10,
				Limit:  5,
			},
			expectedAlbums: []Album{},
		},
		{
			name: "Offset + limit greater than albums length",
			filters: Filters{
				Offset: 3,
				Limit:  5,
			},
			expectedAlbums: []Album{
				{ID: 4, Title: "Your Name.", Artist: "RADWIMPS", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model := NewAlbumsModel()
			albums := model.GetAll(testCase.filters)
			assert.Equal(t, testCase.expectedAlbums, albums)
		})
	}
}

func TestSetAlbumTestData(t *testing.T) {
	t.Parallel()

	testAlbums := []Album{
		{ID: 100, Title: "Test Album 1", Artist: "Test Artist", Description: "Test Description", Image: "test-image-1.jpg"},
		{ID: 101, Title: "Test Album 2", Artist: "Test Artist", Description: "Test Description", Image: "test-image-2.jpg"},
	}

	model := NewAlbumsModel()

	model.SetTestData(testAlbums)

	result := model.GetAll(Filters{Offset: 0, Limit: 10})
	assert.Equal(t, testAlbums, result)
}
