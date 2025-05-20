package usecase

type CreateJamRequest struct {
	UserID   string
	TrackID  string
	Position int64
}

type CreateJamResponse struct {
	RoomID string
	HostID string
}

type JoinJamRequest struct {
	RoomID string
	UserID string
}

type JamMessage struct {
	Type     string
	TrackID  string
	Position int64
	Paused   bool
	UserID   string
	HostID   string
	Users    []string
	Loaded   map[string]bool
}
