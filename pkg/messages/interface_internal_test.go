// Package messages internal tests for interface messages.
package messages

import (
	"testing"
)

// Test_registerInterfaceMessages verifies that all interface rule messages are properly registered.
func Test_registerInterfaceMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "KTN-INTERFACE-001 registered",
			code: "KTN-INTERFACE-001",
		},
	}

	// Iterate through test cases
	for _, test := range tests {
		test := test // Capture range variable
		t.Run(test.name, func(t *testing.T) {
			msg, found := Get(test.code)
			// Verify message is found
			if !found {
				t.Errorf("Get(%q) not found", test.code)
				return
			}
			// Verify code matches
			if msg.Code != test.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, test.code)
			}
			// Verify short description is not empty
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", test.code)
			}
		})
	}
}
