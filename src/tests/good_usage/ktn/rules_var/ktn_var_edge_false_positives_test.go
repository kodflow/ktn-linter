package rules_var

import (
	"testing"

)

// TestValidHTTPCode teste TODO.
//
// Params:
//   - t: contexte de test
func TestValidHTTPCode(t *testing.T) {
	if ValidHTTPCode != 200 {
		t.Errorf("ValidHTTPCode = %v, want 200", ValidHTTPCode)
	}
}

// TestMaxHTTPRetries teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxHTTPRetries(t *testing.T) {
	if MaxHTTPRetries != 5 {
		t.Errorf("MaxHTTPRetries = %v, want 5", MaxHTTPRetries)
	}
}

// TestMaxRetriesLiteral teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxRetriesLiteral(t *testing.T) {
	if MaxRetriesLiteral != 3 {
		t.Errorf("MaxRetriesLiteral = %v, want 3", MaxRetriesLiteral)
	}
}

// TestMaxTimeout teste TODO.
//
// Params:
//   - t: contexte de test
func TestMaxTimeout(t *testing.T) {
	if MaxTimeout != 30 {
		t.Errorf("MaxTimeout = %v, want 30", MaxTimeout)
	}
}
