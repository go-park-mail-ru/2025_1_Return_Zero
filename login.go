package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// @Summary Log in a user
// @Description Authenticate a user using their username/email and password
// @Tags login
// @Accept json
// @Produce json
// @Param user body User true "User credentials (username/email and password)"
// @Success 200 {string} string "Successfully logged in"
// @Failure 400 {string} string "Bad request - invalid input"
// @Failure 401 {string} string "Unauthorized - invalid credentials"
// @Failure 405 {string} string "Method not allowed"
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
}
