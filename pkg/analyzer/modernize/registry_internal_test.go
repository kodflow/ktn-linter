// Internal tests for registry.go - modernize package.
package modernize

import "testing"

// Test_disabled tests that disabled analyzers are filtered out.
//
// Params:
//   - t: testing context
func Test_disabled(t *testing.T) {
	tests := []struct {
		name         string
		disabledName string
	}{
		{
			name:         "newexpr is disabled",
			disabledName: "newexpr",
		},
	}

	analyzers := Analyzers()
	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier que l'analyseur désactivé n'est pas présent
			for _, a := range analyzers {
				// Vérification de la condition
				if a.Name == tt.disabledName {
					t.Errorf("%s should be disabled but was found", tt.disabledName)
				}
			}
		})
	}
}
