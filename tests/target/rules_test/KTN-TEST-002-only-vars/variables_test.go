package onlyvars_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_test/KTN-TEST-002-only-vars"
)

// TestServiceURLConfiguration teste la configuration de l'URL du service.
//
// Params:
//   - t: contexte de test
func TestServiceURLConfiguration(t *testing.T) {
	// Test que les variables sont accessibles
	_ = onlyvars.ServiceURL
	_ = onlyvars.EnableLogging
	_ = onlyvars.ConnectionPoolSize
	t.Log("Variables are accessible")
}
