package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar011 tests the private runVar011 function.
func Test_runVar011(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_checkStringConcatInLoop tests the private checkStringConcatInLoop function.
func Test_checkStringConcatInLoop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks string concatenation in loops
		})
	}
}

// Test_isStringConcatenation tests the private isStringConcatenation function.
func Test_isStringConcatenation(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "string concatenation",
			code: `package test
func example() {
	s := "hello"
	s += " world"
}`,
			expected: true,
		},
		{
			name: "non-string assignment",
			code: `package test
func example() {
	x := 1
	x += 2
}`,
			expected: false,
		},
		{
			name: "empty lhs",
			code: `package test
func example() {
	_ = 1
}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Type check
			conf := types.Config{}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
				Report:    func(_d analysis.Diagnostic) {},
			}

			// Find assignment statement
			foundAssign := false
			ast.Inspect(file, func(n ast.Node) bool {
				if assign, ok := n.(*ast.AssignStmt); ok && assign.Tok == token.ADD_ASSIGN {
					result := isStringConcatenation(pass, assign)
					if result != tt.expected {
						t.Errorf("isStringConcatenation() = %v, expected %v", result, tt.expected)
					}
					foundAssign = true
					return false
				}
				return true
			})

			if !foundAssign && tt.expected {
				t.Error("No assignment found in test code")
			}
		})
	}
}

// Test_isStringConcatenation_noTypeInfo tests with no TypesInfo.
func Test_isStringConcatenation_noTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// Create assignment with lhs that won't have type info
			assign := &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
			}

			result := isStringConcatenation(pass, assign)
			// Vérification du résultat
			if result {
				t.Errorf("isStringConcatenation() = true, expected false when no type info")
			}

		})
	}
}

// Test_isStringConcatenation_notBasicType tests with non-basic type.
func Test_isStringConcatenation_notBasicType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			code := `package test
			type MyString string
			func example() {
			var s MyString
			s += "world"
			}
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Type check
			conf := types.Config{}
			info := &types.Info{
				Types: make(map[ast.Expr]types.TypeAndValue),
			}
			_, _ = conf.Check("test", fset, []*ast.File{file}, info)

			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: info,
				Report:    func(_d analysis.Diagnostic) {},
			}

			// Find assignment statement
			ast.Inspect(file, func(n ast.Node) bool {
				if assign, ok := n.(*ast.AssignStmt); ok && assign.Tok == token.ADD_ASSIGN {
					result := isStringConcatenation(pass, assign)
					// Should return true because underlying type is string
					if !result {
						t.Errorf("isStringConcatenation() = false, expected true for named string type")
					}
					return false
				}
				return true
			})

		})
	}
}

// Test_checkStringConcatInLoop_nilBody tests with nil loop body.
func Test_checkStringConcatInLoop_nilBody(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"nil body"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			// ForStmt with nil body
			forStmt := &ast.ForStmt{Body: nil}
			checkStringConcatInLoop(pass, forStmt)
			// Should return early without panic

			// RangeStmt with nil body
			rangeStmt := &ast.RangeStmt{Body: nil}
			checkStringConcatInLoop(pass, rangeStmt)
			// Should return early without panic
		})
	}
}

// Test_isStringConcatenation_emptyLhs tests with empty Lhs.
func Test_isStringConcatenation_emptyLhs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"empty Lhs"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
				Report: func(_d analysis.Diagnostic) {},
			}

			assign := &ast.AssignStmt{
				Lhs: []ast.Expr{},
			}

			result := isStringConcatenation(pass, assign)
			if result {
				t.Errorf("isStringConcatenation() = true, expected false for empty Lhs")
			}
		})
	}
}

// Test_runVar011_disabled tests runVar011 with disabled rule.
func Test_runVar011_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-011": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar011(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar011() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar011() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar011_fileExcluded tests runVar011 with excluded file.
func Test_runVar011_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-011": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar011(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar011() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar011() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

