// Package messages internal tests for test messages.
package messages

import (
	"testing"
)

// Test_registerTestMessages tests that all test messages are properly registered.
//
// Params:
//   - t *testing.T: Test instance for reporting test results.
func Test_registerTestMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-TEST-001", code: "KTN-TEST-001"},
		{name: "KTN-TEST-002", code: "KTN-TEST-002"},
		{name: "KTN-TEST-003", code: "KTN-TEST-003"},
		{name: "KTN-TEST-004", code: "KTN-TEST-004"},
		{name: "KTN-TEST-005", code: "KTN-TEST-005"},
		{name: "KTN-TEST-006", code: "KTN-TEST-006"},
		{name: "KTN-TEST-007", code: "KTN-TEST-007"},
		{name: "KTN-TEST-008", code: "KTN-TEST-008"},
		{name: "KTN-TEST-009", code: "KTN-TEST-009"},
		{name: "KTN-TEST-010", code: "KTN-TEST-010"},
		{name: "KTN-TEST-011", code: "KTN-TEST-011"},
		{name: "KTN-TEST-012", code: "KTN-TEST-012"},
		{name: "KTN-TEST-013", code: "KTN-TEST-013"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Verify message is registered
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Verify message code matches
			if msg.Code != tt.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify message short description is not empty
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", tt.code)
			}
		})
	}
}
