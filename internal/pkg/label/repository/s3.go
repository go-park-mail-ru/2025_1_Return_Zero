// package repository

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"time"

// 	"image"

// 	_ "image/gif"
// 	_ "image/jpeg"
// 	_ "image/png"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/aws/aws-sdk-go/service/s3/s3manager"
// 	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/customErrors"
// 	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
// 	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/label/domain"
// 	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
// 	"go.uber.org/zap"
// )

// type labelS3Repository struct {
// 	s3         *s3.S3
// 	uploader   *s3manager.Uploader
// 	bucketName string
// 	metrics    *metrics.Metrics
// }

// func NewS3Repository(s3 *s3.S3, bucketName string, metrics *metrics.Metrics) domain.S3Repository {
// 	uploader := s3manager.NewUploaderWithClient(s3)
// 	return &labelS3Repository{
// 		s3:         s3,
// 		bucketName: bucketName,
// 		uploader:   uploader,
// 		metrics:    metrics,
// 	}
// }

// func (r *labelS3Repository) UploadAlbumAvatar(ctx context.Context, albumTitle string, file []byte) (string, error) {
// 	start := time.Now()
// 	logger := loggerPkg.LoggerFromContext(ctx)
// 	date := time.Now()
// 	dateString := date.Format("20060102150405")

// 	_, format, err := image.Decode(bytes.NewReader(file))
// 	logger.Info("format", zap.String("format", format))
// 	if err != nil {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadAlbumAvatar").Inc()
// 		logger.Error("unsupported or invalid image format", zap.Error(err))
// 		return "", customErrors.ErrUnsupportedImageFormatError
// 	}

// 	if format != "jpeg" && format != "png" && format != "gif" {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadAlbumAvatar").Inc()
// 		logger.Error("unsupported image format", zap.String("format", format))
// 		return "", customErrors.ErrUnsupportedImageFormatError
// 	}

// 	// webpBuf := new(bytes.Buffer)
// 	// if err := webp.Encode(webpBuf, img, &webp.Options{Lossless: true}); err != nil {
// 	// 	logger.Error("failed to encode webp", zap.Error(err))
// 	// 	return "", userAvatarFile.ErrFailedToEncodeWebp
// 	// }

// 	fileKey := fmt.Sprintf("/%s-%s.%s", albumTitle, dateString, format)
// 	s3Key := fmt.Sprintf("albums%s", fileKey)

// 	_, err = r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
// 		Bucket:      aws.String(r.bucketName),
// 		Key:         aws.String(s3Key),
// 		Body:        bytes.NewReader(file),
// 		ContentType: aws.String("image/" + format),
// 		ACL:         aws.String("public-read"),
// 	})

// 	if err != nil {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadAlbumAvatar").Inc()
// 		logger.Error("upload failed", zap.Error(err))
// 		return "", customErrors.ErrFailedToUploadImage
// 	}
// 	duration := time.Since(start).Seconds()
// 	r.metrics.DatabaseDuration.WithLabelValues("UploadAlbumAvatar").Observe(duration)

// 	return fmt.Sprintf("https://%s.s3.cloud.ru/albums%s", r.bucketName, fileKey), nil
// }

// func (r *labelS3Repository) UploadTrackAvatar(ctx context.Context, trackTitle string, file []byte) (string, error) {
// 	start := time.Now()
// 	logger := loggerPkg.LoggerFromContext(ctx)
// 	date := time.Now()
// 	dateString := date.Format("20060102150405")

// 	_, format, err := image.Decode(bytes.NewReader(file))
// 	logger.Info("format", zap.String("format", format))
// 	if err != nil {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
// 		logger.Error("unsupported or invalid image format", zap.Error(err))
// 		return "", customErrors.ErrUnsupportedImageFormatError
// 	}

// 	if format != "jpeg" && format != "png" && format != "gif" {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
// 		logger.Error("unsupported image format", zap.String("format", format))
// 		return "", customErrors.ErrUnsupportedImageFormatError
// 	}

// 	fileKey := fmt.Sprintf("/%s-%s.%s", trackTitle, dateString, format)
// 	s3Key := fmt.Sprintf("tracks%s", fileKey)

// 	_, err = r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
// 		Bucket:      aws.String(r.bucketName),
// 		Key:         aws.String(s3Key),
// 		Body:        bytes.NewReader(file),
// 		ContentType: aws.String("image/" + format),
// 		ACL:         aws.String("public-read"),
// 	})

// 	if err != nil {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadTrackAvatar").Inc()
// 		logger.Error("upload failed", zap.Error(err))
// 		return "", customErrors.ErrFailedToUploadImage
// 	}
// 	duration := time.Since(start).Seconds()
// 	r.metrics.DatabaseDuration.WithLabelValues("UploadTrackAvatar").Observe(duration)

// 	return fmt.Sprintf("https://%s.s3.cloud.ru/tracks%s", r.bucketName, fileKey), nil
// }

// func (r *labelS3Repository) UploadTrack(ctx context.Context, trackTitle string, file []byte) (string, error) {
// 	start := time.Now()
// 	logger := loggerPkg.LoggerFromContext(ctx)

// 	fileKey := fmt.Sprintf("%s.mp3", trackTitle)

// 	_, err := r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
// 		Bucket: aws.String(r.bucketName),
// 		Key:    aws.String(fileKey),
// 		Body:   bytes.NewReader(file),
// 	})

// 	if err != nil {
// 		r.metrics.DatabaseErrors.WithLabelValues("UploadTrack").Inc()
// 		logger.Error("upload failed", zap.Error(err))
// 		return "", customErrors.ErrFailedToUploadImage
// 	}
// 	duration := time.Since(start).Seconds()
// 	r.metrics.DatabaseDuration.WithLabelValues("UploadTrack").Observe(duration)

//		return fileKey, nil
//	}
package repository
