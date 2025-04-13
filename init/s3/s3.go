package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
)

func InitS3(cfg config.S3Config) (*s3.S3, error) {
	s3Config := &aws.Config{
		Region:           aws.String(cfg.S3_REGION),
		Credentials:      credentials.NewStaticCredentials(cfg.S3_ACCESS_KEY, cfg.S3_SECRET_KEY, ""),
		Endpoint:         aws.String(cfg.S3_ENDPOINT),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to S3")

	return s3.New(newSession), nil
}
