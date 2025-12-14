// Internal tests for 005.go - struct documentation analyzer.
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

// Test_hasValidDocumentation tests the hasValidDocumentation function.
func Test_hasValidDocumentation(t *testing.T) {
	tests := []struct {
		name       string
		doc        *ast.CommentGroup
		structName string
		want       bool
	}{
		{
			name:       "nil documentation",
			doc:        nil,
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "empty comment list",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "single line doc",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct does something"},
				},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "valid two line doc",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct represents a structure."},
					{Text: "// It provides functionality for testing."},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
		{
			name: "doc not starting with struct name",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// This struct does something."},
					{Text: "// More description here."},
				},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "doc with empty lines",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct is a test structure."},
					{Text: "//"},
					{Text: "// More info here."},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasValidDocumentation(tt.doc, tt.structName, defaultMinStructDocLines)
			// Check result
			if got != tt.want {
				t.Errorf("hasValidDocumentation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runComment005 tests the runComment005 function configuration.
func Test_runComment005(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment005 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer005 is properly configured
			if Analyzer005 == nil {
				t.Error("Analyzer005 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer005.Name != "ktncomment005" {
				t.Errorf("Analyzer005.Name = %q, want %q", Analyzer005.Name, "ktncomment005")
			}
		})
	}
}

// Test_runComment005_ruleDisabled tests behavior when rule is disabled.
func Test_runComment005_ruleDisabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-005": {Enabled: config.Bool(false)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			type MyStruct struct {}`

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
			_, err = runComment005(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment005 failed: %v", err)
			}

			// Should report no errors when rule disabled
			if errorCount != 0 {
				t.Errorf("expected 0 errors when rule disabled, got %d", errorCount)
			}

		})
	}
}

// Test_runComment005_fileExcluded tests behavior when file is excluded.
func Test_runComment005_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Import config package for test
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-005": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			type MyStruct struct {}`

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
			_, err = runComment005(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment005 failed: %v", err)
			}

			// Should report no errors when file excluded
			if errorCount != 0 {
				t.Errorf("expected 0 errors when file excluded, got %d", errorCount)
			}

		})
	}
}

// Test_runComment005_customThreshold tests behavior with custom threshold.
func Test_runComment005_customThreshold(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Configure custom threshold
			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-COMMENT-005": {
						Enabled:   config.Bool(true),
						Threshold: config.Int(3),
					},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			code := `package test
			// MyStruct is a test.
			// Second line.
			type MyStruct struct {}`

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
			_, err = runComment005(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment005 failed: %v", err)
			}

			// Should report error (only 2 lines but needs 3)
			if errorCount != 1 {
				t.Errorf("expected 1 error with custom threshold, got %d", errorCount)
			}

		})
	}
}

// Test_hasValidDocumentation_blockComments tests block-style comments.
func Test_hasValidDocumentation_blockComments(t *testing.T) {
	tests := []struct {
		name       string
		doc        *ast.CommentGroup
		structName string
		want       bool
	}{
		{
			name: "block comment with struct name",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* MyStruct represents a structure. */"},
					{Text: "/* Additional info here. */"},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
		{
			name: "mixed comment styles",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* MyStruct represents a structure. */"},
					{Text: "// More description here."},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasValidDocumentation(tt.doc, tt.structName, defaultMinStructDocLines)
			// Check result
			if got != tt.want {
				t.Errorf("hasValidDocumentation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runComment005_privateStruct tests that private structs are skipped.
func Test_runComment005_privateStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			type privateStruct struct {}`

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
			_, err = runComment005(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment005 failed: %v", err)
			}

			// Should not report error for private struct
			if errorCount != 0 {
				t.Errorf("expected 0 errors for private struct, got %d", errorCount)
			}

		})
	}
}

// Test_runComment005_nonStruct tests that non-struct types are skipped.
func Test_runComment005_nonStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			type MyInt int`

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
			_, err = runComment005(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runComment005 failed: %v", err)
			}

			// Should not report error for non-struct type
			if errorCount != 0 {
				t.Errorf("expected 0 errors for non-struct type, got %d", errorCount)
			}

		})
	}
}
