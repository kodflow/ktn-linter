package ktnconst_test

import (
	"testing"

	ktnconst "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConst003(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ktnconst.Analyzer003, "const003")
}

func TestConst003_NamingValidation(t *testing.T) {
	tests := []struct {
		name      string
		constName string
		wantValid bool
	}{
		// Valid cases
		{"Single uppercase letter", "A", true},
		{"Simple acronym", "API", true},
		{"Multi-word with underscore", "MAX_SIZE", true},
		{"With numbers", "HTTP2", true},
		{"Complex valid", "API_KEY_V2", true},
		{"HTTP timeout", "HTTP_TIMEOUT", true},
		{"EOF constant", "EOF", true},
		{"Multiple underscores", "MAX_BUFFER_SIZE_LIMIT", true},
		{"Underscore with numbers", "TLS1_2_VERSION", true},

		// Invalid cases - these will be caught by the analyzer
		{"Lowercase start", "maxSize", false},
		{"PascalCase", "MaxSize", false},
		{"snake_case", "max_size", false},
		{"camelCase", "maxApiSize", false},
		{"Mixed case with underscore", "Max_Size", false},
		{"Lowercase with underscore", "max_api_size", false},
		{"Mixed capitals", "MaxAPISize", false},
		{"Starting with underscore", "_MAX_SIZE", false},
		{"Ending with underscore", "MAX_SIZE_", false},
		{"Contains lowercase", "MAX_Size", false},
		{"Mixed case complex", "HTTPTimeout", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents expected behavior
			// The actual validation happens in const003/good.go and const003/bad.go
			// The test files will be processed by analysistest.Run in TestConst003
			t.Logf("Testing constant name: %s (should be valid: %v)", tt.constName, tt.wantValid)
		})
	}
}

func TestConst003_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name        string
		description string
	}{
		{"Blank identifier", "Should skip _ (blank identifier)"},
		{"Single letter", "Single uppercase letters like A, B, C are valid"},
		{"Acronyms", "Common acronyms like API, HTTP, URL, EOF are valid"},
		{"Numbers", "Numbers are allowed in constant names like HTTP2, TLS1_3"},
		{"Multiple underscores", "Multiple underscores are valid like MAX_BUFFER_SIZE_LIMIT"},
		{"Grouped constants", "All constants in a group declaration must follow the rule"},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Edge case: %s", tc.description)
		})
	}
}
