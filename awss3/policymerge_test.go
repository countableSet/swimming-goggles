package awss3

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestGetSortedIPAddresses(t *testing.T) {
	var tests = []struct {
		input    string
		expected []string
	}{
		{"{\"Statement\":[{\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":[\"2400:cb00::/32\",\"2405:8100::/32\",\"162.158.0.0/15\"]}}}]}",
			[]string{"162.158.0.0/15", "2400:cb00::/32", "2405:8100::/32"},
		},
	}
	for _, test := range tests {
		resp := s3.GetBucketPolicyOutput{Policy: &test.input}
		result := parsePolicyResponse(&resp)
		resultMap := result.(map[string]interface{})
		addresses := GetSortedIPAddresses(resultMap)
		if len(addresses) != len(test.expected) {
			t.Errorf("Expected %d but got %d", len(test.expected), len(addresses))
		}
		for i, v := range test.expected {
			if addresses[i] != v {
				t.Errorf("Expected %s but got %s", v, addresses[i])
			}
		}
	}
}

func TestMergeIPSliceIntoPolicy(t *testing.T) {
	expected := "{\"Statement\":[{\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":[\"a\",\"b\"]}}}]}"

	policyStr := "{\"Statement\":[{\"Condition\":{\"IpAddress\":{\"aws:SourceIp\":[\"2400:cb00::/32\",\"2405:8100::/32\",\"162.158.0.0/15\"]}}}]}"
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
