package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
	"go.uber.org/zap"
)

type s3Repository struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
}

func NewS3Repository(s3 *s3.S3, bucketName string) *s3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &s3Repository{
		s3:         s3,
		bucketName: bucketName,
		uploader:   uploader}
}

func (r *s3Repository) GetAvatarURL(ctx context.Context, fileKey string) (string, error) {
	logger := middleware.LoggerFromContext(ctx)
	if fileKey == "" {
		logger.Error("fileKey is empty")
		return "", errors.New("empty S3 key")
	}

	return fmt.Sprintf(
		"https://%s.fra1.digitaloceanspaces.com/avatars%s",
		r.bucketName,
		fileKey,
	), nil
}

func (r *s3Repository) UploadUserAvatar(ctx context.Context, username string, fileContent io.Reader) (string, error) {
	logger := middleware.LoggerFromContext(ctx)
	date := time.Now()
	dateString := date.Format("20060102150405")
	fileKey := fmt.Sprintf("/%s-%s.png", username, dateString)
	s3Key := fmt.Sprintf("avatars%s", fileKey)

	_, err := r.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        fileContent,
		ContentType: aws.String("image/png"),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		logger.Error("upload failed", zap.Error(err))
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return fileKey, nil
}

func (r *s3Repository) DeleteUserAvatar(ctx context.Context, fileKey string) error {
	logger := middleware.LoggerFromContext(ctx)
	s3Key := fmt.Sprintf("avatars%s", fileKey)

	_, err := r.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		logger.Error("delete failed", zap.Error(err))
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}
