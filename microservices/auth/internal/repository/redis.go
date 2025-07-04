package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/domain"
	authErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/model/errors"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

const (
	SessionTTL = 24 * time.Hour
)

type authRedisRepository struct {
	redisPool *redis.Pool
	metrics   *metrics.Metrics
}

func NewAuthRedisRepository(redisPool *redis.Pool, metrics *metrics.Metrics) domain.Repository {
	return &authRedisRepository{redisPool: redisPool, metrics: metrics}
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (r *authRedisRepository) CreateSession(ctx context.Context, userID int64) (string, error) {
	start := time.Now()
	conn := r.redisPool.Get()

	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Creating session")
	expiration := int(SessionTTL.Seconds())

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	SID, err := generateSessionID()
	if err != nil {
		logger.Error("failed to generate session ID", zap.Error(err))
		return "", authErrors.NewCreateSessionError("failed to generate session ID: %v", err)
	}

	_, err = redis.DoContext(conn, ctx, "SETEX", SID, expiration, userID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("CreateSession").Inc()
		logger.Error("failed to create session", zap.Error(err))
		return "", authErrors.NewCreateSessionError("failed to create session: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("CreateSession").Observe(duration)
	return SID, nil
}

func (r *authRedisRepository) DeleteSession(ctx context.Context, sessionID string) error {
	start := time.Now()
	conn := r.redisPool.Get()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Deleting session")

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	_, err := redis.DoContext(conn, ctx, "DEL", sessionID)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("DeleteSession").Inc()
		logger.Error("failed to delete session", zap.Error(err))
		return authErrors.NewDeleteSessionError("failed to delete session: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("DeleteSession").Observe(duration)
	return nil
}

func (r *authRedisRepository) GetSession(ctx context.Context, sessionID string) (int64, error) {
	start := time.Now()
	conn := r.redisPool.Get()
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Getting session")

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection:", zap.Error(err))
		}
	}()

	id, err := redis.Int64(redis.DoContext(conn, ctx, "GET", sessionID))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetSession").Inc()
		logger.Error("failed to get session", zap.Error(err))
		return -1, authErrors.NewGetSessionError("failed to get session: %v", err)
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetSession").Observe(duration)
	return id, nil
}
