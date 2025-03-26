package model

type Album struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail_url"`
	ArtistID  uint   `json:"artist_id"`
	Artist    Artist `json:"artist"`
}

type AlbumFilters struct {
	Pagination *Pagination
}
