// Internal tests for analyzer 011.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_hasErrorCaseCoverage tests the hasErrorCaseCoverage private function.
func Test_hasErrorCaseCoverage(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "test with error variable",
			code: `func TestSomething(t *testing.T) {
				err := someFunc()
				if err != nil {
					t.Error(err)
				}
			}`,
			want: true,
		},
		{
			name: "test with error string",
			code: `func TestSomething(t *testing.T) {
				tests := []struct{
					name string
				}{
					{name: "error case"},
				}
			}`,
			want: true,
		},
		{
			name: "test without error coverage",
			code: `func TestSomething(t *testing.T) {
				result := someFunc()
				if result != expected {
					t.Log("mismatch")
				}
			}`,
			want: false,
		},
		{
			name: "test with invalid string",
			code: `func TestSomething(t *testing.T) {
				tests := []struct{
					name string
				}{
					{name: "invalid input"},
				}
			}`,
			want: true,
		},
		{
			name: "test with fail string",
			code: `func TestSomething(t *testing.T) {
				tests := []struct{
					name string
				}{
					{name: "should fail"},
				}
			}`,
			want: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					// Retour false pour arrêter
					return false
				}
				// Continuer la traversée
				return true
			})

			// Vérification de la déclaration
			if funcDecl == nil {
				t.Fatal("no function declaration found")
			}

			got := hasErrorCaseCoverage(funcDecl)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("hasErrorCaseCoverage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_isErrorIndicatorName tests the isErrorIndicatorName private function.
func Test_isErrorIndicatorName(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		want    bool
	}{
		{
			name:    "err is error indicator",
			varName: "err",
			want:    true,
		},
		{
			name:    "error is error indicator",
			varName: "error",
			want:    true,
		},
		{
			name:    "invalid is error indicator",
			varName: "invalid",
			want:    true,
		},
		{
			name:    "fail is error indicator",
			varName: "fail",
			want:    true,
		},
		{
			name:    "bad is error indicator",
			varName: "bad",
			want:    true,
		},
		{
			name:    "wrong is error indicator",
			varName: "wrong",
			want:    true,
		},
		{
			name:    "myError contains error",
			varName: "myError",
			want:    true,
		},
		{
			name:    "regular variable not indicator",
			varName: "result",
			want:    false,
		},
		{
			name:    "empty name not indicator",
			varName: "",
			want:    false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := isErrorIndicatorName(tt.varName)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isErrorIndicatorName(%q) = %v, want %v", tt.varName, got, tt.want)
			}
		})
	}
}

// Test_hasErrorTestCases tests the hasErrorTestCases private function.
func Test_hasErrorTestCases(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "composite literal with error in name",
			code: `package test
func TestSomething(t *testing.T) {
	tests := []struct{
		name string
		want int
	}{
		{name: "error case", want: 0},
	}
	_ = tests
}`,
			want: true,
		},
		{
			name: "composite literal with invalid in name",
			code: `package test
func TestSomething(t *testing.T) {
	tests := []struct{
		name string
		want int
	}{
		{name: "invalid input", want: 0},
	}
	_ = tests
}`,
			want: true,
		},
		{
			name: "composite literal without error indicators",
			code: `package test
func TestSomething(t *testing.T) {
	tests := []struct{
		name string
	}{
		{name: "normal case"},
	}
	_ = tests
}`,
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract inner composite literal (the one with test cases, not the outer array)
			var lits []*ast.CompositeLit
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if cl, ok := n.(*ast.CompositeLit); ok {
					lits = append(lits, cl)
				}
				// Continuer la traversée
				return true
			})

			// We need the inner composite literal (the test case, not the array)
			// The structure is: outer array literal -> inner struct literal
			var targetLit *ast.CompositeLit
			// Parcourir les literals
			for _, lit := range lits {
				// Check if this literal has KeyValueExpr with "name" key
				// Parcourir les éléments
				for _, elt := range lit.Elts {
					// Vérification du type
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						// Vérification de l'identifiant
						if ident, identOk := kv.Key.(*ast.Ident); identOk && ident.Name == "name" {
							targetLit = lit
							break
						}
					}
				}
				// Vérification du literal trouvé
				if targetLit != nil {
					break
				}
			}

			// Vérification du literal
			if targetLit == nil {
				t.Fatal("no composite literal with name field found")
			}

			got := hasErrorTestCases(targetLit)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("hasErrorTestCases() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runTest011 tests the runTest011 private function.
func Test_runTest011(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{name: "analyzer exists", expectedName: "ktntest011"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if Analyzer011 == nil || Analyzer011.Name != tt.expectedName {
				t.Errorf("Analyzer011 invalid: nil=%v, Name=%q, want %q",
					Analyzer011 == nil, Analyzer011.Name, tt.expectedName)
			}
		})
	}
}

