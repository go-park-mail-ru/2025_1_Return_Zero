package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type trackS3Repository struct {
	s3         *s3.S3
	bucketName string
	expiration time.Duration
}

func NewTrackS3Repository(s3 *s3.S3, bucketName string, expiration time.Duration) domain.S3Repository {
	return &trackS3Repository{s3: s3, bucketName: bucketName, expiration: expiration}
}

func (r *trackS3Repository) GetPresignedURL(trackKey string) (string, error) {
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
