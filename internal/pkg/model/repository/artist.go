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
	Listeners   int64  `sql:"listeners_count"`
	Favorites   int64  `sql:"favorites_count"`
}

type ArtistWithRole struct {
	ID    int64  `sql:"id"`
	Title string `sql:"title"`
	Role  string `sql:"role"`
}

type ArtistFilters struct {
	Pagination *Pagination
}
