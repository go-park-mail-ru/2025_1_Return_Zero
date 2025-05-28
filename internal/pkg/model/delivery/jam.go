package delivery

type CreateJamRequest struct {
	TrackID  string `json:"track_id" valid:"required"`
	Position int64  `json:"position" valid:"optional"`
}

type CreateJamResponse struct {
	RoomID string `json:"room_id"`
	HostID string `json:"host_id"`
}

type JamMessage struct {
	Type       string            `json:"type"`
	TrackID    string            `json:"track_id,omitempty"`
	Position   int64             `json:"position"`
	Paused     bool              `json:"paused,omitempty"`
	HostID     string            `json:"host_id,omitempty"`
	UserID     string            `json:"user_id,omitempty"`
	Users      []string          `json:"users,omitempty"`
	Loaded     map[string]bool   `json:"loaded,omitempty"`
	Error      string            `json:"error,omitempty"`
	UserImages map[string]string `json:"user_images,omitempty"`
	UserNames  map[string]string `json:"user_names,omitempty"`
}
