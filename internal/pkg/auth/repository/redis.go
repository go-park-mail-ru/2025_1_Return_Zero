package repository

// sessions map[string - session ID]*model.Session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/gomodule/redigo/redis"
)

const (
	SessionTTL = 24 * time.Hour
)

type AuthRedisRepository struct {
	redisPool *redis.Pool
}

func NewAuthRedisRepository(pool *redis.Pool) auth.Repository {
	repo := &AuthRedisRepository{
		redisPool: pool,
	}

	return repo
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *AuthRedisRepository) CreateSession(ctx context.Context, ID int64) (string, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	SID := generateSessionID()
	expiration := int(SessionTTL.Seconds())
	_, err := redis.DoContext(conn, ctx, "SETEX", SID, expiration, ID)
	if err != nil {
		return "", err
	}
	return SID, nil
}

func (r *AuthRedisRepository) DeleteSession(ctx context.Context, SID string) error {
	conn := r.redisPool.Get()
	defer conn.Close()

	_, err := redis.DoContext(conn, ctx, "DEL", SID)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRedisRepository) GetSession(ctx context.Context, SID string) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()

	id, err := redis.Int64(redis.DoContext(conn, ctx, "GET", SID))
	if err != nil {
		return -1, err
	}
	return id, nil
}
