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
