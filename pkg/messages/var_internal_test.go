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

// Test_registerVarMessages001To018 tests the first batch of VAR messages (001-018).
func Test_registerVarMessages001To018(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "VAR-001", code: "KTN-VAR-001"},
		{name: "VAR-002", code: "KTN-VAR-002"},
		{name: "VAR-003", code: "KTN-VAR-003"},
		{name: "VAR-004", code: "KTN-VAR-004"},
		{name: "VAR-005", code: "KTN-VAR-005"},
		{name: "VAR-006", code: "KTN-VAR-006"},
		{name: "VAR-007", code: "KTN-VAR-007"},
		{name: "VAR-008", code: "KTN-VAR-008"},
		{name: "VAR-009", code: "KTN-VAR-009"},
		{name: "VAR-010", code: "KTN-VAR-010"},
		{name: "VAR-011", code: "KTN-VAR-011"},
		{name: "VAR-012", code: "KTN-VAR-012"},
		{name: "VAR-013", code: "KTN-VAR-013"},
		{name: "VAR-014", code: "KTN-VAR-014"},
		{name: "VAR-015", code: "KTN-VAR-015"},
		{name: "VAR-016", code: "KTN-VAR-016"},
		{name: "VAR-017", code: "KTN-VAR-017"},
		{name: "VAR-018", code: "KTN-VAR-018"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Verify message exists
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Errorf("Short is empty for %q", tt.code)
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Errorf("Verbose is empty for %q", tt.code)
			}
		})
	}
}

// Test_registerVarMessages019To037 tests the second batch of VAR messages (019-037).
func Test_registerVarMessages019To037(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "VAR-019", code: "KTN-VAR-019"},
		{name: "VAR-020", code: "KTN-VAR-020"},
		{name: "VAR-021", code: "KTN-VAR-021"},
		{name: "VAR-022", code: "KTN-VAR-022"},
		{name: "VAR-023", code: "KTN-VAR-023"},
		{name: "VAR-024", code: "KTN-VAR-024"},
		{name: "VAR-025", code: "KTN-VAR-025"},
		{name: "VAR-026", code: "KTN-VAR-026"},
		{name: "VAR-027", code: "KTN-VAR-027"},
		{name: "VAR-028", code: "KTN-VAR-028"},
		{name: "VAR-029", code: "KTN-VAR-029"},
		{name: "VAR-030", code: "KTN-VAR-030"},
		{name: "VAR-031", code: "KTN-VAR-031"},
		{name: "VAR-033", code: "KTN-VAR-033"},
		{name: "VAR-034", code: "KTN-VAR-034"},
		{name: "VAR-035", code: "KTN-VAR-035"},
		{name: "VAR-036", code: "KTN-VAR-036"},
		{name: "VAR-037", code: "KTN-VAR-037"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Verify message exists
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Errorf("Short is empty for %q", tt.code)
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Errorf("Verbose is empty for %q", tt.code)
			}
		})
	}
}
