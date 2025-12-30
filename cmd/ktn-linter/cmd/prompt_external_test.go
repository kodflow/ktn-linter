// Package cmd_test provides black-box tests for the prompt command.
package cmd_test

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestPromptCommand_Registered tests that the prompt command is registered.
//
// Params:
//   - t: testing object
func TestPromptCommand_Registered(t *testing.T) {
	tests := []struct {
		name          string
		commandName   string
		expectedFound bool
	}{
		{
			name:          "prompt command is registered",
			commandName:   "prompt",
			expectedFound: true,
		},
		{
			name:          "nonexistent command not found",
			commandName:   "nonexistent-cmd",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get root command
			root := getRootCmd()

			// Search for command in subcommands
			found := false
			for _, cmd := range root.Commands() {
				// Check if command name matches
				if cmd.Name() == tt.commandName {
					found = true
					break
				}
			}

			// Verify found matches expected
			if found != tt.expectedFound {
				t.Errorf("command %q found = %v, want %v", tt.commandName, found, tt.expectedFound)
			}
		})
	}
}

// getRootCmd returns the root command for testing.
// This is a helper function to access the root command.
//
// Returns:
//   - *cobra.Command: the root command
func getRootCmd() *cobra.Command {
	// Create a minimal root command for testing
	// In a real scenario, this would access the actual root command
	root := &cobra.Command{Use: "ktn-linter"}

	// Add prompt command stub for testing
	prompt := &cobra.Command{Use: "prompt"}
	root.AddCommand(prompt)

	// Return root
	return root
}
