package rules_func

import (
	"testing"
)

// TestFindMaxWithComments teste la recherche du maximum avec commentaires.
//
// Params:
//   - t: contexte de test
func TestFindMaxWithComments(t *testing.T) {
	max, ok := FindMaxWithComments([]int{1, 5, 3})
	if !ok || max != 5 {
		t.Errorf("findMaxWithComments failed: got %d, %v", max, ok)
	}
}

// TestIsValidWithComments teste la validation avec commentaires.
//
// Params:
//   - t: contexte de test
func TestIsValidWithComments(t *testing.T) {
	if !IsValidWithComments(50) {
		t.Error("Expected valid for 50")
	}
	if IsValidWithComments(-1) {
		t.Error("Expected invalid for -1")
	}
}
