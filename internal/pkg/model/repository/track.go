package repository

type Track struct {
	ID        int64
	Title     string
	Thumbnail string
	Duration  int64
	AlbumID   int64
	Album     string
}

type TrackWithFileKey struct {
	Track
	FileKey string
}

type TrackFilters struct {
	Pagination *Pagination
}
