// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package cli

import (
	"fmt"
	"os"
)

const (
	// ExitSuccess means nominal status
	ExitSuccess = iota

	// ExitError means general error
	ExitError

	// ExitBadConnection means failed connection to remote service
	ExitBadConnection

	// ExitBadArgs means invalid argument values were given
	ExitBadArgs = 128
)

// Output prints the specified format message with arguments to stdout.
func Output(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, msg, args...)
}

// ExitWithOutput prints the specified entity and exits program with success.
func ExitWithOutput(msg string, output ...interface{}) {
	fmt.Fprintf(os.Stdout, msg, output...)
	os.Exit(ExitSuccess)
}

// ExitWithSuccess exits program with success without any output.
func ExitWithSuccess() {
	os.Exit(ExitSuccess)
}

// ExitWithError prints the specified error and exits program with the given error code.
func ExitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}

// ExitWithErrorMessage prints the specified message and exits program with the given error code.
func ExitWithErrorMessage(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(ExitError)
}
