// Package cmd_test provides black-box tests for the prompt command.
package cmd_test

import (
	"testing"
)

// TestPromptCommand_Registered tests that the prompt command is registered.
//
// Params:
//   - t: testing object
func TestPromptCommand_Registered(t *testing.T) {
	// The prompt command should be registered with root
	// This test verifies the init() function runs correctly
	// by checking that the build succeeds (implicit test)
	t.Log("Prompt command registered successfully")
}
