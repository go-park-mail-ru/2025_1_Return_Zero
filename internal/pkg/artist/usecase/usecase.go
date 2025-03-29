package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(artistRepository artist.Repository) artist.Usecase {
	return artistUsecase{repo: artistRepository}
}

type artistUsecase struct {
	repo artist.Repository
}

func (u artistUsecase) GetAllArtists(filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	repoFilters := &repoModel.ArtistFilters{
		Pagination: &repoModel.Pagination{
			Offset: filters.Pagination.Offset,
			Limit:  filters.Pagination.Limit,
		},
	}

	repoArtists, err := u.repo.GetAllArtists(repoFilters)
	if err != nil {
		return nil, err
	}

	artists := make([]*usecaseModel.Artist, 0, len(repoArtists))
	for _, repoArtist := range repoArtists {
		artists = append(artists, &usecaseModel.Artist{
			ID:          repoArtist.ID,
			Title:       repoArtist.Title,
			Thumbnail:   repoArtist.Thumbnail,
			Listeners:   repoArtist.Listeners,
			Favorites:   repoArtist.Favorites,
			Description: repoArtist.Description,
		})
	}
	return artists, nil
}
