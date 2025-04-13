package artist

import (
	"errors"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
)

var (
	ErrArtistNotFound = errors.New("artist not found")
)

type Repository interface {
	GetAllArtists(filters *repoModel.ArtistFilters) ([]*repoModel.Artist, error)
	GetArtistByID(id int64) (*repoModel.Artist, error)
	GetArtistTitleByID(id int64) (string, error)
	GetArtistsByTrackID(id int64) ([]*repoModel.ArtistWithRole, error)
	GetArtistStats(id int64) (*repoModel.ArtistStats, error)
	GetArtistsByAlbumID(albumID int64) ([]*repoModel.ArtistWithTitle, error)
}
