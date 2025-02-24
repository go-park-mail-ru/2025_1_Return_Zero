package models

type Models struct {
	Artists *ArtistsModel
	Tracks  *TracksModel
	Albums  *AlbumsModel
}

func NewModels() *Models {
	return &Models{
		Artists: NewArtistsModel(),
		Tracks:  NewTracksModel(),
		Albums:  NewAlbumsModel(),
	}
}
