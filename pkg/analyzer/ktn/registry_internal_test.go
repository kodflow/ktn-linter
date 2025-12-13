// Internal tests for registry in ktn package.
package ktn

import (
	"testing"
)

// Test_categoryAnalyzers tests that categoryAnalyzers returns valid map.
//
// Params:
//   - t: testing context
func Test_categoryAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns valid map"},
	}

	// Iteration over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories := categoryAnalyzers()
			// Check that map is not empty
			if len(categories) == 0 {
				t.Error("categoryAnalyzers() returned empty map")
				return
			}
			// Check all category functions work
			for name, fn := range categories {
				analyzers := fn()
				// Check analyzers are not nil
				if analyzers == nil {
					t.Errorf("category %q returned nil slice", name)
				}
			}
		})
	}
}

// Test_codeToAnalyzerName tests the code to analyzer name conversion.
//
// Params:
//   - t: testing context
func Test_codeToAnalyzerName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{name: "valid func code", code: "KTN-FUNC-001", expected: "ktnfunc001"},
		{name: "valid var code", code: "KTN-VAR-002", expected: "ktnvar002"},
		{name: "valid const code", code: "KTN-CONST-003", expected: "ktnconst003"},
		{name: "valid test code", code: "KTN-TEST-008", expected: "ktntest008"},
		{name: "valid struct code", code: "KTN-STRUCT-001", expected: "ktnstruct001"},
		{name: "invalid prefix", code: "XXX-FUNC-001", expected: ""},
		{name: "missing parts", code: "KTN-FUNC", expected: ""},
		{name: "too many parts", code: "KTN-FUNC-001-extra", expected: ""},
		{name: "empty string", code: "", expected: ""},
		{name: "no prefix", code: "FUNC-001", expected: ""},
	}

	// Iteration over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := codeToAnalyzerName(tt.code)
			// Check result matches expected
			if result != tt.expected {
				t.Errorf("codeToAnalyzerName(%q) = %q, want %q", tt.code, result, tt.expected)
			}
		})
	}
}
