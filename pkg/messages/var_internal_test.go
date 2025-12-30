// Package messages internal tests for var messages.
package messages

import (
	"testing"
)

// Test_registerVarMessages verifies that all var category messages are properly registered.
func Test_registerVarMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-VAR-001", code: "KTN-VAR-001"},
		{name: "KTN-VAR-002", code: "KTN-VAR-002"},
		{name: "KTN-VAR-003", code: "KTN-VAR-003"},
		{name: "KTN-VAR-004", code: "KTN-VAR-004"},
		{name: "KTN-VAR-005", code: "KTN-VAR-005"},
		{name: "KTN-VAR-006", code: "KTN-VAR-006"},
		{name: "KTN-VAR-007", code: "KTN-VAR-007"},
		{name: "KTN-VAR-008", code: "KTN-VAR-008"},
		{name: "KTN-VAR-009", code: "KTN-VAR-009"},
		{name: "KTN-VAR-010", code: "KTN-VAR-010"},
		{name: "KTN-VAR-011", code: "KTN-VAR-011"},
		{name: "KTN-VAR-012", code: "KTN-VAR-012"},
		{name: "KTN-VAR-013", code: "KTN-VAR-013"},
		{name: "KTN-VAR-014", code: "KTN-VAR-014"},
		{name: "KTN-VAR-015", code: "KTN-VAR-015"},
		{name: "KTN-VAR-016", code: "KTN-VAR-016"},
		{name: "KTN-VAR-017", code: "KTN-VAR-017"},
		{name: "KTN-VAR-018", code: "KTN-VAR-018"},
	}

	// Test each rule code in the table.
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Verify message exists
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short description is not empty
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", tt.code)
			}
		})
	}
}
