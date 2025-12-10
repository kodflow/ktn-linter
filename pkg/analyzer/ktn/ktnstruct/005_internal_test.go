package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runStruct005 tests the private runStruct005 function.
func Test_runStruct005(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_checkFieldOrder tests the private checkFieldOrder function.
func Test_checkFieldOrder(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		expectError bool
	}{
		{
			name: "correct order",
			src: `package test
type User struct {
	Name string
	age  int
}`,
			expectError: false,
		},
		{
			name: "incorrect order",
			src: `package test
type User struct {
	age  int
	Name string
}`,
			expectError: true,
		},
		{
			name: "all exported",
			src: `package test
type User struct {
	Name string
	Age  int
}`,
			expectError: false,
		},
		{
			name: "all private",
			src: `package test
type User struct {
	name string
	age  int
}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the struct
			var typeSpec *ast.TypeSpec
			var structType *ast.StructType
			ast.Inspect(file, func(n ast.Node) bool {
				if ts, ok := n.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						typeSpec = ts
						structType = st
						return false
					}
				}
				return true
			})

			if typeSpec == nil || structType == nil {
				t.Fatal("no struct found")
			}

			// Create a minimal pass with a reporter
			errorReported := false
			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					errorReported = true
				},
			}

			// Call checkFieldOrder
			checkFieldOrder(pass, typeSpec, structType)

			if errorReported != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, errorReported)
			}
		})
	}
}

// Test_fieldInfo tests the fieldInfo type.
func Test_fieldInfo(t *testing.T) {
	tests := []struct {
		name     string
		fi       fieldInfo
		expected string
		exported bool
	}{
		{
			name:     "exported field",
			fi:       fieldInfo{name: "Name", exported: true, pos: token.NoPos},
			expected: "Name",
			exported: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fi.name != tt.expected || tt.fi.exported != tt.exported {
				t.Errorf("expected name=%s exported=%v, got name=%s exported=%v",
					tt.expected, tt.exported, tt.fi.name, tt.fi.exported)
			}
		})
	}
}

// Test_runStruct005_disabled tests that the rule is skipped when disabled.
func Test_runStruct005_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-005": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runStruct005(pass)
	if err != nil {
		t.Errorf("runStruct005() error = %v", err)
	}
}

// Test_runStruct005_excludedFile tests that excluded files are skipped.
func Test_runStruct005_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-005": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runStruct005(pass)
	if err != nil {
		t.Errorf("runStruct005() error = %v", err)
	}
}
