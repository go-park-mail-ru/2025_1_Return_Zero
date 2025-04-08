package repository

import (
	"errors"
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
}

type AlbumFilters struct {
	Pagination *Pagination
}

var (
	ErrAlbumNotFound = errors.New("album not found")
)
