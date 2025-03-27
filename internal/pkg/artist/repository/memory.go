package repository

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrArtistNotFound = errors.New("artist not found")
)

type artistMemoryRepository struct {
	mu      sync.RWMutex
	artists map[uint]*repoModel.Artist
}

func NewArtistMemoryRepository() artist.Repository {
	repo := &artistMemoryRepository{
		artists: map[uint]*repoModel.Artist{},
	}

	testArtists := []*repoModel.Artist{
		{ID: 1, Title: "Inabakumori", Thumbnail: "https://i1.sndcdn.com/artworks-000640888066-bwv7e8-t500x500.jpg"},
		{ID: 2, Title: "YOASOBI", Thumbnail: "https://i.scdn.co/image/ab67616100005174bfdd8a29d0c6bc6950055234"},
		{ID: 3, Title: "Kenshi Yonezu", Thumbnail: "https://i.scdn.co/image/ab6761610000e5ebd7ca899f6e53b54976a8594b"},
		{ID: 4, Title: "RADWIMPS", Thumbnail: "https://i.scdn.co/image/ab6761610000e5ebc9d443fb5ced1dd32d106632"},
		{ID: 5, Title: "Official HIGE DANdism", Thumbnail: "https://i.scdn.co/image/ab6761610000e5ebf9f7513528a90d1dde6d3aaa"},
	}

	for _, artist := range testArtists {
		repo.artists[artist.ID] = artist
	}

	return repo
}

func (r *artistMemoryRepository) GetAllArtists(filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error) {
	offset := filters.Pagination.Offset
	limit := filters.Pagination.Limit

	if offset > len(r.artists) {
		return []*repoModel.Artist{}, nil
	}

	if offset+limit > len(r.artists) {
		limit = len(r.artists) - offset
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	artists := make([]*repoModel.Artist, 0, limit)
	for _, artist := range r.artists {
		artists = append(artists, artist)
	}

	if len(artists) == 0 {
		return []*repoModel.Artist{}, nil
	}

	return artists, nil
}

func (r *artistMemoryRepository) GetArtistByID(id uint) (*repoModel.Artist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artist, ok := r.artists[id]
	if !ok {
		return nil, ErrArtistNotFound
	}

	return artist, nil
}
