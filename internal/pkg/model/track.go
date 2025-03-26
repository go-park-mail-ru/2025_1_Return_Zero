package model

type TrackFilters struct {
	Pagination *Pagination
}

// Track represents a music track with its associated album and artist
// @Description A music track entity
type Track struct {
	ID        uint     `json:"id" example:"1" description:"Unique identifier"`
	Title     string   `json:"title" example:"Lagtrain" description:"Track title"`
	Thumbnail string   `json:"thumbnail_url" example:"https://example.com/image.jpg" description:"URL to the track thumbnail"`
	Duration  int      `json:"duration" example:"216" description:"Track duration in seconds"`
	Album     *AlbumDB `json:"album" description:"Associated album"`
	Artist    *Artist  `json:"artist" description:"Associated artist"`
}

// TrackDB represents a music track with its associated album_id and artist_id
// @Description A music track entity
type TrackDB struct {
	ID        uint   `json:"id" example:"1" description:"Unique identifier"`
	Title     string `json:"title" example:"Lagtrain" description:"Track title"`
	Thumbnail string `json:"thumbnail_url" example:"https://example.com/image.jpg" description:"URL to the track thumbnail"`
	Duration  int    `json:"duration" example:"216" description:"Track duration in seconds"`
	AlbumID   uint   `json:"album_id" example:"1" description:"ID of the associated album"`
	ArtistID  uint   `json:"artist_id" example:"1" description:"ID of the associated artist"`
}
