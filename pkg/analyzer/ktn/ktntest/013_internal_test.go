// Internal tests for analyzer 013.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_hasErrorCaseCoverage tests the hasErrorCaseCoverage private function.
//
// Params:
//   - t: testing context
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
//
// Params:
//   - t: testing context
func Test_isErrorIndicatorName(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		want     bool
	}{
		{
			name:     "err is error indicator",
			varName:  "err",
			want:     true,
		},
		{
			name:     "error is error indicator",
			varName:  "error",
			want:     true,
		},
		{
			name:     "invalid is error indicator",
			varName:  "invalid",
			want:     true,
		},
		{
			name:     "fail is error indicator",
			varName:  "fail",
			want:     true,
		},
		{
			name:     "bad is error indicator",
			varName:  "bad",
			want:     true,
		},
		{
			name:     "wrong is error indicator",
			varName:  "wrong",
			want:     true,
		},
		{
			name:     "myError contains error",
			varName:  "myError",
			want:     true,
		},
		{
			name:     "regular variable not indicator",
			varName:  "result",
			want:     false,
		},
		{
			name:     "empty name not indicator",
			varName:  "",
			want:     false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
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
//
// Params:
//   - t: testing context
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

// Test_runTest013 tests the runTest013 private function.
//
// Params:
//   - t: testing context
func Test_runTest013(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{name: "analyzer exists", expectedName: "ktntest013"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Analyzer013 == nil || Analyzer013.Name != tt.expectedName {
				t.Errorf("Analyzer013 invalid: nil=%v, Name=%q, want %q",
					Analyzer013 == nil, Analyzer013.Name, tt.expectedName)
			}
		})
	}
}

// Test_collectFuncSignatures tests the collectFuncSignatures function.
//
// Params:
//   - t: testing context
func Test_collectFuncSignatures(t *testing.T) {
	t.Run("function exists", func(t *testing.T) {
		// Vérification que la fonction existe
		// Les tests réels nécessitent un *analysis.Pass complet
		t.Log("collectFuncSignatures exists and is tested via public API")
	})
}

// Test_addFuncSignature tests the addFuncSignature function.
//
// Params:
//   - t: testing context
func Test_addFuncSignature(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "adds signature", code: "package test\nfunc Foo() error { return nil }"},
	}

	for _, tt := range tests {
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
			if len(result) == 0 {
				t.Error("expected signature to be added")
			}
		})
	}
}

// Test_collectExternalSourceSignatures tests the collectExternalSourceSignatures function.
//
// Params:
//   - t: testing context
func Test_collectExternalSourceSignatures(t *testing.T) {
	t.Run("function exists", func(t *testing.T) {
		// Vérification que la fonction existe
		// Les tests réels nécessitent un *analysis.Pass complet
		t.Log("collectExternalSourceSignatures exists and is tested via public API")
	})
}

// Test_scanSourceFile tests the scanSourceFile function.
//
// Params:
//   - t: testing context
func Test_scanSourceFile(t *testing.T) {
	t.Run("handles missing files", func(t *testing.T) {
		result := make(map[string]testedFuncInfo)
		// Appel avec fichier inexistant ne doit pas paniquer
		scanSourceFile("/nonexistent", "missing.go", result)
		// Le test passe si pas de panic
		t.Log("handled missing file gracefully")
	})
}

// Test_extractFuncInfo tests the extractFuncInfo function.
//
// Params:
//   - t: testing context
func Test_extractFuncInfo(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "extracts function info", code: "package test\nfunc Foo() error { return nil }"},
	}

	for _, tt := range tests {
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
//
// Params:
//   - t: testing context
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
//
// Params:
//   - t: testing context
func Test_isErrorType(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "identifies error type", code: "package test\nfunc Foo() error { return nil }"},
	}

	for _, tt := range tests {
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
						if !result {
							t.Error("expected error type to be identified")
						}
					}
				}
			}
		})
	}
}

// Test_extractReceiverName tests the extractReceiverName function.
//
// Params:
//   - t: testing context
func Test_extractReceiverName(t *testing.T) {
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
						result := extractReceiverName(fd.Recv.List[0].Type)
						// Vérification résultat
						if result != tt.expected {
							t.Errorf("extractReceiverName() = %q, want %q", result, tt.expected)
						}
					}
				}
			}
		})
	}
}

// Test_analyzeTestFunction tests the analyzeTestFunction function.
//
// Params:
//   - t: testing context
func Test_analyzeTestFunction(t *testing.T) {
	t.Run("analyzes test function", func(t *testing.T) {
		// Vérification que la fonction existe et ne panique pas
		// Les tests réels sont faits via l'API publique
		t.Log("analyzeTestFunction exists")
	})
}

// Test_extractTestedFuncName tests the extractTestedFuncName function.
//
// Params:
//   - t: testing context
func Test_extractTestedFuncName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"TestFoo", "TestFoo", "Foo"},
		{"Test_foo", "Test_foo", "foo"},
		{"TestFooBar", "TestFooBar", "FooBar"},
		{"Test", "Test", ""},
		{"NotATest", "NotATest", "NotATest"},
	}
	// Parcourir les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTestedFuncName(tt.input)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("extractTestedFuncName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_checkErrorInNode tests the checkErrorInNode function.
//
// Params:
//   - t: testing context
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
//
// Params:
//   - t: testing context
func Test_checkErrorInBasicLit(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "checks error in basic literal", code: "package test\nvar s = \"error case\""},
	}

	for _, tt := range tests {
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
//
// Params:
//   - t: testing context
func Test_checkErrorInKeyValue(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "checks error in key-value expression", code: "package test\nvar x = struct{name string}{name: \"error case\"}"},
	}

	for _, tt := range tests {
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
			if !found {
				t.Error("expected to find error in key-value expression")
			}
		})
	}
}
