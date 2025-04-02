package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
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

func (u trackUsecase) GetAllTracks(filters *usecaseModel.TrackFilters) ([]*usecaseModel.Track, error) {
	repoFilters := &repoModel.TrackFilters{
		Pagination: &repoModel.Pagination{
			Offset: filters.Pagination.Offset,
			Limit:  filters.Pagination.Limit,
		},
	}
	repoTracks, err := u.trackRepo.GetAllTracks(repoFilters)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		repoArtists, err := u.artistRepo.GetArtistsByTrackID(repoTrack.ID)
		if err != nil {
			return nil, err
		}

		artists := make([]*usecaseModel.TrackArtist, 0, len(repoArtists))
		for _, repoArtist := range repoArtists {
			artists = append(artists, &usecaseModel.TrackArtist{
				ID:    repoArtist.ID,
				Title: repoArtist.Title,
				Role:  repoArtist.Role,
			})
		}

		albumTitle, err := u.albumRepo.GetAlbumTitleByID(repoTrack.AlbumID)
		if err != nil {
			return nil, err
		}

		track := &usecaseModel.Track{
			ID:        repoTrack.ID,
			Title:     repoTrack.Title,
			Thumbnail: repoTrack.Thumbnail,
			Duration:  repoTrack.Duration,
			Artists:   artists,
			Album:     albumTitle,
			AlbumID:   repoTrack.AlbumID,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (u trackUsecase) GetTrackByID(id int64) (*usecaseModel.TrackDetailed, error) {
	repoTrack, err := u.trackRepo.GetTrackByID(id)
	if err != nil {
		return nil, err
	}

	repoArtists, err := u.artistRepo.GetArtistsByTrackID(repoTrack.ID)
	if err != nil {
		return nil, err
	}

	artists := make([]*usecaseModel.TrackArtist, 0, len(repoArtists))
	for _, repoArtist := range repoArtists {
		artists = append(artists, &usecaseModel.TrackArtist{
			ID:    repoArtist.ID,
			Title: repoArtist.Title,
			Role:  repoArtist.Role,
		})
	}

	albumTitle, err := u.albumRepo.GetAlbumTitleByID(repoTrack.AlbumID)
	if err != nil {
		return nil, err
	}

	trackFileUrl, err := u.trackFileRepo.GetPresignedURL(repoTrack.FileKey)
	if err != nil {
		return nil, err
	}

	track := &usecaseModel.TrackDetailed{
		Track: usecaseModel.Track{
			ID:        repoTrack.ID,
			Title:     repoTrack.Title,
			Thumbnail: repoTrack.Thumbnail,
			Duration:  repoTrack.Duration,
			Album:     albumTitle,
			AlbumID:   repoTrack.AlbumID,
			Artists:   artists,
		},
		FileUrl: trackFileUrl,
	}

	return track, nil
}

func (u trackUsecase) GetTracksByArtistID(id int64) ([]*usecaseModel.Track, error) {
	repoTracks, err := u.trackRepo.GetTracksByArtistID(id)
	if err != nil {
		return nil, err
	}

	tracks := make([]*usecaseModel.Track, 0, len(repoTracks))
	for _, repoTrack := range repoTracks {
		repoArtists, err := u.artistRepo.GetArtistsByTrackID(repoTrack.ID)
		if err != nil {
			return nil, err
		}

		artists := make([]*usecaseModel.TrackArtist, 0, len(repoArtists))
		for _, repoArtist := range repoArtists {
			artists = append(artists, &usecaseModel.TrackArtist{
				ID:    repoArtist.ID,
				Title: repoArtist.Title,
				Role:  repoArtist.Role,
			})
		}

		albumTitle, err := u.albumRepo.GetAlbumTitleByID(repoTrack.AlbumID)
		if err != nil {
			return nil, err
		}

		track := &usecaseModel.Track{
			ID:        repoTrack.ID,
			Title:     repoTrack.Title,
			Thumbnail: repoTrack.Thumbnail,
			Duration:  repoTrack.Duration,
			Album:     albumTitle,
			AlbumID:   repoTrack.AlbumID,
			Artists:   artists,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
