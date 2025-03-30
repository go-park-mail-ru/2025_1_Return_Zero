package usecase

type Artist struct {
	ID          int64
	Title       string
	Description string
	Thumbnail   string
}

type ArtistDetailed struct {
	Artist
	Listeners int64 `json:"listeners_count" example:"1000" description:"Number of listeners"`
	Favorites int64 `json:"favorites_count" example:"1000" description:"Number of favorites"`
}

type ArtistFilters struct {
	Pagination *Pagination
}
