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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := readJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check username
	var fieldLogin string
	isRegistred := false
	for _, user := range api.users {
		if user.Username == u.Username {
			isRegistred = true
			fieldLogin = u.Username
			break
		}
	}

	// check email
	if !isRegistred {
		for _, user := range api.users {
			if user.Email == u.Email {
				isRegistred = true
				fieldLogin = u.Email
				break
			}
		}
	}

	if !isRegistred || !checkPasswordHash(u.Password, api.users[fieldLogin].Password) {
		http.Error(w, "Invalid input", http.StatusUnauthorized)
		return
	}
	api.createSession(w, u.ID)
	if err := writeJSON(w, http.StatusOK, "Successfuly logged in", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
