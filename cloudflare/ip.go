package cloudflare

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"gitlab.com/countableset/lambda-s3-cloudflare/util"
)

// https://www.cloudflare.com/ips/
const ipv4Url string = "https://www.cloudflare.com/ips-v4"
const ipv6Url string = "https://www.cloudflare.com/ips-v6"

// GetAllSortedIPAddresses returns a list of both ipv6 and ipv4 addresses
func GetAllSortedIPAddresses() []string {
	ips := append(GetIPv6Addresses(), GetIPv4Addresses()...)
	sort.Strings(ips)
	return ips
}

// GetIPv4Addresses returns a list of v4 addresses
func GetIPv4Addresses() []string {
	return fetchAndParseIPList(ipv4Url)
}

// GetIPv6Addresses returns a list of v6 addresses
func GetIPv6Addresses() []string {
	return fetchAndParseIPList(ipv6Url)
}

func fetchAndParseIPList(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		util.ExitErrorf("Unable to request ip list from url %q, %v", url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ExitErrorf("Unable to parse ip list from url %q, %v", url, err)
	}
	result := strings.Split(string(body), "\n")
	// Chop the last item if it's empty
	if result[len(result)-1] == "" {
		result = append(result[:len(result)-1])
	}
	return result
}
