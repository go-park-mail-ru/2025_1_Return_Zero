package usecase

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
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

func (u trackUsecase) GetAllTracks(filters *model.TrackFilters) ([]*model.Track, error) {
	tracksDB, err := u.trackRepo.GetAllTracks(filters)
	if err != nil {
		return nil, err
	}

	tracks := make([]*model.Track, 0, len(tracksDB))
	for _, trackDB := range tracksDB {
		artist, err := u.artistRepo.GetArtistByID(trackDB.ArtistID)
		if err != nil {
			return nil, err
		}
		track := &model.Track{
			ID:        trackDB.ID,
			Title:     trackDB.Title,
			Thumbnail: trackDB.Thumbnail,
			Duration:  trackDB.Duration,
			Artist:    artist,
		}

		album, err := u.albumRepo.GetAlbumByID(trackDB.AlbumID)
		if err != nil {
			return nil, err
		}
		track.Album = album
		tracks = append(tracks, track)
	}

	return tracks, nil
}
