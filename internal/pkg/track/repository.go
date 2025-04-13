package track

import (
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrStreamNotFound               = errors.New("stream not found")
	ErrFailedToUpdateStreamDuration = errors.New("failed to update stream duration")
	ErrTrackNotFound                = errors.New("track not found")
)

type Repository interface {
	GetAllTracks(filters *repoModel.TrackFilters) ([]*repoModel.Track, error)
	GetTrackByID(id int64) (*repoModel.TrackWithFileKey, error)
	GetTracksByArtistID(id int64) ([]*repoModel.Track, error)
	CreateStream(stream *repoModel.TrackStreamCreateData) (int64, error)
	GetStreamByID(streamID int64) (*repoModel.TrackStream, error)
	UpdateStreamDuration(endedStream *repoModel.TrackStreamUpdateData) error
}
