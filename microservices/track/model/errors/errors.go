package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TrackError struct {
	Code    codes.Code
	Message string
}

func (e *TrackError) Error() string {
	return e.Message
}

func (e *TrackError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewNotFoundError(format string, args ...interface{}) error {
	return &TrackError{
		Code:    codes.NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewInternalError(format string, args ...interface{}) error {
	return &TrackError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewPermissionDeniedError(format string, args ...interface{}) error {
	return &TrackError{
		Code:    codes.PermissionDenied,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewUnsupportedImageFormatError(format string, args ...interface{}) error {
	return &TrackError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewFailedToUploadAvatarError(format string, args ...interface{}) error {
	return &TrackError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrTrackNotFound                = NewNotFoundError("track not found")
	ErrStreamPermissionDenied       = NewPermissionDeniedError("user does not have permission to update this stream")
	ErrStreamNotFound               = NewNotFoundError("stream not found")
	ErrFailedToUpdateStreamDuration = NewInternalError("failed to update stream duration")
)
