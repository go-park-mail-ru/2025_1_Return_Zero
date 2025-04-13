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
	Username string `sql:"username"`
	Email    string `sql:"email"`
	Password string `sql:"password_hash"`

	NewUsername string `sql:"username"`
	NewEmail    string `sql:"email"`
	NewPassword string `sql:"password_hash"`
}

type PrivacySettings struct {
	Username                string `sql:"username"`
	IsPublicPlaylists       bool   `sql:"is_public_playlists"`
	IsPublicMinutesListened bool   `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `sql:"is_public_artists_listened"`
}

type UserAndSettings struct {
	Username                string `sql:"username"`
	Thumbnail               string `sql:"thumbnail_url"`
	IsPublicPlaylists       bool   `sql:"is_public_playlists"`
	IsPublicMinutesListened bool   `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `sql:"is_public_artists_listened"`
}
