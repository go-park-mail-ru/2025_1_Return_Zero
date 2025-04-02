package track

import (
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllTracks(filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTrackByID(id int64) (*usecaseModel.TrackDetailed, error)
	GetTracksByArtistID(id int64) ([]*usecaseModel.Track, error)
}
