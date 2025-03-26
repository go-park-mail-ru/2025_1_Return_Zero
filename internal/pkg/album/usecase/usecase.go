package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func NewUsecase(albumRepository album.Repository) album.Usecase {
	return albumUsecase{repo: albumRepository}
}

type albumUsecase struct {
	repo album.Repository
}

func (u albumUsecase) GetAllAlbums(filters *model.AlbumFilters) ([]*model.Album, error) {
	return u.repo.GetAllAlbums(filters)
}
