package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

const (
	MaxOffset = 10000
	MaxLimit  = 100
)

var (
	ErrInvalidOffset = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit  = errors.New("invalid limit: should be greater than 0")
)

func validatePagination(p *model.Pagination) error {
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

func Pagination(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := &model.Pagination{}

		offset, err := helpers.ReadInt(r.URL.Query(), "offset", 0)
		if err != nil {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		limit, err := helpers.ReadInt(r.URL.Query(), "limit", 10)
		if err != nil {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		pagination.Offset = offset
		pagination.Limit = limit

		err = validatePagination(pagination)
		if err != nil {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), model.PaginationKey, pagination)
		next(w, r.WithContext(ctx))
	}
}

func PaginationFromContext(ctx context.Context) *model.Pagination {
	return ctx.Value(model.PaginationKey).(*model.Pagination)
}
