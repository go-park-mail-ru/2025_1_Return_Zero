package repository

type Artist struct {
	ID          int64  `sql:"id"`
	Title       string `sql:"title"`
	Description string `sql:"description"`
	Thumbnail   string `sql:"thumbnail_url"`
	IsFavorite  bool   `sql:"is_favorite"`
}

type ArtistWithTitle struct {
	ID    int64  `sql:"id"`
	Title string `sql:"title"`
}

type ArtistWithRole struct {
	ID    int64  `sql:"id"`
	Title string `sql:"title"`
	Role  string `sql:"role"`
}

type ArtistStats struct {
	ListenersCount int64 `sql:"listeners_count"`
	FavoritesCount int64 `sql:"favorites_count"`
}

type Pagination struct {
	Offset int64 `sql:"offset"`
	Limit  int64 `sql:"limit"`
}

type Filters struct {
	Pagination *Pagination
}

type ArtistStreamCreateDataList struct {
	ArtistIDs []int64
	UserID    int64
}

type LikeRequest struct {
	ArtistID int64
	UserID   int64
}
