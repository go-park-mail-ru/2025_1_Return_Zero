package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var USER_COUNTER = 0

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// @Summary Register a new user
// @Description Creates a new user if the username and email are unique. Hashes the password, saves the user, and creates a session.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body User true "User data for registration"
// @Success 200 {object} UserToFront "User successfully registered"
// @Failure 400 {string} string "Invalid request"
// @Failure 409 {string} string "User already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /signup [post]
func (api *MyHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := readJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, isRegistred := api.users[u.Username]
	if isRegistred {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}
	for _, user := range api.users {
		if user.Email == u.Email {
			http.Error(w, "User already exist", http.StatusConflict)
			return
		}
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	newUser := &User{
		ID:       uint(USER_COUNTER),
		Username: u.Username,
		Password: hashedPassword,
		Email:    u.Email,
	}
	api.users[u.Username] = newUser
	USER_COUNTER++

	api.createSession(w, newUser.ID)
	sendUser := &UserToFront{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	if err := writeJSON(w, http.StatusOK, sendUser, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
