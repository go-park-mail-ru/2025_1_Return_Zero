package usecase

import "time"

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeEP          AlbumType = "ep"
	AlbumTypeSingle      AlbumType = "single"
	AlbumTypeCompilation AlbumType = "compilation"
)

type AlbumArtist struct {
	ID    int64
	Title string
}

type Album struct {
	ID          int64
	Title       string
	Thumbnail   string
	Type        AlbumType
	ReleaseDate time.Time
	Artists     []*AlbumArtist
	IsLiked     bool
}

type AlbumFilters struct {
	Pagination *Pagination
}

type AlbumLikeRequest struct {
	AlbumID int64
	UserID  int64
	IsLike  bool
}

type CreateTrackRequest struct {
	Title    string
	Track    []byte
}

type CreateAlbumRequest struct {
	ArtistsIDs []int64
	Type       string
	Title      string
	Image      []byte
	Tracks     []*CreateTrackRequest
	LabelID    int64
}
