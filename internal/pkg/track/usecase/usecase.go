package usecase

import (
	"context"

	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	artistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	playlistProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	trackProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/ctxExtractor"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

func NewUsecase(trackClient trackProto.TrackServiceClient, artistClient artistProto.ArtistServiceClient, albumClient albumProto.AlbumServiceClient, playlistClient playlistProto.PlaylistServiceClient, userClient userProto.UserServiceClient) track.Usecase {
	return &trackUsecase{trackClient: trackClient, artistClient: artistClient, albumClient: albumClient, playlistClient: playlistClient, userClient: userClient}
}

type trackUsecase struct {
	trackClient    trackProto.TrackServiceClient
	artistClient   artistProto.ArtistServiceClient
	albumClient    albumProto.AlbumServiceClient
	playlistClient playlistProto.PlaylistServiceClient
	userClient     userProto.UserServiceClient
}

func (u *trackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoFilters := &trackProto.UserIDWithFilters{
		UserId:  &trackProto.UserID{Id: userID},
		Filters: &trackProto.Filters{Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination)},
	}
	protoTracks, err := u.trackClient.GetAllTracks(ctx, protoFilters)
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

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
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
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoTrack, err := u.trackClient.GetTrackByID(ctx, &trackProto.TrackIDWithUserID{TrackId: &trackProto.TrackID{Id: id}, UserId: &trackProto.UserID{Id: userID}})
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackID(ctx, &artistProto.TrackID{Id: id})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitle, err := u.albumClient.GetAlbumTitleByID(ctx, &albumProto.AlbumID{Id: protoTrack.Track.AlbumId})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	trackDetailed := model.TrackDetailedFromProtoToUsecase(protoTrack, protoAlbumTitle, protoArtists)

	return trackDetailed, nil
}

