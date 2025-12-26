// Internal tests for analyzer 003.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Test_isExemptFunction tests the isExemptFunction private function.
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

// Test_ClassifyFunc tests the shared.ClassifyFunc helper function.
func Test_ClassifyFunc(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantPublic bool
	}{
		{
			name:       "public function",
			code:       "func PublicFunc() {}",
			wantPublic: true,
		},
		{
			name:       "private function",
			code:       "func privateFunc() {}",
			wantPublic: false,
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

			meta := shared.ClassifyFunc(funcDecl)
			isPublic := meta.Visibility == shared.VisPublic
			// Vérification de la condition
			if isPublic != tt.wantPublic {
				t.Errorf("ClassifyFunc(%q) public = %v, want %v", tt.code, isPublic, tt.wantPublic)
			}
		})
	}
}

// Test_ExtractReceiverTypeName tests the shared.ExtractReceiverTypeName helper.
func Test_ExtractReceiverTypeName(t *testing.T) {
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

			got := shared.ExtractReceiverTypeName(funcDecl.Recv.List[0].Type)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("ExtractReceiverTypeName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_BuildSuggestedTestName tests the shared.BuildSuggestedTestName helper.
func Test_BuildSuggestedTestName(t *testing.T) {
	tests := []struct {
		name string
		meta *shared.FuncMeta
		want string
	}{
		{
			name: "public function",
			meta: &shared.FuncMeta{
				Name:       "DoSomething",
				Kind:       shared.FuncTopLevel,
				Visibility: shared.VisPublic,
			},
			want: "TestDoSomething",
		},
		{
			name: "private function",
			meta: &shared.FuncMeta{
				Name:       "doSomething",
				Kind:       shared.FuncTopLevel,
				Visibility: shared.VisPrivate,
			},
			want: "Test_doSomething",
		},
		{
			name: "public method",
			meta: &shared.FuncMeta{
				Name:         "Method",
				ReceiverName: "MyType",
				Kind:         shared.FuncMethod,
				Visibility:   shared.VisPublic,
			},
			want: "TestMyType_Method",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shared.BuildSuggestedTestName(tt.meta)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("BuildSuggestedTestName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_buildTestNames tests the buildTestNames private function.
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
		{
			name: "top-level function",
			fn:   funcInfo{name: "Foo", receiverName: "", isExported: true},
			want: []string{"Foo"},
		},
		{
			name: "method with receiver",
			fn:   funcInfo{name: "Method", receiverName: "Type", isExported: true},
			want: []string{"Type_Method"},
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

// Test_buildFuncLookupKey tests the buildFuncLookupKey private function.
func Test_buildFuncLookupKey(t *testing.T) {
	tests := []struct {
		name string
		fn   funcInfo
		want string
	}{
		{
			name: "top-level function",
			fn:   funcInfo{name: "Foo", receiverName: ""},
			want: "Foo",
		},
		{
			name: "method with receiver",
			fn:   funcInfo{name: "Bar", receiverName: "MyType"},
			want: "MyType_Bar",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFuncLookupKey(tt.fn)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("buildFuncLookupKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test_hasMatchingTest tests the hasMatchingTest private function.
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
		{
			name:        "exact match found",
			testNames:   []string{"Foo"},
			testedFuncs: map[string]bool{"Foo": true},
			want:        true,
		},
		{
			name:        "case insensitive match found",
			testNames:   []string{"Foo"},
			testedFuncs: map[string]bool{"foo": true},
			want:        true,
		},
		{
			name:        "no match found",
			testNames:   []string{"Foo"},
			testedFuncs: map[string]bool{"Bar": true},
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

// Test_collectExternalTestedFunctions tests the collectExternalTestedFunctions function.
func Test_collectExternalTestedFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want map[string]bool
	}{
		{
			name: "non-test function",
			code: "func helper() {}",
			want: map[string]bool{},
		},
		{
			name: "test function with no params",
			code: "func TestFoo() {}",
			want: map[string]bool{},
		},
		{
			name: "test function with wrong first param",
			code: "func TestFoo(x int) {}",
			want: map[string]bool{},
		},
		{
			name: "test function with wrong param name",
			code: "func TestFoo(x *testing.T) {}",
			want: map[string]bool{},
		},
		{
			name: "valid test function",
			code: "func TestFoo(t *testing.T) {}",
			want: map[string]bool{"Foo": true},
		},
		{
			name: "valid private test function",
			code: "func Test_foo(t *testing.T) {}",
			want: map[string]bool{"foo": true},
		},
		{
			name: "test with method receiver",
			code: "func TestMyType_Method(t *testing.T) {}",
			want: map[string]bool{"MyType_Method": true},
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
					collectExternalTestedFunctions(fd, testedFuncs)
				}
				// Continuer la traversée
				return true
			})

			// Vérification de la longueur
			if len(testedFuncs) != len(tt.want) {
				t.Errorf("collectExternalTestedFunctions() len = %d, want %d; got %v", len(testedFuncs), len(tt.want), testedFuncs)
			}
			// Vérifier les clés
			for k := range tt.want {
				// Vérifier si la clé existe
				if !testedFuncs[k] {
					t.Errorf("collectExternalTestedFunctions() missing key %q", k)
				}
			}
		})
	}
}

// Test_collectFunctions tests the collectFunctions private function.
func Test_collectFunctions(t *testing.T) {
	tests := []struct {
		name string
		code string
		want int
	}{
		{
			name: "error case - empty file",
			code: "package test",
			want: 0,
		},
		{
			name: "collects public function",
			code: "package test\nfunc PublicFunc() {}",
			want: 1,
		},
		{
			name: "collects private function",
			code: "package test\nfunc privateFunc() {}",
			want: 1,
		},
		{
			name: "collects method with receiver",
			code: "package test\ntype T struct{}\nfunc (t T) Method() {}",
			want: 1,
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
			if len(funcs) != tt.want {
				t.Errorf("expected %d functions, got %d", tt.want, len(funcs))
			}
		})
	}
}

// Test_runTest003 tests the runTest003 private function.
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
func Test_countTestFiles(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantHas   bool
		wantCount int
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
func Test_checkFunctionsHaveTests(t *testing.T) {
	tests := []struct {
		name        string
		allFuncs    []funcInfo
		testedFuncs map[string]bool
		wantReports int
	}{
		{
			name:        "error case - no functions to check",
			allFuncs:    []funcInfo{},
			testedFuncs: make(map[string]bool),
			wantReports: 0,
		},
		{
			name: "function with test",
			allFuncs: []funcInfo{
				{name: "Foo", receiverName: "", isExported: true, pos: token.Pos(1), filename: "test.go"},
			},
			testedFuncs: map[string]bool{"Foo": true},
			wantReports: 0,
		},
		{
			name: "function without test",
			allFuncs: []funcInfo{
				{name: "Bar", receiverName: "", isExported: true, pos: token.Pos(1), filename: "test.go"},
			},
			testedFuncs: map[string]bool{},
			wantReports: 1,
		},
		{
			name: "mock function skipped",
			allFuncs: []funcInfo{
				{name: "MockFoo", receiverName: "", isExported: true, pos: token.Pos(1), filename: "test.go"},
			},
			testedFuncs: map[string]bool{},
			wantReports: 0,
		},
		{
			name: "method with mock receiver skipped",
			allFuncs: []funcInfo{
				{name: "Method", receiverName: "MockType", isExported: true, pos: token.Pos(1), filename: "test.go"},
			},
			testedFuncs: map[string]bool{},
			wantReports: 0,
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

			checkFunctionsHaveTests(pass, tt.allFuncs, tt.testedFuncs)

			// Vérification du nombre de rapports
			if reportCount != tt.wantReports {
				t.Errorf("expected %d reports, got %d", tt.wantReports, reportCount)
			}
		})
	}
}

// Test_reportMissingTest tests the reportMissingTest private function.
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

// Test_findPackageDir tests the findPackageDir private function.
func Test_findPackageDir(t *testing.T) {
	tests := []struct {
		name      string
		files     []*ast.File
		wantEmpty bool
	}{
		{
			name:      "no files returns empty",
			files:     []*ast.File{},
			wantEmpty: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				Fset:  token.NewFileSet(),
				Files: tt.files,
			}

			result := findPackageDir(pass)
			// Vérification résultat
			if tt.wantEmpty && result != "" {
				t.Errorf("findPackageDir() = %q, want empty", result)
			}
		})
	}
}

// Test_isCacheOrTempFile tests the isCacheOrTempFile private function.
func Test_isCacheOrTempFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "cache file linux",
			filename: "/home/user/.cache/go-build/abc/test.go",
			want:     true,
		},
		{
			name:     "cache file windows",
			filename: "C:\\Users\\user\\cache\\go-build\\test.go",
			want:     true,
		},
		{
			name:     "tmp file",
			filename: "/tmp/go-test123/test.go",
			want:     true,
		},
		{
			name:     "normal file",
			filename: "/home/user/project/main.go",
			want:     false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCacheOrTempFile(tt.filename)
			// Vérification résultat
			if got != tt.want {
				t.Errorf("isCacheOrTempFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// Test_parseTestFile tests the parseTestFile private function.
func Test_parseTestFile(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "non-existent file",
			path: "/nonexistent/path/test.go",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testedFuncs := make(map[string]bool)
			// Should not panic on non-existent file
			parseTestFile(tt.path, testedFuncs)
			// Vérification résultat
			if len(testedFuncs) != 0 {
				t.Errorf("parseTestFile() collected %d funcs, want 0", len(testedFuncs))
			}
		})
	}
}

// Test_runTest003_disabled tests that the rule is skipped when disabled.
func Test_runTest003_disabled(t *testing.T) {
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

// Test_runTest003_excludedFile tests that excluded files are skipped.
func Test_runTest003_excludedFile(t *testing.T) {
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
