package model

// Artist represents a music artist
// @Description A music artist entity
type Artist struct {
	ID        uint   `json:"id" example:"1" description:"Unique identifier"`
	Title     string `json:"title" example:"Inabakumori" description:"Artist name"`
	Thumbnail string `json:"thumbnail_url" example:"https://example.com/artist.jpg" description:"URL to the artist thumbnail"`
}

type ArtistFilters struct {
	Pagination *Pagination
}
