// Internal tests for formatter package.
package formatter

import (
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_extractCode tests the extractCode function.
//
// Params:
//   - t: testing context
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

// Test_groupByFile tests the groupByFile method.
//
// Params:
//   - t: testing context
func Test_groupByFile(t *testing.T) {
	tests := []struct {
		name           string
		expectedGroups int
	}{
		{name: "groups diagnostics by file", expectedGroups: 2},
	}

	// Iteration over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// Check groups count
			if len(groups) != tt.expectedGroups {
				t.Errorf("expected %d groups, got %d", tt.expectedGroups, len(groups))
				return
			}
			// Check each group has diagnostics
			for _, group := range groups {
				// Verify group is not empty
				if len(group.Diagnostics) == 0 {
					t.Errorf("group for %s has no diagnostics", group.Filename)
				}
			}
		})
	}
}

// Test_filterAndSortDiagnostics tests the filterAndSortDiagnostics method.
//
// Params:
//   - t: testing context
func Test_filterAndSortDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
	}{
		{name: "filters cache files and sorts", expectedCount: 2},
	}

	// Iteration over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// Check count
			if len(filtered) != tt.expectedCount {
				t.Errorf("expected %d filtered diagnostics, got %d", tt.expectedCount, len(filtered))
				return
			}
			// Check sorting
			if len(filtered) >= 2 {
				pos1 := fset.Position(filtered[0].Pos).Line
				pos2 := fset.Position(filtered[1].Pos).Line
				// Verify sorted order
				if pos1 > pos2 {
					t.Errorf("diagnostics not sorted: %d > %d", pos1, pos2)
				}
			}
		})
	}
}
