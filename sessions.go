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
