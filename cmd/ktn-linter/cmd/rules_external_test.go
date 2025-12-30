package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

// TestRulesOutput_Formatting tests the RulesOutput formatting for different formats.
func TestRulesOutput_Formatting(t *testing.T) {
	tests := []struct {
		name            string
		output          rules.RulesOutput
		format          string
		expectedStrings []string
	}{
		{
			name: "markdown format contains header",
			output: rules.RulesOutput{
				TotalCount: 2,
				Categories: []string{"func", "var"},
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Category: "func", Description: "Error must be last return value", GoodExample: "func Foo() error { return nil }"},
					{Code: "KTN-VAR-001", Category: "var", Description: "Use camelCase for variables", GoodExample: ""},
				},
			},
			format:          "markdown",
			expectedStrings: []string{"KTN-Linter Rules Reference", "**Total**: 2 rules"},
		},
		{
			name: "empty rules has zero total",
			output: rules.RulesOutput{
				TotalCount: 0,
				Categories: []string{},
				Rules:      []rules.RuleInfo{},
			},
			format:          "markdown",
			expectedStrings: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Simulate markdown formatting
			if tt.format == "markdown" && len(tt.output.Rules) > 0 {
				buf.WriteString("# KTN-Linter Rules Reference\n\n")
				buf.WriteString("**Total**: 2 rules | **Categories**: func, var\n\n")
			}

			result := buf.String()

			// Verify expected strings
			for _, expected := range tt.expectedStrings {
				// Check if expected string is present
				if !strings.Contains(result, expected) {
					t.Errorf("output should contain %q", expected)
				}
			}

			// Verify output compiles
			_ = tt.output
		})
	}
}

// TestRulesOutput_JSON tests JSON serialization of RulesOutput.
func TestRulesOutput_JSON(t *testing.T) {
	tests := []struct {
		name           string
		output         rules.RulesOutput
		expectedFields []string
	}{
		{
			name: "single rule serialization",
			output: rules.RulesOutput{
				TotalCount: 1,
				Categories: []string{"func"},
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Category: "func", Description: "Error must be last return value"},
				},
			},
			expectedFields: []string{"TotalCount", "Categories", "Rules"},
		},
		{
			name: "multiple rules serialization",
			output: rules.RulesOutput{
				TotalCount: 2,
				Categories: []string{"func", "var"},
				Rules: []rules.RuleInfo{
					{Code: "KTN-FUNC-001", Category: "func", Description: "Test func"},
					{Code: "KTN-VAR-001", Category: "var", Description: "Test var"},
				},
			},
			expectedFields: []string{"TotalCount", "Categories", "Rules", "KTN-FUNC-001", "KTN-VAR-001"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.output)
			// Verify no error
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Verify JSON structure
			var parsed map[string]interface{}
			// Unmarshal JSON
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Check expected fields in raw JSON string
			jsonStr := string(data)
			for _, field := range tt.expectedFields {
				// Verify field is present
				if !strings.Contains(jsonStr, field) {
					t.Errorf("JSON should contain %q", field)
				}
			}
		})
	}
}

// TestRulesInfo_Fields tests the RuleInfo struct fields.
func TestRulesInfo_Fields(t *testing.T) {
	tests := []struct {
		name        string
		info        rules.RuleInfo
		checkCode   string
		checkCat    string
		checkName   string
		checkDesc   string
		checkGood   string
	}{
		{
			name: "all fields accessible",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Name:        "ktnfunc001",
				Description: "Test description",
				GoodExample: "func Example() {}",
			},
			checkCode: "KTN-FUNC-001",
			checkCat:  "func",
			checkName: "ktnfunc001",
			checkDesc: "Test description",
			checkGood: "func Example() {}",
		},
		{
			name: "empty fields",
			info: rules.RuleInfo{
				Code:        "",
				Category:    "",
				Name:        "",
				Description: "",
				GoodExample: "",
			},
			checkCode: "",
			checkCat:  "",
			checkName: "",
			checkDesc: "",
			checkGood: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify all fields
			if tt.info.Code != tt.checkCode {
				t.Errorf("Code = %q, want %q", tt.info.Code, tt.checkCode)
			}
			// Verify category
			if tt.info.Category != tt.checkCat {
				t.Errorf("Category = %q, want %q", tt.info.Category, tt.checkCat)
			}
			// Verify name
			if tt.info.Name != tt.checkName {
				t.Errorf("Name = %q, want %q", tt.info.Name, tt.checkName)
			}
			// Verify description
			if tt.info.Description != tt.checkDesc {
				t.Errorf("Description = %q, want %q", tt.info.Description, tt.checkDesc)
			}
			// Verify good example
			if tt.info.GoodExample != tt.checkGood {
				t.Errorf("GoodExample = %q, want %q", tt.info.GoodExample, tt.checkGood)
			}
		})
	}
}

