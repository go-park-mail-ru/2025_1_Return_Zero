package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model"
	albumErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/errors"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
)

type AlbumUsecase struct {
	albumRepository domain.Repository
}

func NewAlbumUsecase(albumRepository domain.Repository) domain.Usecase {
	return &AlbumUsecase{albumRepository: albumRepository}
}

func (u *AlbumUsecase) GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters, userID int64) ([]*usecaseModel.Album, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	albums, err := u.albumRepository.GetAllAlbums(ctx, repoFilters, userID)
	if err != nil {
		return nil, err
	}

	return model.AlbumListFromRepositoryToUsecase(albums), nil
}

func (u *AlbumUsecase) GetAlbumByID(ctx context.Context, id int64, userID int64) (*usecaseModel.Album, error) {
	album, err := u.albumRepository.GetAlbumByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return model.AlbumFromRepositoryToUsecase(album), nil
}

func (u *AlbumUsecase) GetAlbumTitleByID(ctx context.Context, id int64) (string, error) {
	title, err := u.albumRepository.GetAlbumTitleByID(ctx, id)
	if err != nil {
		return "", err
	}
	return title, nil
}

func (u *AlbumUsecase) GetAlbumTitleByIDs(ctx context.Context, ids []int64) (*usecaseModel.AlbumTitleMap, error) {
	repoIds := make([]int64, len(ids))
	copy(repoIds, ids)

	albumTitles, err := u.albumRepository.GetAlbumTitleByIDs(ctx, repoIds)
	if err != nil {
		return nil, err
	}

	return model.AlbumTitleMapFromRepositoryToUsecase(albumTitles), nil
}

func (u *AlbumUsecase) GetAlbumsByIDs(ctx context.Context, ids []int64, userID int64) ([]*usecaseModel.Album, error) {
	repoIds := make([]int64, len(ids))
	copy(repoIds, ids)

	albums, err := u.albumRepository.GetAlbumsByIDs(ctx, repoIds, userID)
	if err != nil {
		return nil, err
	}

	return model.AlbumListFromRepositoryToUsecase(albums), nil
}

func (u *AlbumUsecase) CreateStream(ctx context.Context, albumID int64, userID int64) error {
	return u.albumRepository.CreateStream(ctx, albumID, userID)
}

func (u *AlbumUsecase) LikeAlbum(ctx context.Context, request *usecaseModel.LikeRequest) error {
	repoRequest := model.LikeRequestFromUsecaseToRepository(request)

	exists, err := u.albumRepository.CheckAlbumExists(ctx, request.AlbumID)
	if err != nil {
		return err
	}

	if !exists {
		return albumErrors.NewNotFoundError("album not found")
	}
	if request.IsLike {
		err := u.albumRepository.LikeAlbum(ctx, repoRequest)
		if err != nil {
			return err
		}
		return nil
	}

	err = u.albumRepository.UnlikeAlbum(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}
