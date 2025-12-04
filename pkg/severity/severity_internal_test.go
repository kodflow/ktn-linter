// Internal tests for severity.go (white-box testing).
package severity

import (
	"testing"
)

// Test_rulesSeverity tests the internal rulesSeverity map.
func Test_rulesSeverity(t *testing.T) {
	tests := []struct {
		name     string
		ruleCode string
		expected Level
	}{
		{"KTN-VAR-001 is ERROR", "KTN-VAR-001", SEVERITY_ERROR},
		{"KTN-FUNC-001 is ERROR", "KTN-FUNC-001", SEVERITY_ERROR},
		{"KTN-TEST-013 is INFO", "KTN-TEST-013", SEVERITY_INFO},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			level, ok := rulesSeverity[tt.ruleCode]
			// Vérification existence
			if !ok {
				t.Errorf("rule %q not found in rulesSeverity", tt.ruleCode)
				return
			}
			// Vérification niveau
			if level != tt.expected {
				t.Errorf("rulesSeverity[%q] = %v, want %v", tt.ruleCode, level, tt.expected)
			}
		})
	}
}

// Test_rulesSeverityCompleteness tests that all rule categories are covered.
func Test_rulesSeverityCompleteness(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
	}{
		{name: "all categories covered", categories: []string{"COMMENT", "CONST", "VAR", "FUNC", "STRUCT", "TEST"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(rulesSeverity) == 0 {
				t.Error("rulesSeverity should not be empty")
				return
			}

			for _, cat := range tt.categories {
				found := false
				for rule := range rulesSeverity {
					if len(rule) > 4 && rule[4:4+len(cat)] == cat {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("No rules found for category %q", cat)
				}
			}
		})
	}
}