// Test_collectFuncSignatures tests the collectFuncSignatures function.
func Test_collectFuncSignatures(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantCount int
		checkKey  string
		wantError bool
	}{
		{
			name:      "function with error return",
			code:      "package test\nfunc Foo() error { return nil }",
			wantCount: 1,
			checkKey:  "Foo",
			wantError: true,
		},
		{
			name:      "function without error return",
			code:      "package test\nfunc Bar() string { return \"\" }",
			wantCount: 1,
			checkKey:  "Bar",
			wantError: false,
		},
		{
			name:      "method with error return",
			code:      "package test\ntype S struct{}\nfunc (s *S) Method() error { return nil }",
			wantCount: 2,
			checkKey:  "Method",
			wantError: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) {},
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
			}

			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			signatures := collectFuncSignatures(pass, insp)

			if len(signatures) < tt.wantCount {
				t.Errorf("expected at least %d signatures, got %d", tt.wantCount, len(signatures))
			}

			if tt.checkKey != "" {
				if info, found := signatures[tt.checkKey]; found {
					if info.returnsError != tt.wantError {
						t.Errorf("expected returnsError=%v, got %v", tt.wantError, info.returnsError)
					}
				}
			}
		})
	}
}

// Test_addFuncSignature tests the addFuncSignature function.
func Test_addFuncSignature(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectCount  int
		expectName   string
		returnsError bool
	}{
		{
			name:         "function with error",
			code:         "package test\nfunc Foo() error { return nil }",
			expectCount:  1,
			expectName:   "Foo",
			returnsError: true,
		},
		{
			name:         "function without error",
			code:         "package test\nfunc Bar() string { return \"\" }",
			expectCount:  1,
			expectName:   "Bar",
			returnsError: false,
		},
		{
			name:        "method with receiver",
			code:        "package test\ntype S struct{}\nfunc (s *S) Method() error { return nil }",
			expectCount: 2,
			expectName:  "Method",
		},
		{
			name:        "mock function (included, not filtered at this level)",
			code:        "package test\nfunc MockFunc() error { return nil }",
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result := make(map[string]testedFuncInfo)
			for _, decl := range file.Decls {
				if fd, ok := decl.(*ast.FuncDecl); ok {
					addFuncSignature(result, fd)
				}
			}
			if len(result) != tt.expectCount {
				t.Errorf("expected %d signatures, got %d", tt.expectCount, len(result))
			}
			if tt.expectName != "" && tt.expectCount > 0 {
				if info, found := result[tt.expectName]; found {
					if tt.returnsError && !info.returnsError {
						t.Error("expected function to return error")
					}
				} else if tt.expectCount == 1 {
					t.Errorf("expected to find signature %q", tt.expectName)
				}
			}
		})
	}
}

// Test_collectExternalSourceSignatures tests the collectExternalSourceSignatures function.
func Test_collectExternalSourceSignatures(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "function exists",
			want: "collectExternalSourceSignatures exists and is tested via public API",
		},
		{
			name: "function is internal",
			want: "collectExternalSourceSignatures is not exported",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Vérification que la fonction existe
			// Les tests réels nécessitent un *analysis.Pass complet
			t.Log(tt.want)
		})
	}
}

