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
)

func NewUsecase(trackRepository track.Repository, artistRepository artist.Repository, albumRepository album.Repository, trackFileRepository trackFile.Repository) track.Usecase {
	return trackUsecase{trackRepo: trackRepository, artistRepo: artistRepository, albumRepo: albumRepository, trackFileRepo: trackFileRepository}
}

type trackUsecase struct {
	trackRepo     track.Repository
	artistRepo    artist.Repository
	albumRepo     album.Repository
	trackFileRepo trackFile.Repository
}

func (u trackUsecase) GetAllTracks(ctx context.Context, filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := &repoModel.TrackFilters{
		Pagination: model.PaginationFromUsecaseToRepository(filters.Pagination),
	}
	repoTracks, err := u.trackRepo.GetAllTracks(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		repoArtists, err := u.artistRepo.GetArtistsByTrackID(ctx, repoTrack.ID)
		if err != nil {
			return nil, err
		}

		albumTitle, err := u.albumRepo.GetAlbumTitleByID(ctx, repoTrack.AlbumID)
		if err != nil {
			return nil, err
		}

		track := model.TrackFromRepositoryToUsecase(repoTrack, repoArtists, albumTitle)
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

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		repoArtists, err := u.artistRepo.GetArtistsByTrackID(ctx, repoTrack.ID)
		if err != nil {
			return nil, err
		}

		albumTitle, err := u.albumRepo.GetAlbumTitleByID(ctx, repoTrack.AlbumID)
		if err != nil {
			return nil, err
		}

		track := model.TrackFromRepositoryToUsecase(repoTrack, repoArtists, albumTitle)

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
