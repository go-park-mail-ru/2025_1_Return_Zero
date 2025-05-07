package domain

import (
	"context"
	"io"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/repository"
)

type Repository interface {
	CreatePlaylist(ctx context.Context, playlistCreateRequest *repository.CreatePlaylistRequest) (*repository.Playlist, error)
	GetPlaylistByID(ctx context.Context, id int64) (*repository.Playlist, error)
	GetPlaylistWithIsLikedByID(ctx context.Context, id int64, userID int64) (*repository.PlaylistWithIsLiked, error)
	GetCombinedPlaylistsByUserID(ctx context.Context, userID int64) (*repository.PlaylistList, error)
	TrackExistsInPlaylist(ctx context.Context, playlistID int64, trackID int64) (bool, error)
	AddTrackToPlaylist(ctx context.Context, request *repository.AddTrackToPlaylistRequest) error
	RemoveTrackFromPlaylist(ctx context.Context, request *repository.RemoveTrackFromPlaylistRequest) error
	GetPlaylistTrackIds(ctx context.Context, request *repository.GetPlaylistTrackIdsRequest) ([]int64, error)
	UpdatePlaylist(ctx context.Context, request *repository.UpdatePlaylistRequest) (*repository.Playlist, error)
	RemovePlaylist(ctx context.Context, request *repository.RemovePlaylistRequest) error
	GetPlaylistsToAdd(ctx context.Context, request *repository.GetPlaylistsToAddRequest) (*repository.GetPlaylistsToAddResponse, error)
	UpdatePlaylistsPublisityByUserID(ctx context.Context, request *repository.UpdatePlaylistsPublisityByUserIDRequest) error
	LikePlaylist(ctx context.Context, request *repository.LikePlaylistRequest) error
	UnlikePlaylist(ctx context.Context, request *repository.LikePlaylistRequest) error
	GetProfilePlaylists(ctx context.Context, request *repository.GetProfilePlaylistsRequest) (*repository.GetProfilePlaylistsResponse, error)
	SearchPlaylists(ctx context.Context, request *repository.SearchPlaylistsRequest) (*repository.PlaylistList, error)
	CheckExistsPlaylistAndNotDifferentUser(ctx context.Context, playlistID int64, userID int64) (bool, error)
}

type S3Repository interface {
	UploadThumbnail(ctx context.Context, file io.Reader, key string) (string, error)
}
