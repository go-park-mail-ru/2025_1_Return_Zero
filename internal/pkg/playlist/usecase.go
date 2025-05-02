package playlist

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	CreatePlaylist(ctx context.Context, request *usecaseModel.CreatePlaylistRequest) (*usecaseModel.Playlist, error)
	GetCombinedPlaylistsForCurrentUser(ctx context.Context, userID int64) ([]*usecaseModel.Playlist, error)
	AddTrackToPlaylist(ctx context.Context, request *usecaseModel.AddTrackToPlaylistRequest) error
	RemoveTrackFromPlaylist(ctx context.Context, request *usecaseModel.RemoveTrackFromPlaylistRequest) error
	UpdatePlaylist(ctx context.Context, request *usecaseModel.UpdatePlaylistRequest) (*usecaseModel.Playlist, error)
	GetPlaylistByID(ctx context.Context, playlistID int64) (*usecaseModel.Playlist, error)
	RemovePlaylist(ctx context.Context, request *usecaseModel.RemovePlaylistRequest) error
	GetPlaylistsToAdd(ctx context.Context, request *usecaseModel.GetPlaylistsToAddRequest) ([]*usecaseModel.PlaylistWithIsIncludedTrack, error)
}
