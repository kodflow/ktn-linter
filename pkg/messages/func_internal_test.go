// Package messages internal tests for func messages.
package messages

import (
	"testing"
)

// Test_registerFuncMessages validates that all func rule messages are properly registered.
func Test_registerFuncMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "func001", code: "KTN-FUNC-001"},
		{name: "func002", code: "KTN-FUNC-002"},
		{name: "func003", code: "KTN-FUNC-003"},
		{name: "func004", code: "KTN-FUNC-004"},
		{name: "func005", code: "KTN-FUNC-005"},
		{name: "func006", code: "KTN-FUNC-006"},
		{name: "func007", code: "KTN-FUNC-007"},
		{name: "func008", code: "KTN-FUNC-008"},
		{name: "func009", code: "KTN-FUNC-009"},
		{name: "func010", code: "KTN-FUNC-010"},
		{name: "func011", code: "KTN-FUNC-011"},
		{name: "func012", code: "KTN-FUNC-012"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Check if message is found
			msg, found := Get(tt.code)
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Check code matches
			if msg.Code != tt.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, tt.code)
			}
			// Check short message is not empty
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", tt.code)
			}
		})
	}
}
