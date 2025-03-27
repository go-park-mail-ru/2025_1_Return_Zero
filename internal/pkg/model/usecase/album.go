package usecase

type Album struct {
	ID        uint
	Title     string
	Thumbnail string
	Artist    *Artist
}

type AlbumUnpopulated struct {
	ID        uint
	Title     string
	Thumbnail string
	ArtistID  uint
}

type AlbumFilters struct {
	Pagination *Pagination
}
