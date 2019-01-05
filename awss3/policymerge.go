package awss3

import (
	"sort"
)

// GetSortedIPAddresses returns the list of ip address sorts from the policy
func GetSortedIPAddresses(policy map[string]interface{}) []string {
	firstStatement := policy["Statement"].([]interface{})[0].(map[string]interface{})
	condition := firstStatement["Condition"].(map[string]interface{})
	ipAddr := condition["IpAddress"].(map[string]interface{})["aws:SourceIp"].([]interface{})
	ipAddresses := make([]string, len(ipAddr))
	for i, v := range ipAddr {
		ipAddresses[i] = v.(string)
	}
	sort.Strings(ipAddresses)
	return ipAddresses
}

// MergeIPSliceIntoPolicy todo
func MergeIPSliceIntoPolicy(policy map[string]interface{}, ipAddresses []string) interface{} {
	firstStatement := policy["Statement"].([]interface{})[0].(map[string]interface{})
	condition := firstStatement["Condition"].(map[string]interface{})
	condition["IpAddress"].(map[string]interface{})["aws:SourceIp"] = ipAddresses
	return policy
}