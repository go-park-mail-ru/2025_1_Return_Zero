package usecase

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
)

func NewArtistUsecase(artistRepository domain.Repository) domain.Usecase {
	return &artistUsecase{
		artistRepo: artistRepository,
	}
}

type artistUsecase struct {
	artistRepo domain.Repository
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

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.ArtistFilters) (*usecaseModel.ArtistList, error) {
	repoFilters := model.ArtistFiltersFromUsecaseToRepository(filters)
	repoArtists, err := u.artistRepo.GetAllArtists(ctx, repoFilters)
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromRepositoryToUsecase(repoArtists), nil
}

func (u *artistUsecase) GetArtistTitleByID(ctx context.Context, id int64) (string, error) {
	repoTitle, err := u.artistRepo.GetArtistTitleByID(ctx, id)
	if err != nil {
		return "", err
	}
	return repoTitle, nil
}

func (u *artistUsecase) GetArtistsByTrackID(ctx context.Context, id int64) (*usecaseModel.ArtistWithRoleList, error) {
	repoArtists, err := u.artistRepo.GetArtistsByTrackID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithRoleListFromRepositoryToUsecase(repoArtists), nil
}

func (u *artistUsecase) GetArtistsByTrackIDs(ctx context.Context, ids []int64) (*usecaseModel.ArtistWithRoleMap, error) {
	repoArtists, err := u.artistRepo.GetArtistsByTrackIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithRoleMapFromRepositoryToUsecase(repoArtists), nil
}

func (u *artistUsecase) GetArtistsByAlbumID(ctx context.Context, id int64) (*usecaseModel.ArtistWithTitleList, error) {
	repoArtists, err := u.artistRepo.GetArtistsByAlbumID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithTitleListFromRepositoryToUsecase(repoArtists), nil
}

func (u *artistUsecase) GetArtistsByAlbumIDs(ctx context.Context, ids []int64) (*usecaseModel.ArtistWithTitleMap, error) {
	repoArtists, err := u.artistRepo.GetArtistsByAlbumIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithTitleMapFromRepositoryToUsecase(repoArtists), nil
}
