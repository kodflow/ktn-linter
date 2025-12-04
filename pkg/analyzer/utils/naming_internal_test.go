package utils

import "testing"

// TestGetKnownInitialisms tests the functionality of the corresponding implementation.
func TestGetKnownInitialisms(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
	}{
		{name: "contains expected initialisms", expected: []string{"HTTP", "HTTPS", "URL", "API", "JSON", "XML"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialisms := getKnownInitialisms()
			if len(initialisms) == 0 {
				t.Error("getKnownInitialisms() returned empty list")
				return
			}
			for _, exp := range tt.expected {
				found := false
				for _, init := range initialisms {
					if init == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected initialism %q not found in list", exp)
				}
			}
		})
	}
}

// TestTryMatchInitialismPrefix tests the functionality of the corresponding implementation.
func TestTryMatchInitialismPrefix(t *testing.T) {
	initialisms := getKnownInitialisms()

	tests := []struct {
		name              string
		input             string
		expectedRemaining string
		expectedMatch     bool
	}{
		{"match HTTP", "HTTPURL", "URL", true},
		{"match URL", "URLID", "ID", true},
		{"no match", "NOMATCH", "NOMATCH", false},
		{"empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, matched := tryMatchInitialismPrefix(tt.input, initialisms)
			if remaining != tt.expectedRemaining {
				t.Errorf("tryMatchInitialismPrefix(%q) remaining = %q, want %q",
					tt.input, remaining, tt.expectedRemaining)
			}
			if matched != tt.expectedMatch {
				t.Errorf("tryMatchInitialismPrefix(%q) matched = %v, want %v",
					tt.input, matched, tt.expectedMatch)
			}
		})
	}
}
