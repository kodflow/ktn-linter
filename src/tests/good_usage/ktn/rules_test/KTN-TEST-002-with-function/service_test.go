package withfunction

import (
	"testing"

)

// TestProcessData vérifie que ProcessData traite correctement les données.
//
// Params:
//   - t: contexte de test
func TestProcessData(t *testing.T) {
	result := ProcessData("test")
	expected := "processed: test"
	if result != expected {
		t.Errorf("ProcessData() = %v, want %v", result, expected)
	}
}
