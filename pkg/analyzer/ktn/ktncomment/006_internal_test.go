// Internal tests for 006.go - function documentation analyzer.
package ktncomment

import (
	"go/ast"
	"testing"
)

// Test_extractCommentLines tests the extractCommentLines function.
//
// Params:
//   - t: testing context
func Test_extractCommentLines(t *testing.T) {
	tests := []struct {
		name string
		cg   *ast.CommentGroup
		want int
	}{
		{
			name: "single comment",
			cg: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// Hello"},
				},
			},
			want: 1,
		},
		{
			name: "multiple comments",
			cg: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// Line 1"},
					{Text: "// Line 2"},
					{Text: "// Line 3"},
				},
			},
			want: 3,
		},
		{
			name: "block comment ignored",
			cg: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* block comment */"},
				},
			},
			want: 0,
		},
		{
			name: "mixed comments",
			cg: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// Line comment"},
					{Text: "/* block */"},
				},
			},
			want: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCommentLines(tt.cg)
			// Check result length
			if len(got) != tt.want {
				t.Errorf("extractCommentLines() returned %d lines, want %d", len(got), tt.want)
			}
		})
	}
}

// Test_validateDescriptionLine tests the validateDescriptionLine function.
//
// Params:
//   - t: testing context
func Test_validateDescriptionLine(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		funcName string
		wantErr  bool
	}{
		{
			name:     "empty comments",
			comments: []string{},
			funcName: "myFunc",
			wantErr:  true,
		},
		{
			name:     "valid description",
			comments: []string{"// myFunc does something"},
			funcName: "myFunc",
			wantErr:  false,
		},
		{
			name:     "wrong function name",
			comments: []string{"// otherFunc does something"},
			funcName: "myFunc",
			wantErr:  true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescriptionLine(tt.comments, tt.funcName)
			// Check error presence
			if (err != "") != tt.wantErr {
				t.Errorf("validateDescriptionLine() error = %q, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test_validateParamsSection tests the validateParamsSection function.
//
// Params:
//   - t: testing context
func Test_validateParamsSection(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		startIdx int
		wantErr  bool
		wantIdx  int
	}{
		{
			name:     "missing params header",
			comments: []string{"// Something else"},
			startIdx: 0,
			wantErr:  true,
			wantIdx:  0,
		},
		{
			name:     "valid params section",
			comments: []string{"// Params:", "//   - arg: description"},
			startIdx: 0,
			wantErr:  false,
			wantIdx:  2,
		},
		{
			name:     "params header but no items",
			comments: []string{"// Params:"},
			startIdx: 0,
			wantErr:  true,
			wantIdx:  1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, idx := validateParamsSection(tt.comments, tt.startIdx)
			// Check error presence
			if (err != "") != tt.wantErr {
				t.Errorf("validateParamsSection() error = %q, wantErr %v", err, tt.wantErr)
			}
			// Check index
			if idx != tt.wantIdx {
				t.Errorf("validateParamsSection() idx = %d, want %d", idx, tt.wantIdx)
			}
		})
	}
}

// Test_validateReturnsSection tests the validateReturnsSection function.
//
// Params:
//   - t: testing context
func Test_validateReturnsSection(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		startIdx int
		wantErr  bool
		wantIdx  int
	}{
		{
			name:     "missing returns header",
			comments: []string{"// Something else"},
			startIdx: 0,
			wantErr:  true,
			wantIdx:  0,
		},
		{
			name:     "valid returns section",
			comments: []string{"// Returns:", "//   - error: description"},
			startIdx: 0,
			wantErr:  false,
			wantIdx:  2,
		},
		{
			name:     "returns header but no items",
			comments: []string{"// Returns:"},
			startIdx: 0,
			wantErr:  true,
			wantIdx:  1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, idx := validateReturnsSection(tt.comments, tt.startIdx)
			// Check error presence
			if (err != "") != tt.wantErr {
				t.Errorf("validateReturnsSection() error = %q, wantErr %v", err, tt.wantErr)
			}
			// Check index
			if idx != tt.wantIdx {
				t.Errorf("validateReturnsSection() idx = %d, want %d", idx, tt.wantIdx)
			}
		})
	}
}

// Test_validateDocFormat tests the validateDocFormat function.
//
// Params:
//   - t: testing context
func Test_validateDocFormat(t *testing.T) {
	tests := []struct {
		name       string
		comments   []string
		funcName   string
		hasParams  bool
		hasReturns bool
		wantErr    bool
	}{
		{
			name:       "empty comments",
			comments:   []string{},
			funcName:   "myFunc",
			hasParams:  false,
			hasReturns: false,
			wantErr:    true,
		},
		{
			name:       "simple function no params no returns",
			comments:   []string{"// myFunc does something"},
			funcName:   "myFunc",
			hasParams:  false,
			hasReturns: false,
			wantErr:    false,
		},
		{
			name: "function with params",
			comments: []string{
				"// myFunc does something",
				"//",
				"// Params:",
				"//   - arg: description",
			},
			funcName:   "myFunc",
			hasParams:  true,
			hasReturns: false,
			wantErr:    false,
		},
		{
			name: "function with params and returns",
			comments: []string{
				"// myFunc does something",
				"//",
				"// Params:",
				"//   - arg: description",
				"//",
				"// Returns:",
				"//   - error: error description",
			},
			funcName:   "myFunc",
			hasParams:  true,
			hasReturns: true,
			wantErr:    false,
		},
		{
			name:       "missing params section",
			comments:   []string{"// myFunc does something"},
			funcName:   "myFunc",
			hasParams:  true,
			hasReturns: false,
			wantErr:    true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocFormat(tt.comments, tt.funcName, tt.hasParams, tt.hasReturns)
			// Check error presence
			if (err != "") != tt.wantErr {
				t.Errorf("validateDocFormat() error = %q, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test_runComment006 tests the runComment006 function configuration.
//
// Params:
//   - t: testing context
func Test_runComment006(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment006 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer006 is properly configured
			if Analyzer006 == nil {
				t.Error("Analyzer006 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer006.Name != "ktncomment006" {
				t.Errorf("Analyzer006.Name = %q, want %q", Analyzer006.Name, "ktncomment006")
			}
		})
	}
}
