package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(albumRepository album.Repository, artistRepository artist.Repository) album.Usecase {
	return albumUsecase{albumRepo: albumRepository, artistRepo: artistRepository}
}

type albumUsecase struct {
	albumRepo  album.Repository
	artistRepo artist.Repository
}

func (u albumUsecase) GetAllAlbums(filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	repoFilters := &repoModel.AlbumFilters{
		Pagination: &repoModel.Pagination{
			Offset: filters.Pagination.Offset,
			Limit:  filters.Pagination.Limit,
		},
	}
	repoAlbums, err := u.albumRepo.GetAllAlbums(repoFilters)
	if err != nil {
		return nil, err
	}
	albums := make([]*usecaseModel.Album, 0, len(repoAlbums))

	for _, repoAlbum := range repoAlbums {
		repoArtist, err := u.artistRepo.GetArtistByID(repoAlbum.ArtistID)
		if err != nil {
			return nil, err
		}
		album := &usecaseModel.Album{
			ID:        repoAlbum.ID,
			Title:     repoAlbum.Title,
			Thumbnail: repoAlbum.Thumbnail,
			Artist: usecaseModel.AlbumArtist{
				ID:    repoArtist.ID,
				Title: repoArtist.Title,
			},
		}
		albums = append(albums, album)
	}
	return albums, nil
}
