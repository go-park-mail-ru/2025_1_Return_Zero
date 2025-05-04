package domain

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/usecase"
)

type Usecase interface {
	CreatePlaylist(ctx context.Context, playlist *usecase.CreatePlaylistRequest) (*usecase.Playlist, error)
	UploadPlaylistThumbnail(ctx context.Context, playlist *usecase.UploadPlaylistThumbnailRequest) (string, error)
	GetCombinedPlaylistsByUserID(ctx context.Context, request *usecase.GetCombinedPlaylistsByUserIDRequest) (*usecase.PlaylistList, error)
	AddTrackToPlaylist(ctx context.Context, request *usecase.AddTrackToPlaylistRequest) error
	RemoveTrackFromPlaylist(ctx context.Context, request *usecase.RemoveTrackFromPlaylistRequest) error
	GetPlaylistTrackIds(ctx context.Context, request *usecase.GetPlaylistTrackIdsRequest) ([]int64, error)
	UpdatePlaylist(ctx context.Context, request *usecase.UpdatePlaylistRequest) (*usecase.Playlist, error)
	GetPlaylistByID(ctx context.Context, request *usecase.GetPlaylistByIDRequest) (*usecase.PlaylistWithIsLiked, error)
	RemovePlaylist(ctx context.Context, request *usecase.RemovePlaylistRequest) error
	GetPlaylistsToAdd(ctx context.Context, request *usecase.GetPlaylistsToAddRequest) (*usecase.GetPlaylistsToAddResponse, error)
	UpdatePlaylistsPublisityByUserID(ctx context.Context, request *usecase.UpdatePlaylistsPublisityByUserIDRequest) error
	LikePlaylist(ctx context.Context, request *usecase.LikePlaylistRequest) error
	GetProfilePlaylists(ctx context.Context, request *usecase.GetProfilePlaylistsRequest) (*usecase.GetProfilePlaylistsResponse, error)
	SearchPlaylists(ctx context.Context, request *usecase.SearchPlaylistsRequest) (*usecase.PlaylistList, error)
}
