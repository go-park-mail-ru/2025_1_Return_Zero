package repository

// sessions map[string - session ID]*model.Session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type AuthRedisRepository struct {
	redis redis.Conn
}

func NewAuthRedisRepository(redis redis.Conn) auth.Repository {
	repo := &AuthRedisRepository{
		redis: redis,
	}

	return repo
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *AuthRedisRepository) CreateSession(ID int64) string {
	SID := generateSessionID()
	expiration := 24 * 3600 * time.Second 
	r.redis.Do("SETEX", SID, expiration, ID)
	return SID
}

func (r *AuthRedisRepository) DeleteSession(SID string) {
	r.redis.Do("DEL", SID)
}

func (r *AuthRedisRepository) GetSession(SID string) (int64, error) {
	id, err := redis.Int64(r.redis.Do("GET", SID))
	if err != nil {
		return 0, err
	}
	return id, nil
}
