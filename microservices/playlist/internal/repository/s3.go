package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	playlistErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/model/errors"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/playlist/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"go.uber.org/zap"

	_ "image/jpeg"
	_ "image/png"
)

const (
	MaxImageSize = 5 * 1024 * 1024
)

type playlistS3Repository struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
}

func NewPlaylistS3Repository(s3 *s3.S3, bucketName string) domain.S3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &playlistS3Repository{
		s3:         s3,
		bucketName: bucketName,
		uploader:   uploader,
	}
}

func (r *playlistS3Repository) UploadThumbnail(ctx context.Context, file io.Reader, key string) (string, error) {
	logger := loggerPkg.LoggerFromContext(ctx)
	logger.Info("Uploading thumbnail", zap.String("key", key))
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		logger.Error("failed to parse image", "error", err)
		return "", playlistErrors.ErrFailedToParseImage
	}

	if buf.Len() > MaxImageSize {
		logger.Error("image size exceeds 5MB limit", "error", playlistErrors.ErrImageTooBig)
		return "", playlistErrors.ErrImageTooBig
	}

	_, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		logger.Error("failed to parse image", "error", err)
		return "", playlistErrors.ErrFailedToParseImage
	}

	if format != "png" && format != "jpeg" {
		logger.Error("unsupported image format", "error", playlistErrors.ErrUnsupportedImageFormat)
		return "", playlistErrors.ErrUnsupportedImageFormat
	}

	timestampedKey := fmt.Sprintf("%s-%s.%s", key, time.Now().Format("20060102150405"), format)

	s3Key := fmt.Sprintf("playlists/%s", timestampedKey)

	output, err := r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/" + format),
		ACL:         aws.String("public-read"),
		Metadata: map[string]*string{
			"Cache-Control": aws.String("max-age=86400, public"),
		},
	})

	if err != nil {
		logger.Error("failed to upload image", "error", err)
		return "", playlistErrors.NewInternalError("failed to upload image: %v", err)
	}

	outloadedUrl := output.Location

	return outloadedUrl, nil
}
