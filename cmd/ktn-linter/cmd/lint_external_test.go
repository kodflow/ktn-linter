// External tests for the lint command (black-box testing).
package cmd_test

import (
	"os"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
)

// TestLintConstants tests the exported constants from lint.go.
func TestLintConstants(t *testing.T) {
	tests := []struct {
		name     string
		check    func() bool
		errMsg   string
	}{
		{
			name:   "INITIAL_FILE_EDITS_CAP has reasonable value",
			check:  func() bool { return cmd.INITIAL_FILE_EDITS_CAP == 16 },
			errMsg: "expected INITIAL_FILE_EDITS_CAP to be 16",
		},
		{
			name:   "FILE_PERMISSION_RW is 0644",
			check:  func() bool { return cmd.FILE_PERMISSION_RW == os.FileMode(0644) },
			errMsg: "expected FILE_PERMISSION_RW to be 0644",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.errMsg)
			}
		})
	}
}
