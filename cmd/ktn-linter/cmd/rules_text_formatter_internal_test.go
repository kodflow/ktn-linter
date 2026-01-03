// Internal tests for the text rules formatter.
// Note: Public methods (DisplayCategories, DisplayCategoryRules, DisplayRuleDetails)
// are tested through the interface in rules_text_formatter_external_test.go.
// This file is intentionally minimal as the textRulesFormatter has no private methods.
package cmd

import (
	"testing"
)

// Test_textRulesFormatter_struct verifies the struct exists and is usable.
func Test_textRulesFormatter_struct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "formatter struct can be instantiated",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify struct can be created
			f := &textRulesFormatter{}
			// Verify formatter is not nil
			if f == nil {
				t.Error("expected non-nil formatter")
			}
		})
	}
}
