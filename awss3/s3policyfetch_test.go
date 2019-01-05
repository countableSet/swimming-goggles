package awss3

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestParsePolicyResponse(t *testing.T) {
	var tests = []struct {
		input                 string
		expectedVersion       string
		expectedStatementSize int
		expectedSid           string
	}{
		{"{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"PublicReadGetObject\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"s3:GetObject\",\"Resource\":\"arn:aws:s3:::logs.countableset.com/*\",\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":[\"2400:cb00::/32\",\"2405:8100::/32\",\"2405:b500::/32\",\"2606:4700::/32\",\"2803:f800::/32\",\"2c0f:f248::/32\",\"2a06:98c0::/29\",\"103.21.244.0/22\",\"103.22.200.0/22\",\"103.31.4.0/22\",\"104.16.0.0/12\",\"108.162.192.0/18\",\"131.0.72.0/22\",\"141.101.64.0/18\",\"162.158.0.0/15\"]}}}]}",
			"2012-10-17", 1, "PublicReadGetObject"},
	}
	for _, test := range tests {
		resp := s3.GetBucketPolicyOutput{Policy: &test.input}
		result := parsePolicyResponse(&resp)
		if result.Version != test.expectedVersion {
			t.Errorf("Expected %s but got %s", test.expectedVersion, result)
		}
		if len(result.Statement) != test.expectedStatementSize {
			t.Errorf("Expected %d but got %s", test.expectedStatementSize, result)
		}
		if result.Statement[0].Sid != test.expectedSid {
			t.Errorf("Expected %s but got %s", test.expectedSid, result)
		}
	}
}
