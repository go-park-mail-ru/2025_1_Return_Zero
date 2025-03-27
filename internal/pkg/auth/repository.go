package auth

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Repository interface {
	CreateSession(ID uint) string
	DeleteSession(SID string)
	GetSession(SID string) (*model.Session, error)
}