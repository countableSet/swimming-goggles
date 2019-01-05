package cloudflare

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchIP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello\nThere")
	}))
	defer ts.Close()

	result := fetchAndParseIPList(ts.URL)
	if len(result) != 2 {
		t.Errorf("Incorrect result length, got %d, want: %d", 2, len(result))
	}
	if result[0] != "Hello" {
		t.Errorf("Incorrect result value, got %s, want: %s", "Hello", result[0])
	}
	if result[1] != "There" {
		t.Errorf("Incorrect result value, got %s, want: %s", "There", result[1])
	}
}
