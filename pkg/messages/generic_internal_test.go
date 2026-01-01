package messages

import "testing"

// Test_registerGenericMessages tests that all generic messages are properly registered.
func Test_registerGenericMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "KTN-GENERIC-001_registered",
			code: "KTN-GENERIC-001",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Verifier que le message est enregistre
			if !ok {
				t.Errorf("message %s not found in registry", tt.code)
				// Retour anticipe
				return
			}
			// Verifier que le message n'est pas vide
			if msg.Short == "" {
				t.Errorf("message %s has empty Short", tt.code)
			}
			// Verifier que le message verbose n'est pas vide
			if msg.Verbose == "" {
				t.Errorf("message %s has empty Verbose", tt.code)
			}
		})
	}
}
