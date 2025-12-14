// Package messages internal tests for interface messages.
package messages

import (
	"testing"
)

// TestInterfaceMessagesRegistered tests that interface messages are registered.
func TestInterfaceMessagesRegistered(t *testing.T) {
	interfaceRules := []string{
		"KTN-INTERFACE-001",
	}

	// Itération sur les règles
	for _, code := range interfaceRules {
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
