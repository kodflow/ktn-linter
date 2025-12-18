// Mock test file - should be skipped.
package mock_only

import "testing"

// TestMockAPI tests the mock API.
func TestMockAPI(t *testing.T) {
	// Test mock - should be skipped
	m := &MockAPI{}
	result := m.Execute("test")
	// Vérification du résultat
	if result == "" {
		t.Error("empty result")
	}
}