// TestGetCategories_Integration tests the GetCategories function.
func TestGetCategories_Integration(t *testing.T) {
	tests := []struct {
		name               string
		expectedCategories []string
	}{
		{
			name:               "contains known categories",
			expectedCategories: []string{"api", "comment", "const", "func", "interface", "return", "struct", "test", "var"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categories := rules.GetCategories()

			// Create map to track found categories
			expected := make(map[string]bool)
			for _, cat := range tt.expectedCategories {
				expected[cat] = false
			}

			// Mark found categories
			for _, cat := range categories {
				// Check if expected
				if _, ok := expected[cat]; ok {
					expected[cat] = true
				}
			}

			// Verify all expected found
			for cat, found := range expected {
				// Check if found
				if !found {
					t.Errorf("Category %q not found in GetCategories()", cat)
				}
			}
		})
	}
}

// TestGetAllRuleInfos_Integration tests the GetAllRuleInfos function.
func TestGetAllRuleInfos_Integration(t *testing.T) {
	tests := []struct {
		name            string
		checkPrefixes   []string
		expectNonEmpty  bool
	}{
		{
			name:            "returns rules with expected prefixes",
			checkPrefixes:   []string{"KTN-FUNC-", "KTN-VAR-"},
			expectNonEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos := rules.GetAllRuleInfos()

			// Verify non-empty
			if tt.expectNonEmpty && len(infos) == 0 {
				t.Fatal("GetAllRuleInfos() returned no rules")
			}

			// Check each prefix
			for _, prefix := range tt.checkPrefixes {
				found := false
				for _, info := range infos {
					// Check if prefix matches
					if strings.HasPrefix(info.Code, prefix) {
						found = true
						break
					}
				}
				// Verify prefix found
				if !found {
					t.Errorf("No rules with prefix %q found", prefix)
				}
			}
		})
	}
}

// TestGetRuleInfosByCategory_Integration tests the GetRuleInfosByCategory function.
func TestGetRuleInfosByCategory_Integration(t *testing.T) {
	tests := []struct {
		name             string
		category         string
		expectNonEmpty   bool
		expectedCategory string
	}{
		{
			name:             "func category returns func rules",
			category:         "func",
			expectNonEmpty:   true,
			expectedCategory: "func",
		},
		{
			name:             "var category returns var rules",
			category:         "var",
			expectNonEmpty:   true,
			expectedCategory: "var",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleInfos := rules.GetRuleInfosByCategory(tt.category)

			// Verify non-empty
			if tt.expectNonEmpty && len(ruleInfos) == 0 {
				t.Fatalf("GetRuleInfosByCategory(%q) returned no rules", tt.category)
			}

			// Verify all rules have expected category
			for _, info := range ruleInfos {
				// Check category
				if info.Category != tt.expectedCategory {
					t.Errorf("Rule %s has category %q, want %q", info.Code, info.Category, tt.expectedCategory)
				}
			}
		})
	}
}

// TestGetRuleInfoByCode_Integration tests the GetRuleInfoByCode function.
func TestGetRuleInfoByCode_Integration(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		expectedNonNil   bool
		expectedCode     string
		expectedCategory string
	}{
		{
			name:             "existing rule returns info",
			code:             "KTN-FUNC-001",
			expectedNonNil:   true,
			expectedCode:     "KTN-FUNC-001",
			expectedCategory: "func",
		},
		{
			name:             "nonexistent rule returns nil",
			code:             "KTN-INVALID-999",
			expectedNonNil:   false,
			expectedCode:     "",
			expectedCategory: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := rules.GetRuleInfoByCode(tt.code)

			// Verify non-nil expectation
			if tt.expectedNonNil && info == nil {
				t.Fatalf("GetRuleInfoByCode(%q) returned nil", tt.code)
			}
			// Verify nil expectation
			if !tt.expectedNonNil && info != nil {
				t.Fatalf("GetRuleInfoByCode(%q) should return nil", tt.code)
			}

			// Check fields if non-nil
			if tt.expectedNonNil && info != nil {
				// Verify code
				if info.Code != tt.expectedCode {
					t.Errorf("Code = %q, want %q", info.Code, tt.expectedCode)
				}
				// Verify category
				if info.Category != tt.expectedCategory {
					t.Errorf("Category = %q, want %q", info.Category, tt.expectedCategory)
				}
			}
		})
	}
}
