package track

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTrackByID(ctx context.Context, id int64) (*usecaseModel.TrackDetailed, error)
	GetTracksByArtistID(ctx context.Context, id int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error)
	UpdateStreamDuration(ctx context.Context, endedStream *usecaseModel.TrackStreamUpdateData) error
	GetLastListenedTracks(ctx context.Context, userID int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTracksByAlbumID(ctx context.Context, id int64) ([]*usecaseModel.Track, error)
}
