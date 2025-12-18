// Mock test file - should be skipped by analyzer.
package test012

import "testing"

// TestMockService tests mock service.
func TestMockService(t *testing.T) {
	// This test should be skipped because it's in a mock file
	m := &MockService{}
	// Call without assertion (would be passthrough if not skipped)
	m.Execute("test")
}
