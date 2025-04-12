package delivery

type UserToFront struct {
	ID       int64  `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}

type UserDelete struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,alphanum,stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

// RegisterData represents user registration information
// @Description User registration data requiring username (3-20 characters), password (4-25 characters), and valid email (5-30 characters)
type RegisterData struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,alphanum,stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

// LoginData represents user login credentials
// @Description User login data. Either username or email must be provided along with required password (4-25 characters)
type LoginData struct {
	Username string `json:"username" valid:"alphanum,stringlength(3|20)"`
	Password string `json:"password" valid:"required,stringlength(4|25)"`
	Email    string `json:"email" valid:"email,stringlength(5|30)"`
}

// ChangeUserData contains user update information
// @Description Data for user profile update. Requires current credentials and allows new username (3-20 alphanum), new email (5-30 valid format), and new password (4-25 characters)
type ChangeUserData struct {
	Username string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
	Password string `json:"password" valid:"stringlength(4|25)"`

	NewUsername string `json:"new_username" valid:"alphanum,stringlength(3|20)"`
	NewEmail    string `json:"new_email" valid:"email,stringlength(5|30)"`
	NewPassword string `json:"new_password" valid:"stringlength(4|25)"`
}

type PrivacySettings struct {
	Username                string `json:"username" valid:"required,alphanum,stringlength(3|20)"`
	IsPublicPlaylists       bool   `json:"is_public_playlists"`
	IsPublicMinutesListened bool   `json:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `json:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `json:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `json:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `json:"is_public_artists_listened"`
}


type UserAndSettings struct {
	Username                string `json:"username"`
	AvatarUrl               string `json:"avatar_url"`
	IsPublicPlaylists       bool   `json:"is_public_playlists"`
	IsPublicMinutesListened bool   `json:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool   `json:"is_public_favorite_artists"`
	IsPublicTracksListened  bool   `json:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool   `json:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool   `json:"is_public_artists_listened"`
}