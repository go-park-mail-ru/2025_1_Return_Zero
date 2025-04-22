package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
)

type AlbumUsecase struct {
	albumRepository domain.Repository
}

func NewAlbumUsecase(albumRepository domain.Repository) domain.Usecase {
	return &AlbumUsecase{albumRepository: albumRepository}
}

func (u *AlbumUsecase) GetAllAlbums(ctx context.Context, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	albums, err := u.albumRepository.GetAllAlbums(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	return model.AlbumListFromRepositoryToUsecase(albums), nil
}

func (u *AlbumUsecase) GetAlbumByID(ctx context.Context, id int64) (*usecaseModel.Album, error) {
	album, err := u.albumRepository.GetAlbumByID(ctx, id)
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

func (u *AlbumUsecase) GetAlbumsByIDs(ctx context.Context, ids []int64) ([]*usecaseModel.Album, error) {
	repoIds := make([]int64, len(ids))
	copy(repoIds, ids)

	albums, err := u.albumRepository.GetAlbumsByIDs(ctx, repoIds)
	if err != nil {
		return nil, err
	}

	albumsInOrder := make([]*usecaseModel.Album, len(ids))
	for i, id := range ids {
		albumsInOrder[i] = model.AlbumFromRepositoryToUsecase(albums[id])
	}

	return albumsInOrder, nil
}
