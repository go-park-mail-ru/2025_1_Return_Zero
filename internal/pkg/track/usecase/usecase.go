package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

func NewUsecase(trackRepository track.Repository, artistRepository artist.Repository, albumRepository album.Repository) track.Usecase {
	return trackUsecase{trackRepo: trackRepository, artistRepo: artistRepository, albumRepo: albumRepository}
}

type trackUsecase struct {
	trackRepo  track.Repository
	artistRepo artist.Repository
	albumRepo  album.Repository
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
