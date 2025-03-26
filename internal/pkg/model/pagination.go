package model

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

const PaginationKey = "pagination"
