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

type TrackStreamCreateData struct {
	TrackID int64
	UserID  int64
}

// TrackStream represents a stream of a track (more like a listening session)
// @Description A stream of a track
type TrackStream struct {
	ID       int64 `json:"id" example:"1" description:"Unique identifier"`
	TrackID  int64 `json:"track_id" example:"1" description:"Unique identifier of the track"`
	Duration int64 `json:"duration" example:"216" description:"Stream duration in seconds"`
}

// TrackStreamUpdateData represents data that will be sent at the end of stream to update duration of the stream
// @Description an update data for stream
type TrackStreamUpdateData struct {
	Duration int64 `json:"duration" example:"216" description:"Stream duration in seconds" valid:"required,gte=0"`
}

// StreamID represents an id of newly created track stream
// @Description An id of a track stream
type StreamID struct {
	ID int64 `json:"id" example:"1" description:"Unique identifier"`
}

// TrackDetailed represents  a music track with its associated album and artist and also presigned file url
// Description A music track entity with file url
type TrackDetailed struct {
	Track
	FileUrl string `json:"file_url" example:"https://example.com/track.mp3" description:"URL to the track file"`
}
