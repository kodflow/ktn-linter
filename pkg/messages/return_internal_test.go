// Package messages internal tests for return messages.
package messages

import (
	"testing"
)

// Test_registerReturnMessages tests that return messages are properly registered.
func Test_registerReturnMessages(t *testing.T) {
	// Define test cases for return rule codes
	tests := []struct {
		name string
		code string
	}{
		{
			name: "KTN-RETURN-001 registered",
			code: "KTN-RETURN-001",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify message exists
			msg, found := Get(tt.code)
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
