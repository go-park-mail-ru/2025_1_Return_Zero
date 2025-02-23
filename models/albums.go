package models

import (
	"errors"
	"sync"
)

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Image  string `json:"image"`
}

var albums = []Album{
	{ID: 1, Title: "Anticyclone", Artist: "Inabakumori", Image: "https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e"},
	{ID: 2, Title: "THE BOOK", Artist: "YOASOBI", Image: "https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1"},
	{ID: 3, Title: "BOOTLEG", Artist: "Kenshi Yonezu", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
	{ID: 4, Title: "Your Name.", Artist: "RADWIMPS", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
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
		nextID: 0,
	}
}

func (m *AlbumsModel) GetAll(filters Filters) ([]Album, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	offset := filters.offset()
	limit := filters.limit()

	if offset > len(m.albums) {
		return nil, errors.New("offset is greater than the number of albums")
	}

	if offset+limit > len(m.albums) {
		limit = len(m.albums) - offset
	}

	albums := m.albums[offset : offset+limit]
	if len(albums) == 0 {
		return nil, errors.New("no albums found")
	}

	return albums, nil
}
