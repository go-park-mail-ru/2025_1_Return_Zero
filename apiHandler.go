package main

import (
	"sync"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserToFront struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type MyHandler struct {
	sessions map[string]*Session
	users    map[string]*User
	mu       sync.Mutex
}

func NewMyHandler() *MyHandler {
	h := &MyHandler{
		sessions: make(map[string]*Session),
		users:    make(map[string]*User),
	}
	go h.cleanupSessions()
	return h
}
