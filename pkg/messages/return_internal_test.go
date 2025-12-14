// Package messages internal tests for return messages.
package messages

import (
	"testing"
)

// TestReturnMessagesRegistered tests that return messages are registered.
func TestReturnMessagesRegistered(t *testing.T) {
	returnRules := []string{
		"KTN-RETURN-002",
	}

	// Itération sur les règles
	for _, code := range returnRules {
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
