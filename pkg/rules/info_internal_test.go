package rules

import (
	"testing"

	"golang.org/x/tools/go/analysis"
)

func Test_isValidRuleCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "valid code",
			code: "KTN-FUNC-001",
			want: true,
		},
		{
			name: "valid var code",
			code: "KTN-VAR-002",
			want: true,
		},
		{
			name: "missing prefix",
			code: "FUNC-001",
			want: false,
		},
		{
			name: "missing number",
			code: "KTN-FUNC",
			want: false,
		},
		{
			name: "empty code",
			code: "",
			want: false,
		},
		{
			name: "too many parts",
			code: "KTN-FUNC-001-EXTRA",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRuleCode(tt.code)
			if got != tt.want {
				t.Errorf("isValidRuleCode(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func Test_extractCodeWithoutColon(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "code without colon",
			doc:  "KTN-FUNC-001 error must be last",
			want: "KTN-FUNC-001",
		},
		{
			name: "empty doc",
			doc:  "",
			want: "",
		},
		{
			name: "invalid code",
			doc:  "not a valid code",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCodeWithoutColon(tt.doc)
			if got != tt.want {
				t.Errorf("extractCodeWithoutColon(%q) = %q, want %q", tt.doc, got, tt.want)
			}
		})
	}
}

func Test_analyzerToRuleInfo(t *testing.T) {
	// Test with nil analyzer would panic, so skip
	// This tests the conversion logic indirectly through exported functions
}

// Test_analyzersToRuleInfos tests conversion of analyzers to rule info.
func Test_analyzersToRuleInfos(t *testing.T) {
	tests := []struct {
		name    string
		input   []*analysis.Analyzer
		wantLen int
	}{
		{
			name:    "nil slice returns empty",
			input:   nil,
			wantLen: 0,
		},
		{
			name:    "empty slice returns empty",
			input:   []*analysis.Analyzer{},
			wantLen: 0,
		},
		{
			name: "valid analyzers",
			input: []*analysis.Analyzer{
				{Name: "ktnfunc001", Doc: "KTN-FUNC-001: test description"},
			},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzersToRuleInfos(tt.input)
			if len(result) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}
