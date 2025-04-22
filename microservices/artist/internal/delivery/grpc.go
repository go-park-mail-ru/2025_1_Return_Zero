package delivery

import (
	"context"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model"
)

type ArtistService struct {
	artistProto.UnimplementedArtistServiceServer
	artistUsecase domain.Usecase
}

func NewArtistService(artistUsecase domain.Usecase) artistProto.ArtistServiceServer {
	return &ArtistService{
		artistUsecase: artistUsecase,
	}
}

func (s *ArtistService) GetArtistByID(ctx context.Context, req *artistProto.ArtistID) (*artistProto.ArtistDetailed, error) {
	artist, err := s.artistUsecase.GetArtistByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return model.ArtistDetailedFromUsecaseToProto(artist), nil
}

func (s *ArtistService) GetAllArtists(ctx context.Context, req *artistProto.Filters) (*artistProto.ArtistList, error) {
	artists, err := s.artistUsecase.GetAllArtists(ctx, model.ArtistFiltersFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromUsecaseToProto(artists), nil
}

func (s *ArtistService) GetArtistTitleByID(ctx context.Context, req *artistProto.ArtistID) (*artistProto.ArtistTitle, error) {
	artist, err := s.artistUsecase.GetArtistTitleByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &artistProto.ArtistTitle{
		Title: artist,
	}, nil
}

func (s *ArtistService) GetArtistsByTrackID(ctx context.Context, req *artistProto.TrackID) (*artistProto.ArtistWithRoleList, error) {
	artists, err := s.artistUsecase.GetArtistsByTrackID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithRoleListFromUsecaseToProto(artists), nil
}

func (s *ArtistService) GetArtistsByTrackIDs(ctx context.Context, req *artistProto.TrackIDList) (*artistProto.ArtistWithRoleMap, error) {
	artists, err := s.artistUsecase.GetArtistsByTrackIDs(ctx, model.TrackIDListFromProtoToUsecase(req.Ids))
	if err != nil {
		return nil, err
	}
	return model.ArtistWithRoleMapFromUsecaseToProto(artists), nil
}

func (s *ArtistService) GetArtistsByAlbumID(ctx context.Context, req *artistProto.AlbumID) (*artistProto.ArtistWithTitleList, error) {
	artists, err := s.artistUsecase.GetArtistsByAlbumID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return model.ArtistWithTitleListFromUsecaseToProto(artists), nil
}

func (s *ArtistService) GetArtistsByAlbumIDs(ctx context.Context, req *artistProto.AlbumIDList) (*artistProto.ArtistWithTitleMap, error) {
	artists, err := s.artistUsecase.GetArtistsByAlbumIDs(ctx, model.AlbumIDListFromProtoToUsecase(req.Ids))
	if err != nil {
		return nil, err
	}
	return model.ArtistWithTitleMapFromUsecaseToProto(artists), nil
}

func (s *ArtistService) GetAlbumIDsByArtistID(ctx context.Context, req *artistProto.ArtistID) (*artistProto.AlbumIDList, error) {
	albumIDs, err := s.artistUsecase.GetAlbumIDsByArtistID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	albumIDList := make([]*artistProto.AlbumID, 0, len(albumIDs))
	for _, albumID := range albumIDs {
		albumIDList = append(albumIDList, &artistProto.AlbumID{Id: albumID})
	}
	return &artistProto.AlbumIDList{Ids: albumIDList}, nil
}
