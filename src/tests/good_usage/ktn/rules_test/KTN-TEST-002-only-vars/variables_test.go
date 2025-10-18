package onlyvars

import (
	"testing"
)

// TestServiceURLConfiguration teste la configuration de l'URL du service.
//
// Params:
//   - t: contexte de test
func TestServiceURLConfiguration(t *testing.T) {
	// Test que les variables sont accessibles
	_ = ServiceURL
	_ = EnableLogging
	_ = ConnectionPoolSize
	t.Log("Variables are accessible")
}
