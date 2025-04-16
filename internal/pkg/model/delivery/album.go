package delivery

import "time"

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeEP          AlbumType = "ep"
	AlbumTypeCompilation AlbumType = "compilation"
	AlbumTypeSingle      AlbumType = "single"
)

// AlbumArtist represents an artist associated with an album
// @Description An artist associated with an album
type AlbumArtist struct {
	ID    int64  `json:"id" example:"1" description:"Unique identifier"`
	Title string `json:"title" example:"Inabakumori" description:"Artist name"`
}

// Album represents a music album with its associated artist
// @Description A music album entity
type Album struct {
	ID          int64          `json:"id" example:"1" description:"Unique identifier"`
	Title       string         `json:"title" example:"Anticyclone" description:"Album title"`
	Thumbnail   string         `json:"thumbnail_url" example:"https://example.com/album.jpg" description:"URL to the album thumbnail"`
	Artists     []*AlbumArtist `json:"artists" description:"Associated artists"`
	Type        AlbumType      `json:"type" example:"album" description:"Type of the album"`
	ReleaseDate time.Time      `json:"release_date" example:"2021-01-01" description:"Release date of the album"`
}

type AlbumFilters struct {
	Pagination *Pagination
}
