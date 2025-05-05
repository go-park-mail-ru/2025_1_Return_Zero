package usecase

type Playlist struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	Thumbnail string `json:"thumbnail_url"`
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

type LikePlaylistRequest struct {
	UserID     int64
	PlaylistID int64
	IsLike     bool
}

type PlaylistWithIsLiked struct {
	Playlist
	IsLiked bool
}
