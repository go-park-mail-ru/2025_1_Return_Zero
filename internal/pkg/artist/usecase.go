package artist

import (
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllArtists(filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error)
}
