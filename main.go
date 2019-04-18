package main

import (
	"log"

	"github.com/countableset/swimming-goggles/awss3"
	"github.com/countableset/swimming-goggles/cloudflare"
	"github.com/countableset/swimming-goggles/util"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Event struct to parse event json from lambda
type Event struct {
	Buckets []string `json:"buckets"`
}

var svc = s3.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

func lambdaHandler(event Event) (string, error) {
	cloudflareIPAddresses := cloudflare.GetAllSortedIPAddresses()
	for _, bucket := range event.Buckets {
		log.Printf("Running check on bucket %s\n", bucket)
		policy := awss3.GetPolicy(svc, bucket)
		s3IPAddresses := awss3.GetSortedIPAddresses(policy.(map[string]interface{}))
		result := util.TestEqualSlices(s3IPAddresses, cloudflareIPAddresses)
		if result != nil {
			log.Printf("Updating policy for bucket %s...\n", bucket)
			policy = awss3.MergeIPSliceIntoPolicy(policy.(map[string]interface{}), result)
			awss3.UpdatePolicy(svc, bucket, policy)
			log.Printf("Policy updated for bucket %s to %v\n", bucket, policy)
		} else {
			log.Printf("Policy already up-to-date, no changes needed for bucket %s\n", bucket)
		}
	}
	return "Success", nil
}

func main() {
	lambda.Start(lambdaHandler)
}
