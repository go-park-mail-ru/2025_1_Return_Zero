package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AlbumError struct {
	Code    codes.Code
	Message string
}

func (e *AlbumError) Error() string {
	return e.Message
}

func (e *AlbumError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewNotFoundError(format string, args ...interface{}) error {
	return &AlbumError{
		Code:    codes.NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewInternalError(format string, args ...interface{}) error {
	return &AlbumError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrAlbumNotFound = NewNotFoundError("album not found")
)
