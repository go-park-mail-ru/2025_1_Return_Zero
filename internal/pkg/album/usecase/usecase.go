package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func NewUsecase(albumRepository album.Repository, artistRepository artist.Repository) album.Usecase {
	return albumUsecase{albumRepo: albumRepository, artistRepo: artistRepository}
}

type albumUsecase struct {
	albumRepo  album.Repository
	artistRepo artist.Repository
}

func (u albumUsecase) GetAllAlbums(filters *model.AlbumFilters) ([]*model.Album, error) {
	albumsDB, err := u.albumRepo.GetAllAlbums(filters)
	if err != nil {
		return nil, err
	}
	albums := make([]*model.Album, 0, len(albumsDB))

	for _, albumDB := range albumsDB {
		artist, err := u.artistRepo.GetArtistByID(albumDB.ArtistID)
		if err != nil {
			return nil, err
		}
		album := &model.Album{
			ID:        albumDB.ID,
			Title:     albumDB.Title,
			Thumbnail: albumDB.Thumbnail,
			Artist:    artist,
		}
		albums = append(albums, album)
	}
	return albums, nil
}
