package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var USER_COUNTER = 0

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (api *MyHandler) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := readJSON(w, r, &u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, user := range api.users {
		if user.Username == u.Username || user.Email == u.Email {
			http.Error(w, "User already exist", http.StatusConflict)
			return
		}
	}

	hashedPassword, err := hashPassword(u.Password)
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
	if err := writeJSON(w, http.StatusOK, "Successfuly logged in", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
