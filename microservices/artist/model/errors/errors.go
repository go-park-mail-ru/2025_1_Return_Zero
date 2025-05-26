package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArtistError struct {
	Code    codes.Code
	Message string
}

func (e *ArtistError) Error() string {
	return e.Message
}

func (e *ArtistError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewNotFoundError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewInternalError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewBadRequestError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewEmptyS3KeyError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewUnsupportedImageFormatError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewFailedToUploadAvatarError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewConflictError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.AlreadyExists,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewForbiddenError(format string, args ...interface{}) error {
	return &ArtistError{
		Code:    codes.PermissionDenied,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrArtistNotFound = NewNotFoundError("artist not found")
)
