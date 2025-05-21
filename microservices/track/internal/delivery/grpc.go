package delivery

import (
	"context"
	"fmt"

	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TrackService struct {
	trackProto.UnimplementedTrackServiceServer
	trackUsecase domain.Usecase
}

func NewTrackService(trackUsecase domain.Usecase) trackProto.TrackServiceServer {
	return &TrackService{
		trackUsecase: trackUsecase,
	}
}

func (s *TrackService) GetAllTracks(ctx context.Context, req *trackProto.UserIDWithFilters) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.GetAllTracks(ctx, model.FiltersFromProtoToUsecase(req.Filters), req.UserId.Id)
	if err != nil {
		return nil, err
	}
	for i := range tracks {
		fmt.Println(tracks[i].Title)
		fmt.Println(tracks[i].Thumbnail)
		fmt.Println(tracks[i].Duration)
		fmt.Println(tracks[i].AlbumID)
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetTrackByID(ctx context.Context, req *trackProto.TrackIDWithUserID) (*trackProto.TrackDetailed, error) {
	track, err := s.trackUsecase.GetTrackByID(ctx, req.TrackId.Id, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.TrackDetailedFromUsecaseToProto(track), nil
}

func (s *TrackService) CreateStream(ctx context.Context, req *trackProto.TrackStreamCreateData) (*trackProto.StreamID, error) {
	streamID, err := s.trackUsecase.CreateStream(ctx, model.TrackStreamCreateDataFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.StreamIDFromUsecaseToProto(streamID), nil
}

func (s *TrackService) UpdateStreamDuration(ctx context.Context, req *trackProto.TrackStreamUpdateData) (*emptypb.Empty, error) {
	err := s.trackUsecase.UpdateStreamDuration(ctx, model.TrackStreamUpdateDataFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *TrackService) GetLastListenedTracks(ctx context.Context, req *trackProto.UserIDWithFilters) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.GetLastListenedTracks(ctx, req.UserId.Id, model.FiltersFromProtoToUsecase(req.Filters))
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetTracksByIDs(ctx context.Context, req *trackProto.TrackIDList) (*trackProto.TrackList, error) {
	ids, userID := model.TrackIDListFromProtoToUsecase(req)
	tracks, err := s.trackUsecase.GetTracksByIDs(ctx, ids, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetTracksByIDsFiltered(ctx context.Context, req *trackProto.TrackIDListWithFilters) (*trackProto.TrackList, error) {
	ids, filters, userID := model.TrackIDListWithFiltersFromProtoToUsecase(req)
	tracks, err := s.trackUsecase.GetTracksByIDsFiltered(ctx, ids, filters, userID)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetAlbumIDByTrackID(ctx context.Context, req *trackProto.TrackID) (*trackProto.AlbumID, error) {
	albumID, err := s.trackUsecase.GetAlbumIDByTrackID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &trackProto.AlbumID{Id: albumID}, nil
}

func (s *TrackService) GetTracksByAlbumID(ctx context.Context, req *trackProto.AlbumIDWithUserID) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.GetTracksByAlbumID(ctx, req.AlbumId.Id, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetMinutesListenedByUserID(ctx context.Context, req *trackProto.UserID) (*trackProto.MinutesListened, error) {
	minutesListened, err := s.trackUsecase.GetMinutesListenedByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &trackProto.MinutesListened{Minutes: minutesListened}, nil
}

func (s *TrackService) GetTracksListenedByUserID(ctx context.Context, req *trackProto.UserID) (*trackProto.TracksListened, error) {
	tracks, err := s.trackUsecase.GetTracksListenedByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &trackProto.TracksListened{Tracks: tracks}, nil
}

func (s *TrackService) LikeTrack(ctx context.Context, req *trackProto.LikeRequest) (*emptypb.Empty, error) {
	err := s.trackUsecase.LikeTrack(ctx, model.LikeRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *TrackService) GetFavoriteTracks(ctx context.Context, req *trackProto.FavoriteRequest) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.GetFavoriteTracks(ctx, model.FavoriteRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) SearchTracks(ctx context.Context, req *trackProto.Query) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.SearchTracks(ctx, req.Query, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) AddTracksToAlbum(ctx context.Context, req *trackProto.TracksListWithAlbumID) (*trackProto.TrackIdsList, error) {
	tracksList := model.TracksListWithAlbumIDFromProtoToUsecase(req)
	trackIDs, err := s.trackUsecase.AddTracksToAlbum(ctx, tracksList)
	if err != nil {
		return nil, err
	}
	return model.TrackIdsListFromUsecaseToProto(trackIDs), nil
}

func (s *TrackService) DeleteTracksByAlbumID(ctx context.Context, req *trackProto.AlbumID) (*emptypb.Empty, error) {
	err := s.trackUsecase.DeleteTracksByAlbumID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
