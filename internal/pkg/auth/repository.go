package auth

import (
	"context"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Repository interface {
	CreateSession(ctx context.Context, ID int64) (string, error)
	DeleteSession(ctx context.Context, SID string) error
	GetSession(ctx context.Context, SID string) (int64, error)
}
