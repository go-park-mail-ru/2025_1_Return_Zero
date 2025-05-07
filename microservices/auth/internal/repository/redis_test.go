package repository

import (
	"context"
	"testing"
	"time"
	"errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
	"go.uber.org/zap"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
)

func setupTest() (*redis.Pool, *redigomock.Conn, context.Context) {
	conn := redigomock.NewConn()

	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			return nil
		},
	}

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return pool, conn, ctx
}

func TestCreateSession(t *testing.T) {
    pool, conn, ctx := setupTest()
    defer pool.Close()

    repo := NewAuthRedisRepository(pool, metrics.NewMockMetrics())

    conn.Command("SETEX", redigomock.NewAnyData(), 86400, int64(1)).Expect("OK")

    sessionID, err := repo.CreateSession(ctx, 1)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if sessionID == "" {
        t.Fatalf("expected non-empty session ID, got empty string")
    }
}

func TestGetSession(t *testing.T) {
    pool, conn, ctx := setupTest()
    defer pool.Close()

    repo := NewAuthRedisRepository(pool, metrics.NewMockMetrics())
    
    conn.Command("GET", "test-session-id").Expect(int64(1))
    
    userID, err := repo.GetSession(ctx, "test-session-id")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    
    if userID != 1 {
        t.Fatalf("expected user ID 1, got %d", userID)
    }
}

func TestDeleteSession(t *testing.T) {
    pool, conn, ctx := setupTest()
    defer pool.Close()

    repo := NewAuthRedisRepository(pool, metrics.NewMockMetrics())
    
    conn.Command("DEL", "test-session-id").Expect(int64(1))
    
    err := repo.DeleteSession(ctx, "test-session-id")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
}

func TestGetSessionError(t *testing.T) {
    pool, conn, ctx := setupTest()
    defer pool.Close()

    repo := NewAuthRedisRepository(pool, metrics.NewMockMetrics())
    
    conn.Command("GET", "test-session-id").ExpectError(errors.New("redis connection error"))
    
    _, err := repo.GetSession(ctx, "test-session-id")
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}

func TestDeleteSessionError(t *testing.T) {
    pool, conn, ctx := setupTest()
    defer pool.Close()

    repo := NewAuthRedisRepository(pool, metrics.NewMockMetrics())
    
    conn.Command("DEL", "test-session-id").ExpectError(errors.New("redis connection error"))
    
    err := repo.DeleteSession(ctx, "test-session-id")
    if err == nil {
        t.Fatalf("expected error, got nil")
    }
}