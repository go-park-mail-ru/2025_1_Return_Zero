package delivery

// CreatePlaylistRequest
// @Description Create playlist request structure
type CreatePlaylistRequest struct {
	Title     string `form:"title"`
	Thumbnail []byte `form:"thumbnail"`
}

// Playlist
// @Description Playlist structure
type Playlist struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail_url"`
	Username  string `json:"username"`
}

// PlaylistWithIsIncludedTrack
// @Description Playlist with is included track structure
type PlaylistWithIsIncludedTrack struct {
	Playlist
	IsIncluded bool `json:"is_included"`
}

// AddTrackToPlaylistRequest
// @Description Add track to playlist request structure
type AddTrackToPlaylistRequest struct {
	TrackID int64 `json:"track_id"`
}

// UpdatePlaylistRequest
// @Description Update playlist request structure
type UpdatePlaylistRequest struct {
	Title     string `form:"title"`
	Thumbnail []byte `form:"thumbnail"`
}

// PlaylistLikeRequest represents a request to like or unlike an playlist
// @Description A request to like or unlike an playlist. Should be authenticated
type PlaylistLikeRequest struct {
	IsLike bool `json:"value" example:"true" description:"Whether to like or unlike the playlist"`
}

// PlaylistWithIsLiked
// @Description Playlist with is liked structure
type PlaylistWithIsLiked struct {
	Playlist
	IsLiked bool `json:"is_liked"`
}
