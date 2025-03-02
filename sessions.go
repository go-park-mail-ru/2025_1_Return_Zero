package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

type Session struct {
	UserID    uint
	ExpiresAt time.Time
}

func (api *MyHandler) cleanupSessions() {
	for {
		time.Sleep(time.Hour)
		api.mu.Lock()
		now := time.Now()
		for sid, session := range api.sessions {
			if now.After(session.ExpiresAt) {
				delete(api.sessions, sid)
			}
		}
		api.mu.Unlock()
	}
}

func (api *MyHandler) createSession(w http.ResponseWriter, ID uint) {
	SID := generateSessionID()
	expiration := time.Now().Add(24 * time.Hour)
	api.mu.Lock()
	api.sessions[SID] = &Session{
		UserID:    ID,
		ExpiresAt: expiration,
	}
	api.mu.Unlock()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// checkSession verifies the user's session based on the "session_id" cookie.
//
// @Summary Check user session
// @Description Validates the session by checking the "session_id" cookie and retrieving user information.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} UserToFront "User session is valid"
// @Failure 401 {string} string "Invalid cookie or unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /session/check [get]
func (api *MyHandler) checkSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Invalid cookie", http.StatusUnauthorized)
		return
	}

	session, exists := api.sessions[cookie.Value]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var username string
	for _, user := range api.users {
		if user.ID == session.UserID {
			username = user.Username
		}
	}

	user, exists := api.users[username]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := &UserToFront{
		ID:       user.ID,
		Username: username,
		Email:    user.Email,
	}
	if err = writeJSON(w, http.StatusOK, response, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
