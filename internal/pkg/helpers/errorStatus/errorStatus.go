package errorStatus

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	userAvatarFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
)

var mapErrorStatus = map[error]int{
	customErrors.ErrTrackNotFound:                http.StatusNotFound,
	customErrors.ErrStreamNotFound:               http.StatusNotFound,
	customErrors.ErrFailedToUpdateStreamDuration: http.StatusInternalServerError,
	customErrors.ErrStreamPermissionDenied:       http.StatusForbidden,
	user.ErrUsernameExist:                        http.StatusNotFound,
	user.ErrEmailExist:                           http.StatusNotFound,
	user.ErrUserNotFound:                         http.StatusNotFound,
	user.ErrCreateSalt:                           http.StatusNotFound,
	user.ErrWrongPassword:                        http.StatusNotFound,
	auth.ErrSessionNotFound:                      http.StatusNotFound,
	customErrors.ErrInvalidOffset:                http.StatusBadRequest,
	customErrors.ErrInvalidLimit:                 http.StatusBadRequest,
	user.ErrPasswordRequired:                     http.StatusBadRequest,
	userAvatarFile.ErrFailedToUploadAvatar:       http.StatusBadRequest,
	userAvatarFile.ErrUnsupportedImageFormat:     http.StatusBadRequest,
	userAvatarFile.ErrFailedToEncodeWebp:         http.StatusBadRequest,
	userAvatarFile.ErrFailedToParseImage:         http.StatusBadRequest,
	customErrors.ErrArtistNotFound:               http.StatusNotFound,
	customErrors.ErrAlbumNotFound:                http.StatusNotFound,
	customErrors.ErrStreamHistoryUnauthorized:    http.StatusForbidden,
	customErrors.ErrStreamUpdateUnauthorized:     http.StatusForbidden,
	customErrors.ErrStreamCreateUnauthorized:     http.StatusForbidden,
	customErrors.ErrLikeArtistUnauthorized:       http.StatusForbidden,
	customErrors.ErrLikeAlbumUnauthorized:        http.StatusForbidden,
	customErrors.ErrLikeTrackUnauthorized:        http.StatusForbidden,
	customErrors.ErrPlaylistNotFound:             http.StatusNotFound,
	customErrors.ErrPlaylistPermissionDenied:     http.StatusForbidden,
	customErrors.ErrPlaylistBadRequest:           http.StatusBadRequest,
	customErrors.ErrUnsupportedImageFormat:       http.StatusBadRequest,
	customErrors.ErrImageTooBig:                  http.StatusBadRequest,
	customErrors.ErrFailedToParseImage:           http.StatusBadRequest,
	customErrors.ErrFailedToUploadImage:          http.StatusBadRequest,
	customErrors.ErrPlaylistImageNotUploaded:     http.StatusBadRequest,
	customErrors.ErrFailedToCreatePlaylist:       http.StatusBadRequest,
	customErrors.ErrPlaylistUnauthorized:         http.StatusForbidden,
	customErrors.ErrPlaylistDuplicate:            http.StatusBadRequest,
	customErrors.ErrPlaylistTrackNotFound:        http.StatusNotFound,
	customErrors.ErrPlaylistTrackDuplicate:       http.StatusBadRequest,
}

func ErrorStatus(err error) int {
	status, exists := mapErrorStatus[err]
	if !exists {
		return http.StatusInternalServerError
	}
	return status
}
