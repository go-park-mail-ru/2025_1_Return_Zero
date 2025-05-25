package jam

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	CreateJam(ctx context.Context, request *repository.CreateJamRequest) (*repository.CreateJamResponse, error)
	AddUser(ctx context.Context, roomID string, userID string) error
	PauseJam(ctx context.Context, roomID string) error
	GetInitialJamData(ctx context.Context, roomID string) (*repository.JamMessage, error)
	GetHostID(ctx context.Context, roomID string) (string, error)
	LoadTrack(ctx context.Context, roomID string, userID string) error
	CheckAllReadyAndPlay(ctx context.Context, roomID string)
	MarkUserAsReady(ctx context.Context, roomID string, userID string) error
	RemoveUser(ctx context.Context, roomID string, userID string) error
	RemoveJam(ctx context.Context, roomID string) error
	ExistsRoom(ctx context.Context, roomID string) (bool, error)
	SeekJam(ctx context.Context, roomID string, position int64) error
	StoreUserInfo(ctx context.Context, roomID string, userID string, username string, avatarURL string) error
	GetUserInfo(ctx context.Context, roomID string, userID string) (string, string, error)
	SubscribeToJamMessages(ctx context.Context, roomID string) (<-chan []byte, error)
}
