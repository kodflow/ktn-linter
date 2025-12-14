// Package messages internal tests for struct messages.
package messages

import (
	"testing"
)

// TestStructMessagesRegistered tests that struct messages are registered.
func TestStructMessagesRegistered(t *testing.T) {
	structRules := []string{
		"KTN-STRUCT-001",
		"KTN-STRUCT-002",
		"KTN-STRUCT-003",
		"KTN-STRUCT-004",
		"KTN-STRUCT-005",
		"KTN-STRUCT-006",
		"KTN-STRUCT-007",
	}

	// Itération sur les règles
	for _, code := range structRules {
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
