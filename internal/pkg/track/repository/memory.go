package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

type TrackMemoryRepository struct {
	mu     sync.RWMutex
	tracks map[uint]*model.TrackDB
}

func NewTrackMemoryRepository() track.Repository {
	repo := &TrackMemoryRepository{
		tracks: map[uint]*model.TrackDB{},
	}

	testTracks := []*model.TrackDB{
		{ID: 1, Title: "Lagtrain", Thumbnail: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", AlbumID: 1, ArtistID: 1, Duration: 216},
		{ID: 2, Title: "Lost Umbrella", Thumbnail: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg", AlbumID: 1, ArtistID: 1, Duration: 216},
		{ID: 3, Title: "Racing Into The Night", Thumbnail: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", AlbumID: 2, ArtistID: 2, Duration: 216},
		{ID: 4, Title: "Idol", Thumbnail: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg", AlbumID: 2, ArtistID: 2, Duration: 216},
		{ID: 5, Title: "KICK BACK", Thumbnail: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", AlbumID: 3, ArtistID: 3, Duration: 216},
		{ID: 6, Title: "Lemon", Thumbnail: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg", AlbumID: 3, ArtistID: 3, Duration: 216},
	}

	for _, track := range testTracks {
		repo.tracks[track.ID] = track
	}

	return repo
}

func (r *TrackMemoryRepository) GetAllTracks(filters *model.TrackFilters) ([]*model.TrackDB, error) {
	offset := filters.Pagination.Offset
	limit := filters.Pagination.Limit

	if offset > len(r.tracks) {
		return []*model.TrackDB{}, nil
	}

	if offset+limit > len(r.tracks) {
		limit = len(r.tracks) - offset
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	tracks := make([]*model.TrackDB, 0, limit)
	for _, track := range r.tracks {
		tracks = append(tracks, track)
	}

	if len(tracks) == 0 {
		return []*model.TrackDB{}, nil
	}

	return tracks, nil
}
