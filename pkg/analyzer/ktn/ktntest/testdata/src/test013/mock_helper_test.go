// Mock helper test - should be skipped.
package test013_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test013"
)

// TestMockHelper tests mock helper (should be skipped).
func TestMockHelper(t *testing.T) {
	// Test mock helper - should not require error coverage
	result, _ := test013.MockHelper("test")
	// Vérification du résultat
	if result == "" {
		t.Log("empty")
	}
}
