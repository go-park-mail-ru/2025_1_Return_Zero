package track

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Repository interface {
	GetAllTracks(filters *model.TrackFilters) ([]*model.Track, error)
}
