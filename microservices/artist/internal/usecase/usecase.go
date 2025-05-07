package usecase

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model"
	artistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/errors"
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

func (u *artistUsecase) GetArtistByID(ctx context.Context, id int64, userID int64) (*usecaseModel.ArtistDetailed, error) {
	repoArtist, err := u.artistRepo.GetArtistByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	stats, err := u.artistRepo.GetArtistStats(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.ArtistDetailedFromRepositoryToUsecase(repoArtist, stats), nil
}

func (u *artistUsecase) GetAllArtists(ctx context.Context, filters *usecaseModel.Filters, userID int64) (*usecaseModel.ArtistList, error) {
	repoFilters := model.ArtistFiltersFromUsecaseToRepository(filters)
	repoArtists, err := u.artistRepo.GetAllArtists(ctx, repoFilters, userID)
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

func (u *artistUsecase) GetAlbumIDsByArtistID(ctx context.Context, id int64) ([]int64, error) {
	repoAlbumIDs, err := u.artistRepo.GetAlbumIDsByArtistID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoAlbumIDs, nil
}

func (u *artistUsecase) GetTrackIDsByArtistID(ctx context.Context, id int64) ([]int64, error) {
	repoTrackIDs, err := u.artistRepo.GetTrackIDsByArtistID(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoTrackIDs, nil
}

func (u *artistUsecase) CreateStreamsByArtistIDs(ctx context.Context, data *usecaseModel.ArtistStreamCreateDataList) error {
	repoData := model.ArtistStreamCreateDataFromUsecaseToRepository(data)
	err := u.artistRepo.CreateStreamsByArtistIDs(ctx, repoData)
	if err != nil {
		return err
	}
	return nil
}

func (u *artistUsecase) GetArtistsListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	repoArtistsListened, err := u.artistRepo.GetArtistsListenedByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return repoArtistsListened, nil
}

func (u *artistUsecase) LikeArtist(ctx context.Context, request *usecaseModel.LikeRequest) error {
	repoRequest := model.LikeRequestFromUsecaseToRepository(request)

	exists, err := u.artistRepo.CheckArtistExists(ctx, request.ArtistID)
	if err != nil {
		return err
	}

	if !exists {
		return artistErrors.NewNotFoundError("artist not found")
	}
	if request.IsLike {
		err := u.artistRepo.LikeArtist(ctx, repoRequest)
		if err != nil {
			return err
		}
		return nil
	}

	err = u.artistRepo.UnlikeArtist(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *artistUsecase) GetFavoriteArtists(ctx context.Context, filters *usecaseModel.Filters, userID int64) (*usecaseModel.ArtistList, error) {
	repoFilters := model.ArtistFiltersFromUsecaseToRepository(filters)
	repoArtists, err := u.artistRepo.GetFavoriteArtists(ctx, repoFilters, userID)
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromRepositoryToUsecase(repoArtists), nil
}

func (u *artistUsecase) SearchArtists(ctx context.Context, query string, userID int64) (*usecaseModel.ArtistList, error) {
	repoArtists, err := u.artistRepo.SearchArtists(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromRepositoryToUsecase(repoArtists), nil
}
