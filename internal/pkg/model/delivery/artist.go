package delivery

// Artist represents a music artist
// @Description A music artist entity
type Artist struct {
	ID          int64  `json:"id" example:"1" description:"Unique identifier"`
	Title       string `json:"title" example:"Inabakumori" description:"Artist name"`
	Description string `json:"description" example:"Inabakumori is a Japanese artist" description:"Artist description"`
	Thumbnail   string `json:"thumbnail_url" example:"https://example.com/artist.jpg" description:"URL to the artist thumbnail"`
	IsLiked     bool   `json:"is_liked" example:"false" description:"Whether the artist is liked"`
}

// ArtistDetailed represents a detailed music artist entity
// @Description A detailed music artist entity
type ArtistDetailed struct {
	Artist
	Listeners int64 `json:"listeners_count" example:"1000" description:"Number of listeners"`
	Favorites int64 `json:"favorites_count" example:"1000" description:"Number of favorites"`
}

type ArtistFilters struct {
	Pagination *Pagination
}

// LikeRequest represents a request to like or unlike an artist
// @Description A request to like or unlike an artist. Should be authenticated
type ArtistLikeRequest struct {
	IsLike bool `json:"value" example:"true" description:"Whether the artist is liked"`
}
