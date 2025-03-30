package artist

import (
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

type Repository interface {
	GetAllArtists(filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error)
	GetArtistByID(id int64) (*repoModel.Artist, error)
	GetArtistTitleByID(id int64) (string, error)
	GetArtistsByTrackID(id int64) ([]*repoModel.ArtistWithRole, error)
	GetArtistListenersCount(id int64) (int64, error)
	GetArtistFavoritesCount(id int64) (int64, error)
}
