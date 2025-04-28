package usecase

import (
	"time"
)

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
	Type        AlbumType
	Thumbnail   string
	ReleaseDate time.Time
}

type AlbumList struct {
	Albums []*Album
}

type AlbumTitle struct {
	Title string
}

type AlbumTitleMap struct {
	Titles map[int64]*AlbumTitle
}

type Pagination struct {
	Offset int64
	Limit  int64
}

type AlbumFilters struct {
	Pagination *Pagination
}

type AlbumStreamCreateData struct {
	AlbumID int64
	UserID  int64
}
