// Package messages internal tests for struct messages.
package messages

import (
	"testing"
)

// Test_registerStructMessages validates that struct rule messages are properly registered.
func Test_registerStructMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-STRUCT-001 registered", code: "KTN-STRUCT-001"},
		{name: "KTN-STRUCT-002 registered", code: "KTN-STRUCT-002"},
		{name: "KTN-STRUCT-003 registered", code: "KTN-STRUCT-003"},
		{name: "KTN-STRUCT-004 registered", code: "KTN-STRUCT-004"},
		{name: "KTN-STRUCT-005 registered", code: "KTN-STRUCT-005"},
		{name: "KTN-STRUCT-006 registered", code: "KTN-STRUCT-006"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Verify message exists in registry
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Verify message code matches
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
