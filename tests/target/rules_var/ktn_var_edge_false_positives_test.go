package rules_var_test

import "testing"

func TestValidHTTPCode(t *testing.T) {
	if validHTTPCode != 200 {
		t.Errorf("validHTTPCode = %v, want 200", validHTTPCode)
	}
}

func TestMaxHTTPRetries(t *testing.T) {
	if maxHTTPRetries != 5 {
		t.Errorf("maxHTTPRetries = %v, want 5", maxHTTPRetries)
	}
}

func TestMaxRetriesLiteral(t *testing.T) {
	if MaxRetriesLiteral != 3 {
		t.Errorf("MaxRetriesLiteral = %v, want 3", MaxRetriesLiteral)
	}
}

func TestMaxTimeout(t *testing.T) {
	if MaxTimeout != 30 {
		t.Errorf("MaxTimeout = %v, want 30", MaxTimeout)
	}
}
