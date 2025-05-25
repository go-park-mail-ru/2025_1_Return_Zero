package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
)

type Repository interface {
	GetAllTracks(ctx context.Context, filters *repoModel.TrackFilters, userID int64) ([]*repoModel.Track, error)
	GetTrackByID(ctx context.Context, id int64, userID int64) (*repoModel.TrackWithFileKey, error)
	CreateStream(ctx context.Context, stream *repoModel.TrackStreamCreateData) (int64, error)
	GetStreamByID(ctx context.Context, streamID int64) (*repoModel.TrackStream, error)
	UpdateStreamDuration(ctx context.Context, endedStream *repoModel.TrackStreamUpdateData) error
	GetStreamsByUserID(ctx context.Context, userID int64, filters *repoModel.TrackFilters) ([]*repoModel.TrackStream, error)
	GetTracksByIDs(ctx context.Context, ids []int64, userID int64) (map[int64]*repoModel.Track, error)
	GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *repoModel.TrackFilters, userID int64) ([]*repoModel.Track, error)
	GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error)
	GetTracksByAlbumID(ctx context.Context, id int64, userID int64) ([]*repoModel.Track, error)
	GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error)
	GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error)
	LikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error
	CheckTrackExists(ctx context.Context, trackID int64) (bool, error)
	UnlikeTrack(ctx context.Context, likeRequest *repoModel.LikeRequest) error
	GetFavoriteTracks(ctx context.Context, favoriteRequest *repoModel.FavoriteRequest) ([]*repoModel.Track, error)
	SearchTracks(ctx context.Context, query string, userID int64) ([]*repoModel.Track, error)
	GetMostLikedTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error)
	GetMostRecentTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error)
	GetMostListenedLastMonthTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error)
	GetMostLikedLastWeekTracks(ctx context.Context, userID int64) ([]*repoModel.Track, error)
}

type S3Repository interface {
	GetPresignedURL(trackKey string) (string, error)
}
