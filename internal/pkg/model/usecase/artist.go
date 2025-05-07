package usecase

type Artist struct {
	ID          int64
	Title       string
	Description string
	Thumbnail   string
	IsLiked     bool
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
