// Package messages internal tests for const messages.
package messages

import (
	"testing"
)

// TestConstMessagesRegistered tests that const messages are registered.
func TestConstMessagesRegistered(t *testing.T) {
	constRules := []string{
		"KTN-CONST-001",
		"KTN-CONST-002",
		"KTN-CONST-003",
	}

	// Itération sur les règles
	for _, code := range constRules {
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
