package ktn

import (
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TestAllRulesNotNil tests the functionality of the corresponding implementation.
func TestAllRulesNotNil(t *testing.T) {
	// Verify all rule categories are initialized
	if AllRules.Func == nil {
		t.Error("AllRules.Func should not be nil")
	}
	if AllRules.Var == nil {
		t.Error("AllRules.Var should not be nil")
	}
	if AllRules.Struct == nil {
		t.Error("AllRules.Struct should not be nil")
	}
	if AllRules.Interface == nil {
		t.Error("AllRules.Interface should not be nil")
	}
	if AllRules.Const == nil {
		t.Error("AllRules.Const should not be nil")
	}
	if AllRules.Error == nil {
		t.Error("AllRules.Error should not be nil")
	}
	if AllRules.Test == nil {
		t.Error("AllRules.Test should not be nil")
	}
	if AllRules.Alloc == nil {
		t.Error("AllRules.Alloc should not be nil")
	}
	if AllRules.Goroutine == nil {
		t.Error("AllRules.Goroutine should not be nil")
	}
	if AllRules.Pool == nil {
		t.Error("AllRules.Pool should not be nil")
	}
	if AllRules.Mock == nil {
		t.Error("AllRules.Mock should not be nil")
	}
	if AllRules.Method == nil {
		t.Error("AllRules.Method should not be nil")
	}
	if AllRules.Package == nil {
		t.Error("AllRules.Package should not be nil")
	}
	if AllRules.ControlFlow == nil {
		t.Error("AllRules.ControlFlow should not be nil")
	}
	if AllRules.DataStructures == nil {
		t.Error("AllRules.DataStructures should not be nil")
	}
	if AllRules.Ops == nil {
		t.Error("AllRules.Ops should not be nil")
	}
}

// TestGetAllRules tests the functionality of the corresponding implementation.
func TestGetAllRules(t *testing.T) {
	all := GetAllRules()

	// Verify we got a non-nil result
	if all == nil {
		t.Fatal("GetAllRules() should not return nil")
	}

	// Calculate expected total from all categories
	expectedCount := len(AllRules.Func) +
		len(AllRules.Var) +
		len(AllRules.Struct) +
		len(AllRules.Interface) +
		len(AllRules.Const) +
		len(AllRules.Error) +
		len(AllRules.Test) +
		len(AllRules.Alloc) +
		len(AllRules.Goroutine) +
		len(AllRules.Pool) +
		len(AllRules.Mock) +
		len(AllRules.Method) +
		len(AllRules.Package) +
		len(AllRules.ControlFlow) +
		len(AllRules.DataStructures) +
		len(AllRules.Ops)

	if len(all) != expectedCount {
		t.Errorf("GetAllRules() returned %d rules, expected %d", len(all), expectedCount)
	}

	// Verify all analyzers are non-nil
	for i, analyzer := range all {
		if analyzer == nil {
			t.Errorf("Analyzer at index %d is nil", i)
		}
	}
}

// TestGetRulesByCategory tests the functionality of the corresponding implementation.
func TestGetRulesByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected int
	}{
		{"func category", "func", len(AllRules.Func)},
		{"var category", "var", len(AllRules.Var)},
		{"struct category", "struct", len(AllRules.Struct)},
		{"interface category", "interface", len(AllRules.Interface)},
		{"const category", "const", len(AllRules.Const)},
		{"error category", "error", len(AllRules.Error)},
		{"test category", "test", len(AllRules.Test)},
		{"alloc category", "alloc", len(AllRules.Alloc)},
		{"goroutine category", "goroutine", len(AllRules.Goroutine)},
		{"pool category", "pool", len(AllRules.Pool)},
		{"mock category", "mock", len(AllRules.Mock)},
		{"method category", "method", len(AllRules.Method)},
		{"package category", "package", len(AllRules.Package)},
		{"control_flow category", "control_flow", len(AllRules.ControlFlow)},
		{"data_structures category", "data_structures", len(AllRules.DataStructures)},
		{"ops category", "ops", len(AllRules.Ops)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := GetRulesByCategory(tt.category)
			if len(rules) != tt.expected {
				t.Errorf("GetRulesByCategory(%q) returned %d rules, expected %d",
					tt.category, len(rules), tt.expected)
			}

			// Verify all returned analyzers are non-nil
			for i, analyzer := range rules {
				if analyzer == nil {
					t.Errorf("Analyzer at index %d is nil for category %q", i, tt.category)
				}
			}
		})
	}
}

