package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"bytes"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers"
	userAvatarFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile"
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
	logger := helpers.LoggerFromContext(ctx)
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
	logger := helpers.LoggerFromContext(ctx)
	date := time.Now()
	dateString := date.Format("20060102150405")

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, fileContent)
	if err != nil {
		logger.Error("failed to read file content", zap.Error(err))
		return "", userAvatarFile.ErrFailedToParseImage
	}

	_, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	logger.Info("format", zap.String("format", format))
	if err != nil {
		logger.Error("unsupported or invalid image format", zap.Error(err))
		return "", userAvatarFile.ErrUnsupportedImageFormat
	}

	if format != "jpeg" && format != "png" && format != "gif" {
		logger.Error("unsupported image format", zap.String("format", format))
		return "", userAvatarFile.ErrUnsupportedImageFormat
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
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/" + format),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		logger.Error("upload failed", zap.Error(err))
		return "", userAvatarFile.ErrFailedToUploadAvatar
	}

	return fileKey, nil
}
