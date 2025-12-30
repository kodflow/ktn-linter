// Internal tests for API messages.
package messages

import (
	"testing"
)

// Test_registerAPIMessages tests the registerAPIMessages function.
func Test_registerAPIMessages(t *testing.T) {
	tests := []struct {
		name         string
		ruleCode     string
		expectExists bool
	}{
		{
			name:         "KTN-API-001_registered",
			ruleCode:     "KTN-API-001",
			expectExists: true,
		},
	}

	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que le message existe
			msg, exists := Get(tt.ruleCode)
			// Vérifier l'existence
			if exists != tt.expectExists {
				t.Errorf("Get(%q) exists = %v, want %v", tt.ruleCode, exists, tt.expectExists)
				return
			}

			// Vérifier le code
			if tt.expectExists && msg.Code != tt.ruleCode {
				t.Errorf("Get(%q).Code = %q, want %q", tt.ruleCode, msg.Code, tt.ruleCode)
			}

			// Vérifier que Short n'est pas vide
			if tt.expectExists && msg.Short == "" {
				t.Errorf("Get(%q).Short is empty", tt.ruleCode)
			}

			// Vérifier que Verbose n'est pas vide
			if tt.expectExists && msg.Verbose == "" {
				t.Errorf("Get(%q).Verbose is empty", tt.ruleCode)
			}
		})
	}
}