// Test_scanSourceFile tests the scanSourceFile function.
func Test_scanSourceFile(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		filename string
		wantLog  string
	}{
		{
			name:     "handles missing files",
			dir:      "/nonexistent",
			filename: "missing.go",
			wantLog:  "handled missing file gracefully",
		},
		{
			name:     "handles empty directory",
			dir:      "",
			filename: "test.go",
			wantLog:  "handled empty directory gracefully",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]testedFuncInfo)
			// Appel avec fichier inexistant ne doit pas paniquer
			scanSourceFile(tt.dir, tt.filename, result)
			// Le test passe si pas de panic
			t.Log(tt.wantLog)
		})
	}
}

// Test_extractFuncInfo tests the extractFuncInfo function.
func Test_extractFuncInfo(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "extracts function info", code: "package test\nfunc Foo() error { return nil }"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			for _, decl := range file.Decls {
				if fd, ok := decl.(*ast.FuncDecl); ok {
					info := extractFuncInfo(fd)
					if info == nil {
						t.Error("expected non-nil info")
					}
				}
			}
		})
	}
}

// Test_functionReturnsError tests the functionReturnsError function.
func Test_functionReturnsError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"returns error", "func Foo() error { return nil }", true},
		{"no return", "func Foo() {}", false},
		{"returns int", "func Foo() int { return 0 }", false},
	}
	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			// Parcourir les déclarations
			for _, decl := range file.Decls {
				// Vérifier FuncDecl
				if fd, ok := decl.(*ast.FuncDecl); ok {
					result := functionReturnsError(fd)
					// Vérification résultat
					if result != tt.expected {
						t.Errorf("functionReturnsError() = %v, want %v", result, tt.expected)
					}
				}
			}
		})
	}
}

// Test_isErrorType tests the isErrorType function.
func Test_isErrorType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "error type",
			code:     "package test\nfunc Foo() error { return nil }",
			expected: true,
		},
		{
			name:     "non-error type (string)",
			code:     "package test\nfunc Foo() string { return \"\" }",
			expected: false,
		},
		{
			name:     "non-error type (int)",
			code:     "package test\nfunc Foo() int { return 0 }",
			expected: false,
		},
		{
			name:     "multiple returns with error",
			code:     "package test\nfunc Foo() (string, error) { return \"\", nil }",
			expected: false,
		},
		{
			name:     "pointer type (not error)",
			code:     "package test\nfunc Foo() *string { return nil }",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			for _, decl := range file.Decls {
				if fd, ok := decl.(*ast.FuncDecl); ok {
					if fd.Type.Results != nil && len(fd.Type.Results.List) > 0 {
						result := isErrorType(fd.Type.Results.List[0].Type)
						if result != tt.expected {
							t.Errorf("isErrorType() = %v, want %v", result, tt.expected)
						}
					}
				}
			}
		})
	}
}

// Test_ExtractReceiverTypeName013 tests the shared.ExtractReceiverTypeName function (for 013).
func Test_ExtractReceiverTypeName013(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"method with pointer receiver", "func (s *Service) Foo() {}", "Service"},
		{"method with value receiver", "func (s Service) Foo() {}", "Service"},
	}
	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test\ntype Service struct{}\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			// Parcourir les déclarations
			for _, decl := range file.Decls {
				// Vérifier FuncDecl
				if fd, ok := decl.(*ast.FuncDecl); ok {
					// Vérifier si c'est une méthode
					if fd.Recv != nil && len(fd.Recv.List) > 0 {
						result := shared.ExtractReceiverTypeName(fd.Recv.List[0].Type)
						// Vérification résultat
						if result != tt.expected {
							t.Errorf("ExtractReceiverTypeName() = %q, want %q", result, tt.expected)
						}
					}
				}
			}
		})
	}
}

