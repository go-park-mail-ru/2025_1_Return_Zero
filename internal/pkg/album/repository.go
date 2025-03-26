package album

import (
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

type Repository interface {
	GetAllAlbums(filters *model.AlbumFilters) ([]*model.AlbumDB, error)
	GetAlbumByID(id uint) (*model.AlbumDB, error)
}
