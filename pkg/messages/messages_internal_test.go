// Package messages internal tests for private functions.
package messages

import (
	"testing"
)

// Test_registryInitialized tests that registry is initialized with messages.
//
// Params:
//   - t: testing context
func Test_registryInitialized(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-FUNC-001 exists", code: "KTN-FUNC-001"},
		{name: "KTN-VAR-001 exists", code: "KTN-VAR-001"},
		{name: "KTN-CONST-001 exists", code: "KTN-CONST-001"},
		{name: "KTN-COMMENT-001 exists", code: "KTN-COMMENT-001"},
		{name: "KTN-STRUCT-001 exists", code: "KTN-STRUCT-001"},
		{name: "KTN-TEST-004 exists", code: "KTN-TEST-004"},
	}

	// Vérification que le registre n'est pas vide
	if len(registry) == 0 {
		t.Fatal("registry is empty after init()")
	}

	// Itération sur les règles essentielles
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification existence
			if _, ok := registry[tt.code]; !ok {
				t.Errorf("registry missing essential rule %q", tt.code)
			}
		})
	}
}
