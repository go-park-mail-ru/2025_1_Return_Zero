package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type AlbumMemoryRepository struct {
	mu     sync.RWMutex
	albums map[uint]*model.Album
}

func NewAlbumMemoryRepository() album.Repository {
	repo := &AlbumMemoryRepository{
		albums: map[uint]*model.Album{},
	}

	artists := []model.Artist{
		{ID: 1, Title: "Inabakumori", Thumbnail: "https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e"},
		{ID: 2, Title: "YOASOBI", Thumbnail: "https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1"},
		{ID: 3, Title: "Kenshi Yonezu", Thumbnail: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
		{ID: 4, Title: "RADWIMPS", Thumbnail: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
	}

	testAlbums := []*model.Album{
		{ID: 1, Title: "Anticyclone", ArtistID: 1, Artist: artists[0], Thumbnail: "https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e"},
		{ID: 2, Title: "THE BOOK", ArtistID: 2, Artist: artists[1], Thumbnail: "https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1"},
		{ID: 3, Title: "BOOTLEG", ArtistID: 3, Artist: artists[2], Thumbnail: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s"},
		{ID: 4, Title: "Your Name.", ArtistID: 4, Artist: artists[3], Thumbnail: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s"},
	}

	for _, album := range testAlbums {
		repo.albums[album.ID] = album
	}
	return repo
}

func (r *AlbumMemoryRepository) GetAllAlbums(filters *model.AlbumFilters) ([]*model.Album, error) {
	offset := filters.Pagination.Offset
	limit := filters.Pagination.Limit

	if offset > len(r.albums) {
		return []*model.Album{}, nil
	}

	if offset+limit > len(r.albums) {
		limit = len(r.albums) - offset
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	albums := make([]*model.Album, 0, limit)
	for _, album := range r.albums {
		albums = append(albums, album)
	}

	if len(albums) == 0 {
		return []*model.Album{}, nil
	}

	return albums, nil
}
