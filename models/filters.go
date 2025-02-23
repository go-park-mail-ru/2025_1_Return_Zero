package models

import (
	"errors"
	"net/http"
	"strconv"
)

type Filters struct {
	Page     int
	PageSize int
}

func (f Filters) Validate() error {
	if f.Page < 1 {
		return errors.New("invalid page number: should be greater than 0")
	}

	if f.PageSize < 1 {
		return errors.New("invalid page size: should be greater than 0")
	}

	if f.PageSize > 100 {
		return errors.New("invalid page size: should be less than 100")
	}

	return nil
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ParseFromQuery parses the filters from URL query parameters
func ParseFromQuery(r *http.Request) Filters {
	var f Filters

	if page := r.URL.Query().Get("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			f.Page = pageNum
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			f.PageSize = size
		}
	}

	return f
}

// TODO: Probably should add metadata, to give frontend info about numbers of pages and other things
