package repository

type Track struct {
	ID        uint
	Title     string
	Thumbnail string
	Duration  int
	AlbumID   uint
	ArtistID  uint
}

type TrackFilters struct {
	Pagination *Pagination
}
