package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Greatly inspired from the following link.
// https://www.cyphar.com/blog/post/20170412-golang-integration-coverage
func TestRunMain(t *testing.T) {
	var (
		args []string
	)

	for _, arg := range os.Args {
		switch {
		case strings.HasPrefix(arg, "-test"):
		case strings.HasPrefix(arg, "__bypass"):
			args = append(args, strings.TrimPrefix(arg, "__bypass"))
		default:
			args = append(args, arg)
		}
	}
	os.Args = args

	// To retrieve coverage results, os.Exit must NOT be called
	returnCode := mainReturnWithCode()
	fmt.Println("Netorcai return code:", returnCode)
}
