package helpers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/album"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/track"
)

func ErrorStatus(err error) int {
	switch err {
	case track.ErrTrackNotFound:
		return http.StatusNotFound
	case track.ErrStreamNotFound:
		return http.StatusNotFound
	case track.ErrFailedToUpdateStreamDuration:
		return http.StatusInternalServerError
	case track.ErrStreamPermissionDenied:
		return http.StatusForbidden
	case artist.ErrArtistNotFound:
		return http.StatusNotFound
	case album.ErrAlbumNotFound:
		return http.StatusNotFound
	case artist.ErrArtistNotFound:
		return http.StatusNotFound
	case ErrInvalidOffset:
		return http.StatusBadRequest
	case ErrInvalidLimit:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
