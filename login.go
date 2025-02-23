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

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

	isRegistred := false
	for _, user := range api.users {
		if user.Username == u.Username {
			isRegistred = true
			break
		}
	}

	if !isRegistred || !checkPasswordHash(u.Password, api.users[u.Username].Password) {
		http.Error(w, "Invalid input", http.StatusUnauthorized)
		return
	}
	api.createSession(w, u.ID)
}
