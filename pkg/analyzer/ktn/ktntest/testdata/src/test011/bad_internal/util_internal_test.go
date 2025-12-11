package test011_test // want "KTN-TEST-011: le fichier 'util_internal_test.go' doit utiliser 'package test011'"

import "testing"

// TestPrivateUtil - ERREUR: package test011_test dans _internal_test.go
func TestPrivateUtil(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Ne peut pas accéder à privateUtil car on est dans package test011_test
			t.Log("This should fail - wrong package")

		})
	}
}
