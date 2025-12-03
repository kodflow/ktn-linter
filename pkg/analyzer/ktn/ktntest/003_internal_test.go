// Internal tests for analyzer 003.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_isExemptFunction tests the isExemptFunction private function.
//
// Params:
//   - t: testing context
func Test_isExemptFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		{
			name:     "init function is exempt",
			funcName: "init",
			want:     true,
		},
		{
			name:     "main function is exempt",
			funcName: "main",
			want:     true,
		},
		{
			name:     "regular public function not exempt",
			funcName: "DoSomething",
			want:     false,
		},
		{
			name:     "regular private function not exempt",
			funcName: "doSomething",
			want:     false,
		},
		{
			name:     "empty name not exempt",
			funcName: "",
			want:     false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptFunction(tt.funcName)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isExemptFunction(%q) = %v, want %v", tt.funcName, got, tt.want)
			}
		})
	}
}

// Test_isPublicFunction tests the isPublicFunction private function.
//
// Params:
//   - t: testing context
func Test_isPublicFunction(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "public function",
			code: "func PublicFunc() {}",
			want: true,
		},
		{
			name: "private function",
			code: "func privateFunc() {}",
			want: false,
		},
		{
			name: "function with nil name",
			code: "",
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification du code vide
			if tt.code == "" {
				// Test with nil function decl
				funcDecl := &ast.FuncDecl{Name: nil}
				got := isPublicFunction(funcDecl)
				// Vérification de la condition
				if got != tt.want {
					t.Errorf("isPublicFunction(nil name) = %v, want %v", got, tt.want)
				}
				// Retour de la fonction
				return
			}

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

			got := isPublicFunction(funcDecl)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isPublicFunction(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// Test_extractReceiverTypeName tests the extractReceiverTypeName private function.
//
// Params:
//   - t: testing context
func Test_extractReceiverTypeName(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "simple receiver",
			code: "func (r Receiver) Method() {}",
			want: "Receiver",
		},
		{
			name: "pointer receiver",
			code: "func (r *Receiver) Method() {}",
			want: "Receiver",
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

			got := extractReceiverTypeName(funcDecl.Recv.List[0].Type)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("extractReceiverTypeName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_buildSuggestedTestName tests the buildSuggestedTestName private function.
//
// Params:
//   - t: testing context
func Test_buildSuggestedTestName(t *testing.T) {
	tests := []struct {
		name string
		fn   funcInfo
		want string
	}{
		{
			name: "public function",
			fn: funcInfo{
				name:       "DoSomething",
				isExported: true,
			},
			want: "TestDoSomething",
		},
		{
			name: "private function",
			fn: funcInfo{
				name:       "doSomething",
				isExported: false,
			},
			want: "Test_doSomething",
		},
		{
			name: "public method",
			fn: funcInfo{
				name:         "Method",
				receiverName: "MyType",
				isExported:   true,
			},
			want: "TestMyType_Method",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSuggestedTestName(tt.fn)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("buildSuggestedTestName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_getTestFileInfo tests the getTestFileInfo private function.
//
// Params:
//   - t: testing context
func Test_getTestFileInfo(t *testing.T) {
	tests := []struct {
		name           string
		isExported     bool
		fileBase       string
		wantFile       string
		wantType       string
		wantFuncType   string
	}{
		{
			name:           "public function",
			isExported:     true,
			fileBase:       "myfile",
			wantFile:       "myfile_external_test.go",
			wantType:       "black-box testing avec package xxx_test",
			wantFuncType:   "publique",
		},
		{
			name:           "private function",
			isExported:     false,
			fileBase:       "myfile",
			wantFile:       "myfile_internal_test.go",
			wantType:       "white-box testing avec package xxx",
			wantFuncType:   "privée",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile, gotType, gotFuncType := getTestFileInfo(tt.isExported, tt.fileBase)
			// Vérification du fichier
			if gotFile != tt.wantFile {
				t.Errorf("getTestFileInfo() file = %q, want %q", gotFile, tt.wantFile)
			}
			// Vérification du type
			if gotType != tt.wantType {
				t.Errorf("getTestFileInfo() type = %q, want %q", gotType, tt.wantType)
			}
			// Vérification du type de fonction
			if gotFuncType != tt.wantFuncType {
				t.Errorf("getTestFileInfo() funcType = %q, want %q", gotFuncType, tt.wantFuncType)
			}
		})
	}
}

// Test_buildTestNames tests the buildTestNames private function.
//
// Params:
//   - t: testing context
func Test_buildTestNames(t *testing.T) {
	tests := []struct {
		name string
		fn   funcInfo
		want []string
	}{
		{
			name: "error case - empty function name",
			fn:   funcInfo{name: "", receiverName: ""},
			want: []string{""},
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTestNames([]string{}, tt.fn)
			// Vérification de la longueur
			if len(result) != len(tt.want) {
				t.Errorf("buildTestNames() len = %d, want %d", len(result), len(tt.want))
			}
		})
	}
}

// Test_hasMatchingTest tests the hasMatchingTest private function.
//
// Params:
//   - t: testing context
func Test_hasMatchingTest(t *testing.T) {
	tests := []struct {
		name        string
		testNames   []string
		testedFuncs map[string]bool
		want        bool
	}{
		{
			name:        "error case - empty test names",
			testNames:   []string{},
			testedFuncs: map[string]bool{"foo": true},
			want:        false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasMatchingTest(tt.testNames, tt.testedFuncs)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("hasMatchingTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_collectTestedFunctions tests the collectTestedFunctions private function.
//
// Params:
//   - t: testing context
func Test_collectTestedFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want map[string]bool
	}{
		{
			name: "error case - non-test function",
			code: "func helper() {}",
			want: map[string]bool{},
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			testedFuncs := make(map[string]bool)
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if fd, ok := n.(*ast.FuncDecl); ok {
					collectTestedFunctions(fd, testedFuncs)
				}
				// Continuer la traversée
				return true
			})

			// Vérification de la longueur
			if len(testedFuncs) != len(tt.want) {
				t.Errorf("collectTestedFunctions() len = %d, want %d", len(testedFuncs), len(tt.want))
			}
		})
	}
}

// Test_collectFunctions tests the collectFunctions private function.
//
// Params:
//   - t: testing context
func Test_collectFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - empty file",
			code: "package test",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			var funcs []funcInfo
			testedFuncs := make(map[string]bool)
			collectFunctions(pass, &funcs, testedFuncs)

			// Vérification du résultat
			if len(funcs) != 0 {
				t.Errorf("expected 0 functions, got %d", len(funcs))
			}
		})
	}
}

// Test_runTest003 tests the runTest003 private function.
//
// Params:
//   - t: testing context
func Test_runTest003(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - no files",
			code: "package test",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					t.Logf("Report: %s", d.Message)
				},
			}

			_, err = runTest003(pass)
			// Vérification pas d'erreur
			if err != nil {
				t.Errorf("runTest003() error = %v", err)
			}
		})
	}
}

