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
		repoArtist, err := u.artistRepo.GetArtistByID(repoTrack.ArtistID)
		if err != nil {
			return nil, err
		}
		artist := &usecaseModel.Artist{
			ID:        repoArtist.ID,
			Title:     repoArtist.Title,
			Thumbnail: repoArtist.Thumbnail,
		}

		track := &usecaseModel.Track{
			ID:        repoTrack.ID,
			Title:     repoTrack.Title,
			Thumbnail: repoTrack.Thumbnail,
			Duration:  repoTrack.Duration,
			Artist:    artist,
		}

		album, err := u.albumRepo.GetAlbumByID(repoTrack.AlbumID)
		if err != nil {
			return nil, err
		}
		track.Album = &usecaseModel.AlbumUnpopulated{
			ID:        album.ID,
			Title:     album.Title,
			Thumbnail: album.Thumbnail,
			ArtistID:  album.ArtistID,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
