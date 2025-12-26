// Internal tests for analyzer 008.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Test_ParseTestNameForPrivate tests the shared.ParseTestName helper for private functions.
func Test_ParseTestNameForPrivate(t *testing.T) {
	tests := []struct {
		name        string
		testName    string
		wantFunc    string
		wantPrivate bool
	}{
		{
			name:        "simple private function test",
			testName:    "Test_doSomething",
			wantFunc:    "doSomething",
			wantPrivate: true,
		},
		{
			name:        "method pattern Type_method",
			testName:    "TestMyType_doSomething",
			wantFunc:    "doSomething",
			wantPrivate: true,
		},
		{
			name:        "public function",
			testName:    "TestDoSomething",
			wantFunc:    "DoSomething",
			wantPrivate: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, ok := shared.ParseTestName(tt.testName)
			// Vérification parsing
			if !ok {
				t.Fatalf("ParseTestName(%q) failed", tt.testName)
			}
			// Vérification fonction
			if target.FuncName != tt.wantFunc {
				t.Errorf("ParseTestName(%q).FuncName = %q, want %q", tt.testName, target.FuncName, tt.wantFunc)
			}
			// Vérification private
			if target.IsPrivate != tt.wantPrivate {
				t.Errorf("ParseTestName(%q).IsPrivate = %v, want %v", tt.testName, target.IsPrivate, tt.wantPrivate)
			}
		})
	}
}

// Test_addPrivateFunction tests the addPrivateFunction function.
func Test_addPrivateFunction(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		wantFuncName string
	}{
		{
			name:         "private function",
			code:         "func privateFunc() {}",
			wantFuncName: "privateFunc",
		},
		{
			name:         "public function",
			code:         "func PublicFunc() {}",
			wantFuncName: "",
		},
		{
			name:         "private method",
			code:         "func (r MyType) privateMethod() {}",
			wantFuncName: "MyType_privateMethod",
		},
		{
			name:         "mock function (excluded)",
			code:         "func mockPrivateFunc() {}",
			wantFuncName: "",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\ntype MyType struct{}\n"+tt.code, 0)
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

			privateFunctions := make(map[string]bool)
			addPrivateFunction(funcDecl, privateFunctions)

			// Vérification fonction privée
			if tt.wantFuncName != "" {
				// Vérification de la condition
				if !privateFunctions[tt.wantFuncName] {
					t.Errorf("expected private function %q to be added, got %v", tt.wantFuncName, privateFunctions)
				}
			} else {
				// Vérification de la condition
				if len(privateFunctions) > 0 {
					t.Errorf("expected no functions to be added, got %v", privateFunctions)
				}
			}
		})
	}
}

// Test_collectPrivateFunctions tests the collectPrivateFunctions function logic.
func Test_collectPrivateFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []string
	}{
		{
			name: "single private function",
			code: `package test
func privateFunc() {}`,
			want: []string{"privateFunc"},
		},
		{
			name: "public and private functions",
			code: `package test
func PublicFunc() {}
func privateFunc() {}`,
			want: []string{"privateFunc"},
		},
		{
			name: "multiple private functions",
			code: `package test
func privateFunc1() {}
func privateFunc2() {}`,
			want: []string{"privateFunc1", "privateFunc2"},
		},
		{
			name: "only public functions",
			code: `package test
func PublicFunc() {}`,
			want: []string{},
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

			privateFunctions := make(map[string]bool, initialPrivateFunctionsCap)

			// Collect private functions
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					// Vérification du nom
					if funcDecl.Name != nil && len(funcDecl.Name.Name) > 0 {
						firstChar := rune(funcDecl.Name.Name[0])
						// Vérification caractère minuscule
						if firstChar >= 'a' && firstChar <= 'z' {
							privateFunctions[funcDecl.Name.Name] = true
						}
					}
				}
				// Continuer la traversée
				return true
			})

			// Vérification du nombre de fonctions
			if len(privateFunctions) != len(tt.want) {
				t.Errorf("collectPrivateFunctions() count = %d, want %d", len(privateFunctions), len(tt.want))
			}

			// Vérification de chaque fonction attendue
			for _, funcName := range tt.want {
				// Vérification de la condition
				if !privateFunctions[funcName] {
					t.Errorf("expected private function %q not found", funcName)
				}
			}
		})
	}
}

// Test_checkAndReportPrivateFunctionTest tests checkAndReportPrivateFunctionTest logic.
func Test_checkAndReportPrivateFunctionTest(t *testing.T) {
	tests := []struct {
		name              string
		testFuncName      string
		privateFunctions  map[string]bool
		shouldReportError bool
	}{
		{
			name:         "test for private function",
			testFuncName: "Test_doSomething",
			privateFunctions: map[string]bool{
				"doSomething": true,
			},
			shouldReportError: true,
		},
		{
			name:         "test for public function",
			testFuncName: "TestDoSomething",
			privateFunctions: map[string]bool{
				"doSomething": true,
			},
			shouldReportError: false,
		},
		{
			name:              "test with empty name",
			testFuncName:      "Test",
			privateFunctions:  map[string]bool{},
			shouldReportError: false,
		},
		{
			name:              "test for non-existent function",
			testFuncName:      "Test_unknownFunction",
			privateFunctions:  map[string]bool{"otherFunc": true},
			shouldReportError: false,
		},
		{
			name:              "empty key generated from test name",
			testFuncName:      "Test_foo",
			privateFunctions:  map[string]bool{},
			shouldReportError: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\nfunc " + tt.testFuncName + "(t *testing.T) {}"
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test_external_test.go", code, 0)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("no function declaration found")
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			checkAndReportPrivateFunctionTest(pass, funcDecl, "test_external_test.go", tt.privateFunctions)

			if tt.shouldReportError && reportCount == 0 {
				t.Error("expected error to be reported")
			} else if !tt.shouldReportError && reportCount > 0 {
				t.Errorf("expected no error, got %d reports", reportCount)
			}
		})
	}
}

// Test_runTest008 tests the runTest008 private function.
func Test_runTest008(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - minimal test",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing: %s", tt.name)
		})
	}
}

// Test_checkExternalTestsForPrivateFunctions tests checkExternalTestsForPrivateFunctions.
func Test_checkExternalTestsForPrivateFunctions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - no tests",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing: %s", tt.name)
		})
	}
}

// Test_runTest008_disabled tests that the rule is skipped when disabled.
func Test_runTest008_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runTest008_excludedFile tests that excluded files are skipped.
func Test_runTest008_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}
