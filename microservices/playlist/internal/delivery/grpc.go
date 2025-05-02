package delivery

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model"
)

type PlaylistService struct {
	playlistProto.UnimplementedPlaylistServiceServer
	playlistUsecase domain.Usecase
	s3Repository    domain.S3Repository
}

func NewPlaylistService(playlistUsecase domain.Usecase, s3Repository domain.S3Repository) playlistProto.PlaylistServiceServer {
	return &PlaylistService{
		playlistUsecase: playlistUsecase,
		s3Repository:    s3Repository,
	}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, req *playlistProto.CreatePlaylistRequest) (*playlistProto.Playlist, error) {
	playlist, err := s.playlistUsecase.CreatePlaylist(ctx, model.CreatePlaylistRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.PlaylistFromUsecaseToProto(playlist), nil
}

func (s *PlaylistService) UploadPlaylistThumbnail(ctx context.Context, req *playlistProto.UploadPlaylistThumbnailRequest) (*playlistProto.UploadPlaylistThumbnailResponse, error) {
	thumbnail, err := s.playlistUsecase.UploadPlaylistThumbnail(ctx, model.UploadPlaylistThumbnailRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &playlistProto.UploadPlaylistThumbnailResponse{
		Thumbnail: thumbnail,
	}, nil
}

func (s *PlaylistService) GetCombinedPlaylistsByUserID(ctx context.Context, req *playlistProto.GetCombinedPlaylistsByUserIDRequest) (*playlistProto.PlaylistList, error) {
	playlists, err := s.playlistUsecase.GetCombinedPlaylistsByUserID(ctx, model.GetCombinedPlaylistsByUserIDRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.PlaylistListFromUsecaseToProto(playlists), nil
}

func (s *PlaylistService) AddTrackToPlaylist(ctx context.Context, req *playlistProto.AddTrackToPlaylistRequest) (*emptypb.Empty, error) {
	err := s.playlistUsecase.AddTrackToPlaylist(ctx, model.AddTrackToPlaylistRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *PlaylistService) RemoveTrackFromPlaylist(ctx context.Context, req *playlistProto.RemoveTrackFromPlaylistRequest) (*emptypb.Empty, error) {
	err := s.playlistUsecase.RemoveTrackFromPlaylist(ctx, model.RemoveTrackFromPlaylistRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *PlaylistService) GetPlaylistTrackIds(ctx context.Context, req *playlistProto.GetPlaylistTrackIdsRequest) (*playlistProto.GetPlaylistTrackIdsResponse, error) {
	trackIds, err := s.playlistUsecase.GetPlaylistTrackIds(ctx, model.GetPlaylistTrackIdsRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &playlistProto.GetPlaylistTrackIdsResponse{
		TrackIds: trackIds,
	}, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, req *playlistProto.UpdatePlaylistRequest) (*playlistProto.Playlist, error) {
	playlist, err := s.playlistUsecase.UpdatePlaylist(ctx, model.UpdatePlaylistRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.PlaylistFromUsecaseToProto(playlist), nil
}

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, req *playlistProto.GetPlaylistByIDRequest) (*playlistProto.Playlist, error) {
	playlist, err := s.playlistUsecase.GetPlaylistByID(ctx, model.GetPlaylistByIDRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.PlaylistFromUsecaseToProto(playlist), nil
}

func (s *PlaylistService) RemovePlaylist(ctx context.Context, req *playlistProto.RemovePlaylistRequest) (*emptypb.Empty, error) {
	err := s.playlistUsecase.RemovePlaylist(ctx, model.RemovePlaylistRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *PlaylistService) GetPlaylistsToAdd(ctx context.Context, req *playlistProto.GetPlaylistsToAddRequest) (*playlistProto.GetPlaylistsToAddResponse, error) {
	playlists, err := s.playlistUsecase.GetPlaylistsToAdd(ctx, model.GetPlaylistsToAddRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.GetPlaylistsToAddResponseFromUsecaseToProto(playlists), nil
}
