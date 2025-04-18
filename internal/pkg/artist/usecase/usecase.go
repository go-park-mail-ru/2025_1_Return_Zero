package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(artistRepository artist.Repository) artist.Usecase {
	return &artistUsecase{
		artistRepo: artistRepository,
	}
}

type artistUsecase struct {
	artistRepo artist.Repository
}

func (u *artistUsecase) GetArtistByID(ctx context.Context, id int64) (*usecaseModel.ArtistDetailed, error) {
	repoArtist, err := u.artistRepo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stats, err := u.artistRepo.GetArtistStats(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.ArtistDetailedFromRepositoryToUsecase(repoArtist, stats), nil
}

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error) {
	repoFilters := &repoModel.ArtistFilters{
		Pagination: model.PaginationFromUsecaseToRepository(filters.Pagination),
	}

	repoArtists, err := u.artistRepo.GetAllArtists(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	return model.ArtistsFromRepositoryToUsecase(repoArtists), nil
}
