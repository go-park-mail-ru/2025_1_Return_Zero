package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/playlist"

	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

func NewUsecase(playlistClient *playlistProto.PlaylistServiceClient, userClient *userProto.UserServiceClient) playlist.Usecase {
	return &playlistUsecase{playlistClient: playlistClient, userClient: userClient}
}

type playlistUsecase struct {
	playlistClient *playlistProto.PlaylistServiceClient
	userClient     *userProto.UserServiceClient
}

func (u *playlistUsecase) CreatePlaylist(ctx context.Context, request *usecaseModel.CreatePlaylistRequest) (*usecaseModel.Playlist, error) {
	thumbnail, err := (*u.playlistClient).UploadPlaylistThumbnail(ctx, model.UploadPlaylistThumbnailRequestFromUsecaseToProto(request.Title, request.Thumbnail))
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	playlist, err := (*u.playlistClient).CreatePlaylist(ctx, model.CreatePlaylistRequestFromUsecaseToProto(request, thumbnail.GetThumbnail()))
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	user, err := (*u.userClient).GetUserByID(ctx, &userProto.UserID{
		Id: playlist.GetUserId(),
	})
	if err != nil {
		return nil, err
	}

	return model.PlaylistFromProtoToUsecase(playlist, user.Username), nil
}

func (u *playlistUsecase) GetCombinedPlaylistsForCurrentUser(ctx context.Context, userID int64) ([]*usecaseModel.Playlist, error) {
	request := &playlistProto.GetCombinedPlaylistsByUserIDRequest{
		UserId: userID,
	}

	playlists, err := (*u.playlistClient).GetCombinedPlaylistsByUserID(ctx, request)
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	usecasePlaylists := make([]*usecaseModel.Playlist, len(playlists.GetPlaylists()))
	for i, playlist := range playlists.GetPlaylists() {
		user, err := (*u.userClient).GetUserByID(ctx, &userProto.UserID{
			Id: playlist.GetUserId(),
		})
		if err != nil {
			return nil, err
		}
		usecasePlaylists[i] = model.PlaylistFromProtoToUsecase(playlist, user.Username)
	}

	return usecasePlaylists, nil
}

func (u *playlistUsecase) AddTrackToPlaylist(ctx context.Context, request *usecaseModel.AddTrackToPlaylistRequest) error {
	_, err := (*u.playlistClient).AddTrackToPlaylist(ctx, model.AddTrackToPlaylistRequestFromUsecaseToProto(request))
	if err != nil {
		return customErrors.HandlePlaylistGRPCError(err)
	}
	return nil
}

func (u *playlistUsecase) RemoveTrackFromPlaylist(ctx context.Context, request *usecaseModel.RemoveTrackFromPlaylistRequest) error {
	_, err := (*u.playlistClient).RemoveTrackFromPlaylist(ctx, model.RemoveTrackFromPlaylistRequestFromUsecaseToProto(request))
	if err != nil {
		return customErrors.HandlePlaylistGRPCError(err)
	}
	return nil
}

func (u *playlistUsecase) UpdatePlaylist(ctx context.Context, request *usecaseModel.UpdatePlaylistRequest) (*usecaseModel.Playlist, error) {
	thumbnail := ""
	if request.Thumbnail != nil {
		thumbnailObject, err := (*u.playlistClient).UploadPlaylistThumbnail(ctx, model.UploadPlaylistThumbnailRequestFromUsecaseToProto(request.Title, request.Thumbnail))
		if err != nil {
			return nil, customErrors.HandlePlaylistGRPCError(err)
		}
		thumbnail = thumbnailObject.GetThumbnail()
	}

	playlist, err := (*u.playlistClient).UpdatePlaylist(ctx, model.UpdatePlaylistRequestFromUsecaseToProto(request, thumbnail))
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	user, err := (*u.userClient).GetUserByID(ctx, &userProto.UserID{
		Id: playlist.GetUserId(),
	})
	if err != nil {
		return nil, err
	}

	return model.PlaylistFromProtoToUsecase(playlist, user.Username), nil
}

func (u *playlistUsecase) GetPlaylistByID(ctx context.Context, playlistID int64) (*usecaseModel.Playlist, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	user, err := (*u.userClient).GetUserByID(ctx, &userProto.UserID{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	playlist, err := (*u.playlistClient).GetPlaylistByID(ctx, &playlistProto.GetPlaylistByIDRequest{
		Id:     playlistID,
		UserId: userID,
	})
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	return model.PlaylistFromProtoToUsecase(playlist, user.Username), nil
}

func (u *playlistUsecase) RemovePlaylist(ctx context.Context, request *usecaseModel.RemovePlaylistRequest) error {
	_, err := (*u.playlistClient).RemovePlaylist(ctx, model.RemovePlaylistRequestFromUsecaseToProto(request))
	if err != nil {
		return customErrors.HandlePlaylistGRPCError(err)
	}
	return nil
}

func (u *playlistUsecase) GetPlaylistsToAdd(ctx context.Context, request *usecaseModel.GetPlaylistsToAddRequest) ([]*usecaseModel.PlaylistWithIsIncludedTrack, error) {
	playlists, err := (*u.playlistClient).GetPlaylistsToAdd(ctx, model.GetPlaylistsToAddRequestFromUsecaseToProto(request))
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	user, err := (*u.userClient).GetUserByID(ctx, &userProto.UserID{
		Id: request.UserID,
	})
	if err != nil {
		return nil, err
	}

	return model.GetPlaylistsToAddResponseFromProtoToUsecase(playlists, user.Username), nil
}
