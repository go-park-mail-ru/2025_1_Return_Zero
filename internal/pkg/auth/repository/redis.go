package repository

// sessions map[string - session ID]*model.Session

import (
	"context"
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

const (
	SessionTTL = 24 * time.Hour
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

func (r *AuthRedisRepository) CreateSession(ctx context.Context, ID int64) (string, error) {
	SID := generateSessionID()
	expiration := int(SessionTTL.Seconds())
	_, err := redis.DoContext(r.redis, ctx, "SETEX", SID, expiration, ID)
	if err != nil {
		return "", err
	}
	return SID, nil
}

func (r *AuthRedisRepository) DeleteSession(ctx context.Context, SID string) error {
	_, err := redis.DoContext(r.redis, ctx, "DEL", SID)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRedisRepository) GetSession(ctx context.Context, SID string) (int64, error) {
	id, err := redis.Int64(redis.DoContext(r.redis, ctx, "GET", SID))
	if err != nil {
		return -1, err
	}
	return id, nil
}
