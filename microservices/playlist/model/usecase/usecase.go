package usecase

import (
	"io"
)

type Playlist struct {
	ID        int64
	Title     string
	UserID    int64
	Thumbnail string
}

type PlaylistList struct {
	Playlists []*Playlist
}

type CreatePlaylistRequest struct {
	Title     string
	UserID    int64
	Thumbnail string
}

type UploadPlaylistThumbnailRequest struct {
	Title     string
	Thumbnail io.Reader
}

type GetCombinedPlaylistsByUserIDRequest struct {
	UserID int64
}

type AddTrackToPlaylistRequest struct {
	UserID     int64
	PlaylistID int64
	TrackID    int64
}

type RemoveTrackFromPlaylistRequest struct {
	UserID     int64
	PlaylistID int64
	TrackID    int64
}

type GetPlaylistTrackIdsRequest struct {
	UserID     int64
	PlaylistID int64
}

type UpdatePlaylistRequest struct {
	UserID     int64
	PlaylistID int64
	Title      string
	Thumbnail  string
}

type GetPlaylistByIDRequest struct {
	UserID     int64
	PlaylistID int64
}

type RemovePlaylistRequest struct {
	UserID     int64
	PlaylistID int64
}

type GetPlaylistsToAddRequest struct {
	UserID  int64
	TrackID int64
}

type PlaylistWithIsIncludedTrack struct {
	Playlist   *Playlist
	IsIncluded bool
}

type GetPlaylistsToAddResponse struct {
	Playlists []*PlaylistWithIsIncludedTrack
}
