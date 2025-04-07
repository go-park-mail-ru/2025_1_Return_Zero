package repository

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Repository struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
	expiration time.Duration
}

func NewS3Repository(s3 *s3.S3, bucketName string, expiration time.Duration) *s3Repository {
	uploader := s3manager.NewUploaderWithClient(s3)
	return &s3Repository{s3: s3, bucketName: bucketName,
		expiration: expiration, uploader: uploader}
}

func (r *s3Repository) GetPresignedURL(fileKey string) (string, error) {
	if fileKey == "" {
		return "", errors.New("empty S3 key")
	}
	path := fmt.Sprintf("avatars/%s", fileKey)
	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(path),
	}

	req, _ := r.s3.GetObjectRequest(input)
	url, err := req.Presign(r.expiration)
	if err != nil {
		return "", err
	}
	
	return url, nil
}

func (r *s3Repository) UploadUserAvatar(username string, fileContent io.Reader) (string, error) {
	fileKey := fmt.Sprintf("avatars/%s.png", username)

	_, err := r.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fileKey),
		Body:        fileContent,
		ContentType: aws.String("image/png"),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.png", username), nil
}
