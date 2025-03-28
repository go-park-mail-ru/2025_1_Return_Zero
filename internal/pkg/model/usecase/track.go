package usecase

type Track struct {
	ID        uint
	Title     string
	Thumbnail string
	Duration  int
	AlbumID   uint
	ArtistID  uint
	Album     string
	Artist    string
}

type TrackFilters struct {
	Pagination *Pagination
}
