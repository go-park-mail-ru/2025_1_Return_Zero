package domain

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
)

type Usecase interface {
	GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters, userID int64) ([]*usecaseModel.Track, error)
	GetTrackByID(ctx context.Context, id int64, userID int64) (*usecaseModel.TrackDetailed, error)
	CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error)
	UpdateStreamDuration(ctx context.Context, endedStream *usecaseModel.TrackStreamUpdateData) error
	GetLastListenedTracks(ctx context.Context, userID int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error)
	GetTracksByIDs(ctx context.Context, ids []int64, userID int64) ([]*usecaseModel.Track, error)
	GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *usecaseModel.TrackFilters, userID int64) ([]*usecaseModel.Track, error)
	GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error)
	GetTracksByAlbumID(ctx context.Context, id int64, userID int64) ([]*usecaseModel.Track, error)
	GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error)
	GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error)
	LikeTrack(ctx context.Context, likeRequest *usecaseModel.LikeRequest) error
	GetFavoriteTracks(ctx context.Context, favoriteRequest *usecaseModel.FavoriteRequest) ([]*usecaseModel.Track, error)
	SearchTracks(ctx context.Context, query string, userID int64) ([]*usecaseModel.Track, error)
	AddTracksToAlbum(ctx context.Context, tracksList *usecaseModel.TracksListWithAlbumID) ([]int64, error)
	DeleteTracksByAlbumID(ctx context.Context, albumID int64) error
	GetMostLikedTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error)
	GetMostRecentTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error)
	GetMostListenedLastMonthTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error)
	GetMostLikedLastWeekTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error)
}
