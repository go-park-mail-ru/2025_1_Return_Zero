package domain

import (
	"context"

	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
)

type Usecase interface {
	GetAllArtists(ctx context.Context, filters *usecaseModel.Filters, userID int64) (*usecaseModel.ArtistList, error)
	GetArtistByID(ctx context.Context, id int64, userID int64) (*usecaseModel.ArtistDetailed, error)
	GetArtistTitleByID(ctx context.Context, id int64) (string, error)
	GetArtistsByTrackID(ctx context.Context, id int64) (*usecaseModel.ArtistWithRoleList, error)
	GetArtistsByTrackIDs(ctx context.Context, trackIDs []int64) (*usecaseModel.ArtistWithRoleMap, error)
	GetArtistsByAlbumID(ctx context.Context, albumID int64) (*usecaseModel.ArtistWithTitleList, error)
	GetArtistsByAlbumIDs(ctx context.Context, albumIDs []int64) (*usecaseModel.ArtistWithTitleMap, error)
	GetAlbumIDsByArtistID(ctx context.Context, id int64) ([]int64, error)
	GetTrackIDsByArtistID(ctx context.Context, id int64) ([]int64, error)
	CreateStreamsByArtistIDs(ctx context.Context, data *usecaseModel.ArtistStreamCreateDataList) error
	GetArtistsListenedByUserID(ctx context.Context, userID int64) (int64, error)
	LikeArtist(ctx context.Context, request *usecaseModel.LikeRequest) error
	GetFavoriteArtists(ctx context.Context, filters *usecaseModel.Filters, userID int64) (*usecaseModel.ArtistList, error)
	SearchArtists(ctx context.Context, query string, userID int64) (*usecaseModel.ArtistList, error)
	CreateArtist(ctx context.Context, artist *usecaseModel.ArtistLoad) (*usecaseModel.Artist, error)
	EditArtist(ctx context.Context, artist *usecaseModel.ArtistEdit) (*usecaseModel.Artist, error)
	GetArtistsLabelID(ctx context.Context, filters *usecaseModel.Filters, labelID int64) (*usecaseModel.ArtistList, error)
	DeleteArtist(ctx context.Context, artist *usecaseModel.ArtistDelete) error
	ConnectArtists(ctx context.Context, artistIDs []int64, albumID int64, trackIDs []int64) error
}
