package repository

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type AuthMemoryRepository struct {
	mu       sync.RWMutex
	sessions map[string]*model.Session
}

func NewAuthMemoryRepository() auth.Repository {
	repo := &AuthMemoryRepository{
		sessions: make(map[string]*model.Session),
	}

	return repo
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *AuthMemoryRepository) CreateSession(ID uint) string {
	SID := generateSessionID()
	expiration := time.Now().Add(24 * time.Hour)
	session := &model.Session{
		UserID:    ID,
		ExpiresAt: expiration,
	}
	r.mu.Lock()
	r.sessions[SID] = session
	r.mu.Unlock()
	return SID
}

func (r *AuthMemoryRepository) DeleteSession(SID string) {
	r.mu.Lock()
	delete(r.sessions, SID)
	r.mu.Unlock()
}

func (r *AuthMemoryRepository) GetSession(SID string) (*model.Session, error) {
	r.mu.RLock()
	session, ok := r.sessions[SID]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}
