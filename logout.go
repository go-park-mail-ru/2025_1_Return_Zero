package main

import (
	"net/http"
	"time"
)

// @Summary Log out a user
// @Description Terminate the user's session and clear the session cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Successfully logged out"
// @Failure 401 {string} string "Unauthorized - no valid session"
// @Failure 500 {string} string "Internal server error"
// @Router /logout [post]
func (api *MyHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
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
