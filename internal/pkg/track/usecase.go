package track

import (
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllTracks(filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
}
