package models

import (
	"time"
)

type User struct {
	ID       uint   `json:"-"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" valid:"email"`
}

type Session struct {
	UserID    uint
	ExpiresAt time.Time
}
