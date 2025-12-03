// External tests for the lint command (black-box testing).
package cmd_test

import (
	"os"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
)

// TestLintConstants tests the exported constants from lint.go.
func TestLintConstants(t *testing.T) {
	// Test INITIAL_FILE_EDITS_CAP
	t.Run("INITIAL_FILE_EDITS_CAP has reasonable value", func(t *testing.T) {
		// Vérification de la valeur
		if cmd.INITIAL_FILE_EDITS_CAP != 16 {
			t.Errorf("expected 16, got %d", cmd.INITIAL_FILE_EDITS_CAP)
		}
	})

	// Test FILE_PERMISSION_RW
	t.Run("FILE_PERMISSION_RW is 0644", func(t *testing.T) {
		const EXPECTED os.FileMode = 0644
		// Vérification de la valeur
		if cmd.FILE_PERMISSION_RW != EXPECTED {
			t.Errorf("expected %v, got %v", EXPECTED, cmd.FILE_PERMISSION_RW)
		}
	})
}
