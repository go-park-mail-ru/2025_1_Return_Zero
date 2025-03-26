package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

func NewUsecase(trackRepository track.Repository) track.Usecase {
	return trackUsecase{repo: trackRepository}
}

type trackUsecase struct {
	repo track.Repository
}

func (u trackUsecase) GetAllTracks(filters *model.TrackFilters) ([]*model.Track, error) {
	return u.repo.GetAllTracks(filters)
}
