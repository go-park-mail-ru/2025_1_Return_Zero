package helpers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

var mapErrorStatus = map[error]int{
	track.ErrTrackNotFound:                http.StatusNotFound,
	track.ErrStreamNotFound:               http.StatusNotFound,
	track.ErrFailedToUpdateStreamDuration: http.StatusInternalServerError,
	track.ErrStreamPermissionDenied:       http.StatusForbidden,
	artist.ErrArtistNotFound:              http.StatusNotFound,
	album.ErrAlbumNotFound:                http.StatusNotFound,
	user.ErrUsernameExist:                 http.StatusNotFound,
	user.ErrEmailExist:                    http.StatusNotFound,
	user.ErrUserNotFound:                  http.StatusNotFound,
	user.ErrCreateSalt:                    http.StatusNotFound,
	user.ErrWrongPassword:                 http.StatusNotFound,
	auth.ErrSessionNotFound:               http.StatusNotFound,
	ErrInvalidOffset:                      http.StatusBadRequest,
	ErrInvalidLimit:                       http.StatusBadRequest,
	user.ErrPasswordRequired:                   http.StatusBadRequest,
}

func ErrorStatus(err error) int {
	status, exists := mapErrorStatus[err]
	if !exists {
		return http.StatusInternalServerError
	}
	return status
}
