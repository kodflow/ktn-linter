// Package messages internal tests for var messages.
package messages

import (
	"testing"
)

// TestVarMessagesRegistered tests that var messages are registered.
func TestVarMessagesRegistered(t *testing.T) {
	varRules := []string{
		"KTN-VAR-001",
		"KTN-VAR-002",
		"KTN-VAR-003",
		"KTN-VAR-004",
		"KTN-VAR-005",
		"KTN-VAR-006",
		"KTN-VAR-007",
		"KTN-VAR-008",
		"KTN-VAR-009",
		"KTN-VAR-010",
		"KTN-VAR-011",
		"KTN-VAR-012",
		"KTN-VAR-013",
		"KTN-VAR-014",
		"KTN-VAR-015",
		"KTN-VAR-016",
		"KTN-VAR-017",
		"KTN-VAR-018",
	}

	// Itération sur les règles
	for _, code := range varRules {
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
