package usecase

type AlbumArtist struct {
	ID    uint
	Title string
}

type Album struct {
	ID        uint
	Title     string
	Thumbnail string
	Artist    AlbumArtist
}

type AlbumFilters struct {
	Pagination *Pagination
}
