package errorStatus

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
)

var mapErrorStatus = map[error]int{
	customErrors.ErrTrackNotFound:                http.StatusNotFound,
	customErrors.ErrStreamNotFound:               http.StatusNotFound,
	customErrors.ErrFailedToUpdateStreamDuration: http.StatusInternalServerError,
	customErrors.ErrStreamPermissionDenied:       http.StatusForbidden,

	customErrors.ErrUserExist:        http.StatusConflict,
	customErrors.ErrUserNotFound:     http.StatusNotFound,
	customErrors.ErrCreateSalt:       http.StatusInternalServerError,
	customErrors.ErrWrongPassword:    http.StatusUnauthorized,
	customErrors.ErrPasswordRequired: http.StatusBadRequest,

	customErrors.ErrCreateSession: http.StatusInternalServerError,
	customErrors.ErrGetSession:    http.StatusUnauthorized,
	customErrors.ErrDeleteSession: http.StatusInternalServerError,

	customErrors.ErrInvalidOffset:            http.StatusBadRequest,
	customErrors.ErrInvalidLimit:             http.StatusBadRequest,
	customErrors.ErrPasswordRequired:         http.StatusBadRequest,
	customErrors.ErrUnsupportedImageFormat:   http.StatusBadRequest,
	customErrors.ErrFailedToParseImage:       http.StatusBadRequest,
	customErrors.ErrArtistNotFound:           http.StatusNotFound,
	customErrors.ErrAlbumNotFound:            http.StatusNotFound,
	customErrors.ErrUnauthorized:             http.StatusForbidden,
	customErrors.ErrPlaylistNotFound:         http.StatusNotFound,
	customErrors.ErrPlaylistPermissionDenied: http.StatusForbidden,
	customErrors.ErrPlaylistDuplicate:        http.StatusConflict,
	customErrors.ErrPlaylistTrackNotFound:    http.StatusNotFound,
	customErrors.ErrPlaylistTrackDuplicate:   http.StatusConflict,
	customErrors.ErrPlaylistImageNotUploaded: http.StatusBadRequest,
	customErrors.ErrPlaylistBadRequest:       http.StatusBadRequest,
	customErrors.ErrPlaylistUnauthorized:     http.StatusUnauthorized,

	customErrors.ErrLableExist: http.StatusBadRequest,
}

func ErrorStatus(err error) int {
	status, exists := mapErrorStatus[err]
	if !exists {
		return http.StatusInternalServerError
	}
	return status
}
