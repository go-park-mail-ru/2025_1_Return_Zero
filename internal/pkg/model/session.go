package model

import (
	"time"
)

type Session struct {
	UserID    uint
	ExpiresAt time.Time
}