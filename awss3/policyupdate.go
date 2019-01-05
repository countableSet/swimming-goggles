package awss3

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"gitlab.com/countableset/lambda-s3-cloudflare/util"
)

// UpdatePolicy updates the policy of the bucket with the policy
func UpdatePolicy(svc *s3.S3, bucket string, policy interface{}) {
	policyEncoded := encodedPolicy(policy)
	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policyEncoded),
	}

	_, err := svc.PutBucketPolicy(input)
	if err != nil {
		util.ExitErrorf("Unable to update bucket policy %v", err)
	}
}

func encodedPolicy(policy interface{}) string {
	encoded, err := json.Marshal(policy)
	if err != nil {
		util.ExitErrorf("Error occured during json marshal %v", err)
	}
	return string(encoded)
}
