package usecase

type Artist struct {
	ID          int64
	Title       string
	Description string
	Thumbnail   string
	IsFavorite  bool
	LabelID     int64
}

type ArtistList struct {
	Artists []*Artist
}

type ArtistDetailed struct {
	*Artist
	ListenersCount int64
	FavoritesCount int64
}

type ArtistWithTitle struct {
	ID    int64
	Title string
}

type ArtistWithTitleList struct {
	Artists []*ArtistWithTitle
}

type ArtistWithTitleMap struct {
	Artists map[int64]*ArtistWithTitleList
}

type ArtistWithRole struct {
	ID    int64
	Title string
	Role  string
}

type ArtistWithRoleList struct {
	Artists []*ArtistWithRole
}

type ArtistWithRoleMap struct {
	Artists map[int64]*ArtistWithRoleList
}

type Pagination struct {
	Offset int64
	Limit  int64
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
