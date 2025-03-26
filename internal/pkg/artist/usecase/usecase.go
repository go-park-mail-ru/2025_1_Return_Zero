package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func NewUsecase(artistRepository artist.Repository) artist.Usecase {
	return artistUsecase{repo: artistRepository}
}

type artistUsecase struct {
	repo artist.Repository
}

func (u artistUsecase) GetAllArtists(filters *model.ArtistFilters) ([]*model.Artist, error) {
	return u.repo.GetAllArtists(filters)
}
