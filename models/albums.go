package models

import (
	"sync"
)

type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

var albums = []Album{
	{ID: 1, Title: "Anticyclone", Artist: "Inabakumori", Description: "Single", Image: "https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e"},
	{ID: 2, Title: "THE BOOK", Artist: "YOASOBI", Description: "Single", Image: "https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1"},
	{ID: 3, Title: "BOOTLEG", Artist: "Kenshi Yonezu", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
	{ID: 4, Title: "Your Name.", Artist: "RADWIMPS", Description: "Ep", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
}

type AlbumsModel struct {
	albums []Album
	mutex  sync.RWMutex
	nextID int
}

func NewAlbumsModel() *AlbumsModel {
	return &AlbumsModel{
		// TODO: Change to empty album object list after RK1
		albums: albums,
		mutex:  sync.RWMutex{},
		// TODO: Change to 0 after RK1
		nextID: 4,
	}
}

func (m *AlbumsModel) GetAll(filters Filters) []Album {
	offset := filters.Offset
	limit := filters.Limit

	if offset > len(m.albums) {
		return []Album{}
	}

	if offset+limit > len(m.albums) {
		limit = len(m.albums) - offset
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	albums := m.albums[offset : offset+limit]
	if len(albums) == 0 {
		return []Album{}
	}

	return albums
}

// Only for testing purposes
func (m *AlbumsModel) SetTestData(testAlbums []Album) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.albums = testAlbums
}