// TestGetRulesByCategoryUnknown tests the functionality of the corresponding implementation.
func TestGetRulesByCategoryUnknown(t *testing.T) {
	tests := []string{
		"unknown",
		"invalid",
		"",
		"FUNC", // case sensitive
		"Func", // case sensitive
	}

	for _, category := range tests {
		t.Run("unknown_"+category, func(t *testing.T) {
			rules := GetRulesByCategory(category)
			if rules != nil {
				t.Errorf("GetRulesByCategory(%q) should return nil for unknown category, got %d rules",
					category, len(rules))
			}
		})
	}
}

// TestGetRulesByCategoryReturnsCorrectRules tests the functionality of the corresponding implementation.
func TestGetRulesByCategoryReturnsCorrectRules(t *testing.T) {
	// Verify that GetRulesByCategory returns the same slice as AllRules
	tests := []struct {
		category string
		expected []*analysis.Analyzer
	}{
		{"func", AllRules.Func},
		{"var", AllRules.Var},
		{"struct", AllRules.Struct},
		{"interface", AllRules.Interface},
		{"const", AllRules.Const},
		{"error", AllRules.Error},
		{"test", AllRules.Test},
		{"alloc", AllRules.Alloc},
		{"goroutine", AllRules.Goroutine},
		{"pool", AllRules.Pool},
		{"mock", AllRules.Mock},
		{"method", AllRules.Method},
		{"package", AllRules.Package},
		{"control_flow", AllRules.ControlFlow},
		{"data_structures", AllRules.DataStructures},
		{"ops", AllRules.Ops},
	}

	for _, tt := range tests {
		t.Run("verify_"+tt.category, func(t *testing.T) {
			rules := GetRulesByCategory(tt.category)

			// Check that we get the same number of rules
			if len(rules) != len(tt.expected) {
				t.Errorf("GetRulesByCategory(%q) returned %d rules, expected %d",
					tt.category, len(rules), len(tt.expected))
			}

			// Check that the rules match
			for i := range rules {
				if i >= len(tt.expected) {
					break
				}
				if rules[i] != tt.expected[i] {
					t.Errorf("GetRulesByCategory(%q)[%d] does not match AllRules.%s[%d]",
						tt.category, i, tt.category, i)
				}
			}
		})
	}
}

// TestAllRulesContainValidAnalyzers tests the functionality of the corresponding implementation.
func TestAllRulesContainValidAnalyzers(t *testing.T) {
	// Verify that all analyzers have valid names
	all := GetAllRules()

	namesSeen := make(map[string]bool)

	for _, analyzer := range all {
		if analyzer.Name == "" {
			t.Error("Found analyzer with empty name")
		}

		// Check for duplicates
		if namesSeen[analyzer.Name] {
			t.Errorf("Duplicate analyzer name found: %s", analyzer.Name)
		}
		namesSeen[analyzer.Name] = true

		// Verify Run function is not nil
		if analyzer.Run == nil {
			t.Errorf("Analyzer %s has nil Run function", analyzer.Name)
		}
	}
}

// TestGetAllRulesConsistency tests the functionality of the corresponding implementation.
func TestGetAllRulesConsistency(t *testing.T) {
	// Call GetAllRules multiple times and verify we get the same count
	firstCall := GetAllRules()
	secondCall := GetAllRules()

	if len(firstCall) != len(secondCall) {
		t.Errorf("GetAllRules() returned different counts: %d vs %d",
			len(firstCall), len(secondCall))
	}
}
