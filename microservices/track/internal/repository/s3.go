package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	trackErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/errors"
)

type trackS3Repository struct {
	s3              *s3.S3
	trackBucketName string
	imageBucketName string
	uploader        *s3manager.Uploader
	expiration      time.Duration
	metrics         *metrics.Metrics
}

func NewTrackS3Repository(s3 *s3.S3, trackBucketName string, imageBucketName string, expiration time.Duration, metrics *metrics.Metrics) domain.S3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &trackS3Repository{
		s3:              s3,
		trackBucketName: trackBucketName,
		imageBucketName: imageBucketName,
		expiration:      expiration,
		uploader:        uploader,
		metrics:         metrics,
	}
}

func (r *trackS3Repository) GetPresignedURL(trackKey string) (string, error) {
	start := time.Now()
	// Вот тут сомнительно, возможно стоит кидать ошибку, если ключа нет
	// В бд ключ не может быть пустым, но на всякий случай
	if trackKey == "" {
		r.metrics.DatabaseErrors.WithLabelValues("GetPresignedURL").Inc()
		return "", nil
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.trackBucketName),
		Key:    aws.String(trackKey),
	}

	req, _ := r.s3.GetObjectRequest(input)
	url, err := req.Presign(r.expiration)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("GetPresignedURL").Inc()
		return "", err
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("GetPresignedURL").Observe(duration)
	return url, nil
}

func (r *trackS3Repository) UploadTrack(ctx context.Context, fileKey string, file []byte) error {
	start := time.Now()

	if fileKey == "" {
		r.metrics.DatabaseErrors.WithLabelValues("UploadTrack").Inc()
		return fmt.Errorf("empty file key provided")
	}

	input := &s3manager.UploadInput{
		Bucket: aws.String(r.trackBucketName),
		Key:    aws.String(fmt.Sprintf("%s.mp3", fileKey)),
		Body:   bytes.NewReader(file),
	}

	_, err := r.uploader.UploadWithContext(ctx, input)
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadTrack").Inc()
		return fmt.Errorf("failed to upload file: %w", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UploadTrack").Observe(duration)

	return nil
}

func (r *trackS3Repository) UploadTrackAvatar(ctx context.Context, trackTitle string, file []byte) (string, error) {
	start := time.Now()
	logger := loggerPkg.LoggerFromContext(ctx)
	date := time.Now()
	dateString := date.Format("20060102150405")

	_, format, err := image.Decode(bytes.NewReader(file))
	fmt.Println("FORMAT TRACK", format)
	logger.Info("format", zap.String("format", format))
	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
		logger.Error("unsupported or invalid image format", zap.Error(err))
		return "", trackErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	if format != "jpeg" && format != "png" && format != "gif" {
		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
		logger.Error("unsupported image format", zap.String("format", format))
		return "", trackErrors.NewUnsupportedImageFormatError("unsupported or invalid image format")
	}

	// webpBuf := new(bytes.Buffer)
	// if err := webp.Encode(webpBuf, img, &webp.Options{Lossless: true}); err != nil {
	// 	logger.Error("failed to encode webp", zap.Error(err))
	// 	return "", userAvatarFile.ErrFailedToEncodeWebp
	// }

	fileKey := fmt.Sprintf("/%s-%s.%s", trackTitle, dateString, format)
	s3Key := fmt.Sprintf("tracks%s", fileKey)

	_, err = r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(r.imageBucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/" + format),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
		logger.Error("upload failed", zap.Error(err))
		return "", trackErrors.NewFailedToUploadAvatarError("failed to upload avatar")
	}
	duration := time.Since(start).Seconds()
	r.metrics.DatabaseDuration.WithLabelValues("UploadTrackAvatar").Observe(duration)

	return fmt.Sprintf("https://%s.s3.cloud.ru/tracks%s", r.imageBucketName, fileKey), nil
}
