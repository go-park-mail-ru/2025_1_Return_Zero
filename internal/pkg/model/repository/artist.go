package repository

type Artist struct {
	ID        uint
	Title     string
	Thumbnail string
}

type ArtistFilters struct {
	Pagination *Pagination
}
