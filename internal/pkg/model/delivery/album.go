package delivery

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeEP          AlbumType = "ep"
	AlbumTypeCompilation AlbumType = "compilation"
	AlbumTypeSingle      AlbumType = "single"
)

// Album represents a music album with its associated artist
// @Description A music album entity
type Album struct {
	ID        int64     `json:"id" example:"1" description:"Unique identifier"`
	Title     string    `json:"title" example:"Anticyclone" description:"Album title"`
	Thumbnail string    `json:"thumbnail_url" example:"https://example.com/album.jpg" description:"URL to the album thumbnail"`
	Artist    string    `json:"artist" description:"Associated artist"`
	ArtistID  int64     `json:"artist_id" example:"1" description:"Unique identifier of the associated artist"`
	Type      AlbumType `json:"type" example:"album" description:"Type of the album"`
}

type AlbumFilters struct {
	Pagination *Pagination
}
