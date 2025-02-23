package models

import (
	"errors"
	"sync"
)

type Track struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Image  string `json:"image"`
}

var tracks = []Track{
	{ID: 1, Name: "Lagtrain", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg"},
	{ID: 2, Name: "Lost Umbrella", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg"},
	{ID: 3, Name: "Racing Into The Night", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg"},
	{ID: 4, Name: "Idol", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg"},
	{ID: 5, Name: "Monster", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-ztyGtBiqtACBb5zy-WtrLrg-t500x500.jpg"},
	{ID: 6, Name: "KICK BACK", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg"},
	{ID: 7, Name: "Lemon", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg"},
	{ID: 8, Name: "Peace Sign", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000482219301-jrnq0h-t500x500.jpg"},
	{ID: 9, Name: "Sparkle", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000452912388-ft13zk-t1080x1080.jpg"},
	{ID: 10, Name: "Nandemonaiya", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000230768346-878y9o-t500x500.jpg"},
	{ID: 11, Name: "Suzume", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-OR55dgkv9l0JHg6J-NUMaSQ-t500x500.jpg"},
	{ID: 12, Name: "Pretender", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-000644002372-j1fgr1-t500x500.jpg"},
	{ID: 13, Name: "Mixed Nuts", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-68ZsJYMEYjCHMEpM-z4UHxg-t500x500.jpg"},
	{ID: 14, Name: "Cry Baby", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-G0RPyB0xahP2CyHW-4H1THQ-t500x500.jpg"},
	{ID: 15, Name: "Dream Lantern", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000350712186-6xaoo7-t500x500.jpg"},
	{ID: 16, Name: "Zenzenzense", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000189644938-tywci0-t1080x1080.jpg"},
	{ID: 17, Name: "Shinigami", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-Z0nrZzzmeWrfD6ny-iVaI8w-t500x500.jpg"},
	{ID: 18, Name: "Gunjo", Artist: "YOASOBI", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpTz3Ys6eJleSr2shfdB2BMq15WsipNe4rgQ&s"},
	{ID: 19, Name: "Tabun", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-dumxejUZ4jURPErm-xUFVFw-t500x500.jpg"},
	{ID: 20, Name: "Ghost City Tokyo", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-ssoxHlQypZXAQKap-tEfJ6A-t500x500.jpg"},
}

type TracksModel struct {
	tracks []Track
	mutex  sync.RWMutex
	nextID int
}

func NewTracksModel() *TracksModel {
	return &TracksModel{
		// TODO: Change to empty track object list after RK1
		tracks: tracks,
		mutex:  sync.RWMutex{},
		nextID: 0,
	}
}

func (m *TracksModel) GetAll(filters Filters) ([]Track, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	offset := filters.offset()
	limit := filters.limit()

	if offset > len(m.tracks) {
		return nil, errors.New("offset is greater than the number of tracks")
	}

	if offset+limit > len(m.tracks) {
		limit = len(m.tracks) - offset
	}

	tracks := m.tracks[offset : offset+limit]
	if len(tracks) == 0 {
		return nil, errors.New("no tracks found")
	}

	return tracks, nil
}
