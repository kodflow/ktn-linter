// Package messages internal tests for func messages.
package messages

import (
	"testing"
)

// TestFuncMessagesRegistered tests that func messages are registered.
func TestFuncMessagesRegistered(t *testing.T) {
	funcRules := []string{
		"KTN-FUNC-001",
		"KTN-FUNC-002",
		"KTN-FUNC-003",
		"KTN-FUNC-004",
		"KTN-FUNC-005",
		"KTN-FUNC-006",
		"KTN-FUNC-007",
		"KTN-FUNC-008",
		"KTN-FUNC-009",
		"KTN-FUNC-010",
		"KTN-FUNC-011",
		"KTN-FUNC-012",
	}

	// Itération sur les règles
	for _, code := range funcRules {
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
