package utils

import "testing"

func TestIsAllCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"all caps", "HTTP", true},
		{"mixed case", "Http", false},
		{"lowercase", "http", false},
		{"empty string", "", false},
		{"numbers only", "123", false},
		{"caps with numbers", "HTTP2", true},
		{"underscore with caps", "HTTP_URL", true},
		{"single lowercase", "h", false},
		{"single uppercase", "H", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAllCaps(tt.input)
			if got != tt.expected {
				t.Errorf("IsAllCaps(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsMixedCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid MixedCaps", "MyVariable", true},
		{"valid mixedCaps", "myVariable", true},
		{"snake_case", "my_variable", false},
		{"ALL_CAPS non-initialism", "MY_VAR", false},
		{"valid initialism", "HTTPServer", true},
		{"empty string", "", false},
		{"single char", "x", true},
		{"numbers", "var123", true},
		{"all caps valid", "HTTP", true},
		{"all caps invalid", "ABC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMixedCaps(tt.input)
			if got != tt.expected {
				t.Errorf("IsMixedCaps(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidInitialism(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid single initialism", "HTTP", true},
		{"valid double initialism", "HTTPURL", true},
		{"valid triple initialism", "HTTPURLID", true},
		{"invalid with underscore", "HTTP_URL", false},
		{"invalid mixed", "HTTPserver", false}, // server is not an initialism
		{"not an initialism", "NOTINIT", false},
		{"empty", "", false},
		{"partial match", "HTTPS", true},
		{"lowercase", "http", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidInitialism(tt.input)
			if got != tt.expected {
				t.Errorf("IsValidInitialism(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetKnownInitialisms(t *testing.T) {
	initialisms := getKnownInitialisms()

	// Verify some expected initialismsexist
	expected := []string{"HTTP", "HTTPS", "URL", "API", "JSON", "XML"}
	for _, exp := range expected {
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

	// Verify list is not empty
	if len(initialisms) == 0 {
		t.Error("getKnownInitialisms() returned empty list")
	}
}

func TestTryMatchInitialismPrefix(t *testing.T) {
	initialisms := getKnownInitialisms()

	tests := []struct {
		name            string
		input           string
		expectedRemaining string
		expectedMatch   bool
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
