package main

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func (api *MyHandler) cleanupSessions() {
	for {
		time.Sleep(24 * time.Hour)
		api.mu.Lock()
		for sid := range api.sessions {
			delete(api.sessions, sid)
		}
		api.mu.Unlock()
	}
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
