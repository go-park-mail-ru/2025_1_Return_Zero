package track

import (
	"errors"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

var (
	ErrStreamPermissionDenied = errors.New("user does not have permission to update this stream")
)

type Usecase interface {
	GetAllTracks(filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTrackByID(id int64) (*usecaseModel.TrackDetailed, error)
	GetTracksByArtistID(id int64) ([]*usecaseModel.Track, error)
	CreateStream(stream *usecaseModel.TrackStreamCreateData) (int64, error)
	UpdateStreamDuration(endedStream *usecaseModel.TrackStreamUpdateData) error
}
