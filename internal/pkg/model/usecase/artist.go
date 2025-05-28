package usecase

type Artist struct {
	ID          int64
	Title       string
	Description string
	Thumbnail   string
	IsLiked     bool
	LabelID     int64
}

type ArtistDetailed struct {
	Artist
	Listeners int64
	Favorites int64
}

type ArtistFilters struct {
	Pagination *Pagination
}

type ArtistLikeRequest struct {
	ArtistID int64
	UserID   int64
	IsLike   bool
}

type ArtistLoad struct {
	Title   string
	Image   []byte
	LabelID int64
}

type ArtistEdit struct {
	ArtistID int64
	NewTitle string
	Image    []byte
	LabelID  int64
}

type ArtistDelete struct {
	ArtistID int64
	LabelID  int64
}
