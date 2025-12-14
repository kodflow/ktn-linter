// Package messages internal tests for comment messages.
package messages

import (
	"testing"
)

// TestCommentMessagesRegistered tests that comment messages are registered.
func TestCommentMessagesRegistered(t *testing.T) {
	commentRules := []string{
		"KTN-COMMENT-001",
		"KTN-COMMENT-002",
		"KTN-COMMENT-003",
		"KTN-COMMENT-004",
		"KTN-COMMENT-005",
		"KTN-COMMENT-006",
		"KTN-COMMENT-007",
	}

	// Itération sur les règles
	for _, code := range commentRules {
		t.Run(code, func(t *testing.T) {
			msg, found := Get(code)
			// Vérification existence
			if !found {
				t.Errorf("Get(%q) not found", code)
				return
			}
			// Vérification code
			if msg.Code != code {
				t.Errorf("msg.Code = %q, want %q", msg.Code, code)
			}
			// Vérification short non vide
			if msg.Short == "" {
				t.Errorf("msg.Short is empty for %q", code)
			}
		})
	}
}
