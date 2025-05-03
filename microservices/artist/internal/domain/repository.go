package domain

import (
	"context"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
)

type Repository interface {
	GetAllArtists(ctx context.Context, filters *repoModel.Filters, userID int64) ([]*repoModel.Artist, error)
	GetArtistByID(ctx context.Context, id int64, userID int64) (*repoModel.Artist, error)
	GetArtistTitleByID(ctx context.Context, id int64) (string, error)
	GetArtistsByTrackID(ctx context.Context, id int64) ([]*repoModel.ArtistWithRole, error)
	GetArtistsByTrackIDs(ctx context.Context, trackIDs []int64) (map[int64][]*repoModel.ArtistWithRole, error)
	GetArtistStats(ctx context.Context, id int64) (*repoModel.ArtistStats, error)
	GetArtistsByAlbumID(ctx context.Context, albumID int64) ([]*repoModel.ArtistWithTitle, error)
	GetArtistsByAlbumIDs(ctx context.Context, albumIDs []int64) (map[int64][]*repoModel.ArtistWithTitle, error)
	GetAlbumIDsByArtistID(ctx context.Context, id int64) ([]int64, error)
	GetTrackIDsByArtistID(ctx context.Context, id int64) ([]int64, error)
	CreateStreamsByArtistIDs(ctx context.Context, data *repoModel.ArtistStreamCreateDataList) error
	GetArtistsListenedByUserID(ctx context.Context, userID int64) (int64, error)
	LikeArtist(ctx context.Context, request *repoModel.LikeRequest) error
	UnlikeArtist(ctx context.Context, request *repoModel.LikeRequest) error
	CheckArtistExists(ctx context.Context, id int64) (bool, error)
	GetFavoriteArtists(ctx context.Context, filters *repoModel.Filters, userID int64) ([]*repoModel.Artist, error)
}
