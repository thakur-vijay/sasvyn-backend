package storage

import (
	"context"
	"fmt"
	"os"

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

func DeleteObject(ctx context.Context, client *s3.Client, key string) error {
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
		Key:    aws.String(key),
	})

	return err
}
