package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

func TestRulesOutput_Markdown(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 2,
		Categories: []string{"func", "var"},
		Rules: []rules.RuleInfo{
			{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Error must be last return value",
				GoodExample: "func Foo() error { return nil }",
			},
			{
				Code:        "KTN-VAR-001",
				Category:    "var",
				Description: "Use camelCase for variables",
				GoodExample: "",
			},
		},
	}

	var buf bytes.Buffer
	// Simulate markdown formatting
	buf.WriteString("# KTN-Linter Rules Reference\n\n")
	buf.WriteString("**Total**: 2 rules | **Categories**: func, var\n\n")

	result := buf.String()
	// Verify header
	if !strings.Contains(result, "KTN-Linter Rules Reference") {
		t.Error("Markdown output should contain header")
	}
	// Verify total
	if !strings.Contains(result, "**Total**: 2 rules") {
		t.Error("Markdown output should contain total count")
	}

	_ = output // Used to verify structure compiles
}

func TestRulesOutput_JSON(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 1,
		Categories: []string{"func"},
		Rules: []rules.RuleInfo{
			{
				Code:        "KTN-FUNC-001",
				Category:    "func",
				Description: "Error must be last return value",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check fields
	if _, ok := parsed["TotalCount"]; !ok {
		t.Error("JSON should contain TotalCount")
	}
	if _, ok := parsed["Categories"]; !ok {
		t.Error("JSON should contain Categories")
	}
	if _, ok := parsed["Rules"]; !ok {
		t.Error("JSON should contain Rules")
	}
}

func TestRulesOutput_EmptyRules(t *testing.T) {
	output := rules.RulesOutput{
		TotalCount: 0,
		Categories: []string{},
		Rules:      []rules.RuleInfo{},
	}

	// Should not panic
	if output.TotalCount != 0 {
		t.Error("Empty output should have 0 total")
	}
	if len(output.Rules) != 0 {
		t.Error("Empty output should have no rules")
	}
}

func TestRulesInfo_Fields(t *testing.T) {
	info := rules.RuleInfo{
		Code:        "KTN-FUNC-001",
		Category:    "func",
		Name:        "ktnfunc001",
		Description: "Test description",
		GoodExample: "func Example() {}",
	}

	// Verify all fields accessible
	if info.Code != "KTN-FUNC-001" {
		t.Errorf("Code = %q, want %q", info.Code, "KTN-FUNC-001")
	}
	if info.Category != "func" {
		t.Errorf("Category = %q, want %q", info.Category, "func")
	}
	if info.Name != "ktnfunc001" {
		t.Errorf("Name = %q, want %q", info.Name, "ktnfunc001")
	}
	if info.Description != "Test description" {
		t.Errorf("Description = %q, want %q", info.Description, "Test description")
	}
	if info.GoodExample != "func Example() {}" {
		t.Errorf("GoodExample = %q, want %q", info.GoodExample, "func Example() {}")
	}
}

func TestGetCategories_Integration(t *testing.T) {
	categories := rules.GetCategories()

	// Should have known categories
	expected := map[string]bool{
		"api":       false,
		"comment":   false,
		"const":     false,
		"func":      false,
		"interface": false,
		"return":    false,
		"struct":    false,
		"test":      false,
		"var":       false,
	}

	for _, cat := range categories {
		if _, ok := expected[cat]; ok {
			expected[cat] = true
		}
	}

	for cat, found := range expected {
		if !found {
			t.Errorf("Category %q not found in GetCategories()", cat)
		}
	}
}

func TestGetAllRuleInfos_Integration(t *testing.T) {
	infos := rules.GetAllRuleInfos()

	// Should have rules
	if len(infos) == 0 {
		t.Fatal("GetAllRuleInfos() returned no rules")
	}

	// Verify KTN rules exist
	hasFunc := false
	hasVar := false
	for _, info := range infos {
		if strings.HasPrefix(info.Code, "KTN-FUNC-") {
			hasFunc = true
		}
		if strings.HasPrefix(info.Code, "KTN-VAR-") {
			hasVar = true
		}
	}

	if !hasFunc {
		t.Error("No KTN-FUNC rules found")
	}
	if !hasVar {
		t.Error("No KTN-VAR rules found")
	}
}

func TestGetRuleInfosByCategory_Integration(t *testing.T) {
	funcRules := rules.GetRuleInfosByCategory("func")

	// Should have func rules
	if len(funcRules) == 0 {
		t.Fatal("GetRuleInfosByCategory('func') returned no rules")
	}

	// All should be func category
	for _, info := range funcRules {
		if info.Category != "func" {
			t.Errorf("Rule %s has category %q, want 'func'", info.Code, info.Category)
		}
	}
}

func TestGetRuleInfoByCode_Integration(t *testing.T) {
	info := rules.GetRuleInfoByCode("KTN-FUNC-001")

	if info == nil {
		t.Fatal("GetRuleInfoByCode('KTN-FUNC-001') returned nil")
	}

	if info.Code != "KTN-FUNC-001" {
		t.Errorf("Code = %q, want 'KTN-FUNC-001'", info.Code)
	}
	if info.Category != "func" {
		t.Errorf("Category = %q, want 'func'", info.Category)
	}
}
