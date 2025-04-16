package usecase

type Artist struct {
	ID          int64
	Title       string
	Description string
	Thumbnail   string
}

type ArtistDetailed struct {
	Artist
	Listeners int64
	Favorites int64
}

type ArtistFilters struct {
	Pagination *Pagination
}
