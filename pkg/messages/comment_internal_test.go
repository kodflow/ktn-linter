// Package messages internal tests for comment messages.
package messages

import (
	"testing"
)

// Test_registerCommentMessages tests the registerCommentMessages function.
func Test_registerCommentMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-COMMENT-001 registered", code: "KTN-COMMENT-001"},
		{name: "KTN-COMMENT-002 registered", code: "KTN-COMMENT-002"},
		{name: "KTN-COMMENT-003 registered", code: "KTN-COMMENT-003"},
		{name: "KTN-COMMENT-004 registered", code: "KTN-COMMENT-004"},
		{name: "KTN-COMMENT-005 registered", code: "KTN-COMMENT-005"},
		{name: "KTN-COMMENT-006 registered", code: "KTN-COMMENT-006"},
		{name: "KTN-COMMENT-007 registered", code: "KTN-COMMENT-007"},
	}

	// Itération sur les règles
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := Get(tt.code)
			// Vérification existence
			if !found {
				t.Errorf("Get(%q) not found", tt.code)
				return
			}
			// Vérification code
			if msg.Code != tt.code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, tt.code)
			}
			// Vérification short non vide
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", tt.code)
			}
		})
	}
}
