package models

import (
	"sync"
)

type Artist struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

var artists = []Artist{
	{ID: 1, Title: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-000640888066-bwv7e8-t500x500.jpg"},
	{ID: 2, Title: "YOASOBI", Image: "https://i.scdn.co/image/ab67616100005174bfdd8a29d0c6bc6950055234"},
	{ID: 3, Title: "Kenshi Yonezu", Image: "https://i.scdn.co/image/ab6761610000e5ebd7ca899f6e53b54976a8594b"},
	{ID: 4, Title: "RADWIMPS", Image: "https://i.scdn.co/image/ab6761610000e5ebc9d443fb5ced1dd32d106632"},
	{ID: 5, Title: "Official HIGE DANdism", Image: "https://i.scdn.co/image/ab6761610000e5ebf9f7513528a90d1dde6d3aaa"},
}

type ArtistsModel struct {
	artists []Artist
	mutex   sync.RWMutex
	nextID  int
}

func NewArtistsModel() *ArtistsModel {
	return &ArtistsModel{
		// TODO: Change to empty artist object list after RK1
		artists: artists,
		mutex:   sync.RWMutex{},
		// TODO: Change to 0 after RK1
		nextID: 5,
	}
}

func (m *ArtistsModel) GetAll(filters Filters) []Artist {
	offset := filters.Offset
	limit := filters.Limit

	if offset > len(m.artists) {
		return []Artist{}
	}

	if offset+limit > len(m.artists) {
		limit = len(m.artists) - offset
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	artists := m.artists[offset : offset+limit]

	if len(artists) == 0 {
		return []Artist{}
	}

	return artists
}
