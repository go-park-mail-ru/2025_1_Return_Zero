package delivery

import (
	"context"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AlbumService struct {
	albumProto.UnimplementedAlbumServiceServer
	albumUsecase domain.Usecase
}

func NewAlbumService(albumUsecase domain.Usecase) albumProto.AlbumServiceServer {
	return &AlbumService{
		albumUsecase: albumUsecase,
	}
}

func (s *AlbumService) GetAllAlbums(ctx context.Context, req *albumProto.FiltersWithUserID) (*albumProto.AlbumList, error) {
	userID := req.UserId.Id
	albums, err := s.albumUsecase.GetAllAlbums(ctx, model.AlbumFiltersFromProtoToUsecase(req.Filters), userID)
	if err != nil {
		return nil, err
	}
	return model.AlbumListFromUsecaseToProto(albums), nil
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, req *albumProto.AlbumIDWithUserID) (*albumProto.Album, error) {
	album, err := s.albumUsecase.GetAlbumByID(ctx, req.AlbumId.Id, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.AlbumFromUsecaseToProto(album), nil
}

func (s *AlbumService) GetAlbumTitleByID(ctx context.Context, req *albumProto.AlbumID) (*albumProto.AlbumTitle, error) {
	albumTitle, err := s.albumUsecase.GetAlbumTitleByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &albumProto.AlbumTitle{
		Title: albumTitle,
	}, nil
}

func (s *AlbumService) GetAlbumTitleByIDs(ctx context.Context, req *albumProto.AlbumIDList) (*albumProto.AlbumTitleMap, error) {
	ids := make([]int64, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = id.Id
	}
	albumTitles, err := s.albumUsecase.GetAlbumTitleByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return model.AlbumTitleMapFromUsecaseToProto(albumTitles), nil
}

func (s *AlbumService) GetAlbumsByIDs(ctx context.Context, req *albumProto.AlbumIDListWithUserID) (*albumProto.AlbumList, error) {
	ids := make([]int64, len(req.Ids.Ids))
	for i, id := range req.Ids.Ids {
		ids[i] = id.Id
	}
	albums, err := s.albumUsecase.GetAlbumsByIDs(ctx, ids, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.AlbumListFromUsecaseToProto(albums), nil
}

func (s *AlbumService) CreateStream(ctx context.Context, req *albumProto.AlbumStreamCreateData) (*emptypb.Empty, error) {
	err := s.albumUsecase.CreateStream(ctx, req.AlbumId.Id, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AlbumService) LikeAlbum(ctx context.Context, req *albumProto.LikeRequest) (*emptypb.Empty, error) {
	err := s.albumUsecase.LikeAlbum(ctx, model.LikeRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AlbumService) GetFavoriteAlbums(ctx context.Context, req *albumProto.FiltersWithUserID) (*albumProto.AlbumList, error) {
	albums, err := s.albumUsecase.GetFavoriteAlbums(ctx, model.AlbumFiltersFromProtoToUsecase(req.Filters), req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.AlbumListFromUsecaseToProto(albums), nil
}
