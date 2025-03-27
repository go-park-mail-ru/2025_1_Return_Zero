package delivery

type TrackFilters struct {
	Pagination *Pagination
}

// TrackAlbum represents an album of a track
// @Description An album of a track entity
type TrackAlbum struct {
	ID    uint   `json:"id" example:"1" description:"Unique identifier"`
	Title string `json:"title" example:"Lagtrain" description:"Album title"`
}

// TrackArtist represents an artist of a track
// @Description An artist of a track entity
type TrackArtist struct {
	ID    uint   `json:"id" example:"1" description:"Unique identifier"`
	Title string `json:"title" example:"Lagtrain" description:"Artist title"`
}

// Track represents a music track with its associated album and artist
// @Description A music track entity
type Track struct {
	ID        uint        `json:"id" example:"1" description:"Unique identifier"`
	Title     string      `json:"title" example:"Lagtrain" description:"Track title"`
	Thumbnail string      `json:"thumbnail_url" example:"https://example.com/image.jpg" description:"URL to the track thumbnail"`
	Duration  int         `json:"duration" example:"216" description:"Track duration in seconds"`
	Album     TrackAlbum  `json:"album" description:"Associated album"`
	Artist    TrackArtist `json:"artist" description:"Associated artist"`
}
