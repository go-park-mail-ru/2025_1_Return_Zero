package usecase

import (
	"io"
)

type User struct {
	ID        int64
	Email     string
	Username  string
	Avatar    io.Reader
	Password  string
	AvatarUrl string
	LabelID   int64
}

type ChangeUserData struct {
	Password    string
	NewUsername string
	NewEmail    string
	NewPassword string
}

type PrivacySettings struct {
	Username                string
	IsPublicPlaylists       bool
	IsPublicMinutesListened bool
	IsPublicFavoriteArtists bool
	IsPublicTracksListened  bool
	IsPublicFavoriteTracks  bool
	IsPublicArtistsListened bool
}

type UserPrivacy struct {
	IsPublicPlaylists       bool
	IsPublicMinutesListened bool
	IsPublicFavoriteArtists bool
	IsPublicTracksListened  bool
	IsPublicFavoriteTracks  bool
	IsPublicArtistsListened bool
}

type UserStatistics struct {
	MinutesListened int64
	TracksListened  int64
	ArtistsListened int64
}

type UserFullData struct {
	Username   string
	Email      string
	AvatarUrl  string
	Privacy    *UserPrivacy
	Statistics *UserStatistics
}

type UserChangeSettings struct {
	Privacy     *UserPrivacy
	Password    string
	NewUsername string
	NewEmail    string
	NewPassword string
}

type ChangeSettings struct {
	Password    string
	NewUsername string
	NewEmail    string
	NewPassword string
}

type Label struct {
	Id      int64
	Name    string
	Members []string
}