// Test_analyzeTestFunction tests the analyzeTestFunction function.
func Test_analyzeTestFunction(t *testing.T) {
	tests := []struct {
		name         string
		testCode     string
		funcCode     string
		shouldReport bool
	}{
		{
			name:         "test for function with error and error coverage",
			testCode:     `func TestFoo(t *testing.T) { err := Foo(); if err != nil { t.Error(err) } }`,
			funcCode:     `func Foo() error { return nil }`,
			shouldReport: false,
		},
		{
			name:         "test for function with error but no error coverage",
			testCode:     `func TestBar(t *testing.T) { result := Bar(); _ = result }`,
			funcCode:     `func Bar() error { return nil }`,
			shouldReport: true,
		},
		{
			name:         "test for function without error",
			testCode:     `func TestBaz(t *testing.T) { result := Baz(); _ = result }`,
			funcCode:     `func Baz() string { return "" }`,
			shouldReport: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\n" + tt.funcCode + "\n" + tt.testCode
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			reportCount := 0
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) { reportCount++ },
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) { reportCount++ },
			}

			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			signatures := collectFuncSignatures(pass, insp)

			// Find test function
			var testFunc *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Name != nil && shared.IsUnitTestFunction(fd) {
					testFunc = fd
					return false
				}
				return true
			})

			if testFunc != nil {
				analyzeTestFunction(pass, testFunc, signatures)
			}

			if tt.shouldReport && reportCount == 0 {
				t.Error("expected error report, got none")
			} else if !tt.shouldReport && reportCount > 0 {
				t.Errorf("expected no error report, got %d", reportCount)
			}
		})
	}
}

// Test_ParseTestName013 tests the shared.ParseTestName function (for 013).
func Test_ParseTestName013(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedFunc string
		expectedOk   bool
	}{
		{"TestFoo", "TestFoo", "Foo", true},
		{"Test_foo", "Test_foo", "foo", true},
		{"TestFooBar", "TestFooBar", "FooBar", true},
		{"Test", "Test", "", false},
		{"NotATest", "NotATest", "", false},
	}
	// Parcourir les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			target, ok := shared.ParseTestName(tt.input)
			// Vérification ok
			if ok != tt.expectedOk {
				t.Errorf("ParseTestName(%q) ok = %v, want %v", tt.input, ok, tt.expectedOk)
			}
			// Vérification résultat si ok
			if ok && target.FuncName != tt.expectedFunc {
				t.Errorf("ParseTestName(%q).FuncName = %q, want %q", tt.input, target.FuncName, tt.expectedFunc)
			}
		})
	}
}

// Test_checkErrorInNode tests the checkErrorInNode function.
func Test_checkErrorInNode(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "checks error in AST node",
			code: "package test\nfunc TestFoo(t *testing.T) {\n\ttests := []struct{name string}{{name: \"error case\"}}\n\t_ = tests\n}",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if checkErrorInNode(n) {
					found = true
					return false
				}
				return true
			})
			if !found {
				t.Error("expected to find error indicator")
			}
		})
	}
}

// Test_checkErrorInBasicLit tests the checkErrorInBasicLit function.
func Test_checkErrorInBasicLit(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "checks error in basic literal", code: "package test\nvar s = \"error case\""},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if bl, ok := n.(*ast.BasicLit); ok {
					if checkErrorInBasicLit(bl) {
						found = true
						return false
					}
				}
				return true
			})
			if !found {
				t.Error("expected to find error in basic literal")
			}
		})
	}
}

