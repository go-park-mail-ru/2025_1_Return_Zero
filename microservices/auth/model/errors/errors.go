package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthError struct {
	Code    codes.Code
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func (e *AuthError) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func NewCreateSessionError(format string, args ...interface{}) *AuthError {
	return &AuthError{
		Code:    codes.Unavailable,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewDeleteSessionError(format string, args ...interface{}) *AuthError {
	return &AuthError{
		Code:    codes.Unavailable,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewGetSessionError(format string, args ...interface{}) *AuthError {
	return &AuthError{
		Code:    codes.Unavailable,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrCreateSession = NewCreateSessionError("failed to create session")
	ErrDeleteSession = NewDeleteSessionError("failed to delete session")
	ErrGetSession    = NewGetSessionError("failed to get session")
)
