package main

import (
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*User
	mu       sync.Mutex
}

func NewMyHandler() *MyHandler {
	h := &MyHandler{
		sessions: make(map[string]uint, 10),
		users:    make(map[string]*User),
	}
	go h.cleanupSessions()
	return h
}

func (api *MyHandler) createSession(w http.ResponseWriter, ID uint) {
	SID := generateSessionID()
	api.mu.Lock()
	api.sessions[SID] = ID
	api.mu.Unlock()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}
