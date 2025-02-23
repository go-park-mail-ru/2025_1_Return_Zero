package main

import (
	"net/http"
	"regexp"
)

func passwordValidation(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile("[A-Z]").MatchString(password)
	hasDigit := regexp.MustCompile("[0-9]").MatchString(password)

	return hasUpper && hasDigit
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
		if user.Username == u.Username {
			http.Error(w, "User already exist", http.StatusConflict)
			return
		}
	}

	if !passwordValidation(u.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
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
	}
	api.users[u.Username] = newUser
	USER_COUNTER++

	api.createSession(w, newUser.ID)
}
