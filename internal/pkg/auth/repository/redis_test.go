package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	userID := int64(123)

	conn.Command("SETEX", redigomock.NewAnyData(), int(SessionTTL.Seconds()), userID).Expect("OK")

	sessionID, err := repo.CreateSession(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestCreateSessionError(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	userID := int64(123)

	conn.Command("SETEX", redigomock.NewAnyData(), int(SessionTTL.Seconds()), userID).ExpectError(errors.New("redis error"))

	sessionID, err := repo.CreateSession(ctx, userID)

	assert.Error(t, err)
	assert.Empty(t, sessionID)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestDeleteSession(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	sessionID := "test-session-id"

	conn.Command("DEL", sessionID).Expect(int64(1))

	err := repo.DeleteSession(ctx, sessionID)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestDeleteSessionError(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	sessionID := "test-session-id"

	conn.Command("DEL", sessionID).ExpectError(errors.New("redis error"))

	err := repo.DeleteSession(ctx, sessionID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestGetSession(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	sessionID := "test-session-id"
	userID := int64(123)

	conn.Command("GET", sessionID).Expect(userID)

	resultID, err := repo.GetSession(ctx, sessionID)

	assert.NoError(t, err)
	assert.Equal(t, userID, resultID)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestGetSessionError(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := NewAuthRedisRepository(pool)
	ctx := context.Background()
	sessionID := "test-session-id"

	conn.Command("GET", sessionID).ExpectError(errors.New("redis error"))

	resultID, err := repo.GetSession(ctx, sessionID)

	assert.Error(t, err)
	assert.Equal(t, int64(-1), resultID)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func TestGenerateSessionID(t *testing.T) {
	sessionID := generateSessionID()
	assert.NotEmpty(t, sessionID)

	anotherSessionID := generateSessionID()
	assert.NotEqual(t, sessionID, anotherSessionID)
}
