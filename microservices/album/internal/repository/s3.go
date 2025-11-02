package repository

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/internal/domain"
	albumErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/errors"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"go.uber.org/zap"
)

type artistS3Repository struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
	metrics    *metrics.Metrics
}

func NewS3Repository(s3 *s3.S3, bucketName string, metrics *metrics.Metrics) domain.S3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &artistS3Repository{
		s3:         s3,
		bucketName: bucketName,
		uploader:   uploader,
		metrics:    metrics,
	}
}

func (r *artistS3Repository) GetAvatarURL(ctx context.Context, fileKey string) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	if fileKey == "" {
		r.metrics.DatabaseErrors.WithLabelValues("GetAvatarURL").Inc()
		logger.Error("fileKey is empty")
		return "", albumErrors.NewEmptyS3KeyError("fileKey is empty")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAvatarURL").Observe(duration)
	return fmt.Sprintf(
		"https://%s.s3.cloud.ru/albums%s",
		r.bucketName,
		fileKey,
	), nil
}

func (r *artistS3Repository) UploadAlbumAvatar(ctx context.Context, artistTitle string, file []byte) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	date := time.Now()
	dateString := date.Format("20060102150405")

	_, format, err := image.Decode(bytes.NewReader(file))
	fmt.Println("ALBUM FORMAT: ", format)
	logger.Info("format", zap.String("format", format))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadAlbumAvatar").Inc()
		logger.Error("unsupported or invalid image format", zap.Error(err))
		return "", albumErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	if format != "jpeg" && format != "png" && format != "gif" {
		r.metrics.DatabaseErrors.WithLabelValues("UploadAlbumAvatar").Inc()
		logger.Error("unsupported image format", zap.String("format", format))
		return "", albumErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	// webpBuf := new(bytes.Buffer)
	// if err := webp.Encode(webpBuf, img, &webp.Options{Lossless: true}); err != nil {
	// 	logger.Error("failed to encode webp", zap.Error(err))
	// 	return "", userAvatarFile.ErrFailedToEncodeWebp
	// }

	fileKey := fmt.Sprintf("/%s-%s.%s", artistTitle, dateString, format)
	s3Key := fmt.Sprintf("albums%s", fileKey)

	_, err = r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/" + format),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadArtistAvatar").Inc()
		logger.Error("upload failed", zap.Error(err))
		return "", albumErrors.NewFailedToUploadAvatarError("failed to upload avatar")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UploadArtistAvatar").Observe(duration)

	return fmt.Sprintf("https://%s.s3.cloud.ru/albums%s", r.bucketName, fileKey), nil
}
