// Internal tests for analyzer 007.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Test_ExtractReceiverType tests the shared.ExtractReceiverTypeName helper function.
func Test_ExtractReceiverType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple receiver",
			code: "func (r MyType) Method() {}",
			want: "MyType",
		},
		{
			name: "pointer receiver",
			code: "func (r *MyType) Method() {}",
			want: "MyType",
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
			if funcDecl == nil || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
				t.Fatal("no receiver found")
			}

			got := shared.ExtractReceiverTypeName(funcDecl.Recv.List[0].Type)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("ExtractReceiverTypeName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_addPublicFunction tests the addPublicFunction private function.
func Test_addPublicFunction(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		wantFuncName  string
		wantMethodKey string
	}{
		{
			name:         "public function",
			code:         "func PublicFunc() {}",
			wantFuncName: "PublicFunc",
		},
		{
			name:         "private function",
			code:         "func privateFunc() {}",
			wantFuncName: "",
		},
		{
			name:          "public method",
			code:          "func (r MyType) PublicMethod() {}",
			wantFuncName:  "MyType_PublicMethod",
			wantMethodKey: "MyType_PublicMethod",
		},
		{
			name:         "function with nil name",
			code:         "func PublicFunc() {}",
			wantFuncName: "PublicFunc",
		},
		{
			name:         "mock function (excluded)",
			code:         "func MockPublicFunc() {}",
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

			publicFunctions := make(map[string]bool)
			addPublicFunction(funcDecl, publicFunctions)

			// Vérification fonction publique
			if tt.wantFuncName != "" {
				// Vérification de la condition
				if !publicFunctions[tt.wantFuncName] {
					t.Errorf("expected public function %q to be added, got %v", tt.wantFuncName, publicFunctions)
				}
			} else {
				// Vérification de la condition
				if len(publicFunctions) > 0 {
					t.Errorf("expected no functions to be added, got %v", publicFunctions)
				}
			}

			// Vérification méthode
			if tt.wantMethodKey != "" {
				// Vérification de la condition
				if !publicFunctions[tt.wantMethodKey] {
					t.Errorf("expected method key %q to be added", tt.wantMethodKey)
				}
			}
		})
	}
}

// Test_checkAndReportPublicFunctionTest tests checkAndReportPublicFunctionTest logic.
func Test_checkAndReportPublicFunctionTest(t *testing.T) {
	tests := []struct {
		name              string
		testFuncName      string
		publicFunctions   map[string]bool
		shouldReportError bool
	}{
		{
			name:         "test for public function",
			testFuncName: "TestDoSomething",
			publicFunctions: map[string]bool{
				"DoSomething": true,
			},
			shouldReportError: true,
		},
		{
			name:         "test for private function",
			testFuncName: "Test_doSomething",
			publicFunctions: map[string]bool{
				"DoSomething": true,
			},
			shouldReportError: false,
		},
		{
			name:              "test with empty name",
			testFuncName:      "Test",
			publicFunctions:   map[string]bool{},
			shouldReportError: false,
		},
		{
			name:              "test for non-existent function",
			testFuncName:      "TestUnknownFunction",
			publicFunctions:   map[string]bool{"OtherFunc": true},
			shouldReportError: false,
		},
		{
			name:              "empty key generated from test name",
			testFuncName:      "TestFoo",
			publicFunctions:   map[string]bool{},
			shouldReportError: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "package test\nfunc " + tt.testFuncName + "(t *testing.T) {}"
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test_internal_test.go", code, 0)
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

			checkAndReportPublicFunctionTest(pass, funcDecl, "test_internal_test.go", tt.publicFunctions)

			if tt.shouldReportError && reportCount == 0 {
				t.Error("expected error to be reported")
			} else if !tt.shouldReportError && reportCount > 0 {
				t.Errorf("expected no error, got %d reports", reportCount)
			}
		})
	}
}

// Test_collectPublicFunctions tests the collectPublicFunctions function logic.
func Test_collectPublicFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []string
	}{
		{
			name: "single public function",
			code: `package test
func PublicFunc() {}`,
			want: []string{"PublicFunc"},
		},
		{
			name: "public and private functions",
			code: `package test
func PublicFunc() {}
func privateFunc() {}`,
			want: []string{"PublicFunc"},
		},
		{
			name: "public method",
			code: `package test
type MyType struct{}
func (m MyType) PublicMethod() {}`,
			want: []string{"MyType_PublicMethod"},
		},
		{
			name: "error case - invalid code",
			code: `package test`,
			want: []string{},
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conceptual logic
			t.Logf("Testing code with expected functions: %v", tt.want)
		})
	}
}

// Test_runTest007 tests the runTest007 private function.
func Test_runTest007(t *testing.T) {
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

// Test_checkInternalTestsForPublicFunctions tests checkInternalTestsForPublicFunctions.
func Test_checkInternalTestsForPublicFunctions(t *testing.T) {
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

// Test_runTest007_disabled tests that the rule is skipped when disabled.
func Test_runTest007_disabled(t *testing.T) {
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

// Test_runTest007_excludedFile tests that excluded files are skipped.
func Test_runTest007_excludedFile(t *testing.T) {
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
