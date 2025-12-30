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
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
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
	tests := []struct {
		name           string
		infos          []rules.RuleInfo
		checkIndex     int
		wantHasExample bool
	}{
		{
			name: "valid rule has example",
			infos: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Category: "func"},
			},
			checkIndex:     0,
			wantHasExample: true,
		},
		{
			name: "fake rule has no example",
			infos: []rules.RuleInfo{
				{Code: "KTN-FAKE-999", Category: "fake"},
			},
			checkIndex:     0,
			wantHasExample: false,
		},
		{
			name: "empty code has no example",
			infos: []rules.RuleInfo{
				{Code: "", Category: ""},
			},
			checkIndex:     0,
			wantHasExample: false,
		},
		{
			name: "multiple rules preserves length",
			infos: []rules.RuleInfo{
				{Code: "KTN-FUNC-001", Category: "func"},
				{Code: "KTN-FUNC-002", Category: "func"},
			},
			checkIndex:     0,
			wantHasExample: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			enriched := rules.LoadGoodExamples(tt.infos)

			// Check length preserved
			if len(enriched) != len(tt.infos) {
				t.Errorf("LoadGoodExamples() returned %d items, want %d", len(enriched), len(tt.infos))
			}

			// Check example presence at index
			if tt.checkIndex < len(enriched) {
				hasExample := enriched[tt.checkIndex].GoodExample != ""
				if hasExample != tt.wantHasExample {
					t.Errorf("enriched[%d].GoodExample empty = %v, want %v", tt.checkIndex, !hasExample, !tt.wantHasExample)
				}
			}
		})
	}
}

func TestInvalidCodeError(t *testing.T) {
	tests := []struct {
		name            string
		code            string
		wantContains    string
		wantErrContains string
	}{
		{
			name:            "missing prefix",
			code:            "INVALID",
			wantContains:    "INVALID",
			wantErrContains: "invalid rule code",
		},
		{
			name:            "missing number",
			code:            "KTN-FUNC",
			wantContains:    "KTN-FUNC",
			wantErrContains: "invalid rule code",
		},
		{
			name:            "empty code",
			code:            "",
			wantContains:    "",
			wantErrContains: "invalid rule code",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			_, err := rules.GetTestdataPath(tt.code)
			if err == nil {
				t.Fatal("Expected error for invalid code")
			}

			errStr := err.Error()
			if tt.wantContains != "" && !strings.Contains(errStr, tt.wantContains) {
				t.Errorf("Error message should contain %q: %s", tt.wantContains, errStr)
			}
			if !strings.Contains(errStr, tt.wantErrContains) {
				t.Errorf("Error message should contain %q: %s", tt.wantErrContains, errStr)
			}
		})
	}
}

func TestNewInvalidCodeError(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		reason     string
		wantCode   string
		wantReason string
	}{
		{
			name:       "basic error",
			code:       "BAD-CODE",
			reason:     "invalid format",
			wantCode:   "BAD-CODE",
			wantReason: "invalid format",
		},
		{
			name:       "empty reason",
			code:       "TEST",
			reason:     "",
			wantCode:   "TEST",
			wantReason: "",
		},
		{
			name:       "empty code",
			code:       "",
			reason:     "missing code",
			wantCode:   "",
			wantReason: "missing code",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			err := rules.NewInvalidCodeError(tt.code, tt.reason)

			if err.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", err.Code, tt.wantCode)
			}
			if err.Reason != tt.wantReason {
				t.Errorf("Reason = %q, want %q", err.Reason, tt.wantReason)
			}

			errStr := err.Error()
			if errStr == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}

func TestInvalidCodeError_Error(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		reason       string
		wantContains string
	}{
		{
			name:         "basic error with reason",
			code:         "BAD",
			reason:       "invalid",
			wantContains: "invalid",
		},
		{
			name:         "error contains code",
			code:         "TEST-CODE",
			reason:       "some reason",
			wantContains: "TEST-CODE",
		},
		{
			name:         "empty reason still works",
			code:         "EMPTY",
			reason:       "",
			wantContains: "EMPTY",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			err := &rules.InvalidCodeError{Code: tt.code, Reason: tt.reason}
			result := err.Error()

			if result == "" {
				t.Error("Error() should not return empty string")
			}
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("Error() = %q, should contain %q", result, tt.wantContains)
			}
		})
	}
}
