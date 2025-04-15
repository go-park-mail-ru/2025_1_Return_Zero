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
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"go.uber.org/zap"
)

var (
	ErrSessionNotFound = errors.New("session not found")
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
	logger := helpers.LoggerFromContext(ctx)
	conn := r.redisPool.Get()
	defer conn.Close()

	SID := generateSessionID()
	expiration := int(SessionTTL.Seconds())
	_, err := redis.DoContext(conn, ctx, "SETEX", SID, expiration, ID)
	if err != nil {
		logger.Error("failed to create session", zap.Error(err))
		return "", err
	}
	return SID, nil
}

func (r *AuthRedisRepository) DeleteSession(ctx context.Context, SID string) error {
	logger := helpers.LoggerFromContext(ctx)
	conn := r.redisPool.Get()
	defer conn.Close()

	_, err := redis.DoContext(conn, ctx, "DEL", SID)
	if err != nil {
		logger.Error("failed to delete session", zap.Error(err))
		return err
	}
	return nil
}

func (r *AuthRedisRepository) GetSession(ctx context.Context, SID string) (int64, error) {
	logger := helpers.LoggerFromContext(ctx)
	conn := r.redisPool.Get()
	defer conn.Close()

	id, err := redis.Int64(redis.DoContext(conn, ctx, "GET", SID))
	if err != nil {
		if err == redis.ErrNil {
			logger.Error("session not found", zap.String("SID", SID))
			return -1, ErrSessionNotFound
		}
		logger.Error("failed to get session", zap.Error(err))
		return -1, err
	}
	return id, nil
}
