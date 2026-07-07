// Package urlsigner generates time-limited presigned URLs for objects stored in
// S3-compatible buckets (AWS S3 or Cloudflare R2).
package urlsigner

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GeneratePresignedURL returns a URL that grants temporary read access to objectKey in
// bucketName, valid for lifetimeSecs seconds.
//
// Credentials and region are read from AMZ_ACCESS_KEY_ID, AMZ_SECRET_ACCESS_KEY, and
// AMZ_REGION, matching the rest of the AWS-facing code in this repo. To target Cloudflare
// R2 (or any other S3-compatible endpoint) instead of AWS S3, additionally set
// AWS_ENDPOINT_URL_S3 to that provider's endpoint (e.g. https://<account_id>.r2.cloudflarestorage.com)
// — no other code change is needed to move a bucket between providers.
func GeneratePresignedURL(ctx context.Context, bucketName, objectKey string, lifetimeSecs int64) (string, error) {
	client, err := newS3Client(ctx)
	if err != nil {
		return "", err
	}

	presignClient := s3.NewPresignClient(client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs) * time.Second
	})
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func newS3Client(ctx context.Context) (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		os.Getenv("AMZ_ACCESS_KEY_ID"),
		os.Getenv("AMZ_SECRET_ACCESS_KEY"),
		"",
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(os.Getenv("AMZ_REGION")),
	)
	if err != nil {
		return nil, err
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL_S3")

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			// R2 (and most non-AWS S3-compatible providers) don't support
			// AWS's virtual-hosted-style bucket addressing on a custom endpoint.
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	}), nil
}
