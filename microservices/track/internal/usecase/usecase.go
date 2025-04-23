package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
	"go.uber.org/zap"
)

type TrackUsecase struct {
	trackRepo domain.Repository
	s3Repo    domain.S3Repository
}

func NewTrackUsecase(trackRepo domain.Repository, s3Repo domain.S3Repository) domain.Usecase {
	return &TrackUsecase{trackRepo: trackRepo, s3Repo: s3Repo}
}

func (u *TrackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	repoTracks, err := u.trackRepo.GetAllTracks(ctx, repoFilters)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetTrackByID(ctx context.Context, id int64) (*usecaseModel.TrackDetailed, error) {
	repoTrack, err := u.trackRepo.GetTrackByID(ctx, id)
	if err != nil {
		return nil, err
	}

	trackFileUrl, err := u.s3Repo.GetPresignedURL(repoTrack.FileKey)
	if err != nil {
		return nil, err
	}

	return model.TrackDetailedFromRepositoryToUsecase(repoTrack, trackFileUrl), nil
}

func (u *TrackUsecase) CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error) {
	repoStream := model.TrackStreamCreateDataFromUsecaseToRepository(stream)
	repoStreamID, err := u.trackRepo.CreateStream(ctx, repoStream)
	if err != nil {
		return 0, err
	}
	return repoStreamID, nil
}

func (u *TrackUsecase) UpdateStreamDuration(ctx context.Context, stream *usecaseModel.TrackStreamUpdateData) error {
	logger := loggerPkg.LoggerFromContext(ctx)
	existingStream, err := u.trackRepo.GetStreamByID(ctx, stream.StreamID)
	if err != nil {
		return err
	}

	if existingStream.UserID != stream.UserID {
		logger.Warn("updating stream doesn't belong to user", zap.Error(err))
		return customErrors.ErrStreamNotFound
	}

	repoStream := model.TrackStreamUpdateDataFromUsecaseToRepository(stream)
	err = u.trackRepo.UpdateStreamDuration(ctx, repoStream)
	if err != nil {
		return err
	}
	return nil
}

func (u *TrackUsecase) GetLastListenedTracks(ctx context.Context, userID int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	repoStreams, err := u.trackRepo.GetStreamsByUserID(ctx, userID, repoFilters)
	if err != nil {
		return nil, err
	}
	streamIDs := make([]int64, len(repoStreams))
	for i, stream := range repoStreams {
		streamIDs[i] = stream.ID
	}
	repoTrackIDs := make([]int64, len(repoStreams))
	for i, stream := range repoStreams {
		repoTrackIDs[i] = stream.TrackID
	}
	repoTracks, err := u.trackRepo.GetTracksByIDs(ctx, repoTrackIDs)
	if err != nil {
		return nil, err
	}

	usecaseTracks := make([]*usecaseModel.Track, len(repoTracks))
	for i, id := range streamIDs {
		usecaseTracks[i] = model.TrackFromRepositoryToUsecase(repoTracks[id])
	}

	return usecaseTracks, nil
}

func (u *TrackUsecase) GetTracksByIDs(ctx context.Context, ids []int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetTracksByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	usecaseTracks := make([]*usecaseModel.Track, len(repoTracks))
	for i, id := range ids {
		usecaseTracks[i] = model.TrackFromRepositoryToUsecase(repoTracks[id])
	}

	return usecaseTracks, nil
}

func (u *TrackUsecase) GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	repoTracks, err := u.trackRepo.GetTracksByIDsFiltered(ctx, ids, repoFilters)
	if err != nil {
		return nil, err
	}

	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}
