package customErrors

import "errors"

var (
	ErrArtistNotFound               = errors.New("artist not found")
	ErrInvalidOffset                = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit                 = errors.New("invalid limit: should be greater than 0")
	ErrAlbumNotFound                = errors.New("album not found")
	ErrStreamNotFound               = errors.New("stream not found")
	ErrFailedToUpdateStreamDuration = errors.New("failed to update stream duration")
	ErrTrackNotFound                = errors.New("track not found")
	ErrStreamPermissionDenied       = errors.New("user does not have permission to update this stream")
)
