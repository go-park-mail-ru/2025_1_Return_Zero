package usecase

import (
	"context"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model"
	playlistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/errors"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/usecase"
)

type PlaylistUsecase struct {
	playlistRepo domain.Repository
	s3Repo       domain.S3Repository
}

func NewPlaylistUsecase(playlistRepo domain.Repository, s3Repo domain.S3Repository) domain.Usecase {
	return &PlaylistUsecase{playlistRepo: playlistRepo, s3Repo: s3Repo}
}

func (u *PlaylistUsecase) CreatePlaylist(ctx context.Context, playlist *usecaseModel.CreatePlaylistRequest) (*usecaseModel.Playlist, error) {
	repoCreatePlaylistRequest := model.CreatePlaylistRequestFromUsecaseToRepository(playlist)
	repoPlaylist, err := u.playlistRepo.CreatePlaylist(ctx, repoCreatePlaylistRequest)
	if err != nil {
		return nil, err
	}
	return model.PlaylistFromRepositoryToUsecase(repoPlaylist), nil
}

func (u *PlaylistUsecase) UploadPlaylistThumbnail(ctx context.Context, playlist *usecaseModel.UploadPlaylistThumbnailRequest) (string, error) {
	thumbnail, err := u.s3Repo.UploadThumbnail(ctx, playlist.Thumbnail, playlist.Title)
	if err != nil {
		return "", err
	}
	return thumbnail, nil
}

func (u *PlaylistUsecase) GetCombinedPlaylistsByUserID(ctx context.Context, request *usecaseModel.GetCombinedPlaylistsByUserIDRequest) (*usecaseModel.PlaylistList, error) {
	repoPlaylistList, err := u.playlistRepo.GetCombinedPlaylistsByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return model.PlaylistListFromRepositoryToUsecase(repoPlaylistList), nil
}

func (u *PlaylistUsecase) AddTrackToPlaylist(ctx context.Context, request *usecaseModel.AddTrackToPlaylistRequest) error {
	repoRequest := model.AddTrackToPlaylistRequestFromUsecaseToRepository(request)
	err := u.playlistRepo.AddTrackToPlaylist(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *PlaylistUsecase) RemoveTrackFromPlaylist(ctx context.Context, request *usecaseModel.RemoveTrackFromPlaylistRequest) error {
	repoRequest := model.RemoveTrackFromPlaylistRequestFromUsecaseToRepository(request)
	err := u.playlistRepo.RemoveTrackFromPlaylist(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *PlaylistUsecase) GetPlaylistTrackIds(ctx context.Context, request *usecaseModel.GetPlaylistTrackIdsRequest) ([]int64, error) {
	repoPlaylist, err := u.playlistRepo.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		return nil, err
	}

	if repoPlaylist.UserID != request.UserID && !repoPlaylist.IsPublic {
		return nil, playlistErrors.ErrPlaylistPermissionDenied
	}

	repoRequest := model.GetPlaylistTrackIdsRequestFromUsecaseToRepository(request)
	trackIds, err := u.playlistRepo.GetPlaylistTrackIds(ctx, repoRequest)
	if err != nil {
		return nil, err
	}
	return trackIds, nil
}

func (u *PlaylistUsecase) UpdatePlaylist(ctx context.Context, request *usecaseModel.UpdatePlaylistRequest) (*usecaseModel.Playlist, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	repoRequest := model.UpdatePlaylistRequestFromUsecaseToRepository(request)

	playlist, err := u.playlistRepo.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		return nil, err
	}

	if playlist.UserID != request.UserID {
		logger.Warn("playlist permission denied", "playlist_id", request.PlaylistID, "user_id", request.UserID)
		return nil, playlistErrors.ErrPlaylistPermissionDenied
	}

	repoPlaylist, err := u.playlistRepo.UpdatePlaylist(ctx, repoRequest)
	if err != nil {
		return nil, err
	}
	return model.PlaylistFromRepositoryToUsecase(repoPlaylist), nil
}

func (u *PlaylistUsecase) GetPlaylistByID(ctx context.Context, request *usecaseModel.GetPlaylistByIDRequest) (*usecaseModel.Playlist, error) {
	repoPlaylist, err := u.playlistRepo.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		return nil, err
	}

	if repoPlaylist.UserID != request.UserID && !repoPlaylist.IsPublic {
		return nil, playlistErrors.ErrPlaylistPermissionDenied
	}

	return model.PlaylistFromRepositoryToUsecase(repoPlaylist), nil
}

func (u *PlaylistUsecase) RemovePlaylist(ctx context.Context, request *usecaseModel.RemovePlaylistRequest) error {
	playlist, err := u.playlistRepo.GetPlaylistByID(ctx, request.PlaylistID)
	if err != nil {
		return err
	}

	if playlist.UserID != request.UserID {
		return playlistErrors.ErrPlaylistPermissionDenied
	}

	repoRequest := model.RemovePlaylistRequestFromUsecaseToRepository(request)
	err = u.playlistRepo.RemovePlaylist(ctx, repoRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *PlaylistUsecase) GetPlaylistsToAdd(ctx context.Context, request *usecaseModel.GetPlaylistsToAddRequest) (*usecaseModel.GetPlaylistsToAddResponse, error) {
	repoRequest := model.GetPlaylistsToAddRequestFromUsecaseToRepository(request)
	repoResponse, err := u.playlistRepo.GetPlaylistsToAdd(ctx, repoRequest)
	if err != nil {
		return nil, err
	}
	return model.GetPlaylistsToAddResponseFromRepositoryToUsecase(repoResponse), nil
}
