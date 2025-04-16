package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
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

func (u albumUsecase) GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	repoFilters := &repoModel.AlbumFilters{
		Pagination: model.PaginationFromUsecaseToRepository(filters.Pagination),
	}

	repoAlbums, err := u.albumRepo.GetAllAlbums(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	albumIDs := make([]int64, 0, len(repoAlbums))
	for _, repoAlbum := range repoAlbums {
		albumIDs = append(albumIDs, repoAlbum.ID)
	}

	repoArtists, err := u.artistRepo.GetArtistsByAlbumIDs(ctx, albumIDs)
	if err != nil {
		return nil, err
	}

	albums := make([]*usecaseModel.Album, 0, len(repoAlbums))
	for _, repoAlbum := range repoAlbums {
		usecaseAlbum := model.AlbumFromRepositoryToUsecase(repoAlbum, repoArtists[repoAlbum.ID])
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}

func (u albumUsecase) GetAlbumsByArtistID(ctx context.Context, artistID int64) ([]*usecaseModel.Album, error) {
	repoAlbums, err := u.albumRepo.GetAlbumsByArtistID(ctx, artistID)
	if err != nil {
		return nil, err
	}

	albumIDs := make([]int64, 0, len(repoAlbums))
	for _, repoAlbum := range repoAlbums {
		albumIDs = append(albumIDs, repoAlbum.ID)
	}

	repoArtists, err := u.artistRepo.GetArtistsByAlbumIDs(ctx, albumIDs)
	if err != nil {
		return nil, err
	}

	albums := make([]*usecaseModel.Album, 0, len(repoAlbums))
	for _, repoAlbum := range repoAlbums {
		usecaseAlbum := model.AlbumFromRepositoryToUsecase(repoAlbum, repoArtists[repoAlbum.ID])
		albums = append(albums, usecaseAlbum)
	}
	return albums, nil
}
