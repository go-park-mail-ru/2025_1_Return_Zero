package delivery

// Album represents a music album with its associated artist
// @Description A music album entity
type Album struct {
	ID        uint   `json:"id" example:"1" description:"Unique identifier"`
	Title     string `json:"title" example:"Anticyclone" description:"Album title"`
	Thumbnail string `json:"thumbnail_url" example:"https://example.com/album.jpg" description:"URL to the album thumbnail"`
	Artist    string `json:"artist" description:"Associated artist"`
	ArtistID  uint   `json:"artist_id" example:"1" description:"Unique identifier of the associated artist"`
}

type AlbumFilters struct {
	Pagination *Pagination
}
