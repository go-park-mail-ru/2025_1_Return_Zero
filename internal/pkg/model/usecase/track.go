package usecase

type TrackArtist struct {
	ID    int64
	Title string
	Role  string
}

type Track struct {
	ID        int64
	Title     string
	Thumbnail string
	Duration  int64
	AlbumID   int64
	Album     string
	Artists   []*TrackArtist
}

type TrackDetailed struct {
	Track
	FileUrl string
}

type TrackFilters struct {
	Pagination *Pagination
}
