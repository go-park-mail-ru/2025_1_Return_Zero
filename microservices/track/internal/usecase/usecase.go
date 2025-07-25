package usecase

import (
	"bytes"
	"context"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/usecase"
	"github.com/hajimehoshi/go-mp3"
	"go.uber.org/zap"
)

type TrackUsecase struct {
	trackRepo domain.Repository
	s3Repo    domain.S3Repository
}

func NewTrackUsecase(trackRepo domain.Repository, s3Repo domain.S3Repository) domain.Usecase {
	return &TrackUsecase{trackRepo: trackRepo, s3Repo: s3Repo}
}

func (u *TrackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters, userID int64) ([]*usecaseModel.Track, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	repoTracks, err := u.trackRepo.GetAllTracks(ctx, repoFilters, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetTrackByID(ctx context.Context, id int64, userID int64) (*usecaseModel.TrackDetailed, error) {
	repoTrack, err := u.trackRepo.GetTrackByID(ctx, id, userID)
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
		logger.Warn("updating stream doesn't belong to user", zap.Error(trackErrors.ErrStreamPermissionDenied))
		return trackErrors.ErrStreamPermissionDenied
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

	if len(repoStreams) == 0 {
		return []*usecaseModel.Track{}, nil
	}

	streamIDs := make([]int64, len(repoStreams))
	for i, stream := range repoStreams {
		streamIDs[i] = stream.ID
	}
	repoTrackIDs := make([]int64, len(repoStreams))
	for i, stream := range repoStreams {
		repoTrackIDs[i] = stream.TrackID
	}

	repoTracks, err := u.trackRepo.GetTracksByIDs(ctx, repoTrackIDs, userID)
	if err != nil {
		return nil, err
	}

	usecaseTracks := make([]*usecaseModel.Track, 0, len(repoStreams))
	for _, stream := range repoStreams {
		track, exists := repoTracks[stream.TrackID]
		if exists {
			usecaseTracks = append(usecaseTracks, model.TrackFromRepositoryToUsecase(track))
		}
	}

	return usecaseTracks, nil
}

func (u *TrackUsecase) GetTracksByIDs(ctx context.Context, ids []int64, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetTracksByIDs(ctx, ids, userID)
	if err != nil {
		return nil, err
	}

	usecaseTracks := make([]*usecaseModel.Track, 0, len(ids))
	for _, id := range ids {
		track, exists := repoTracks[id]
		if exists {
			usecaseTracks = append(usecaseTracks, model.TrackFromRepositoryToUsecase(track))
		}
	}

	return usecaseTracks, nil
}

func (u *TrackUsecase) GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *usecaseModel.TrackFilters, userID int64) ([]*usecaseModel.Track, error) {
	repoFilters := model.FiltersFromUsecaseToRepository(filters)
	repoTracks, err := u.trackRepo.GetTracksByIDsFiltered(ctx, ids, repoFilters, userID)
	if err != nil {
		return nil, err
	}

	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error) {
	repoAlbumID, err := u.trackRepo.GetAlbumIDByTrackID(ctx, id)
	if err != nil {
		return 0, err
	}
	return repoAlbumID, nil
}

func (u *TrackUsecase) GetTracksByAlbumID(ctx context.Context, id int64, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetTracksByAlbumID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	usecaseTracks := make([]*usecaseModel.Track, len(repoTracks))
	for i, repoTrack := range repoTracks {
		usecaseTracks[i] = model.TrackFromRepositoryToUsecase(repoTrack)
	}
	return usecaseTracks, nil
}

func (u *TrackUsecase) GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	repoMinutesListened, err := u.trackRepo.GetMinutesListenedByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return repoMinutesListened, nil
}

func (u *TrackUsecase) GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	repoTracksListened, err := u.trackRepo.GetTracksListenedByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return repoTracksListened, nil
}

func (u *TrackUsecase) LikeTrack(ctx context.Context, request *usecaseModel.LikeRequest) error {
	repoRequest := model.LikeRequestFromUsecaseToRepository(request)

	exists, err := u.trackRepo.CheckTrackExists(ctx, request.TrackID)
	if err != nil {
		return err
	}

	if !exists {
		return trackErrors.NewNotFoundError("track not found")
	}
	if request.IsLike {
		err := u.trackRepo.LikeTrack(ctx, repoRequest)
		if err != nil {
			return err
		}
		return nil
	}

	err = u.trackRepo.UnlikeTrack(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *TrackUsecase) GetFavoriteTracks(ctx context.Context, favoriteRequest *usecaseModel.FavoriteRequest) ([]*usecaseModel.Track, error) {
	repoRequest := model.FavoriteRequestFromUsecaseToRepository(favoriteRequest)
	repoTracks, err := u.trackRepo.GetFavoriteTracks(ctx, repoRequest)
	if err != nil {
		return nil, err
	}

	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) SearchTracks(ctx context.Context, query string, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.SearchTracks(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func getMp3Duration(fileData []byte) (int64, error) {
	reader := bytes.NewReader(fileData)

	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		return 0, err
	}

	samples := decoder.Length() / 4
	sampleRate := decoder.SampleRate()
	duration := time.Duration(float64(samples) / float64(sampleRate) * float64(time.Second))

	return int64(duration.Seconds()), nil
}

func (u *TrackUsecase) AddTracksToAlbum(ctx context.Context, tracksList *usecaseModel.TracksListWithAlbumID) ([]int64, error) {
	var tracksToUpload []*usecaseModel.Track
	if len(tracksList.Cover) == 0 {
		return nil, trackErrors.NewUnsupportedImageFormatError("пустое изображение обложки")
	}
	trackCover, err := u.s3Repo.UploadTrackAvatar(ctx, tracksList.Tracks[0].Title, tracksList.Cover)
	if err != nil {
		return nil, err
	}
	for _, track := range tracksList.Tracks {
		err := u.s3Repo.UploadTrack(ctx, track.Title, track.File)
		if err != nil {
			return nil, err
		}
		duration, err := getMp3Duration(track.File)
		if err != nil {
			return nil, err
		}
		tracksToUpload = append(tracksToUpload, &usecaseModel.Track{
			Title:     track.Title,
			Thumbnail: trackCover,
			AlbumID:   tracksList.AlbumID,
			Duration:  duration,
		})
	}
	trackListRepo := model.TrackListFromUsecaseToRepository(tracksToUpload)

	trackIDs, err := u.trackRepo.AddTracksToAlbum(ctx, trackListRepo)
	if err != nil {
		return nil, err
	}
	return trackIDs, nil
}

func (u *TrackUsecase) DeleteTracksByAlbumID(ctx context.Context, albumID int64) error {
	return u.trackRepo.DeleteTracksByAlbumID(ctx, albumID)
}

func (u *TrackUsecase) GetMostLikedTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetMostLikedTracks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetMostRecentTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetMostRecentTracks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetMostListenedLastMonthTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetMostListenedLastMonthTracks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}

func (u *TrackUsecase) GetMostLikedLastWeekTracks(ctx context.Context, userID int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetMostLikedLastWeekTracks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromRepositoryToUsecase(repoTracks), nil
}
