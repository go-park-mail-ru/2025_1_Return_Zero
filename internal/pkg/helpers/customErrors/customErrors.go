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
