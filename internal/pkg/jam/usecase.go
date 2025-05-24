package jam

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	CreateJam(ctx context.Context, request *usecase.CreateJamRequest) (*usecase.CreateJamResponse, error)
	JoinJam(ctx context.Context, request *usecase.JoinJamRequest) (*usecase.JamMessage, error)
	HandleClientMessage(ctx context.Context, roomID string, userID string, m *usecase.JamMessage) error
	LeaveJam(ctx context.Context, roomID string, userID string) error
}
