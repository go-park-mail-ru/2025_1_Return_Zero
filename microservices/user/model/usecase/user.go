package usecase

type RegisterData struct {
	Username string
	Email    string
	Password string
}

type UserFront struct {
	Username  string
	Email     string
	Thumbnail string
	Id        int64
}

type LoginData struct {
	Username string
	Email    string
	Password string
}

type UserDelete struct {
	Username string
	Email    string
	Password string
}

type ChangeUserData struct {
	Password    string
	NewUsername string
	NewEmail    string
	NewPassword string
}

type PrivacySettings struct {
	IsPublicPlaylists       bool
	IsPublicMinutesListened bool
	IsPublicFavoriteArtists bool
	IsPublicTracksListened  bool
	IsPublicFavoriteTracks  bool
	IsPublicArtistsListened bool
}

type UserFullData struct {
	Username  string
	Thumbnail string
	Email     string
	Privacy   *PrivacySettings
}

type CreateLabelRequest struct {
	Name   string

}

type Label struct {
	ID      int64
	Name    string
	Members []string
}
