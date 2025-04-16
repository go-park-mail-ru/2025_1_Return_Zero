package repository

import (
	"errors"
)

type Genre struct {
	ID   int64  `sql:"id"`
	Name string `sql:"name"`
}

var (
	ErrGenreNotFound = errors.New("genre not found")
)
