package models

import (
	"errors"
)

const (
	MaxOffset = 10000
	MaxLimit  = 100
)

type Filters struct {
	Offset int
	Limit  int
}

func (f *Filters) Validate() error {

	if f.Offset > MaxOffset {
		f.Offset = MaxOffset
	}

	if f.Offset < 0 {
		f.Offset = 0
		return errors.New("invalid offset: should be greater than 0")
	}

	if f.Limit > MaxLimit {
		f.Limit = MaxLimit
	}

	if f.Limit < 0 {
		f.Limit = 0
		return errors.New("invalid limit: should be greater than 0")
	}

	return nil
}

// TODO: Probably should add metadata, to give frontend info about numbers of pages and other things
