package main

import (
	"gitlab.com/countableset/lambda-s3-cloudflare/cloudflare"
	"gitlab.com/countableset/lambda-s3-cloudflare/util"

	"gitlab.com/countableset/lambda-s3-cloudflare/awss3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// TODO
	// configure creds profile?
	// configure bucket name?
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "logtest")},
	)
	svc := s3.New(sess)
	bucket := "logs.countableset.com"
	policy := awss3.GetPolicy(svc, bucket)
	s3IPAddresses := awss3.GetSortedIPAddresses(policy.(map[string]interface{}))
	cloudflareIPAddresses := cloudflare.GetAllSortedIPAddresses()
	result := util.TestEqualSlices(s3IPAddresses, cloudflareIPAddresses)
	if result != nil {
		policy = awss3.MergeIPSliceIntoPolicy(policy.(map[string]interface{}), result)
		awss3.UpdatePolicy(svc, bucket, policy)
	}
}
