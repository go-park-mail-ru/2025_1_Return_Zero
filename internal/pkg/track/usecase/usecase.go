package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/trackFile"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

func NewUsecase(trackRepository track.Repository, artistRepository artist.Repository, albumRepository album.Repository, trackFileRepository trackFile.Repository, userRepository user.Repository) track.Usecase {
	return trackUsecase{trackRepo: trackRepository, artistRepo: artistRepository, albumRepo: albumRepository, trackFileRepo: trackFileRepository, userRepo: userRepository}
}

type trackUsecase struct {
	trackRepo     track.Repository
	artistRepo    artist.Repository
	albumRepo     album.Repository
	trackFileRepo trackFile.Repository
	userRepo      user.Repository
}

func (u trackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := &repoModel.TrackFilters{
		Pagination: model.PaginationFromUsecaseToRepository(filters.Pagination),
	}
	repoTracks, err := u.trackRepo.GetAllTracks(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	trackIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		trackIDs = append(trackIDs, repoTrack.ID)
	}

	albumIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		albumIDs = append(albumIDs, repoTrack.AlbumID)
	}

	repoArtists, err := u.artistRepo.GetArtistsByTrackIDs(ctx, trackIDs)
	if err != nil {
		return nil, err
	}

	albumTitles, err := u.albumRepo.GetAlbumTitleByIDs(ctx, albumIDs)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		track := model.TrackFromRepositoryToUsecase(repoTrack, repoArtists[repoTrack.ID], albumTitles[repoTrack.AlbumID])
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u trackUsecase) GetTrackByID(ctx context.Context, id int64) (*usecaseModel.TrackDetailed, error) {
	repoTrack, err := u.trackRepo.GetTrackByID(ctx, id)
	if err != nil {
		return nil, err
	}

	repoArtists, err := u.artistRepo.GetArtistsByTrackID(ctx, repoTrack.ID)
	if err != nil {
		return nil, err
	}

	albumTitle, err := u.albumRepo.GetAlbumTitleByID(ctx, repoTrack.AlbumID)
	if err != nil {
		return nil, err
	}

	trackFileUrl, err := u.trackFileRepo.GetPresignedURL(repoTrack.FileKey)
	if err != nil {
		return nil, err
	}

	trackDetailed := model.TrackDetailedFromRepositoryToUsecase(repoTrack, repoArtists, albumTitle, trackFileUrl)

	return trackDetailed, nil
}

func (u trackUsecase) GetTracksByArtistID(ctx context.Context, id int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetTracksByArtistID(ctx, id)
	if err != nil {
		return nil, err
	}

	trackIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		trackIDs = append(trackIDs, repoTrack.ID)
	}

	albumIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		albumIDs = append(albumIDs, repoTrack.AlbumID)
	}

	albumTitles, err := u.albumRepo.GetAlbumTitleByIDs(ctx, albumIDs)
	if err != nil {
		return nil, err
	}

	repoArtists, err := u.artistRepo.GetArtistsByTrackIDs(ctx, trackIDs)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		track := model.TrackFromRepositoryToUsecase(repoTrack, repoArtists[repoTrack.ID], albumTitles[repoTrack.AlbumID])
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u trackUsecase) CreateStream(ctx context.Context, stream *usecaseModel.TrackStreamCreateData) (int64, error) {
	repoTrackStreamCreateData := model.TrackStreamCreateDataFromUsecaseToRepository(stream)
	streamID, err := u.trackRepo.CreateStream(ctx, repoTrackStreamCreateData)
	if err != nil {
		return 0, err
	}

	return streamID, nil
}

func (u trackUsecase) UpdateStreamDuration(ctx context.Context, endedStream *usecaseModel.TrackStreamUpdateData) error {
	repoTrackStream, err := u.trackRepo.GetStreamByID(ctx, endedStream.StreamID)
	if err != nil {
		return err
	}

	if repoTrackStream.UserID != endedStream.UserID {
		return track.ErrStreamPermissionDenied
	}

	err = u.trackRepo.UpdateStreamDuration(ctx, model.TrackStreamUpdateDataFromUsecaseToRepository(endedStream))
	if err != nil {
		return err
	}

	return nil
}

func (u trackUsecase) GetLastListenedTracks(ctx context.Context, username string, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := &repoModel.TrackFilters{
		Pagination: model.PaginationFromUsecaseToRepository(filters.Pagination),
	}
	logger := middleware.LoggerFromContext(ctx)
	userID, err := u.userRepo.GetIDByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	repoStreams, err := u.trackRepo.GetStreamsByUserID(ctx, userID, repoFilters)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(repoStreams))
	for _, stream := range repoStreams {
		ids = append(ids, stream.TrackID)
	}

	repoTracks, err := u.trackRepo.GetTracksByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	trackIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		trackIDs = append(trackIDs, repoTrack.ID)
	}

	albumIDs := make([]int64, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		albumIDs = append(albumIDs, repoTrack.AlbumID)
	}

	albumTitles, err := u.albumRepo.GetAlbumTitleByIDs(ctx, albumIDs)
	if err != nil {
		return nil, err
	}

	repoArtists, err := u.artistRepo.GetArtistsByTrackIDs(ctx, trackIDs)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		track := model.TrackFromRepositoryToUsecase(repoTrack, repoArtists[repoTrack.ID], albumTitles[repoTrack.AlbumID])
		tracks = append(tracks, track)
	}

	return tracks, nil
}
