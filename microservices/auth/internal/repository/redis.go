package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	authErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/model/errors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/domain"
	"github.com/gomodule/redigo/redis"
)

const (
	SessionTTL = 24 * time.Hour
)

type authRedisRepository struct {
	redisPool *redis.Pool
}

func NewAuthRedisRepository(redisPool *redis.Pool) domain.Repository {
	return &authRedisRepository{redisPool: redisPool}
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *authRedisRepository) CreateSession(ctx context.Context, userID int64) (string, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	SID := generateSessionID()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating session")
	expiration := int(SessionTTL.Seconds())
	_, err := redis.DoContext(conn, ctx, "SETEX", SID, expiration, userID)
	if err != nil {
		return "", authErrors.NewCreateSessionError("failed to create session: %v", err)
	}
	return SID, nil
}

func (r *authRedisRepository) DeleteSession(ctx context.Context, sessionID string) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Deleting session")
	_, err := redis.DoContext(conn, ctx, "DEL", sessionID)
	if err != nil {
		return authErrors.NewDeleteSessionError("failed to delete session: %v", err)
	}
	return nil
}

func (r *authRedisRepository) GetSession(ctx context.Context, sessionID string) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting session")
	id, err := redis.Int64(redis.DoContext(conn, ctx, "GET", sessionID))
	if err != nil {
		return -1, authErrors.NewGetSessionError("failed to get session: %v", err)
	}
	return id, nil
}
