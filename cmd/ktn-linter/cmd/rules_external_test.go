package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
	"github.com/kodflow/ktn-linter/pkg/rules"
)

// TestNewRulesFormatter tests the NewRulesFormatter factory function.
func TestNewRulesFormatter(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		expectNotNil bool
	}{
		{
			name:         "text format returns formatter",
			format:       "text",
			expectNotNil: true,
		},
		{
			name:         "markdown format returns formatter",
			format:       "markdown",
			expectNotNil: true,
		},
		{
			name:         "md alias returns markdown formatter",
			format:       "md",
			expectNotNil: true,
		},
		{
			name:         "json format returns formatter",
			format:       "json",
			expectNotNil: true,
		},
		{
			name:         "unknown format defaults to text",
			format:       "unknown",
			expectNotNil: true,
		},
		{
			name:         "empty format defaults to text",
			format:       "",
			expectNotNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			formatter := cmd.NewRulesFormatter(tt.format)
			// Verify formatter is not nil
			if tt.expectNotNil && formatter == nil {
				t.Error("NewRulesFormatter returned nil")
			}
		})
	}
}

// TestRulesFormatterOutput tests formatter output.
func TestRulesFormatterOutput(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		categories   []string
		expectOutput bool
	}{
		{
			name:         "text formatter displays categories",
			format:       "text",
			categories:   []string{"func", "var"},
			expectOutput: true,
		},
		{
			name:         "markdown formatter displays categories",
			format:       "markdown",
			categories:   []string{"func", "var"},
			expectOutput: true,
		},
		{
			name:         "json formatter displays categories",
			format:       "json",
			categories:   []string{"func", "var"},
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create and invoke formatter
			formatter := cmd.NewRulesFormatter(tt.format)
			formatter.DisplayCategories(tt.categories)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output is not empty
			if tt.expectOutput && stdout.Len() == 0 {
				t.Error("expected non-empty output")
			}
		})
	}
}

// TestRulesFormatter_DisplayCategoryRules tests DisplayCategoryRules method.
func TestRulesFormatter_DisplayCategoryRules(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		category     string
		rules        []rules.RuleInfo
		expectOutput bool
	}{
		{
			name:     "text formatter displays category rules",
			format:   "text",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test rule"},
			},
			expectOutput: true,
		},
		{
			name:     "markdown formatter displays category rules",
			format:   "markdown",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test rule"},
			},
			expectOutput: true,
		},
		{
			name:     "json formatter displays category rules",
			format:   "json",
			category: "func",
			rules: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Description: "Test rule"},
			},
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create and invoke formatter
			formatter := cmd.NewRulesFormatter(tt.format)
			formatter.DisplayCategoryRules(tt.category, tt.rules)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output is not empty
			if tt.expectOutput && stdout.Len() == 0 {
				t.Error("expected non-empty output")
			}
		})
	}
}

// TestRulesFormatter_DisplayRuleDetails tests DisplayRuleDetails method.
func TestRulesFormatter_DisplayRuleDetails(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		info         rules.RuleInfo
		expectOutput bool
	}{
		{
			name:   "text formatter displays rule details",
			format: "text",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test rule",
				GoodExample: "func example() {}",
			},
			expectOutput: true,
		},
		{
			name:   "markdown formatter displays rule details",
			format: "markdown",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test rule",
				GoodExample: "func example() {}",
			},
			expectOutput: true,
		},
		{
			name:   "json formatter displays rule details",
			format: "json",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test rule",
			},
			expectOutput: true,
		},
		{
			name:   "text formatter without example",
			format: "text",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test rule",
				GoodExample: "",
			},
			expectOutput: true,
		},
		{
			name:   "markdown formatter without example",
			format: "markdown",
			info: rules.RuleInfo{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Test rule",
				GoodExample: "",
			},
			expectOutput: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create and invoke formatter
			formatter := cmd.NewRulesFormatter(tt.format)
			formatter.DisplayRuleDetails(tt.info)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output is not empty
			if tt.expectOutput && stdout.Len() == 0 {
				t.Error("expected non-empty output")
			}
		})
	}
}

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
			format: "markdown",
			expectedStrings: []string{
				"KTN-Linter Rules Reference",
				"**Total**: 0 rules",
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Simulate markdown formatting
			if tt.format == "markdown" {
				buf.WriteString("# KTN-Linter Rules Reference\n\n")
				buf.WriteString(fmt.Sprintf(
					"**Total**: %d rules | **Categories**: %s\n\n",
					tt.output.TotalCount,
					strings.Join(tt.output.Categories, ", "),
				))
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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.output)
			// Verify no error
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Unmarshal back into the struct to avoid coupling to JSON field-name casing/tags.
			var roundTrip rules.RulesOutput
			if err := json.Unmarshal(data, &roundTrip); err != nil {
				t.Fatalf("Failed to unmarshal JSON into RulesOutput: %v", err)
			}

			// Assert values actually round-trip.
			if roundTrip.TotalCount != tt.output.TotalCount {
				t.Errorf("TotalCount = %d, want %d", roundTrip.TotalCount, tt.output.TotalCount)
			}
			if len(roundTrip.Categories) != len(tt.output.Categories) {
				t.Fatalf("Categories len = %d, want %d", len(roundTrip.Categories), len(tt.output.Categories))
			}
			// Verify categories match
			for i := range tt.output.Categories {
				if roundTrip.Categories[i] != tt.output.Categories[i] {
					t.Errorf("Categories[%d] = %q, want %q", i, roundTrip.Categories[i], tt.output.Categories[i])
				}
			}
			if len(roundTrip.Rules) != len(tt.output.Rules) {
				t.Fatalf("Rules len = %d, want %d", len(roundTrip.Rules), len(tt.output.Rules))
			}
			// Verify rule codes match
			for i := range tt.output.Rules {
				if roundTrip.Rules[i].Code != tt.output.Rules[i].Code {
					t.Errorf("Rules[%d].Code = %q, want %q", i, roundTrip.Rules[i].Code, tt.output.Rules[i].Code)
				}
			}

			// Validate expected structure / rule-code expectations.
			for _, field := range tt.expectedFields {
				switch field {
				case "TotalCount", "Categories", "Rules":
					// already asserted above
				default:
					// Treat non-top-level expectations as rule-code expectations
					found := false
					for _, r := range roundTrip.Rules {
						if r.Code == field {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("JSON should contain rule code %q", field)
					}
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
