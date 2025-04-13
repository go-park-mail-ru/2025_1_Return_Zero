package repository

import "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"

func PaginationFromUsecaseToRepository(usecasePagination *usecase.Pagination) *Pagination {
	return &Pagination{
		Offset: usecasePagination.Offset,
		Limit:  usecasePagination.Limit,
	}
}

type Pagination struct {
	Offset int
	Limit  int
}
