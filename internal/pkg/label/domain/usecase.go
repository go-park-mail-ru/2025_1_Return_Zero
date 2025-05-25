package domain

import (
	"context"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
)

type Usecase interface {
	CreateLabel(ctx context.Context, label *usecaseModel.Label) (*usecaseModel.Label, error)
	GetLabel(ctx context.Context, id int64) (*usecaseModel.Label, error)
	CreateArtist(ctx context.Context, artist *usecaseModel.ArtistLoad) (*usecaseModel.Artist, error)
	EditArtist(ctx context.Context, artist *usecaseModel.ArtistEdit) (*usecaseModel.Artist, error)
	GetArtists(ctx context.Context, labelID int64, filters *usecaseModel.ArtistFilters) ([]*usecaseModel.Artist, error)
	DeleteArtist(ctx context.Context, artist *usecaseModel.ArtistDelete) error
	CreateAlbum(ctx context.Context, album *usecaseModel.CreateAlbumRequest) (int64, string, error)
	UpdateLabel(ctx context.Context, labelID int64, toAdd, toRemove []string) error
	DeleteAlbum(ctx context.Context, albumID, labelID int64) error
	GetAlbumsByLabelID(ctx context.Context, labelID int64, filters *usecaseModel.AlbumFilters) ([]*usecaseModel.Album, error)
}
