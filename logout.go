package main

import (
	"net/http"
	"time"
)

func (api *MyHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, "No cookie", http.StatusUnauthorized)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, "No cookie", http.StatusUnauthorized)
		return
	}

	api.mu.Lock()
	delete(api.sessions, session.Value)
	api.mu.Unlock()

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	if err := writeJSON(w, http.StatusOK, "Successfuly logged out", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
