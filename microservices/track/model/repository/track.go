package repository

type Track struct {
	ID        int64
	Title     string
	Thumbnail string
	Duration  int64
	AlbumID   int64
}

type TrackStreamCreateData struct {
	TrackID int64
	UserID  int64
}

type TrackStreamUpdateData struct {
	StreamID int64
	Duration int64
}

type TrackStream struct {
	ID       int64
	UserID   int64
	TrackID  int64
	Duration int64
}

type TrackWithFileKey struct {
	Track
	FileKey string
}

type Pagination struct {
	Limit  int64
	Offset int64
}

type TrackFilters struct {
	Pagination *Pagination
}
