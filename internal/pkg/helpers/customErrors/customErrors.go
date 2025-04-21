package customErrors

import "errors"

var (
	ErrArtistNotFound = errors.New("artist not found")
	ErrInvalidOffset  = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit   = errors.New("invalid limit: should be greater than 0")
)
