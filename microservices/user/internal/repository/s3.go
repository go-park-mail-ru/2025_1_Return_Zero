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
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/domain"
	userErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/errors"
	"go.uber.org/zap"
)

type userS3Repository struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
	metrics    *metrics.Metrics
}

func NewS3Repository(s3 *s3.S3, bucketName string, metrics *metrics.Metrics) domain.S3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &userS3Repository{
		s3:         s3,
		bucketName: bucketName,
		uploader:   uploader,
		metrics:    metrics,
	}
}

func (r *userS3Repository) GetAvatarURL(ctx context.Context, fileKey string) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	if fileKey == "" {
		r.metrics.DatabaseErrors.WithLabelValues("GetAvatarURL").Inc()
		logger.Error("fileKey is empty")
		return "", userErrors.NewEmptyS3KeyError("fileKey is empty")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetAvatarURL").Observe(duration)
	return fmt.Sprintf(
		"https://%s.fra1.digitaloceanspaces.com/avatars%s",
		r.bucketName,
		fileKey,
	), nil
}

func (r *userS3Repository) UploadUserAvatar(ctx context.Context, username string, file []byte) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	date := time.Now()
	dateString := date.Format("20060102150405")

	_, format, err := image.Decode(bytes.NewReader(file))
	logger.Info("format", zap.String("format", format))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadUserAvatar").Inc()
		logger.Error("unsupported or invalid image format", zap.Error(err))
		return "", userErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	if format != "jpeg" && format != "png" && format != "gif" {
		r.metrics.DatabaseErrors.WithLabelValues("UploadUserAvatar").Inc()
		logger.Error("unsupported image format", zap.String("format", format))
		return "", userErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	// webpBuf := new(bytes.Buffer)
	// if err := webp.Encode(webpBuf, img, &webp.Options{Lossless: true}); err != nil {
	// 	logger.Error("failed to encode webp", zap.Error(err))
	// 	return "", userAvatarFile.ErrFailedToEncodeWebp
	// }

	fileKey := fmt.Sprintf("/%s-%s.%s", username, dateString, format)
	s3Key := fmt.Sprintf("avatars%s", fileKey)

	_, err = r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/" + format),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadUserAvatar").Inc()
		logger.Error("upload failed", zap.Error(err))
		return "", userErrors.NewFailedToUploadAvatarError("failed to upload avatar")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UploadUserAvatar").Observe(duration)

	return fileKey, nil
}
