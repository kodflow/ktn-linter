package test008

import "testing"

// Testinitialize teste initialize (fonction privée)
func Testinitialize(t *testing.T) {
	result := initialize()
	if result != "initialized" {
		t.Errorf("Expected 'initialized', got %s", result)
	}
}
