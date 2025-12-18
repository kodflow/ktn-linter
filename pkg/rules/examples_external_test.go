package rules_test

import (
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/rules"
)

func TestGetTestdataPath(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    string
		wantErr bool
	}{
		{
			name:    "func rule",
			code:    "KTN-FUNC-001",
			want:    "pkg/analyzer/ktn/ktnfunc/testdata/src/func001",
			wantErr: false,
		},
		{
			name:    "var rule",
			code:    "KTN-VAR-002",
			want:    "pkg/analyzer/ktn/ktnvar/testdata/src/var002",
			wantErr: false,
		},
		{
			name:    "struct rule",
			code:    "KTN-STRUCT-003",
			want:    "pkg/analyzer/ktn/ktnstruct/testdata/src/struct003",
			wantErr: false,
		},
		{
			name:    "invalid prefix",
			code:    "INVALID-001",
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing number",
			code:    "KTN-FUNC",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty code",
			code:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rules.GetTestdataPath(tt.code)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("GetTestdataPath(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTestdataPath(%q) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestLoadGoodExample(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantEmpty  bool
		wantSubstr string
	}{
		{
			name:       "existing rule",
			code:       "KTN-FUNC-001",
			wantEmpty:  false,
			wantSubstr: "package func001",
		},
		{
			name:      "non-existent rule",
			code:      "KTN-FAKE-999",
			wantEmpty: true,
		},
		{
			name:      "invalid code",
			code:      "INVALID",
			wantEmpty: true,
		},
		{
			name:      "empty code",
			code:      "",
			wantEmpty: true,
		},
		{
			name:      "code with wrong format",
			code:      "KTN-FUNC",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.LoadGoodExample(tt.code)
			gotEmpty := got == ""
			if gotEmpty != tt.wantEmpty {
				t.Errorf("LoadGoodExample(%q) empty = %v, want %v", tt.code, gotEmpty, tt.wantEmpty)
			}
			if tt.wantSubstr != "" && !strings.Contains(got, tt.wantSubstr) {
				t.Errorf("LoadGoodExample(%q) does not contain %q", tt.code, tt.wantSubstr)
			}
		})
	}
}

func TestLoadGoodExamples(t *testing.T) {
	// Create sample infos
	infos := []rules.RuleInfo{
		{Code: "KTN-FUNC-001", Category: "func"},
		{Code: "KTN-FUNC-002", Category: "func"},
		{Code: "KTN-FAKE-999", Category: "fake"},
		{Code: "", Category: ""},
	}

	// Load examples
	enriched := rules.LoadGoodExamples(infos)

	// Should have same length
	if len(enriched) != len(infos) {
		t.Errorf("LoadGoodExamples() returned %d items, want %d", len(enriched), len(infos))
	}

	// First two should have examples
	if enriched[0].GoodExample == "" {
		t.Error("LoadGoodExamples() should populate KTN-FUNC-001 example")
	}

	// Last one (fake) should be empty
	if enriched[2].GoodExample != "" {
		t.Error("LoadGoodExamples() should leave KTN-FAKE-999 example empty")
	}

	// Empty code should also be empty
	if enriched[3].GoodExample != "" {
		t.Error("LoadGoodExamples() should leave empty code example empty")
	}
}

func TestInvalidCodeError(t *testing.T) {
	_, err := rules.GetTestdataPath("INVALID")
	if err == nil {
		t.Fatal("Expected error for invalid code")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "INVALID") {
		t.Errorf("Error message should contain code: %s", errStr)
	}
	if !strings.Contains(errStr, "invalid rule code") {
		t.Errorf("Error message should describe the error: %s", errStr)
	}
}
