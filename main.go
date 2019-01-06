package main

import (
	"log"

	"gitlab.com/countableset/lambda-s3-cloudflare/cloudflare"
	"gitlab.com/countableset/lambda-s3-cloudflare/util"

	"gitlab.com/countableset/lambda-s3-cloudflare/awss3"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket string = "logs.countableset.com"

var svc = s3.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

func lambdaHandler() (string, error) {
	policy := awss3.GetPolicy(svc, bucket)
	s3IPAddresses := awss3.GetSortedIPAddresses(policy.(map[string]interface{}))
	cloudflareIPAddresses := cloudflare.GetAllSortedIPAddresses()
	result := util.TestEqualSlices(s3IPAddresses, cloudflareIPAddresses)
	if result != nil {
		log.Println("Updating policy...")
		policy = awss3.MergeIPSliceIntoPolicy(policy.(map[string]interface{}), result)
		awss3.UpdatePolicy(svc, bucket, policy)
		log.Printf("Policy updated to %v\n", policy)
	} else {
		log.Println("Policy already up-to-date, no changes needed!")
	}
	return "Success", nil
}

func main() {
	lambda.Start(lambdaHandler)
}
