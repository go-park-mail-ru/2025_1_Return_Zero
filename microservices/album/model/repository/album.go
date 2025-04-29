package repository

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
	ID          int64     `sql:"id"`
	Title       string    `sql:"title"`
	Type        AlbumType `sql:"type"`
	Thumbnail   string    `sql:"thumbnail_url"`
	ReleaseDate time.Time `sql:"release_date"`
	IsFavorite  bool      `sql:"is_favorite"`
}

type Pagination struct {
	Offset int64 `sql:"offset"`
	Limit  int64 `sql:"limit"`
}

type AlbumFilters struct {
	Pagination *Pagination
}

type AlbumStreamCreateData struct {
	AlbumID int64
	UserID  int64
}

type LikeRequest struct {
	AlbumID int64
	UserID  int64
}
