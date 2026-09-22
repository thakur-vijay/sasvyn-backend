package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewR2Client() *s3.Client {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	endpoint := fmt.Sprintf(
		"https://%s.r2.cloudflarestorage.com",
		accountID,
	)

	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		),
		Region: "auto",
	}

	return s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
	})
}

func CreateImageURL(
	ctx context.Context,
	presignClient *s3.PresignClient,
	imgKey string,
) (string, error) {
	if imgKey == "" {
		return "", nil
	}

	bucket := os.Getenv("R2_BUCKET_NAME")

	presigned, err := presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(imgKey),
		},
		func(options *s3.PresignOptions) {
			options.Expires = 15 * time.Minute
		},
	)
	if err != nil {
		return "", err
	}

	return presigned.URL, nil
}
