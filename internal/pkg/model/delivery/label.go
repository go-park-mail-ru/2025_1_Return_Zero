package delivery

type CreateArtistRequest struct {
	Title   string `json:"title"`
	Image   []byte `json:"image"`
	LabelID int64  `json:"label_id"`
}

type EditArtistRequest struct {
	ArtistID int64  `json:"artist_id"`
	NewTitle string `json:"new_title,omitempty"`
	Image    []byte `json:"image,omitempty"`
	LabelID  int64  `json:"label_id"`
}

type DeleteArtistRequest struct {
	ArtistID int64 `json:"artist_id"`
	LabelID  int64 `json:"-"`
}

type CreateTrackRequest struct {
	Title string `json:"title"`
	Track []byte `json:"track"`
}

type CreateAlbumRequest struct {
	ArtistsIDs []int64               `json:"artists_ids"`
	Type       string                `json:"type"`
	Title      string                `json:"title"`
	Image      []byte                `json:"image"`
	Tracks     []*CreateTrackRequest `json:"tracks"`
	LabelID    int64                 `json:"label_id"`
}

type EditLabelRequest struct {
	LabelID  int64    `json:"label_id"`
	ToAdd    []string `json:"to_add,omitempty"`
	ToRemove []string `json:"to_remove,omitempty"`
}

type DeleteAlbumRequest struct {
	AlbumID int64 `json:"album_id"`
	LabelID int64 `json:"label_id"`
}
