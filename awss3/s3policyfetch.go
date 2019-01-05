package awss3

import (
	"encoding/json"

	"gitlab.com/countableset/lambda-s3-cloudflare/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Policy stucture for s3 policy response
type Policy struct {
	Version   string `json:"Version"`
	Statement []struct {
		Sid       string `json:"Sid"`
		Condition struct {
			IPAddress struct {
				SourceIP []string `json:"aws:SourceIp"`
			} `json:"IpAddress"`
		} `json:"Condition"`
	} `json:"Statement"`
}

// GetPolicy fetches the bucket policy from aws and parses it
func GetPolicy(svc *s3.S3, bucket string) Policy {
	policyOutput := getPolicyReponseFromAws(svc, bucket)
	return parsePolicyResponse(policyOutput)
}

func getPolicyReponseFromAws(svc *s3.S3, bucket string) *s3.GetBucketPolicyOutput {
	resp, err := svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		util.ExitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}
	return resp
}

func parsePolicyResponse(output *s3.GetBucketPolicyOutput) Policy {
	p := Policy{}
	json.Unmarshal([]byte(*output.Policy), &p)
	return p
}
