// Helper test - this should be skipped for KTN-TEST-006 because package ends with _test.
package test006_test

import "testing"

// TestHelper tests the helper.
func TestHelper(t *testing.T) {
	// Test helper
	result := Helper("test")
	// Vérification du résultat
	if result == "" {
		t.Error("empty result")
	}
}
