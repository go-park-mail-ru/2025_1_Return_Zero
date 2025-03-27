package repository

type Album struct {
	ID        uint
	Title     string
	Thumbnail string
	ArtistID  uint
}

type AlbumFilters struct {
	Pagination *Pagination
}
