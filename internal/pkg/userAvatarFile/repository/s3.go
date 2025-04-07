package repository

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (r *s3Repository) GetAvatarURL(fileKey string) (string, error) {
	if fileKey == "" {
		return "", errors.New("empty S3 key")
	}

	return fmt.Sprintf(
		"https://%s.fra1.digitaloceanspaces.com/avatars%s",
		r.bucketName,
		fileKey,
	), nil
}

func (r *s3Repository) UploadUserAvatar(username string, fileContent io.Reader) (string, error) {
	fileKey := fmt.Sprintf("/%s.png", username)
	s3Key := fmt.Sprintf("avatars%s", fileKey)

	_, err := r.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        fileContent,
		ContentType: aws.String("image/png"),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return fileKey, nil
}
