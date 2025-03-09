package models

import (
	"sync"
)

type Track struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Image  string `json:"image"`
	Album  string `json:"album"`
}

var tracks = []Track{
	{ID: 1, Title: "Lagtrain", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", Album: "Anticyclone"},
	{ID: 2, Title: "Lost Umbrella", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg", Album: "Anticyclone"},
	{ID: 3, Title: "Racing Into The Night", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", Album: "THE BOOK"},
	{ID: 4, Title: "Idol", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg", Album: "THE BOOK"},
	{ID: 5, Title: "Monster", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-ztyGtBiqtACBb5zy-WtrLrg-t500x500.jpg", Album: "THE BOOK"},
	{ID: 6, Title: "KICK BACK", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", Album: "BOOTLEG"},
	{ID: 7, Title: "Lemon", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg", Album: "BOOTLEG"},
	{ID: 8, Title: "Peace Sign", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000482219301-jrnq0h-t500x500.jpg", Album: "BOOTLEG"},
	{ID: 9, Title: "Sparkle", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000452912388-ft13zk-t1080x1080.jpg", Album: "Your Name."},
	{ID: 10, Title: "Nandemonaiya", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000230768346-878y9o-t500x500.jpg", Album: "Your Name."},
	{ID: 11, Title: "Suzume", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-OR55dgkv9l0JHg6J-NUMaSQ-t500x500.jpg", Album: "Your Name."},
	{ID: 12, Title: "Pretender", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-000644002372-j1fgr1-t500x500.jpg", Album: "Your Name."},
	{ID: 13, Title: "Mixed Nuts", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-68ZsJYMEYjCHMEpM-z4UHxg-t500x500.jpg", Album: "Your Name."},
	{ID: 14, Title: "Cry Baby", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-G0RPyB0xahP2CyHW-4H1THQ-t500x500.jpg", Album: "Your Name."},
	{ID: 15, Title: "Dream Lantern", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000350712186-6xaoo7-t500x500.jpg", Album: "Your Name."},
	{ID: 16, Title: "Zenzenzense", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000189644938-tywci0-t1080x1080.jpg", Album: "Your Name."},
	{ID: 17, Title: "Shinigami", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-Z0nrZzzmeWrfD6ny-iVaI8w-t500x500.jpg", Album: "BOOTLEG"},
	{ID: 18, Title: "Gunjo", Artist: "YOASOBI", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpTz3Ys6eJleSr2shfdB2BMq15WsipNe4rgQ&s", Album: "BOOTLEG"},
	{ID: 19, Title: "Tabun", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-dumxejUZ4jURPErm-xUFVFw-t500x500.jpg", Album: "BOOTLEG"},
	{ID: 20, Title: "Ghost City Tokyo", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-ssoxHlQypZXAQKap-tEfJ6A-t500x500.jpg", Album: "BOOTLEG"},
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
		// TODO: Change to 20 after RK1
		nextID: 20,
	}
}

func (m *TracksModel) GetAll(filters Filters) []Track {
	offset := filters.Offset
	limit := filters.Limit

	if offset > len(m.tracks) {
		return []Track{}
	}

	if offset+limit > len(m.tracks) {
		limit = len(m.tracks) - offset
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	tracks := m.tracks[offset : offset+limit]

	if len(tracks) == 0 {
		return []Track{}
	}

	return tracks
}

// Only for testing purposes
func (m *TracksModel) SetTestData(testTracks []Track) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tracks = testTracks
}