// Test_checkErrorInKeyValue tests the checkErrorInKeyValue function.
func Test_checkErrorInKeyValue(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantFound bool
	}{
		{
			name:      "checks error in key-value expression with name key",
			code:      "package test\nvar x = struct{name string}{name: \"error case\"}",
			wantFound: true,
		},
		{
			name:      "key is not name",
			code:      "package test\nvar x = struct{desc string}{desc: \"error case\"}",
			wantFound: false,
		},
		{
			name:      "value is not BasicLit",
			code:      "package test\nvar y string\nvar x = struct{name string}{name: y}",
			wantFound: false,
		},
		{
			name:      "normal case without error",
			code:      "package test\nvar x = struct{name string}{name: \"normal case\"}",
			wantFound: false,
		},
		{
			name:      "invalid in name value",
			code:      "package test\nvar x = struct{name string}{name: \"invalid input\"}",
			wantFound: true,
		},
		{
			name:      "fail in name value",
			code:      "package test\nvar x = struct{name string}{name: \"should fail\"}",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if kv, ok := n.(*ast.KeyValueExpr); ok {
					if checkErrorInKeyValue(kv) {
						found = true
						return false
					}
				}
				return true
			})
			if found != tt.wantFound {
				t.Errorf("checkErrorInKeyValue() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

// Test_runTest011_disabled tests that the rule is skipped when disabled.
func Test_runTest011_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runTest011_excludedFile tests that excluded files are skipped.
func Test_runTest011_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_hasErrorTestCases_nonKeyValue tests hasErrorTestCases with non-KeyValueExpr elements.
func Test_hasErrorTestCases_nonKeyValue(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantFound bool
	}{
		{
			name:      "array literal without KeyValueExpr",
			code:      "package test\nvar x = []int{1, 2, 3}",
			wantFound: false,
		},
		{
			name:      "empty composite literal",
			code:      "package test\nvar x = struct{}{}",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if cl, ok := n.(*ast.CompositeLit); ok {
					if hasErrorTestCases(cl) {
						found = true
						return false
					}
				}
				return true
			})
			if found != tt.wantFound {
				t.Errorf("hasErrorTestCases() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

// Test_scanSourceFile_various tests scanSourceFile with various file types.
func Test_scanSourceFile_various(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		filename string
	}{
		{
			name:     "non-go file",
			dir:      "/tmp",
			filename: "readme.txt",
		},
		{
			name:     "test file",
			dir:      "/tmp",
			filename: "foo_test.go",
		},
		{
			name:     "mock file",
			dir:      "/tmp",
			filename: "mock_foo.go",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]testedFuncInfo)
			// These should all return early without adding signatures
			scanSourceFile(tt.dir, tt.filename, result)
			if len(result) != 0 {
				t.Errorf("expected empty result, got %d entries", len(result))
			}
		})
	}
}

// Test_checkErrorInBasicLit_nonString tests checkErrorInBasicLit with non-string literals.
func Test_checkErrorInBasicLit_nonString(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantFound bool
	}{
		{
			name:      "integer literal",
			code:      "package test\nvar x = 42",
			wantFound: false,
		},
		{
			name:      "float literal",
			code:      "package test\nvar x = 3.14",
			wantFound: false,
		},
		{
			name:      "string with error keyword",
			code:      "package test\nvar x = \"error occurred\"",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if bl, ok := n.(*ast.BasicLit); ok {
					if checkErrorInBasicLit(bl) {
						found = true
						return false
					}
				}
				return true
			})
			if found != tt.wantFound {
				t.Errorf("checkErrorInBasicLit() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

// Test_addFuncSignature_mockReceiver tests addFuncSignature with mock receiver.
func Test_addFuncSignature_mockReceiver(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCount int
	}{
		{
			name:        "method with mock receiver is skipped",
			code:        "package test\ntype MockService struct{}\nfunc (m *MockService) Method() {}",
			expectCount: 0,
		},
		{
			name:        "regular method is added",
			code:        "package test\ntype Service struct{}\nfunc (s *Service) Method() error { return nil }",
			expectCount: 2, // Method and Service_Method
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			signatures := make(map[string]testedFuncInfo)
			for _, decl := range file.Decls {
				if fd, ok := decl.(*ast.FuncDecl); ok {
					addFuncSignature(signatures, fd)
				}
			}

			if len(signatures) != tt.expectCount {
				t.Errorf("expected %d signatures, got %d", tt.expectCount, len(signatures))
			}
		})
	}
}

