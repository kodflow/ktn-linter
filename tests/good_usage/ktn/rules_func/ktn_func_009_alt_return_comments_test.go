package rules_func_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/good_usage/rules_func"
)

// TestFindMaxWithComments teste la recherche du maximum avec commentaires.
//
// Params:
//   - t: contexte de test
func TestFindMaxWithComments(t *testing.T) {
	max, ok := rules_func.FindMaxWithComments([]int{1, 5, 3})
	if !ok || max != 5 {
		t.Errorf("findMaxWithComments failed: got %d, %v", max, ok)
	}
}

// TestIsValidWithComments teste la validation avec commentaires.
//
// Params:
//   - t: contexte de test
func TestIsValidWithComments(t *testing.T) {
	if !rules_func.IsValidWithComments(50) {
		t.Error("Expected valid for 50")
	}
	if rules_func.IsValidWithComments(-1) {
		t.Error("Expected invalid for -1")
	}
}
