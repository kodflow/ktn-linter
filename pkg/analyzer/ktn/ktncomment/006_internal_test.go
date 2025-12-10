// Internal tests for 006.go - function documentation analyzer.
package ktncomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
		{
			name:     "params section with blank line after",
			comments: []string{"// Params:", "//   - arg: description", "//"},
			startIdx: 0,
			wantErr:  false,
			wantIdx:  3,
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
		{
			name: "function with only returns",
			comments: []string{
				"// myFunc does something",
				"//",
				"// Returns:",
				"//   - error: error description",
			},
			funcName:   "myFunc",
			hasParams:  false,
			hasReturns: true,
			wantErr:    false,
		},
		{
			name: "missing returns section",
			comments: []string{
				"// myFunc does something",
				"//",
			},
			funcName:   "myFunc",
			hasParams:  false,
			hasReturns: true,
			wantErr:    true,
		},
		{
			name: "multiline description before params",
			comments: []string{
				"// myFunc does something",
				"// with multiple lines",
				"// of description",
				"//",
				"// Params:",
				"//   - arg: description",
			},
			funcName:   "myFunc",
			hasParams:  true,
			hasReturns: false,
			wantErr:    false,
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

// Test_runComment006_ruleDisabled tests behavior when rule is disabled.
//
// Params:
//   - t: testing context
func Test_runComment006_ruleDisabled(t *testing.T) {
	// Import config package for test
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-COMMENT-006": {Enabled: config.Bool(false)},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	code := `package test
func myFunc() {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report no errors when rule disabled
	if errorCount != 0 {
		t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
	}
}

// Test_runComment006_fileExcluded tests behavior when file is excluded.
//
// Params:
//   - t: testing context
func Test_runComment006_fileExcluded(t *testing.T) {
	// Import config package for test
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-COMMENT-006": {
				Enabled: config.Bool(true),
				Exclude: []string{"*.go"},
			},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	code := `package test
func myFunc() {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report no errors when file excluded
	if errorCount != 0 {
		t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
	}
}

// Test_validateDocFormat_skipDescriptionLines tests skipping description lines.
//
// Params:
//   - t: testing context
func Test_validateDocFormat_skipDescriptionLines(t *testing.T) {
	tests := []struct {
		name       string
		comments   []string
		funcName   string
		hasParams  bool
		hasReturns bool
		wantErr    bool
	}{
		{
			name: "direct params without blank line",
			comments: []string{
				"// myFunc does something",
				"// Params:",
				"//   - arg: description",
			},
			funcName:   "myFunc",
			hasParams:  true,
			hasReturns: false,
			wantErr:    false,
		},
		{
			name: "direct returns without blank line",
			comments: []string{
				"// myFunc does something",
				"// Returns:",
				"//   - error: error description",
			},
			funcName:   "myFunc",
			hasParams:  false,
			hasReturns: true,
			wantErr:    false,
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

// Test_runComment006_testFile tests runComment006 with test files.
//
// Params:
//   - t: testing context
func Test_runComment006_testFile(t *testing.T) {
	cfg := config.DefaultConfig()
	config.Set(cfg)
	defer config.Reset()

	code := `package test
func TestExample() {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example_test.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report no errors for test files
	if errorCount != 0 {
		t.Errorf("expected 0 errors for test files, got %d", errorCount)
	}
}

// Test_runComment006_testFunction tests runComment006 with test functions.
//
// Params:
//   - t: testing context
func Test_runComment006_testFunction(t *testing.T) {
	cfg := config.DefaultConfig()
	config.Set(cfg)
	defer config.Reset()

	code := `package example
func BenchmarkExample() {}
func ExampleUsage() {}
func FuzzTarget() {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report no errors for test/benchmark/example/fuzz functions
	if errorCount != 0 {
		t.Errorf("expected 0 errors for test functions, got %d", errorCount)
	}
}

// Test_runComment006_missingDoc tests runComment006 with missing documentation.
//
// Params:
//   - t: testing context
func Test_runComment006_missingDoc(t *testing.T) {
	cfg := config.DefaultConfig()
	config.Set(cfg)
	defer config.Reset()

	code := `package example
func myFunc() {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report error for missing documentation
	if errorCount == 0 {
		t.Error("expected error for missing documentation, got 0")
	}
}

// Test_runComment006_invalidDocFormat tests runComment006 with invalid doc format.
//
// Params:
//   - t: testing context
func Test_runComment006_invalidDocFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	config.Set(cfg)
	defer config.Reset()

	code := `package example
// myFunc does something
// but is missing the required format
func myFunc(x int) int { return x }`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", code, parser.ParseComments)
	// Check parsing success
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	pass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Run inspect analyzer
	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{file},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)
	pass.ResultOf[inspect.Analyzer] = inspectResult

	errorCount := 0
	pass.Report = func(d analysis.Diagnostic) {
		errorCount++
	}

	// Run analyzer
	_, err = runComment006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runComment006 failed: %v", err)
	}

	// Should report error for invalid doc format
	if errorCount == 0 {
		t.Error("expected error for invalid documentation format, got 0")
	}
}
