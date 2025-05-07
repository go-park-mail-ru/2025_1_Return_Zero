package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type trackS3Repository struct {
	s3         *s3.S3
	bucketName string
	expiration time.Duration
	metrics    *metrics.Metrics
}

func NewTrackS3Repository(s3 *s3.S3, bucketName string, expiration time.Duration, metrics *metrics.Metrics) domain.S3Repository {
	return &trackS3Repository{s3: s3, bucketName: bucketName, expiration: expiration}
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
		Bucket: aws.String(r.bucketName),
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
