package usecase

type TrackAlbum struct {
	ID    uint
	Title string
}

type TrackArtist struct {
	ID    uint
	Title string
}

type Track struct {
	ID        uint
	Title     string
	Thumbnail string
	Duration  int
	Album     TrackAlbum
	Artist    TrackArtist
}

type TrackFilters struct {
	Pagination *Pagination
}
