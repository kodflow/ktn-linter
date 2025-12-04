// Package func004_special tests special functions handling.
package func004_special

import "github.com/spf13/cobra"

// main is the program entry point - should be ignored.
func main() {
	// Entry point
	run(nil, nil)
}

// init is called automatically - should be ignored.
func init() {
	// Auto-called
}

// run is used as a callback - should not trigger error.
func run(_ *cobra.Command, _ []string) error {
	// Callback function
	return nil
}

// helper is assigned to a variable - should not trigger error.
func helper() {
	// Assigned to var
}

// helperFunc holds the helper reference.
var helperFunc = helper
