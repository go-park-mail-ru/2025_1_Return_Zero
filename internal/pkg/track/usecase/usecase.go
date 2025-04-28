package usecase

import (
	"context"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

func NewUsecase(trackClient *trackProto.TrackServiceClient, artistClient *artistProto.ArtistServiceClient, albumClient *albumProto.AlbumServiceClient) track.Usecase {
	return &trackUsecase{trackClient: trackClient, artistClient: artistClient, albumClient: albumClient}
}

type trackUsecase struct {
	trackClient  *trackProto.TrackServiceClient
	artistClient *artistProto.ArtistServiceClient
	albumClient  *albumProto.AlbumServiceClient
}

func (u *trackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	protoFilters := &trackProto.Filters{
		Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination),
	}
	protoTracks, err := (*u.trackClient).GetAllTracks(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		albumIDs = append(albumIDs, protoTrack.AlbumId)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitles, err := (*u.albumClient).GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	tracks := make([]*usecaseModel.Track, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		track := model.TrackFromProtoToUsecase(protoTrack, protoAlbumTitles.Titles[protoTrack.AlbumId], protoArtists.Artists[protoTrack.Id])
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u *trackUsecase) GetTrackByID(ctx context.Context, id int64) (*usecaseModel.TrackDetailed, error) {
	protoTrack, err := (*u.trackClient).GetTrackByID(ctx, &trackProto.TrackID{Id: id})
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByTrackID(ctx, &artistProto.TrackID{Id: id})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitle, err := (*u.albumClient).GetAlbumTitleByID(ctx, &albumProto.AlbumID{Id: protoTrack.Track.AlbumId})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	trackDetailed := model.TrackDetailedFromProtoToUsecase(protoTrack, protoAlbumTitle, protoArtists)

	return trackDetailed, nil
}

func (u *trackUsecase) GetTracksByArtistID(ctx context.Context, id int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	protoFilters := &trackProto.Filters{
		Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination),
	}

	artistTrackIDs, err := (*u.artistClient).GetTrackIDsByArtistID(ctx, &artistProto.ArtistID{Id: id})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoTracks, err := (*u.trackClient).GetTracksByIDsFiltered(ctx, &trackProto.TrackIDListWithFilters{Ids: model.TrackIDListFromArtistToTrackProto(artistTrackIDs), Filters: protoFilters})
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		albumIDs = append(albumIDs, protoTrack.AlbumId)
	}

	protoAlbumTitles, err := (*u.albumClient).GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	tracks := make([]*usecaseModel.Track, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		track := model.TrackFromProtoToUsecase(protoTrack, protoAlbumTitles.Titles[protoTrack.AlbumId], protoArtists.Artists[protoTrack.Id])
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u *trackUsecase) CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error) {
	protoTrackStreamCreateData := model.TrackStreamCreateDataFromUsecaseToProto(stream)
	streamID, err := (*u.trackClient).CreateStream(ctx, protoTrackStreamCreateData)
	if err != nil {
		return 0, customErrors.HandleTrackGRPCError(err)
	}

	albumID, err := (*u.trackClient).GetAlbumIDByTrackID(ctx, &trackProto.TrackID{Id: stream.TrackID})
	if err != nil {
		return 0, customErrors.HandleTrackGRPCError(err)
	}

	artists, err := (*u.artistClient).GetArtistsByTrackID(ctx, &artistProto.TrackID{Id: stream.TrackID})
	if err != nil {
		return 0, customErrors.HandleArtistGRPCError(err)
	}

	artistIDs := make([]int64, 0, len(artists.Artists))
	for _, artist := range artists.Artists {
		artistIDs = append(artistIDs, artist.Id)
	}

	_, err = (*u.artistClient).CreateStreamsByArtistIDs(ctx, model.ArtistStreamCreateDataListFromUsecaseToProto(stream.UserID, artistIDs))
	if err != nil {
		return 0, customErrors.HandleArtistGRPCError(err)
	}

	_, err = (*u.albumClient).CreateStream(ctx, &albumProto.AlbumStreamCreateData{
		AlbumId: &albumProto.AlbumID{Id: albumID.Id},
		UserId:  &albumProto.UserID{Id: stream.UserID},
	})
	if err != nil {
		return 0, customErrors.HandleAlbumGRPCError(err)
	}

	return streamID.Id, nil
}

func (u *trackUsecase) UpdateStreamDuration(ctx context.Context, endedStream *usecaseModel.TrackStreamUpdateData) error {
	protoTrackStreamUpdateData := model.TrackStreamUpdateDataFromUsecaseToProto(endedStream)
	_, err := (*u.trackClient).UpdateStreamDuration(ctx, protoTrackStreamUpdateData)
	if err != nil {
		return customErrors.HandleTrackGRPCError(err)
	}

	return nil
}

func (u *trackUsecase) GetLastListenedTracks(ctx context.Context, userID int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {

	protoUserIDWithFilters := &trackProto.UserIDWithFilters{
		UserId:  &trackProto.UserID{Id: userID},
		Filters: &trackProto.Filters{Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination)},
	}

	protoTracks, err := (*u.trackClient).GetLastListenedTracks(ctx, protoUserIDWithFilters)
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		albumIDs = append(albumIDs, protoTrack.AlbumId)
	}

	albumTitles, err := (*u.albumClient).GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	tracks := make([]*usecaseModel.Track, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		track := model.TrackFromProtoToUsecase(protoTrack, albumTitles.Titles[protoTrack.AlbumId], protoArtists.Artists[protoTrack.Id])
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u *trackUsecase) GetTracksByAlbumID(ctx context.Context, id int64) ([]*usecaseModel.Track, error) {
	protoAlbumID := &trackProto.AlbumID{Id: id}
	protoTracks, err := (*u.trackClient).GetTracksByAlbumID(ctx, protoAlbumID)
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumTitle, err := (*u.albumClient).GetAlbumTitleByID(ctx, &albumProto.AlbumID{Id: id})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := (*u.artistClient).GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	tracks := make([]*usecaseModel.Track, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		track := model.TrackFromProtoToUsecase(protoTrack, albumTitle, protoArtists.Artists[protoTrack.Id])
		tracks = append(tracks, track)
	}

	return tracks, nil
}
