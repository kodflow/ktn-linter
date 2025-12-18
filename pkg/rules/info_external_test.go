package rules_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

func TestExtractRuleCode(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "valid KTN code with description",
			doc:  "KTN-FUNC-001: error must be last return value",
			want: "KTN-FUNC-001",
		},
		{
			name: "valid KTN code without description",
			doc:  "KTN-VAR-002",
			want: "KTN-VAR-002",
		},
		{
			name: "non-KTN rule",
			doc:  "modernize: use modern Go idioms",
			want: "",
		},
		{
			name: "empty doc",
			doc:  "",
			want: "",
		},
		{
			name: "KTN prefix only",
			doc:  "KTN-FUNC: missing number",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ExtractRuleCode(tt.doc)
			if got != tt.want {
				t.Errorf("ExtractRuleCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "valid doc with description",
			doc:  "KTN-FUNC-001: error must be last return value",
			want: "error must be last return value",
		},
		{
			name: "doc without colon",
			doc:  "KTN-FUNC-001 error must be last",
			want: "KTN-FUNC-001 error must be last",
		},
		{
			name: "empty doc",
			doc:  "",
			want: "",
		},
		{
			name: "colon at end",
			doc:  "KTN-FUNC-001:",
			want: "KTN-FUNC-001:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ExtractDescription(tt.doc)
			if got != tt.want {
				t.Errorf("ExtractDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractCategory(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "func category",
			code: "KTN-FUNC-001",
			want: "func",
		},
		{
			name: "var category",
			code: "KTN-VAR-002",
			want: "var",
		},
		{
			name: "uppercase category normalized",
			code: "KTN-STRUCT-003",
			want: "struct",
		},
		{
			name: "invalid code",
			code: "INVALID",
			want: "",
		},
		{
			name: "empty code",
			code: "",
			want: "",
		},
		{
			name: "KTN prefix but missing number",
			code: "KTN-FUNC",
			want: "",
		},
		{
			name: "KTN prefix but too many parts",
			code: "KTN-FUNC-001-EXTRA",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ExtractCategory(tt.code)
			if got != tt.want {
				t.Errorf("ExtractCategory() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetAllRuleInfos(t *testing.T) {
	infos := rules.GetAllRuleInfos()

	// Should have multiple rules
	if len(infos) == 0 {
		t.Error("GetAllRuleInfos() returned empty slice")
	}

	// Should be sorted by code
	for i := 1; i < len(infos); i++ {
		if infos[i-1].Code >= infos[i].Code {
			t.Errorf("Rules not sorted: %s >= %s", infos[i-1].Code, infos[i].Code)
		}
	}

	// Each rule should have required fields
	for _, info := range infos {
		if info.Code == "" {
			t.Error("Rule with empty code found")
		}
		if info.Category == "" {
			t.Errorf("Rule %s has empty category", info.Code)
		}
		if info.Name == "" {
			t.Errorf("Rule %s has empty name", info.Code)
		}
	}
}

func TestGetRuleInfosByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantMin  int
	}{
		{
			name:     "func category",
			category: "func",
			wantMin:  5,
		},
		{
			name:     "var category",
			category: "var",
			wantMin:  5,
		},
		{
			name:     "unknown category",
			category: "unknown",
			wantMin:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos := rules.GetRuleInfosByCategory(tt.category)
			if len(infos) < tt.wantMin {
				t.Errorf("GetRuleInfosByCategory(%q) = %d rules, want >= %d", tt.category, len(infos), tt.wantMin)
			}

			// All returned rules should be in requested category
			for _, info := range infos {
				if info.Category != tt.category {
					t.Errorf("Rule %s has category %q, want %q", info.Code, info.Category, tt.category)
				}
			}
		})
	}
}

func TestGetRuleInfoByCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantNil bool
	}{
		{
			name:    "existing rule",
			code:    "KTN-FUNC-001",
			wantNil: false,
		},
		{
			name:    "non-existent rule",
			code:    "KTN-FAKE-999",
			wantNil: true,
		},
		{
			name:    "invalid format",
			code:    "INVALID",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := rules.GetRuleInfoByCode(tt.code)
			gotNil := info == nil
			if gotNil != tt.wantNil {
				t.Errorf("GetRuleInfoByCode(%q) nil = %v, want %v", tt.code, gotNil, tt.wantNil)
			}
			if info != nil && info.Code != tt.code {
				t.Errorf("GetRuleInfoByCode(%q).Code = %q", tt.code, info.Code)
			}
		})
	}
}

func TestGetCategories(t *testing.T) {
	categories := rules.GetCategories()

	// Should have multiple categories
	if len(categories) == 0 {
		t.Error("GetCategories() returned empty slice")
	}

	// Check expected categories exist
	expected := map[string]bool{
		"func":   false,
		"var":    false,
		"struct": false,
		"const":  false,
	}

	for _, cat := range categories {
		if _, ok := expected[cat]; ok {
			expected[cat] = true
		}
	}

	for cat, found := range expected {
		if !found {
			t.Errorf("Expected category %q not found", cat)
		}
	}

	// Should be sorted
	for i := 1; i < len(categories); i++ {
		if categories[i-1] >= categories[i] {
			t.Errorf("Categories not sorted: %s >= %s", categories[i-1], categories[i])
		}
	}
}
