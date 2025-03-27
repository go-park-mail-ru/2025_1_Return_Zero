package delivery

// Album represents a music album with its associated artist
// @Description A music album entity
type Album struct {
	ID        uint    `json:"id" example:"1" description:"Unique identifier"`
	Title     string  `json:"title" example:"Anticyclone" description:"Album title"`
	Thumbnail string  `json:"thumbnail_url" example:"https://example.com/album.jpg" description:"URL to the album thumbnail"`
	Artist    *Artist `json:"artist" description:"Associated artist"`
}

type AlbumUnpopulated struct {
	ID        uint   `json:"id" example:"1" description:"Unique identifier"`
	Title     string `json:"title" example:"Anticyclone" description:"Album title"`
	Thumbnail string `json:"thumbnail_url" example:"https://example.com/album.jpg" description:"URL to the album thumbnail"`
	ArtistID  uint   `json:"artist_id" example:"1" description:"ID of the associated artist"`
}

type AlbumFilters struct {
	Pagination *Pagination
}
