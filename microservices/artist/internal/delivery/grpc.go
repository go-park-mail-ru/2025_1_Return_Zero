package delivery

import (
	"context"

	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *ArtistService) GetArtistByID(ctx context.Context, req *artistProto.ArtistIDWithUserID) (*artistProto.ArtistDetailed, error) {
	userID := req.UserId.Id
	artist, err := s.artistUsecase.GetArtistByID(ctx, req.ArtistId.Id, userID)
	if err != nil {
		return nil, err
	}
	return model.ArtistDetailedFromUsecaseToProto(artist), nil
}

func (s *ArtistService) GetAllArtists(ctx context.Context, req *artistProto.FiltersWithUserID) (*artistProto.ArtistList, error) {
	userID := req.UserId.Id
	artists, err := s.artistUsecase.GetAllArtists(ctx, model.ArtistFiltersFromProtoToUsecase(req.Filters), userID)
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

func (s *ArtistService) GetTrackIDsByArtistID(ctx context.Context, req *artistProto.ArtistID) (*artistProto.TrackIDList, error) {
	trackIDs, err := s.artistUsecase.GetTrackIDsByArtistID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	trackIDList := make([]*artistProto.TrackID, 0, len(trackIDs))
	for _, trackID := range trackIDs {
		trackIDList = append(trackIDList, &artistProto.TrackID{Id: trackID})
	}
	return &artistProto.TrackIDList{Ids: trackIDList}, nil
}

func (s *ArtistService) CreateStreamsByArtistIDs(ctx context.Context, req *artistProto.ArtistStreamCreateDataList) (*emptypb.Empty, error) {
	err := s.artistUsecase.CreateStreamsByArtistIDs(ctx, model.ArtistStreamCreateDataFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ArtistService) GetArtistsListenedByUserID(ctx context.Context, req *artistProto.UserID) (*artistProto.ArtistListened, error) {
	artistsListened, err := s.artistUsecase.GetArtistsListenedByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &artistProto.ArtistListened{ArtistsListened: artistsListened}, nil
}

func (s *ArtistService) LikeArtist(ctx context.Context, req *artistProto.LikeRequest) (*emptypb.Empty, error) {
	err := s.artistUsecase.LikeArtist(ctx, model.LikeRequestFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ArtistService) GetFavoriteArtists(ctx context.Context, req *artistProto.FiltersWithUserID) (*artistProto.ArtistList, error) {
	artists, err := s.artistUsecase.GetFavoriteArtists(ctx, model.ArtistFiltersFromProtoToUsecase(req.Filters), req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromUsecaseToProto(artists), nil
}

func (s *ArtistService) SearchArtists(ctx context.Context, req *artistProto.Query) (*artistProto.ArtistList, error) {
	artists, err := s.artistUsecase.SearchArtists(ctx, req.Query, req.UserId.Id)
	if err != nil {
		return nil, err
	}
	return model.ArtistListFromUsecaseToProto(artists), nil
}
