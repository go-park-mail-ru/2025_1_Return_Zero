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
	IsLiked   bool
}

type TrackStreamUpdateData struct {
	StreamID int64
	UserID   int64
	Duration int64
}

type TrackStreamCreateData struct {
	TrackID int64
	UserID  int64
}

type TrackStream struct {
	ID       int64
	TrackID  int64
	Duration int64
}

type TrackDetailed struct {
	Track
	FileUrl string
}

type TrackFilters struct {
	Pagination *Pagination
}

type TrackLikeRequest struct {
	TrackID int64
	IsLike  bool
	UserID  int64
}