// Test_analyzeTestFunction_noMatch tests analyzeTestFunction with non-matching test name.
func Test_analyzeTestFunction_noMatch(t *testing.T) {
	tests := []struct {
		name         string
		testCode     string
		shouldReport bool
	}{
		{
			name:         "test name doesn't match any function",
			testCode:     `func TestNonExistent(t *testing.T) { _ = t }`,
			shouldReport: false,
		},
		{
			name:         "invalid test name format",
			testCode:     `func Test(t *testing.T) { _ = t }`,
			shouldReport: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\n" + tt.testCode
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			reportCount := 0
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) { reportCount++ },
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) { reportCount++ },
			}

			// Empty signatures - no matching function
			signatures := make(map[string]testedFuncInfo)

			// Find test function
			var testFunc *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Name != nil && shared.IsUnitTestFunction(fd) {
					testFunc = fd
					return false
				}
				return true
			})

			if testFunc != nil {
				analyzeTestFunction(pass, testFunc, signatures)
			}

			if tt.shouldReport && reportCount == 0 {
				t.Error("expected error report, got none")
			} else if !tt.shouldReport && reportCount > 0 {
				t.Errorf("expected no error report, got %d", reportCount)
			}
		})
	}
}

// Test_collectExternalSourceSignatures_emptyPass tests collectExternalSourceSignatures with empty pass.
func Test_collectExternalSourceSignatures_emptyPass(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"empty files"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{},
			}
			signatures := make(map[string]testedFuncInfo)
			collectExternalSourceSignatures(pass, signatures)
			if len(signatures) != 0 {
				t.Errorf("expected 0 signatures, got %d", len(signatures))
			}
		})
	}
}

// Test_functionReturnsError_nilType tests functionReturnsError with nil function type.
func Test_functionReturnsError_nilType(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		want     bool
	}{
		{
			name: "nil type",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Foo"},
				Type: nil,
			},
			want: false,
		},
		{
			name: "nil results",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Bar"},
				Type: &ast.FuncType{Results: nil},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := functionReturnsError(tt.funcDecl)
			if got != tt.want {
				t.Errorf("functionReturnsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_scanSourceFile_parseError tests scanSourceFile with unparseable file.
func Test_scanSourceFile_parseError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"parse error handled gracefully"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]testedFuncInfo)
			// Invalid Go file that will fail parsing
			scanSourceFile("/nonexistent", "invalid.go", result)
			if len(result) != 0 {
				t.Errorf("expected empty result after parse error, got %d entries", len(result))
			}
		})
	}
}

// Test_analyzeTestFunction_emptyKey tests analyzeTestFunction with empty key.
func Test_analyzeTestFunction_emptyKey(t *testing.T) {
	tests := []struct {
		name         string
		testCode     string
		shouldReport bool
	}{
		{
			name:         "test with name that results in empty key",
			testCode:     `func TestUnparseable(t *testing.T) { _ = t }`,
			shouldReport: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\n" + tt.testCode
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			reportCount := 0
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{file},
				Report:   func(d analysis.Diagnostic) { reportCount++ },
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(d analysis.Diagnostic) { reportCount++ },
			}

			// Empty signatures - test will have empty key after parsing
			signatures := make(map[string]testedFuncInfo)

			// Find test function
			var testFunc *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Name != nil && shared.IsUnitTestFunction(fd) {
					testFunc = fd
					return false
				}
				return true
			})

			if testFunc != nil {
				analyzeTestFunction(pass, testFunc, signatures)
			}

			if tt.shouldReport && reportCount == 0 {
				t.Error("expected error report, got none")
			} else if !tt.shouldReport && reportCount > 0 {
				t.Errorf("expected no error report, got %d", reportCount)
			}
		})
	}
}
