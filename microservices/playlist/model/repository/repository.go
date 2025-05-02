package repository

type Playlist struct {
	ID        int64  `sql:"id"`
	Title     string `sql:"title"`
	Thumbnail string `sql:"thumbnail"`
	UserID    int64  `sql:"user_id"`
	IsPublic  bool   `sql:"is_public"`
}

type PlaylistList struct {
	Playlists []*Playlist
}

type CreatePlaylistRequest struct {
	Title     string `sql:"title"`
	UserID    int64  `sql:"user_id"`
	Thumbnail string `sql:"thumbnail"`
}

type GetToAddByUserIdRequest struct {
	UserID     int64
	PlaylistID int64
	TrackID    int64
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
