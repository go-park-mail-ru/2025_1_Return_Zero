package delivery

import (
	"context"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model"
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

func (s *AlbumService) GetAllAlbums(ctx context.Context, req *albumProto.Filters) (*albumProto.AlbumList, error) {
	albums, err := s.albumUsecase.GetAllAlbums(ctx, model.AlbumFiltersFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.AlbumListFromUsecaseToProto(albums), nil
}

func (s *AlbumService) GetAlbumByID(ctx context.Context, req *albumProto.AlbumID) (*albumProto.Album, error) {
	album, err := s.albumUsecase.GetAlbumByID(ctx, req.Id)
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

func (s *AlbumService) GetAlbumsByIDs(ctx context.Context, req *albumProto.AlbumIDList) (*albumProto.AlbumList, error) {
	ids := make([]int64, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = id.Id
	}
	albums, err := s.albumUsecase.GetAlbumsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return model.AlbumListFromUsecaseToProto(albums), nil
}
