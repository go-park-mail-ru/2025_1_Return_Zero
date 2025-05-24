package usecase

import (
	"context"
	"errors"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"go.uber.org/zap"
)

type Usecase struct {
	jamRepository jam.Repository
}

func NewUsecase(jamRepository jam.Repository) *Usecase {
	return &Usecase{
		jamRepository: jamRepository,
	}
}

func (u *Usecase) CreateJam(ctx context.Context, request *usecase.CreateJamRequest) (*usecase.CreateJamResponse, error) {
	repoRequest := &repository.CreateJamRequest{
		UserID:   request.UserID,
		TrackID:  request.TrackID,
		Position: request.Position,
	}
	jamResponse, err := u.jamRepository.CreateJam(ctx, repoRequest)
	if err != nil {
		return nil, err
	}
	return &usecase.CreateJamResponse{
		RoomID: jamResponse.RoomID,
		HostID: jamResponse.HostID,
	}, nil
}

func (u *Usecase) JoinJam(ctx context.Context, request *usecase.JoinJamRequest) (*usecase.JamMessage, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	exists, err := u.jamRepository.ExistsRoom(ctx, request.RoomID)
	if err != nil {
		logger.Error("failed to check if room exists", zap.Error(err))
		return nil, err
	}
	if !exists {
		return nil, errors.New("room not found")
	}
	hostID, err := u.jamRepository.GetHostID(ctx, request.RoomID)
	if err != nil {
		logger.Error("failed to get host id", zap.Error(err))
		return nil, err
	}

	if hostID != request.UserID {
		err = u.jamRepository.AddUser(ctx, request.RoomID, request.UserID)
		if err != nil {
			logger.Error("failed to add user", zap.Error(err))
			return nil, err
		}

		err = u.jamRepository.PauseJam(ctx, request.RoomID)
		if err != nil {
			logger.Error("failed to pause jam", zap.Error(err))
			return nil, err
		}
	}

	repoJamData, err := u.jamRepository.GetInitialJamData(ctx, request.RoomID)
	if err != nil {
		logger.Error("failed to get initial jam data", zap.Error(err))
		return nil, err
	}

	jamData := model.JamMessageFromRepositoryToUsecase(repoJamData)

	return jamData, nil
}

func (u *Usecase) HandleClientMessage(ctx context.Context, roomID string, userID string, m *usecase.JamMessage) error {
	if m.Type == "jam:closed" {
		return errors.New("jam closed")
	}

	hostID, err := u.jamRepository.GetHostID(ctx, roomID)
	if err != nil {
		return err
	}

	isHost := hostID == userID

	switch m.Type {
	case "host:load":
		if !isHost {
			return errors.New("not host")
		}
		err := u.jamRepository.LoadTrack(ctx, roomID, m.TrackID)
		if err != nil {
			return err
		}
		err = u.jamRepository.PauseJam(ctx, roomID)
		if err != nil {
			return err
		}
	case "client:ready":
		err := u.jamRepository.MarkUserAsReady(ctx, roomID, userID)
		if err != nil {
			return err
		}
		u.jamRepository.CheckAllReadyAndPlay(ctx, roomID)
	case "host:play":
		if !isHost {
			return errors.New("not host")
		}
		u.jamRepository.CheckAllReadyAndPlay(ctx, roomID)
	case "host:pause":
		if !isHost {
			return errors.New("not host")
		}
		err := u.jamRepository.PauseJam(ctx, roomID)
		if err != nil {
			return err
		}
	case "host:seek":
		if !isHost {
			return errors.New("not host")
		}
		err := u.jamRepository.SeekJam(ctx, roomID, m.Position)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) LeaveJam(ctx context.Context, roomID string, userID string) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	hostID, err := u.jamRepository.GetHostID(ctx, roomID)
	if err != nil {
		logger.Error("failed to get host id", zap.Error(err))
		return err
	}

	if hostID == userID {
		err = u.jamRepository.RemoveJam(ctx, roomID)
		if err != nil {
			logger.Error("failed to remove jam", zap.Error(err))
			return err
		}
		return nil
	}

	err = u.jamRepository.RemoveUser(ctx, roomID, userID)
	if err != nil {
		logger.Error("failed to remove user", zap.Error(err))
		return err
	}

	u.jamRepository.CheckAllReadyAndPlay(ctx, roomID)
	return nil
}
