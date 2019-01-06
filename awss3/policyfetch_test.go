package awss3

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestParsePolicyResponse(t *testing.T) {
	input := `{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Sid": "PublicReadGetObject",
			"Effect": "Allow",
			"Principal": "*",
			"Action": "s3:GetObject",
			"Resource": "arn:aws:s3:::logs.countableset.com/*",
			"Condition": {
			  "IpAddress": {
				"aws:SourceIp": [
				  "2400:cb00::/32",
				  "162.158.0.0/15"
				]
			  }
			}
		  }
		]
	  }`
	expectedVersion := "2012-10-17"
	expectedStatementSize := 1
	expectedSid := "PublicReadGetObject"

	resp := s3.GetBucketPolicyOutput{Policy: &input}
	result := parsePolicyResponse(&resp)
	resultMap := result.(map[string]interface{})
	if resultMap["Version"].(string) != expectedVersion {
		t.Errorf("Expected %s but got %s", expectedVersion, result)
	}
	if len(resultMap["Statement"].([]interface{})) != expectedStatementSize {
		t.Errorf("Expected %d but got %s", expectedStatementSize, result)
	}
	if resultMap["Statement"].([]interface{})[0].(map[string]interface{})["Sid"] != expectedSid {
		t.Errorf("Expected %s but got %s", expectedSid, result)
	}
}
