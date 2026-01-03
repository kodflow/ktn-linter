// Internal tests for the orchestrator.
package orchestrator

import (
	"bytes"
	"testing"
)

// TestOrchestrator_runSingleModule tests the runSingleModule method.
func TestOrchestrator_runSingleModule(t *testing.T) {
	tests := []struct {
		name        string
		patterns    []string
		expectError bool
	}{
		{
			name:        "run with invalid pattern",
			patterns:    []string{"./nonexistent/package"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			orch := NewOrchestrator(&buf, false)

			_, err := orch.runSingleModule("", tt.patterns, Options{})

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
