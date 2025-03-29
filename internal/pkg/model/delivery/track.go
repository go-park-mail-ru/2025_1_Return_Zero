package delivery

// TrackArtist represents an artist associated with a track
// @Description An artist associated with a track
type TrackArtist struct {
	ID    int64  `json:"id" example:"1" description:"Unique identifier"`
	Title string `json:"title" example:"Lagtrain" description:"Track title"`
	Role  string `json:"role" example:"Main artist" description:"Role of the artist"`
}

type TrackFilters struct {
	Pagination *Pagination
}

// Track represents a music track with its associated album and artist
// @Description A music track entity
type Track struct {
	ID        int64          `json:"id" example:"1" description:"Unique identifier"`
	Title     string         `json:"title" example:"Lagtrain" description:"Track title"`
	Thumbnail string         `json:"thumbnail_url" example:"https://example.com/image.jpg" description:"URL to the track thumbnail"`
	Duration  int64          `json:"duration" example:"216" description:"Track duration in seconds"`
	AlbumID   int64          `json:"album_id" example:"1" description:"Unique identifier of the associated album"`
	Album     string         `json:"album" description:"Associated album"`
	Artists   []*TrackArtist `json:"artists" description:"Associated artists"`
}
