package repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        int64  `sql:"id"`
	Username  string `sql:"username"`
	Password  string `sql:"password_hash"`
	Email     string `sql:"email"`
	Thumbnail string `sql:"thumbnail_url"`
}

type ChangeUserData struct {
	Password    string
	NewUsername string `sql:"username"`
	NewEmail    string `sql:"email"`
	NewPassword string
}

type UserAndSettings struct {
	Username                string `sql:"username"`
	Thumbnail               string `sql:"thumbnail_url"`
	Email                   string `sql:"email"`
	IsPublicPlaylists       bool   `sql:"is_public_playlists"`
	IsPublicMinutesListened bool   `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `sql:"is_public_artists_listened"`
}

type UserSettings struct {
	Username                string `sql:"username"`
	Thumbnail               string `sql:"thumbnail_url"`
	Email                   string `sql:"email"`
	IsPublicPlaylists       bool   `sql:"is_public_playlists"`
	IsPublicMinutesListened bool   `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `sql:"is_public_artists_listened"`
}

type UserChangeSettings struct {
	IsPublicPlaylists       bool `sql:"is_public_playlists"`
	IsPublicMinutesListened bool `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool `sql:"is_public_artists_listened"`

	Password    string
	NewUsername string `sql:"username"`
	NewEmail    string `sql:"email"`
	NewPassword string
}

type UserStats struct {
	MinutesListened int64 `sql:"minutes_listened"`
	TracksListened  int64 `sql:"tracks_listened"`
	ArtistsListened int64 `sql:"artists_listened"`
}

type UserPrivacySettings struct {
	IsPublicPlaylists       bool `sql:"is_public_playlists"`
	IsPublicMinutesListened bool `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool `sql:"is_public_artists_listened"`
}

type UserFullData struct {
	Username   string `sql:"username"`
	Thumbnail  string `sql:"thumbnail_url"`
	Email      string `sql:"email"`
	Privacy    *UserPrivacySettings
	Statistics *UserStats
}
