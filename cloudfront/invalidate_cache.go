package cloudfront

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func InvalidateCache(paths []string) (*cloudfront.CreateInvalidationOutput, error) {
	sesh := session.Must(session.NewSession())
	cf := cloudfront.New(sesh)

	awsPaths := pathsToAWSPaths(paths)

	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(os.Getenv("DISTRIBUTION_ID")),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format(time.RFC3339)),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(int64(len(awsPaths))),
				Items:    awsPaths,
			},
		},
	}

	result, err := cf.CreateInvalidation(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func pathsToAWSPaths(paths []string) []*string {
	var awsPaths []*string
	for _, path := range paths {
		awsPaths = append(awsPaths, aws.String(path))
	}

	return awsPaths
}
