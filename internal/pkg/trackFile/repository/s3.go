package repository

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Repository struct {
	s3         *s3.S3
	bucketName string
	expiration time.Duration
}

func NewS3Repository(s3 *s3.S3, bucketName string, expiration time.Duration) *s3Repository {
	return &s3Repository{s3: s3, bucketName: bucketName, expiration: expiration}
}

func (r *s3Repository) GetPresignedURL(trackKey string) (string, error) {
	// Вот тут сомнительно, возможно стоит кидать ошибку, если ключа нет
	// В бд ключ не может быть пустым, но на всякий случай
	if trackKey == "" {
		return "", nil
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(trackKey),
	}

	req, _ := r.s3.GetObjectRequest(input)
	url, err := req.Presign(r.expiration)
	if err != nil {
		return "", err
	}

	return url, nil
}
