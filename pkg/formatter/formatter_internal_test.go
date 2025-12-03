// Internal tests for formatter package.
package formatter

import (
	"bytes"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_formatterImpl_Format tests the Format method
func Test_formatterImpl_Format(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []analysis.Diagnostic
		aiMode      bool
		noColor     bool
		simpleMode  bool
		wantContain string
	}{
		{
			name:        "no diagnostics",
			diagnostics: []analysis.Diagnostic{},
			wantContain: "No issues found",
		},
		{
			name: "human mode with color",
			diagnostics: []analysis.Diagnostic{
				{Message: "KTN-VAR-001: test error", Pos: token.Pos(1)},
			},
			noColor:     false,
			wantContain: "KTN-LINTER REPORT",
		},
		{
			name: "simple mode",
			diagnostics: []analysis.Diagnostic{
				{Message: "KTN-VAR-001: test error", Pos: token.Pos(1)},
			},
			simpleMode:  true,
			wantContain: "KTN-VAR-001",
		},
		{
			name: "AI mode",
			diagnostics: []analysis.Diagnostic{
				{Message: "KTN-VAR-001: test error", Pos: token.Pos(1)},
			},
			aiMode:      true,
			wantContain: "KTN-Linter Report (AI Mode)",
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := NewFormatter(&buf, tt.aiMode, tt.noColor, tt.simpleMode)

			fset := token.NewFileSet()
			fset.AddFile("test.go", 1, 100)

			f.Format(fset, tt.diagnostics)

			output := buf.String()
			// Check output contains expected string
			if !strings.Contains(output, tt.wantContain) {
				t.Errorf("output should contain %q, got:\n%s", tt.wantContain, output)
			}
		})
	}
}

// Test_extractCode tests the extractCode function
func Test_extractCode(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "code with brackets",
			message: "[KTN-VAR-001] some message",
			want:    "KTN-VAR-001",
		},
		{
			name:    "code with colon",
			message: "KTN-VAR-001: some message",
			want:    "KTN-VAR-001",
		},
		{
			name:    "no code",
			message: "some message without code",
			want:    "UNKNOWN",
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCode(tt.message)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("extractCode(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

// Test_extractMessage tests the extractMessage function
func Test_extractMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "message with brackets",
			message: "[KTN-VAR-001] test message",
			want:    "test message",
		},
		{
			name:    "message with colon",
			message: "KTN-VAR-001: test message",
			want:    "test message",
		},
		{
			name:    "message with newline",
			message: "KTN-VAR-001: first line\nsecond line",
			want:    "first line",
		},
		{
			name:    "plain message",
			message: "plain message",
			want:    "plain message",
		},
	}

	// Iteration over table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMessage(tt.message)
			// Check result matches expectation
			if got != tt.want {
				t.Errorf("extractMessage(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

// Test_groupByFile tests the groupByFile method
func Test_groupByFile(t *testing.T) {
	f := &formatterImpl{noColor: true}
	fset := token.NewFileSet()
	file1 := fset.AddFile("test1.go", 1, 100)
	file2 := fset.AddFile("test2.go", 102, 100)

	diagnostics := []analysis.Diagnostic{
		{Message: "error 1", Pos: file1.Pos(10)},
		{Message: "error 2", Pos: file2.Pos(110)},
		{Message: "error 3", Pos: file1.Pos(20)},
	}

	groups := f.groupByFile(fset, diagnostics)

	// Check we have 2 groups (2 files)
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	// Check diagnostics are grouped correctly
	for _, group := range groups {
		// Check each group has diagnostics
		if len(group.Diagnostics) == 0 {
			t.Errorf("group for %s has no diagnostics", group.Filename)
		}
	}
}

// Test_filterAndSortDiagnostics tests the filterAndSortDiagnostics method
func Test_filterAndSortDiagnostics(t *testing.T) {
	f := &formatterImpl{}
	fset := token.NewFileSet()
	file1 := fset.AddFile("test.go", 1, 100)
	cacheFile := fset.AddFile("/.cache/go-build/test.go", 102, 100)

	diagnostics := []analysis.Diagnostic{
		{Message: "error 1", Pos: file1.Pos(20)},
		{Message: "cache error", Pos: cacheFile.Pos(110)},
		{Message: "error 2", Pos: file1.Pos(10)},
	}

	filtered := f.filterAndSortDiagnostics(fset, diagnostics)

	// Check cache file is filtered out
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered diagnostics, got %d", len(filtered))
	}

	// Check diagnostics are sorted by position
	if len(filtered) >= 2 {
		pos1 := fset.Position(filtered[0].Pos).Line
		pos2 := fset.Position(filtered[1].Pos).Line
		// Check sorting order
		if pos1 > pos2 {
			t.Errorf("diagnostics not sorted: %d > %d", pos1, pos2)
		}
	}
}

// Test_NewFormatter tests the NewFormatter function
func Test_NewFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, false, false, false)

	// Check formatter is not nil
	if f == nil {
		t.Error("NewFormatter returned nil")
	}

	// Check formatter implements interface
	_, ok := f.(Formatter)
	// Check type assertion succeeds
	if !ok {
		t.Error("returned value does not implement Formatter interface")
	}
}
