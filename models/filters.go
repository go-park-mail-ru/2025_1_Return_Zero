package models

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	MaxOffset = 10000
	MaxLimit  = 100
)

type Filters struct {
	Offset int
	Limit  int
}

func (f Filters) Validate() error {

	if f.Offset > MaxOffset {
		f.Offset = MaxOffset
	}

	if f.Offset < 0 {
		return errors.New("invalid offset: should be greater than 0")
	}

	if f.Limit > MaxLimit {
		f.Limit = MaxLimit
	}

	if f.Limit < 0 {
		return errors.New("invalid limit: should be greater than 0")
	}

	return nil
}

// ParseFromQuery parses the filters from URL query parameters
func ParseFromQuery(r *http.Request) Filters {
	var f Filters

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if offsetNum, err := strconv.Atoi(offset); err == nil {
			f.Offset = offsetNum
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil {
			f.Limit = limitNum
		}
	}

	return f
}

// TODO: Probably should add metadata, to give frontend info about numbers of pages and other things
