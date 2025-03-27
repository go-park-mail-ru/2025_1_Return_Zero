package artist

import (
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	GetAllArtists(filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error)
	GetArtistByID(id uint) (*repoModel.Artist, error)
}
