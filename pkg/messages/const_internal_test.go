// Package messages internal tests for const messages.
package messages

import (
	"testing"
)

// Test_registerConstMessages tests the registerConstMessages function.
func Test_registerConstMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-CONST-001 registered", code: "KTN-CONST-001"},
		{name: "KTN-CONST-002 registered", code: "KTN-CONST-002"},
		{name: "KTN-CONST-003 registered", code: "KTN-CONST-003"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			if msg.Code != tt.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, tt.code)
			}
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", tt.code)
			}
		})
	}
}
