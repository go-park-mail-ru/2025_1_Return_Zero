package customErrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrArtistNotFound               = errors.New("artist not found")
	ErrInvalidOffset                = errors.New("invalid offset: should be greater than 0")
	ErrInvalidLimit                 = errors.New("invalid limit: should be greater than 0")
	ErrAlbumNotFound                = errors.New("album not found")
	ErrStreamNotFound               = errors.New("stream not found")
	ErrFailedToUpdateStreamDuration = errors.New("failed to update stream duration")
	ErrTrackNotFound                = errors.New("track not found")
	ErrStreamPermissionDenied       = errors.New("user does not have permission to update this stream")
	ErrStream                       = errors.New("stream not found")
	ErrStreamHistoryUnauthorized    = errors.New("unauthorized users can't save to stream history")
	ErrStreamUpdateUnauthorized     = errors.New("user does not have permission to update this stream")
	ErrStreamCreateUnauthorized     = errors.New("unauthorized users can't create stream")
	ErrLikeArtistUnauthorized       = errors.New("unauthorized users can't like artist")
	ErrLikeAlbumUnauthorized        = errors.New("unauthorized users can't like album")
	ErrLikeTrackUnauthorized        = errors.New("unauthorized users can't like track")
	ErrPlaylistNotFound             = errors.New("playlist not found")
	ErrPlaylistPermissionDenied     = errors.New("user does not have permission to update this playlist")
	ErrPlaylistBadRequest           = errors.New("invalid playlist request")
	ErrUnsupportedImageFormat       = errors.New("unsupported image format: only JPEG and PNG are allowed")
	ErrImageTooBig                  = errors.New("image size exceeds 5MB limit")
	ErrFailedToParseImage           = errors.New("failed to parse image")
	ErrFailedToUploadImage          = errors.New("failed to upload image")
	ErrFailedToCreatePlaylist       = errors.New("failed to create playlist")
	ErrPlaylistUnauthorized         = errors.New("unauthorized users can't create playlist")
	ErrPlaylistImageNotUploaded     = errors.New("playlist image not uploaded")
	ErrPlaylistDuplicate            = errors.New("playlist with this title by you already exists")
	ErrPlaylistTrackNotFound        = errors.New("track not found in playlist")
	ErrPlaylistTrackDuplicate       = errors.New("track already in playlist")
)

func HandleAlbumGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrAlbumNotFound
	case codes.Internal:
		return errors.New("internal server error: " + st.Message())
	default:
		return err
	}
}

func HandleArtistGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrArtistNotFound
	case codes.Internal:
		return errors.New("internal server error: " + st.Message())
	default:
		return err
	}
}

func HandleTrackGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		switch st.Message() {
		case "track not found":
			return ErrTrackNotFound
		case "stream not found":
			return ErrStreamNotFound
		default:
			return err
		}
	case codes.PermissionDenied:
		return ErrStreamPermissionDenied
	case codes.Internal:
		switch st.Message() {
		case "failed to update stream duration":
			return ErrFailedToUpdateStreamDuration
		default:
			return errors.New("internal server error: " + st.Message())
		}
	default:
		return err
	}

}

func HandlePlaylistGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		switch st.Message() {
		case "playlist not found":
			return ErrPlaylistNotFound
		case "track not found in playlist":
			return ErrPlaylistTrackNotFound
		default:
			return err
		}
	case codes.PermissionDenied:
		return ErrPlaylistPermissionDenied
	case codes.InvalidArgument:
		switch st.Message() {
		case "invalid playlist request":
			return ErrPlaylistBadRequest
		case "unsupported image format: only JPEG and PNG are allowed":
			return ErrUnsupportedImageFormat
		case "image size exceeds 5MB limit":
			return ErrImageTooBig
		case "failed to parse image":
			return ErrFailedToParseImage
		case "failed to upload image":
			return ErrFailedToUploadImage
		default:
			return err
		}
	case codes.AlreadyExists:
		switch st.Message() {
		case "playlist with this title by you already exists":
			return ErrPlaylistDuplicate
		case "track already in playlist":
			return ErrPlaylistTrackDuplicate
		default:
			return err
		}
	case codes.Internal:
		return errors.New("internal server error: " + st.Message())
	default:
		return err
	}
}
