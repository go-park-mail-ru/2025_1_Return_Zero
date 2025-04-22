package errorStatus

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
	userAvatarFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
)

var mapErrorStatus = map[error]int{
	track.ErrTrackNotFound:                   http.StatusNotFound,
	track.ErrStreamNotFound:                  http.StatusNotFound,
	track.ErrFailedToUpdateStreamDuration:    http.StatusInternalServerError,
	track.ErrStreamPermissionDenied:          http.StatusForbidden,
	user.ErrUsernameExist:                    http.StatusNotFound,
	user.ErrEmailExist:                       http.StatusNotFound,
	user.ErrUserNotFound:                     http.StatusNotFound,
	user.ErrCreateSalt:                       http.StatusNotFound,
	user.ErrWrongPassword:                    http.StatusNotFound,
	auth.ErrSessionNotFound:                  http.StatusNotFound,
	customErrors.ErrInvalidOffset:            http.StatusBadRequest,
	customErrors.ErrInvalidLimit:             http.StatusBadRequest,
	user.ErrPasswordRequired:                 http.StatusBadRequest,
	userAvatarFile.ErrFailedToUploadAvatar:   http.StatusBadRequest,
	userAvatarFile.ErrUnsupportedImageFormat: http.StatusBadRequest,
	userAvatarFile.ErrFailedToEncodeWebp:     http.StatusBadRequest,
	userAvatarFile.ErrFailedToParseImage:     http.StatusBadRequest,
	customErrors.ErrArtistNotFound:           http.StatusNotFound,
	customErrors.ErrAlbumNotFound:            http.StatusNotFound,
}

func ErrorStatus(err error) int {
	status, exists := mapErrorStatus[err]
	if !exists {
		return http.StatusInternalServerError
	}
	return status
}
