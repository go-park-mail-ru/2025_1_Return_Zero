package errorStatus

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	userAvatarFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
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
	customErrors.ErrGetSession:    http.StatusInternalServerError,
	customErrors.ErrDeleteSession: http.StatusInternalServerError,

	customErrors.ErrInvalidOffset:             http.StatusBadRequest,
	customErrors.ErrInvalidLimit:              http.StatusBadRequest,
	customErrors.ErrPasswordRequired:          http.StatusBadRequest,
	userAvatarFile.ErrFailedToUploadAvatar:    http.StatusBadRequest,
	userAvatarFile.ErrUnsupportedImageFormat:  http.StatusBadRequest,
	userAvatarFile.ErrFailedToEncodeWebp:      http.StatusBadRequest,
	userAvatarFile.ErrFailedToParseImage:      http.StatusBadRequest,
	customErrors.ErrArtistNotFound:            http.StatusNotFound,
	customErrors.ErrAlbumNotFound:             http.StatusNotFound,
	customErrors.ErrStreamHistoryUnauthorized: http.StatusForbidden,
	customErrors.ErrStreamUpdateUnauthorized:  http.StatusForbidden,
	customErrors.ErrStreamCreateUnauthorized:  http.StatusForbidden,
	customErrors.ErrLikeArtistUnauthorized:    http.StatusForbidden,
	customErrors.ErrLikeAlbumUnauthorized:     http.StatusForbidden,
	customErrors.ErrLikeTrackUnauthorized:     http.StatusForbidden,
}

func ErrorStatus(err error) int {
	status, exists := mapErrorStatus[err]
	if !exists {
		return http.StatusInternalServerError
	}
	return status
}
