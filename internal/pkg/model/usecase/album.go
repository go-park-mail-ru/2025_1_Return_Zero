package usecase

type Album struct {
	ID        uint
	Title     string
	Thumbnail string
	Artist    string
	ArtistID  uint
}

type AlbumFilters struct {
	Pagination *Pagination
}
