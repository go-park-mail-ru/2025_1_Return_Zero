package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(artistRepository artist.Repository) artist.Usecase {
	return artistUsecase{
		artistRepo: artistRepository,
	}
}

type artistUsecase struct {
	artistRepo artist.Repository
}

func (u artistUsecase) GetArtistByID(id int64) (*usecaseModel.ArtistDetailed, error) {
	repoArtist, err := u.artistRepo.GetArtistByID(id)
	if err != nil {
		return nil, err
	}

	stats, err := u.artistRepo.GetArtistStats(id)
	if err != nil {
		return nil, err
	}

	return &usecaseModel.ArtistDetailed{
		Artist: usecaseModel.Artist{
			ID:          repoArtist.ID,
			Title:       repoArtist.Title,
			Thumbnail:   repoArtist.Thumbnail,
			Description: repoArtist.Description,
		},
		Listeners: stats.ListenersCount,
		Favorites: stats.FavoritesCount,
	}, nil
}

func (u artistUsecase) GetAllArtists(filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	repoFilters := &repoModel.ArtistFilters{
		Pagination: &repoModel.Pagination{
			Offset: filters.Pagination.Offset,
			Limit:  filters.Pagination.Limit,
		},
	}

	repoArtists, err := u.artistRepo.GetAllArtists(repoFilters)
	if err != nil {
		return nil, err
	}

	artists := make([]*usecaseModel.Artist, 0, len(repoArtists))
	for _, repoArtist := range repoArtists {
		artists = append(artists, &usecaseModel.Artist{
			ID:          repoArtist.ID,
			Title:       repoArtist.Title,
			Thumbnail:   repoArtist.Thumbnail,
			Description: repoArtist.Description,
		})
	}
	return artists, nil
}
