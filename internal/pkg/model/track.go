package model

type TrackFilters struct {
	Pagination *Pagination
}

type Track struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Thumbnail string  `json:"thumbnail_url"`
	Album     *Album  `json:"album"`
	Artist    *Artist `json:"artist"`
}
