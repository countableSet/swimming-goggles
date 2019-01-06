package awss3

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestGetSortedIPAddresses(t *testing.T) {
	input := `{
		"Statement": [
		  {
			"Action": "s3:GetObject",
			"Condition": {
			  "IpAddress": {
				"aws:SourceIp": ["2400:cb00::/32", "2405:8100::/32", "162.158.0.0/15"]
			  }
			}
		  }
		]
	  }
	  `
	expected := []string{"162.158.0.0/15", "2400:cb00::/32", "2405:8100::/32"}

	resp := s3.GetBucketPolicyOutput{Policy: &input}
	result := parsePolicyResponse(&resp)
	resultMap := result.(map[string]interface{})
	addresses := GetSortedIPAddresses(resultMap)
	if len(addresses) != len(expected) {
		t.Errorf("Expected %d but got %d", len(expected), len(addresses))
	}
	for i, v := range expected {
		if addresses[i] != v {
			t.Errorf("Expected %s but got %s", v, addresses[i])
		}
	}
}

func TestMergeIPSliceIntoPolicy(t *testing.T) {
	expected := "{\"Statement\":[{\"Action\":\"s3:GetObject\",\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":[\"a\",\"b\"]}}}]}"

	policyStr := `{
		"Statement": [
		  {
			"Action": "s3:GetObject",
			"Condition": {
			  "IpAddress": {
				"aws:SourceIp": ["2400:cb00::/32", "2405:8100::/32", "162.158.0.0/15"]
			  }
			}
		  }
		]
	  }
	  `
	resp := s3.GetBucketPolicyOutput{Policy: &policyStr}
	policy := parsePolicyResponse(&resp)
	policyMap := policy.(map[string]interface{})
	result := MergeIPSliceIntoPolicy(policyMap, []string{"a", "b"})

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Errorf("Error occured during json marshal %v", err)
	}
	if string(encoded) != expected {
		t.Errorf("Expected %s but got %s", expected, string(encoded))
	}
}

func TestFindBlock(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`{
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
					"aws:SourceIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  }
			]
		  }`},
		{`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Sid": "VisualEditor",
				"Effect": "Allow",
				"Action": "rout53:Edit",
				"Resource": "arn:aws:*"
			  },
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*",
				"Condition": {
				  "IpAddress": {
					"aws:SourceIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  }
			]
		  }`},
		{`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*",
				"Condition": {
				  "Forbidden": true
				}
			  },
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*",
				"Condition": {
				  "IpAddress": {
					"aws:SourceIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  }
			]
		  }
		  `},
		{`{
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
					"aws:TargetIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  },
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*",
				"Condition": {
				  "IpAddress": {
					"aws:SourceIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  }
			]
		  }
		  `},
		{`{
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
					"aws:SourceIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  },
			  {
				"Sid": "VisualEditor",
				"Effect": "Allow",
				"Action": "rout53:Edit",
				"Resource": "arn:aws:*"
			  }
			]
		  }
		  `},
	}
	for _, test := range tests {
		resp := s3.GetBucketPolicyOutput{Policy: &test.input}
		policy := parsePolicyResponse(&resp)

		_, c := findBlock(policy)
		if reflect.TypeOf(*c).Kind() != reflect.Slice {
			t.Errorf("Expected type of slice")
		}
		addresses := (*c).([]interface{})
		if len(addresses) != 2 {
			t.Errorf("Expected %d but got %d", 2, len(addresses))
		}
		if addresses[0].(string) != "2400:cb00::/32" {
			t.Errorf("Expected %s but got %s", "2400:cb00::/32", addresses[0].(string))
		}
		if addresses[1].(string) != "162.158.0.0/15" {
			t.Errorf("Expected %s but got %s", "162.158.0.0/15", addresses[1].(string))
		}
	}
}

func TestFindNothingInBlock(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`{
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
					"aws:TargetIp": ["2400:cb00::/32", "162.158.0.0/15"]
				  }
				}
			  }
			]
		  }`},
		{`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*",
				"Condition": {
				  "Forbidden": true
				}
			  }
			]
		  }
		  `},
		{`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::logs.countableset.com/*"
			  }
			]
		  }
		  `},
		{`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "route53:Get",
				"Resource": "arn:aws:*"
			  }
			]
		  }
		  `},
	}
	for _, test := range tests {
		resp := s3.GetBucketPolicyOutput{Policy: &test.input}
		policy := parsePolicyResponse(&resp)

		_, result := findBlock(policy)
		if result != nil {
			t.Errorf("Expected type of nil")
		}
	}
}
