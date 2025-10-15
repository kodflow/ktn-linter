package withfunction_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_test/KTN-TEST-002-with-function"
)

// TestProcessData vérifie que ProcessData traite correctement les données.
//
// Params:
//   - t: contexte de test
func TestProcessData(t *testing.T) {
	result := withfunction.ProcessData("test")
	expected := "processed: test"
	if result != expected {
		t.Errorf("ProcessData() = %v, want %v", result, expected)
	}
}
