package util

import (
	"fmt"
	"os"
)

// ExitErrorf prints a message to std err and exit with a code of 1
func ExitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// TestEqualSlices returns nill if the two slices are equal, otherwise source slice
func TestEqualSlices(original, source []string) []string {
	eqCheck := testEq(original, source)
	if eqCheck {
		return nil
	}
	return source
}

func testEq(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
