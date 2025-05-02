package usecase

type Playlist struct {
	ID        int64
	Title     string
	Username  string
	Thumbnail string
}

type PlaylistWithIsIncludedTrack struct {
	Playlist
	IsIncluded bool
}

type CreatePlaylistRequest struct {
	UserID    int64
	Title     string
	Thumbnail []byte
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

type UpdatePlaylistRequest struct {
	UserID     int64
	PlaylistID int64
	Title      string
	Thumbnail  []byte
}

type RemovePlaylistRequest struct {
	UserID     int64
	PlaylistID int64
}

type GetPlaylistsToAddRequest struct {
	UserID  int64
	TrackID int64
}
