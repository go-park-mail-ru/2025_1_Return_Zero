package delivery

// UserToFront represents user data
// @Description User data
type UserToFront struct {
	ID       int64  `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
	IsLabel  bool   `json:"is_label"`
}

type UserDelete struct {
	Username string `json:"username" valid:"required,matches(^[a-zA-Z0-9_]+$),stringlength(3|20)"`
	Password string `json:"password" valid:"required,matches(^[a-zA-Z0-9_]+$),stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

// RegisterData represents user registration information
// @Description User registration data requiring username (3-20 characters), password (4-25 characters), and valid email (5-30 characters)
type RegisterData struct {
	Username string `json:"username" valid:"required,matches(^[a-zA-Z0-9_]+$),stringlength(3|20)"`
	Password string `json:"password" valid:"required,matches(^[a-zA-Z0-9_]+$),stringlength(4|25)"`
	Email    string `json:"email" valid:"required,email,stringlength(5|30)"`
}

// LoginData represents user login credentials
// @Description User login data. Either username or email must be provided along with required password (4-25 characters)
type LoginData struct {
	Username string `json:"username" valid:"matches(^[a-zA-Z0-9_]+$),stringlength(3|20)"`
	Password string `json:"password" valid:"required,matches(^[a-zA-Z0-9_]+$),stringlength(4|25)"`
	Email    string `json:"email" valid:"email,stringlength(5|30)"`
}

type Privacy struct {
	IsPublicPlaylists       bool `json:"is_public_playlists"`
	IsPublicMinutesListened bool `json:"is_public_minutes_listened"`
	IsPublicFavoriteArtists bool `json:"is_public_favorite_artists"`
	IsPublicTracksListened  bool `json:"is_public_tracks_listened"`
	IsPublicFavoriteTracks  bool `json:"is_public_favorite_tracks"`
	IsPublicArtistsListened bool `json:"is_public_artists_listened"`
}

type Statistics struct {
	MinutesListened int64 `json:"minutes_listened"`
	TracksListened  int64 `json:"tracks_listened"`
	ArtistsListened int64 `json:"artists_listened"`
}

type UserFullData struct {
	Username   string      `json:"username"`
	Email      string      `json:"email,omitempty"`
	AvatarUrl  string      `json:"avatar_url"`
	Privacy    *Privacy    `json:"privacy,omitempty"`
	Statistics *Statistics `json:"statistics,omitempty"`
}

type UserChangeSettings struct {
	Privacy     *Privacy `json:"privacy"`
	Password    string   `json:"password" valid:"stringlength(4|25)"`
	NewUsername string   `json:"new_username" valid:"matches(^[a-zA-Z0-9_]+$),stringlength(3|20)"`
	NewEmail    string   `json:"new_email" valid:"email,stringlength(5|30)"`
	NewPassword string   `json:"new_password" valid:"matches(^[a-zA-Z0-9_]+$),stringlength(4|25)"`
}

type AvatarURL struct {
	AvatarUrl string `json:"avatar_url"`
}

type Label struct {
	Id        int64    `json:"id,omitempty"`
	Usernames []string `json:"usernames"`
	LabelName string   `json:"label_name"`
}
