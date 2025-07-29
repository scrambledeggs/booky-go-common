package cloudfronthelpers

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func InvalidateCacheStatus(id string) (*cloudfront.GetInvalidationOutput, error) {
	sesh := session.Must(session.NewSession())
	cf := cloudfront.New(sesh)

	out, err := cf.GetInvalidation(&cloudfront.GetInvalidationInput{
		DistributionId: aws.String(os.Getenv("DISTRIBUTION_ID")),
		Id:             &id,
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
