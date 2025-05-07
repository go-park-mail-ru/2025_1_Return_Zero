package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlaylistError struct {
	Code    codes.Code
	Message string
}

func (e *PlaylistError) Error() string {
	return e.Message
}

func (e *PlaylistError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewNotFoundError(format string, args ...interface{}) error {
	return &PlaylistError{
		Code:    codes.NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewInternalError(format string, args ...interface{}) error {
	return &PlaylistError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewBadRequestError(format string, args ...interface{}) error {
	return &PlaylistError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewAlreadyExistsError(format string, args ...interface{}) error {
	return &PlaylistError{
		Code:    codes.AlreadyExists,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewPermissionDeniedError(format string, args ...interface{}) error {
	return &PlaylistError{
		Code:    codes.PermissionDenied,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrPlaylistNotFound         = NewNotFoundError("playlist not found")
	ErrPlaylistPermissionDenied = NewPermissionDeniedError("user does not have permission to update this playlist")
	ErrUnsupportedImageFormat   = NewBadRequestError("unsupported image format: only JPEG and PNG are allowed")
	ErrImageTooBig              = NewBadRequestError("image size exceeds 5MB limit")
	ErrFailedToParseImage       = NewInternalError("failed to parse image")
	ErrFailedToUploadImage      = NewInternalError("failed to upload image")
	ErrPlaylistDuplicate        = NewAlreadyExistsError("playlist with this title by you already exists")
	ErrPlaylistTrackNotFound    = NewNotFoundError("track not found in playlist")
	ErrPlaylistTrackDuplicate   = NewAlreadyExistsError("track already in playlist")
)
