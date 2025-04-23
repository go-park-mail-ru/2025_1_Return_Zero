package domain

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
)

type Usecase interface {
	GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTrackByID(ctx context.Context, id int64) (*usecaseModel.TrackDetailed, error)
	CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error)
	UpdateStreamDuration(ctx context.Context, endedStream *usecaseModel.TrackStreamUpdateData) error
	GetLastListenedTracks(ctx context.Context, userID int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTracksByIDs(ctx context.Context, ids []int64) ([]*usecaseModel.Track, error)
	GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
}
