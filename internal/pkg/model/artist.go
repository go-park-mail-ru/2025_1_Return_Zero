package model

type Artist struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail_url"`
}

type ArtistFilters struct {
	Pagination *Pagination
}