func (u *trackUsecase) GetTracksByArtistID(ctx context.Context, id int64, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	artistTrackIDs, err := u.artistClient.GetTrackIDsByArtistID(ctx, &artistProto.ArtistID{Id: id})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoFilters := &trackProto.TrackIDListWithFilters{
		Ids:     model.TrackIDListFromArtistToTrackProto(artistTrackIDs, userID),
		Filters: &trackProto.Filters{Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination)},
	}

	protoTracks, err := u.trackClient.GetTracksByIDsFiltered(ctx, protoFilters)
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

	protoAlbumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
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
	streamID, err := u.trackClient.CreateStream(ctx, protoTrackStreamCreateData)
	if err != nil {
		return 0, customErrors.HandleTrackGRPCError(err)
	}

	albumID, err := u.trackClient.GetAlbumIDByTrackID(ctx, &trackProto.TrackID{Id: stream.TrackID})
	if err != nil {
		return 0, customErrors.HandleTrackGRPCError(err)
	}

	artists, err := u.artistClient.GetArtistsByTrackID(ctx, &artistProto.TrackID{Id: stream.TrackID})
	if err != nil {
		return 0, customErrors.HandleArtistGRPCError(err)
	}

	artistIDs := make([]int64, 0, len(artists.Artists))
	for _, artist := range artists.Artists {
		artistIDs = append(artistIDs, artist.Id)
	}

	_, err = u.artistClient.CreateStreamsByArtistIDs(ctx, model.ArtistStreamCreateDataListFromUsecaseToProto(stream.UserID, artistIDs))
	if err != nil {
		return 0, customErrors.HandleArtistGRPCError(err)
	}

	_, err = u.albumClient.CreateStream(ctx, &albumProto.AlbumStreamCreateData{
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
	_, err := u.trackClient.UpdateStreamDuration(ctx, protoTrackStreamUpdateData)
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

	protoTracks, err := u.trackClient.GetLastListenedTracks(ctx, protoUserIDWithFilters)
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

	albumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
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
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoAlbumID := &trackProto.AlbumIDWithUserID{
		AlbumId: &trackProto.AlbumID{Id: id},
		UserId:  &trackProto.UserID{Id: userID},
	}
	protoTracks, err := u.trackClient.GetTracksByAlbumID(ctx, protoAlbumID)
	if err != nil {
		return nil, customErrors.HandleTrackGRPCError(err)
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumTitle, err := u.albumClient.GetAlbumTitleByID(ctx, &albumProto.AlbumID{Id: id})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
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

func (u *trackUsecase) LikeTrack(ctx context.Context, request *usecaseModel.TrackLikeRequest) error {
	protoRequest := model.TrackLikeRequestFromUsecaseToProto(request)
	_, err := u.trackClient.LikeTrack(ctx, protoRequest)
	if err != nil {
		return customErrors.HandleTrackGRPCError(err)
	}
	return nil
}

func (u *trackUsecase) GetPlaylistTracks(ctx context.Context, id int64) ([]*usecaseModel.Track, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoPlaylistTrackIds, err := u.playlistClient.GetPlaylistTrackIds(ctx, &playlistProto.GetPlaylistTrackIdsRequest{
		PlaylistId: id,
		UserId:     userID,
	})
	if err != nil {
		return nil, customErrors.HandlePlaylistGRPCError(err)
	}

	if len(protoPlaylistTrackIds.TrackIds) == 0 {
		return make([]*usecaseModel.Track, 0), nil
	}

	trackIDList := make([]*trackProto.TrackID, 0, len(protoPlaylistTrackIds.TrackIds))
	for _, trackID := range protoPlaylistTrackIds.TrackIds {
		trackIDList = append(trackIDList, &trackProto.TrackID{Id: trackID})
	}

	protoTracks, err := u.trackClient.GetTracksByIDs(ctx, &trackProto.TrackIDList{Ids: trackIDList, UserId: &trackProto.UserID{Id: userID}})
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

	albumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
	if err != nil {
		return nil, customErrors.HandleAlbumGRPCError(err)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
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

func (u *trackUsecase) GetFavoriteTracks(ctx context.Context, filters *usecaseModel.TrackFilters, username string) ([]*usecaseModel.Track, error) {
	profileUserID, err := u.userClient.GetIDByUsername(ctx, &userProto.Username{Username: username})
	if err != nil {
		return nil, customErrors.HandleUserGRPCError(err)
	}

	requestUserID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		requestUserID = -1
	}

	profilePrivacy, err := u.userClient.GetUserPrivacyByID(ctx, &userProto.UserID{Id: profileUserID.Id})
	if err != nil {
		return nil, customErrors.HandleUserGRPCError(err)
	}

	if !profilePrivacy.IsPublicFavoriteTracks && requestUserID != profileUserID.Id {
		return make([]*usecaseModel.Track, 0), nil
	}

	protoFilters := &trackProto.FavoriteRequest{
		ProfileUserId: &trackProto.UserID{Id: profileUserID.Id},
		RequestUserId: &trackProto.UserID{Id: requestUserID},
		Filters:       &trackProto.Filters{Pagination: model.PaginationFromUsecaseToTrackProto(filters.Pagination)},
	}

	protoTracks, err := u.trackClient.GetFavoriteTracks(ctx, protoFilters)
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

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
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

func (u *trackUsecase) SearchTracks(ctx context.Context, query string) ([]*usecaseModel.Track, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	protoTracks, err := u.trackClient.SearchTracks(ctx, &trackProto.Query{Query: query, UserId: &trackProto.UserID{Id: userID}})
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

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
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

func (u *trackUsecase) GetSelectionTracks(ctx context.Context, selection string) ([]*usecaseModel.Track, error) {
	userID, exists := ctxExtractor.UserFromContext(ctx)
	if !exists {
		userID = -1
	}

	var protoTracks *trackProto.TrackList
	var err error
	switch selection {
	case "most-liked":
		protoTracks, err = u.trackClient.GetMostLikedTracks(ctx, &trackProto.UserID{Id: userID})
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}
	case "most-recent":
		protoTracks, err = u.trackClient.GetMostRecentTracks(ctx, &trackProto.UserID{Id: userID})
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}
	case "most-listened-last-month":
		protoTracks, err = u.trackClient.GetMostListenedLastMonthTracks(ctx, &trackProto.UserID{Id: userID})
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}
	case "most-liked-last-week":
		protoTracks, err = u.trackClient.GetMostLikedLastWeekTracks(ctx, &trackProto.UserID{Id: userID})
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}
	case "top-chart":
		protoTracks, err = u.getMostListenedFromMostListenedArtists(ctx, userID)
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}
	default:
		return nil, customErrors.ErrInvalidSelection
	}

	trackIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		trackIDs = append(trackIDs, protoTrack.Id)
	}

	albumIDs := make([]int64, 0, len(protoTracks.Tracks))
	for _, protoTrack := range protoTracks.Tracks {
		albumIDs = append(albumIDs, protoTrack.AlbumId)
	}

	protoArtists, err := u.artistClient.GetArtistsByTrackIDs(ctx, &artistProto.TrackIDList{Ids: model.TrackIdsFromUsecaseToArtistProto(trackIDs)})
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	protoAlbumTitles, err := u.albumClient.GetAlbumTitleByIDs(ctx, &albumProto.AlbumIDList{Ids: model.AlbumIdsFromUsecaseToAlbumProto(albumIDs)})
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

func (u *trackUsecase) getMostListenedFromMostListenedArtists(ctx context.Context, userID int64) (*trackProto.TrackList, error) {
	protoFilters := &artistProto.FiltersWithUserID{
		Filters: &artistProto.Filters{
			Pagination: &artistProto.Pagination{
				Offset: 0,
				Limit:  5, // TOP 5
			},
		},
		UserId: &artistProto.UserID{Id: -1},
	}

	mostListenedArtists, err := u.artistClient.GetAllArtists(ctx, protoFilters)
	if err != nil {
		return nil, customErrors.HandleArtistGRPCError(err)
	}

	var combinedProtoTracks trackProto.TrackList
	for _, artist := range mostListenedArtists.Artists {
		artistTrackIDs, err := u.artistClient.GetTrackIDsByArtistID(ctx, &artistProto.ArtistID{Id: artist.Id})
		if err != nil {
			return nil, customErrors.HandleArtistGRPCError(err)
		}

		protoFilters := &trackProto.TrackIDListWithFilters{
			Ids:     model.TrackIDListFromArtistToTrackProto(artistTrackIDs, userID),
			Filters: &trackProto.Filters{Pagination: model.PaginationFromUsecaseToTrackProto(&usecaseModel.Pagination{Offset: 0, Limit: 3})},
		}

		protoTracks, err := u.trackClient.GetTracksByIDsFiltered(ctx, protoFilters)
		if err != nil {
			return nil, customErrors.HandleTrackGRPCError(err)
		}

		combinedProtoTracks.Tracks = append(combinedProtoTracks.Tracks, protoTracks.Tracks...)
	}

	return &combinedProtoTracks, nil
}
