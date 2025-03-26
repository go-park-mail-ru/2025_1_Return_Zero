package model

type TrackFilters struct {
	Pagination *Pagination
}

// Track represents a music track with its associated album and artist
// @Description A music track entity
type Track struct {
	ID        uint    `json:"id" example:"1" description:"Unique identifier"`
	Title     string  `json:"title" example:"Lagtrain" description:"Track title"`
	Thumbnail string  `json:"thumbnail_url" example:"https://example.com/image.jpg" description:"URL to the track thumbnail"`
	Album     *Album  `json:"album" description:"Associated album"`
	Artist    *Artist `json:"artist" description:"Associated artist"`
}
