package helpers

import (
	"errors"
	"net/http"

	deliveryModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
)

const (
	MaxOffset = 10000
	MaxLimit  = 100
)

var (
	ErrInvalidOffset = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit  = errors.New("invalid limit: should be greater than 0")
)

func validatePagination(p *deliveryModel.Pagination) error {
	if p.Offset > MaxOffset {
		p.Offset = MaxOffset
	}

	if p.Offset < 0 {
		p.Offset = 0
		return ErrInvalidOffset
	}

	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}

	if p.Limit < 0 {
		p.Limit = 0
		return ErrInvalidLimit
	}

	return nil

}

func GetPagination(r *http.Request) (*deliveryModel.Pagination, error) {
	pagination := &deliveryModel.Pagination{}

	offset, err := ReadInt(r.URL.Query(), "offset", 0)
	if err != nil {
		return nil, err
	}

	limit, err := ReadInt(r.URL.Query(), "limit", 10)
	if err != nil {
		return nil, err
	}

	pagination.Offset = offset
	pagination.Limit = limit

	err = validatePagination(pagination)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
