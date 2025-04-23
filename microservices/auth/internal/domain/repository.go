package domain

import (
	"context"
)

type Repository interface {
	CreateSession(ctx context.Context, userID int64) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (int64, error)
}
