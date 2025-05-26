package repository

type CreateJamRequest struct {
	UserID   string `json:"user_id"`
	TrackID  string `json:"track_id"`
	Position int64  `json:"position"`
}

type CreateJamResponse struct {
	RoomID string `json:"room_id"`
	HostID string `json:"host_id"`
}

type JoinJamRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

type JamMessage struct {
	Type       string            `json:"type"`
	TrackID    string            `json:"track_id,omitempty"`
	Position   int64             `json:"position"`
	Paused     bool              `json:"paused,omitempty"`
	UserID     string            `json:"user_id,omitempty"`
	HostID     string            `json:"host_id,omitempty"`
	Users      []string          `json:"users,omitempty"`
	Loaded     map[string]bool   `json:"loaded,omitempty"`
	UserImages map[string]string `json:"user_images,omitempty"`
	UserNames  map[string]string `json:"user_names,omitempty"`
}
