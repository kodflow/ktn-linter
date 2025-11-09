package test009

import "testing"

// TestExport teste Export (fonction publique) dans internal - ERREUR! // want "KTN-TEST-009: le test 'TestExport' dans 'wrong_internal_test.go' teste une fonction publique 'Export'"
func TestExport(t *testing.T) {
	result := Export()
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}
