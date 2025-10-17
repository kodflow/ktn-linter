package rules_var_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/good_usage/rules_var"
)

// TestValidHTTPCode teste TODO.
//
// Params:
//   - t: contexte de test
func TestValidHTTPCode(t *testing.T) {
	if rules_var.ValidHTTPCode != 200 {
		t.Errorf("rules_var.ValidHTTPCode = %v, want 200", rules_var.ValidHTTPCode)
	}
}

// TestMaxHTTPRetries teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxHTTPRetries(t *testing.T) {
	if rules_var.MaxHTTPRetries != 5 {
		t.Errorf("rules_var.MaxHTTPRetries = %v, want 5", rules_var.MaxHTTPRetries)
	}
}

// TestMaxRetriesLiteral teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxRetriesLiteral(t *testing.T) {
	if rules_var.MaxRetriesLiteral != 3 {
		t.Errorf("rules_var.MaxRetriesLiteral = %v, want 3", rules_var.MaxRetriesLiteral)
	}
}

// TestMaxTimeout teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxTimeout(t *testing.T) {
	if rules_var.MaxTimeout != 30 {
		t.Errorf("rules_var.MaxTimeout = %v, want 30", rules_var.MaxTimeout)
	}
}
