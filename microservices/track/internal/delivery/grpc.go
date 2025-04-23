package delivery

import (
	"context"

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

func (s *TrackService) GetAllTracks(ctx context.Context, req *trackProto.Filters) (*trackProto.TrackList, error) {
	tracks, err := s.trackUsecase.GetAllTracks(ctx, model.FiltersFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetTrackByID(ctx context.Context, req *trackProto.TrackID) (*trackProto.TrackDetailed, error) {
	track, err := s.trackUsecase.GetTrackByID(ctx, req.Id)
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
	ids := make([]int64, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = id.Id
	}
	tracks, err := s.trackUsecase.GetTracksByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}

func (s *TrackService) GetTracksByIDsFiltered(ctx context.Context, req *trackProto.TrackIDListWithFilters) (*trackProto.TrackList, error) {
	ids, filters := model.TrackIDListWithFiltersFromProtoToUsecase(req)
	tracks, err := s.trackUsecase.GetTracksByIDsFiltered(ctx, ids, filters)
	if err != nil {
		return nil, err
	}
	return model.TrackListFromUsecaseToProto(tracks), nil
}
