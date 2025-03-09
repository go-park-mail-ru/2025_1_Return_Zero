package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// @Summary User login
// @Description Authenticates a user based on provided credentials (either username+password or email+password).
// @Tags auth
// @Accept json
// @Produce json
// @Param request body User true "User credentials"
// @Success 200 {object} UserToFront
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (api *MyHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	var user *User
	var isRegistered bool
	if err := readJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if u.Username == "" {
		for _, rangeUser := range api.users {
			if rangeUser.Email == u.Email {
				user = rangeUser
				isRegistered = true
				break
			}
		}
	} else {
		user, isRegistered = api.users[u.Username]
	}

	if !isRegistered || !checkPasswordHash(u.Password, user.Password) {
		http.Error(w, "Invalid input", http.StatusUnauthorized)
		return
	}
	api.createSession(w, u.ID)

	sendUser := &UserToFront{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}

	if err := writeJSON(w, http.StatusOK, sendUser, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
