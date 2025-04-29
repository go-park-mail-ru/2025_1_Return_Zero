package customErrors

import "errors"

var (
	ErrArtistNotFound   = errors.New("artist not found")
	ErrInvalidOffset    = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit     = errors.New("invalid limit: should be greater than 0")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExist        = errors.New("user with this username or email already exists")
	ErrCreateSalt       = errors.New("failed to create salt")
	ErrWrongPassword    = errors.New("wrong password")
	ErrUsernameExist    = errors.New("username already exists")
	ErrEmailExist       = errors.New("email already exists")
	ErrPasswordRequired = errors.New("password is required")
)
