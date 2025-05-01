package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserError struct {
	Code    codes.Code
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

func (e *UserError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewNotFoundError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewUserExistError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.AlreadyExists,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewCreateSaltError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewWrongPasswordError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.Unauthenticated,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewPasswordRequierdError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewEmptyS3KeyError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewUnsupportedImageFormatError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.InvalidArgument,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewFailedToUploadAvatarError(format string, args ...interface{}) error {
	return &UserError{
		Code:    codes.Internal,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrUserNotFound           = NewNotFoundError("user not found")
	ErrUserExist              = NewUserExistError("user already exist")
	ErrCreateSalt             = NewCreateSaltError("failed to create salt")
	ErrWrongPassword          = NewWrongPasswordError("wrong password")
	ErrPasswordRequierd       = NewPasswordRequierdError("password required")
	ErrEmptyS3Key             = NewEmptyS3KeyError("s3 key is empty")
	ErrUnsupportedImageFormat = NewUnsupportedImageFormatError("unsupported image format")
	ErrFailedToUploadAvatar   = NewFailedToUploadAvatarError("failed to upload avatar")
)
