package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestCheckFuncParamsNilParams teste le cas où funcDecl.Type.Params est nil.
//
// Params:
//   - t: instance de test
func TestCheckFuncParamsNilParams(t *testing.T) {
	code := `package test
// externalFunc est définie ailleurs.
func externalFunc()`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	hasError := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			hasError = true
		},
	}

	_, _ = analyzer.FuncAnalyzer.Run(pass)

	if hasError {
		t.Error("Expected no error for function without params")
	}
}

// TestCheckGodocQualityEdgeCases teste les cas limites de checkGodocQuality.
//
// Params:
//   - t: instance de test
func TestCheckGodocQualityEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
		errorCode   string
	}{
		{
			name: "Godoc commence par le nom de la fonction (OK)",
			code: `package test
// simpleFunc fait quelque chose.
func simpleFunc() {
	x := 1
	x++
}`,
			expectError: false,
		},
		{
			name: "Godoc ne commence PAS par le nom de la fonction (KTN-FUNC-002)",
			code: `package test
// Cette fonction fait quelque chose.
func badFunc() {
	x := 1
	x++
}`,
			expectError: true,
			errorCode:   "KTN-FUNC-002",
		},
		{
			name: "Fonction sans params ni returns (OK)",
			code: `package test
// simpleFunc fait quelque chose.
func simpleFunc() {
	x := 1
	x++
}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			hasError := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(diag analysis.Diagnostic) {
					if tt.errorCode == "" || strings.Contains(diag.Message, tt.errorCode) {
						hasError = true
					}
				},
			}

			analyzer.FuncAnalyzer.Run(pass)

			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, hasError)
			}
		})
	}
}
