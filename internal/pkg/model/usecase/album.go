package usecase

import "time"

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeEP          AlbumType = "ep"
	AlbumTypeSingle      AlbumType = "single"
	AlbumTypeCompilation AlbumType = "compilation"
)

type Album struct {
	ID          int64
	Title       string
	Thumbnail   string
	Artist      string
	ArtistID    int64
	Type        AlbumType
	ReleaseDate time.Time
}

type AlbumFilters struct {
	Pagination *Pagination
}
