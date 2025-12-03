// Tests with proper assertions.
package test013_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test013"
)

// TestProcessData teste ProcessData avec des assertions.
//
// Params:
//   - t: contexte de test
func TestProcessData(t *testing.T) {
	// Test avec assertion
	result := test013.ProcessData("hello")
	// Vérification avec assertion
	if result != "processed:hello" {
		t.Errorf("ProcessData() = %v, want %v", result, "processed:hello")
	}
}

// TestGetCount teste GetCount avec des assertions.
//
// Params:
//   - t: contexte de test
func TestGetCount(t *testing.T) {
	// Test avec assertion
	got := test013.GetCount()
	// Vérification avec assertion
	if got != 42 {
		t.Fatalf("GetCount() = %d, want 42", got)
	}
}
