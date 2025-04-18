package track

import (
	"context"
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrStreamNotFound               = errors.New("stream not found")
	ErrFailedToUpdateStreamDuration = errors.New("failed to update stream duration")
	ErrTrackNotFound                = errors.New("track not found")
)

type Repository interface {
	GetAllTracks(ctx context.Context, filters *repoModel.TrackFilters) ([]*repoModel.Track, error)
	GetTrackByID(ctx context.Context, id int64) (*repoModel.TrackWithFileKey, error)
	GetTracksByArtistID(ctx context.Context, id int64, filters *repoModel.TrackFilters) ([]*repoModel.Track, error)
	CreateStream(ctx context.Context, stream *repoModel.TrackStreamCreateData) (int64, error)
	GetStreamByID(ctx context.Context, streamID int64) (*repoModel.TrackStream, error)
	UpdateStreamDuration(ctx context.Context, endedStream *repoModel.TrackStreamUpdateData) error
	GetStreamsByUserID(ctx context.Context, userID int64, filters *repoModel.TrackFilters) ([]*repoModel.TrackStream, error)
	GetTracksByIDs(ctx context.Context, ids []int64) (map[int64]*repoModel.Track, error)
}
