package artist

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Repository interface {
	GetAllArtists(filters *model.ArtistFilters) ([]*model.Artist, error)
}
