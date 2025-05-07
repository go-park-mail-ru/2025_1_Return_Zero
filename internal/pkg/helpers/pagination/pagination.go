package pagination

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	customErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	query "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/query"
	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
)

func validatePagination(p *deliveryModel.Pagination, cfg *config.PaginationConfig) error {
	if p.Offset > cfg.MaxOffset {
		p.Offset = cfg.MaxOffset
	}

	if p.Offset < 0 {
		p.Offset = 0
		return customErrors.ErrInvalidOffset
	}

	if p.Limit > cfg.MaxLimit {
		p.Limit = cfg.MaxLimit
	}

	if p.Limit < 0 {
		p.Limit = 0
		return customErrors.ErrInvalidLimit
	}

	return nil

}

func GetPagination(r *http.Request, cfg *config.PaginationConfig) (*deliveryModel.Pagination, error) {
	pagination := &deliveryModel.Pagination{}

	offset, err := query.ReadInt(r.URL.Query(), "offset", cfg.DefaultOffset)
	if err != nil {
		return nil, customErrors.ErrInvalidOffset
	}

	limit, err := query.ReadInt(r.URL.Query(), "limit", cfg.DefaultLimit)
	if err != nil {
		return nil, customErrors.ErrInvalidLimit
	}

	pagination.Offset = offset
	pagination.Limit = limit

	err = validatePagination(pagination, cfg)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
