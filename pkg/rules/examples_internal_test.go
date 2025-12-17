package rules

import (
	"testing"
)

func Test_parseRuleCode(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		wantCategory string
		wantNumber   string
		wantErr      bool
	}{
		{
			name:         "valid func code",
			code:         "KTN-FUNC-001",
			wantCategory: "func",
			wantNumber:   "001",
			wantErr:      false,
		},
		{
			name:         "valid var code",
			code:         "KTN-VAR-002",
			wantCategory: "var",
			wantNumber:   "002",
			wantErr:      false,
		},
		{
			name:    "missing prefix",
			code:    "FUNC-001",
			wantErr: true,
		},
		{
			name:    "invalid format",
			code:    "KTN-FUNC",
			wantErr: true,
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, number, err := parseRuleCode(tt.code)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("parseRuleCode(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if category != tt.wantCategory {
					t.Errorf("parseRuleCode(%q) category = %q, want %q", tt.code, category, tt.wantCategory)
				}
				if number != tt.wantNumber {
					t.Errorf("parseRuleCode(%q) number = %q, want %q", tt.code, number, tt.wantNumber)
				}
			}
		})
	}
}

func Test_fileExists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: "examples.go",
			want: true,
		},
		{
			name: "non-existing file",
			path: "nonexistent.go",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func Test_findProjectRoot(t *testing.T) {
	root := findProjectRoot()
	// Should return a non-empty path
	if root == "" {
		t.Error("findProjectRoot() returned empty string")
	}
}

func TestNewInvalidCodeError(t *testing.T) {
	err := NewInvalidCodeError("BAD-CODE", "invalid format")

	if err.Code != "BAD-CODE" {
		t.Errorf("Code = %q, want 'BAD-CODE'", err.Code)
	}
	if err.Reason != "invalid format" {
		t.Errorf("Reason = %q, want 'invalid format'", err.Reason)
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}
}

// TestInvalidCodeError_Error tests the Error method on InvalidCodeError.
func TestInvalidCodeError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    *InvalidCodeError
		expect string
	}{
		{
			name:   "basic error",
			err:    &InvalidCodeError{Code: "BAD", Reason: "invalid"},
			expect: "invalid",
		},
		{
			name:   "empty reason",
			err:    &InvalidCodeError{Code: "TEST", Reason: ""},
			expect: "TEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result == "" {
				t.Error("Error() should not return empty string")
			}
		})
	}
}
