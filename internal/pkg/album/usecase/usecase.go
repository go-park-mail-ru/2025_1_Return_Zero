package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/genre"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(albumRepository album.Repository, artistRepository artist.Repository, genreRepository genre.Repository) album.Usecase {
	return albumUsecase{albumRepo: albumRepository, artistRepo: artistRepository, genreRepo: genreRepository}
}

type albumUsecase struct {
	albumRepo  album.Repository
	artistRepo artist.Repository
	genreRepo  genre.Repository
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
		repoArtists, err := u.artistRepo.GetArtistsByAlbumID(repoAlbum.ID)
		if err != nil {
			return nil, err
		}

		albumArtists := make([]*usecaseModel.AlbumArtist, 0, len(repoArtists))
		for _, repoArtist := range repoArtists {
			albumArtists = append(albumArtists, &usecaseModel.AlbumArtist{
				ID:    repoArtist.ID,
				Title: repoArtist.Title,
			})
		}

		albumType := usecaseModel.AlbumType(repoAlbum.Type)

		album := &usecaseModel.Album{
			ID:          repoAlbum.ID,
			Title:       repoAlbum.Title,
			Thumbnail:   repoAlbum.Thumbnail,
			Type:        albumType,
			ReleaseDate: repoAlbum.ReleaseDate,
			Artists:     albumArtists,
		}
		albums = append(albums, album)
	}
	return albums, nil
}

func (u albumUsecase) GetAlbumsByArtistID(artistID int64) ([]*usecaseModel.Album, error) {
	repoAlbums, err := u.albumRepo.GetAlbumsByArtistID(artistID)
	if err != nil {
		return nil, err
	}

	albums := make([]*usecaseModel.Album, 0, len(repoAlbums))

	for _, repoAlbum := range repoAlbums {
		repoArtists, err := u.artistRepo.GetArtistsByAlbumID(repoAlbum.ID)
		if err != nil {
			return nil, err
		}

		albumArtists := make([]*usecaseModel.AlbumArtist, 0, len(repoArtists))
		for _, repoArtist := range repoArtists {
			albumArtists = append(albumArtists, &usecaseModel.AlbumArtist{
				ID:    repoArtist.ID,
				Title: repoArtist.Title,
			})
		}

		albumType := usecaseModel.AlbumType(repoAlbum.Type)

		album := &usecaseModel.Album{
			ID:          repoAlbum.ID,
			Title:       repoAlbum.Title,
			Thumbnail:   repoAlbum.Thumbnail,
			Type:        albumType,
			ReleaseDate: repoAlbum.ReleaseDate,
			Artists:     albumArtists,
		}
		albums = append(albums, album)
	}
	return albums, nil
}
