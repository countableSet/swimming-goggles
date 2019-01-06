package awss3

import (
	"sort"

	"gitlab.com/countableset/lambda-s3-cloudflare/util"
)

// GetSortedIPAddresses returns the list of ip address sorts from the policy
func GetSortedIPAddresses(policy map[string]interface{}) []string {
	_, c := findBlock(policy)
	if c == nil {
		util.ExitErrorf("Unable to find ip address field in policy %v", policy)
	}
	addressesInterface := (*c).([]interface{})
	ipAddresses := make([]string, len(addressesInterface))
	for i, v := range addressesInterface {
		ipAddresses[i] = v.(string)
	}
	sort.Strings(ipAddresses)
	return ipAddresses
}

// MergeIPSliceIntoPolicy returns an updated policy with the merges ip addresses for the parameter
func MergeIPSliceIntoPolicy(policy map[string]interface{}, ipAddresses []string) interface{} {
	p, c := findBlock(policy)
	if c == nil {
		util.ExitErrorf("Unable to find ip address field in policy %v", policy)
	}
	(*p).(map[string]interface{})["aws:SourceIp"] = ipAddresses
	return policy
}

func findBlock(root interface{}) (*interface{}, *interface{}) {
	conditions := []func(map[string]interface{}) (interface{}, bool){
		func(field map[string]interface{}) (interface{}, bool) {
			item, ok := field["Statement"]
			return item, ok
		},
		func(field map[string]interface{}) (interface{}, bool) {
			if action, ok := field["Action"]; ok && action == "s3:GetObject" {
				item, ok := field["Condition"]
				return item, ok
			}
			return nil, false
		},
		func(field map[string]interface{}) (interface{}, bool) {
			item, ok := field["IpAddress"]
			return item, ok
		},
		func(field map[string]interface{}) (interface{}, bool) {
			item, ok := field["aws:SourceIp"]
			return item, ok
		},
	}
	return testBlock(nil, &root, 0, &conditions)
}

func testBlock(
	parent *interface{},
	branch *interface{},
	index int,
	conditions *[]func(map[string]interface{}) (interface{}, bool),
) (*interface{}, *interface{}) {
	if index < len(*conditions) {
		switch branchType := (*branch).(type) {
		case []interface{}:
			for _, item := range branchType {
				p, result := testBlock(branch, &item, index, conditions)
				if result != nil {
					return p, result
				}
			}
		case map[string]interface{}:
			if field, ok := (*conditions)[index](branchType); ok {
				return testBlock(branch, &field, index+1, conditions)
			} else {
				return nil, nil
			}
		default:
			util.ExitErrorf("Invalid type of branch=%v", branch)
		}
		return nil, nil
	}
	return parent, branch
}
