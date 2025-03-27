package usecase

type Track struct {
	ID        uint
	Title     string
	Thumbnail string
	Duration  int
	Album     *AlbumUnpopulated
	Artist    *Artist
}

type TrackFilters struct {
	Pagination *Pagination
}
