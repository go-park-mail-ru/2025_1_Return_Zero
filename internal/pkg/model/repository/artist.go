package repository

import "errors"

var (
	ErrArtistNotFound = errors.New("artist not found")
)

type Artist struct {
	ID          int64  `sql:"id"`
	Title       string `sql:"title"`
	Description string `sql:"description"`
	Thumbnail   string `sql:"thumbnail_url"`
}

type ArtistWithRole struct {
	ID    int64  `sql:"id"`
	Title string `sql:"title"`
	Role  string `sql:"role"`
}

type ArtistFilters struct {
	Pagination *Pagination
}
