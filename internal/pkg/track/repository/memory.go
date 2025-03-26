package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

type TrackMemoryRepository struct {
	mu     sync.RWMutex
	tracks map[uint]*model.Track
}

func NewTrackMemoryRepository() track.Repository {
	repo := &TrackMemoryRepository{
		tracks: map[uint]*model.Track{},
	}

	inabakumori := &model.Artist{ID: 1, Title: "Inabakumori", Thumbnail: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg"}
	yoasobi := &model.Artist{ID: 2, Title: "YOASOBI", Thumbnail: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg"}
	kenshi := &model.Artist{ID: 3, Title: "Kenshi Yonezu", Thumbnail: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg"}

	anticyclone := &model.Album{ID: 1, Title: "Anticyclone", Thumbnail: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", ArtistID: 1, Artist: *inabakumori}
	theBook := &model.Album{ID: 2, Title: "THE BOOK", Thumbnail: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", ArtistID: 2, Artist: *yoasobi}
	bootleg := &model.Album{ID: 3, Title: "BOOTLEG", Thumbnail: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", ArtistID: 3, Artist: *kenshi}

	testTracks := []*model.Track{
		{ID: 1, Title: "Lagtrain", Thumbnail: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", Album: anticyclone, Artist: inabakumori},
		{ID: 2, Title: "Lost Umbrella", Thumbnail: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg", Album: anticyclone, Artist: inabakumori},
		{ID: 3, Title: "Racing Into The Night", Thumbnail: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", Album: theBook, Artist: yoasobi},
		{ID: 4, Title: "Idol", Thumbnail: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg", Album: theBook, Artist: yoasobi},
		{ID: 5, Title: "KICK BACK", Thumbnail: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", Album: bootleg, Artist: kenshi},
		{ID: 6, Title: "Lemon", Thumbnail: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg", Album: bootleg, Artist: kenshi},
	}

	for _, track := range testTracks {
		repo.tracks[track.ID] = track
	}

	return repo
}

func (r *TrackMemoryRepository) GetAllTracks(filters *model.TrackFilters) ([]*model.Track, error) {
	offset := filters.Pagination.Offset
	limit := filters.Pagination.Limit

	if offset > len(r.tracks) {
		return []*model.Track{}, nil
	}

	if offset+limit > len(r.tracks) {
		limit = len(r.tracks) - offset
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	tracks := make([]*model.Track, 0, limit)
	for _, track := range r.tracks {
		tracks = append(tracks, track)
	}

	if len(tracks) == 0 {
		return []*model.Track{}, nil
	}

	return tracks, nil
}
