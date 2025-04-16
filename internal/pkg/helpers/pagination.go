package helpers

import (
	"errors"
	"net/http"

	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
)

var (
	ErrInvalidOffset = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit  = errors.New("invalid limit: should be greater than 0")
)

func validatePagination(p *deliveryModel.Pagination, cfg *deliveryModel.PaginationConfig) error {
	if p.Offset > cfg.MaxOffset {
		p.Offset = cfg.MaxOffset
	}

	if p.Offset < 0 {
		p.Offset = 0
		return ErrInvalidOffset
	}

	if p.Limit > cfg.MaxLimit {
		p.Limit = cfg.MaxLimit
	}

	if p.Limit < 0 {
		p.Limit = 0
		return ErrInvalidLimit
	}

	return nil

}

func GetPagination(r *http.Request, cfg *deliveryModel.PaginationConfig) (*deliveryModel.Pagination, error) {
	pagination := &deliveryModel.Pagination{}

	offset, err := ReadInt(r.URL.Query(), "offset", cfg.DefaultOffset)
	if err != nil {
		return nil, err
	}

	limit, err := ReadInt(r.URL.Query(), "limit", cfg.DefaultLimit)
	if err != nil {
		return nil, err
	}

	pagination.Offset = offset
	pagination.Limit = limit

	err = validatePagination(pagination, cfg)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