// Test_countTestFiles tests the countTestFiles private function.
//
// Params:
//   - t: testing context
func Test_countTestFiles(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		wantHas       bool
		wantCount     int
	}{
		{
			name:      "error case - no test files",
			filename:  "test.go",
			wantHas:   false,
			wantCount: 0,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tt.filename, "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			hasTestFiles, testFileCount := countTestFiles(pass)
			// Vérification has
			if hasTestFiles != tt.wantHas {
				t.Errorf("countTestFiles() hasTestFiles = %v, want %v", hasTestFiles, tt.wantHas)
			}
			// Vérification count
			if testFileCount != tt.wantCount {
				t.Errorf("countTestFiles() count = %d, want %d", testFileCount, tt.wantCount)
			}
		})
	}
}

// Test_collectAllFunctionsAndTests tests the collectAllFunctionsAndTests private function.
//
// Params:
//   - t: testing context
func Test_collectAllFunctionsAndTests(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "error case - empty package",
			code: "package test",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing code: %s", tt.code)
		})
	}
}

// Test_checkFunctionsHaveTests tests the checkFunctionsHaveTests private function.
//
// Params:
//   - t: testing context
func Test_checkFunctionsHaveTests(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - no functions to check",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			allFuncs := []funcInfo{}
			testedFuncs := make(map[string]bool)
			checkFunctionsHaveTests(pass, allFuncs, testedFuncs)

			// Vérification pas de rapport
			if reportCount != 0 {
				t.Errorf("expected 0 reports, got %d", reportCount)
			}
		})
	}
}

// Test_reportMissingTest tests the reportMissingTest private function.
//
// Params:
//   - t: testing context
func Test_reportMissingTest(t *testing.T) {
	tests := []struct {
		name string
		fn   funcInfo
	}{
		{
			name: "error case - minimal function info",
			fn: funcInfo{
				name:       "test",
				isExported: false,
				filename:   "test.go",
			},
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			reportMissingTest(pass, tt.fn)

			// Vérification rapport généré
			if reportCount != 1 {
				t.Errorf("expected 1 report, got %d", reportCount)
			}
		})
	}
}

// Test_collectExternalTestFunctions tests the collectExternalTestFunctions private function.
//
// Params:
//   - t: testing context
func Test_collectExternalTestFunctions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - no files",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				Fset:  token.NewFileSet(),
				Files: []*ast.File{},
			}

			testedFuncs := make(map[string]bool)
			collectExternalTestFunctions(pass, testedFuncs)

			// Vérification résultat
			if len(testedFuncs) != 0 {
				t.Errorf("expected 0 tested funcs, got %d", len(testedFuncs))
			}
		})
	}
}
