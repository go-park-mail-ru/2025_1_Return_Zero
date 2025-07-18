package repository

type RegisterData struct {
	Username string `sql:"username"`
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

type User struct {
	ID        int64  `sql:"id"`
	Username  string `sql:"username"`
	Email     string `sql:"email"`
	Thumbnail string `sql:"thumbnail"`
	LabelId   int64  `sql:"is_label"`
}

type LoginData struct {
	Username string `sql:"username"`
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

type UserDelete struct {
	Username string `sql:"username"`
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

type ChangeUserData struct {
	Password string `sql:"password"`

	NewUsername string `sql:"new_username"`
	NewEmail    string `sql:"new_email"`
	NewPassword string `sql:"new_password"`
}

type PrivacySettings struct {
	IsPublicPlaylists       bool `sql:"is_public_playlists"`
	IsPublicMinutesListened bool `sql:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool `sql:"is_public_favorite_artists"`
	IsPublicTracksListened  bool `sql:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool `sql:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool `sql:"is_public_artists_listened"`
}

type UserFullData struct {
	Username  string `sql:"username"`
	Thumbnail string `sql:"thumbnail_url"`
	Email     string `sql:"email"`
	Privacy   *PrivacySettings
}

type Label struct {
	ID      int64   `sql:"id"`
	Name    string  `sql:"name"`
	Members []int64 `sql:"members"`
}
