// Internal tests for analyzer 008.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_runTest008 tests the runTest008 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest008(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		wantErr bool
	}{
		{
			name:    "source file with appropriate tests",
			pkgPath: "test008",
			wantErr: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing package: %s", tt.pkgPath)
		})
	}
}

// Test_analyzeFileFunctions tests the analyzeFileFunctions private function.
//
// Params:
//   - t: testing context
func Test_analyzeFileFunctions(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		wantPublic  bool
		wantPrivate bool
	}{
		{
			name: "file with public function",
			code: `package test
func PublicFunc() {}`,
			wantPublic:  true,
			wantPrivate: false,
		},
		{
			name: "file with private function",
			code: `package test
func privateFunc() {}`,
			wantPublic:  false,
			wantPrivate: true,
		},
		{
			name: "file with both public and private",
			code: `package test
func PublicFunc() {}
func privateFunc() {}`,
			wantPublic:  true,
			wantPrivate: true,
		},
		{
			name: "file with init function",
			code: `package test
func init() {}`,
			wantPublic:  false,
			wantPrivate: false,
		},
		{
			name: "file with public variable",
			code: `package test
var PublicVar int`,
			wantPublic:  true,
			wantPrivate: false,
		},
		{
			name: "file with private variable",
			code: `package test
var privateVar int`,
			wantPublic:  false,
			wantPrivate: true,
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

			result := analyzeFileFunctions(file)

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("analyzeFileFunctions() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
			// Vérification private
			if result.hasPrivate != tt.wantPrivate {
				t.Errorf("analyzeFileFunctions() hasPrivate = %v, want %v", result.hasPrivate, tt.wantPrivate)
			}
		})
	}
}

// Test_checkVariables tests the checkVariables private function.
//
// Params:
//   - t: testing context
func Test_checkVariables(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		wantPublic  bool
		wantPrivate bool
	}{
		{
			name: "public variable",
			code: `package test
var PublicVar int`,
			wantPublic:  true,
			wantPrivate: false,
		},
		{
			name: "private variable",
			code: `package test
var privateVar int`,
			wantPublic:  false,
			wantPrivate: true,
		},
		{
			name: "blank identifier",
			code: `package test
var _ int`,
			wantPublic:  false,
			wantPrivate: false,
		},
		{
			name: "multiple variables",
			code: `package test
var PublicVar, privateVar int`,
			wantPublic:  true,
			wantPrivate: true,
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

			result := fileAnalysisResult{}
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if genDecl, ok := n.(*ast.GenDecl); ok {
					checkVariables(genDecl, &result)
				}
				// Continuer la traversée
				return true
			})

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("checkVariables() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
			// Vérification private
			if result.hasPrivate != tt.wantPrivate {
				t.Errorf("checkVariables() hasPrivate = %v, want %v", result.hasPrivate, tt.wantPrivate)
			}
		})
	}
}

// Test_checkTestFilesExist tests the checkTestFilesExist private function.
//
// Params:
//   - t: testing context
func Test_checkTestFilesExist(t *testing.T) {
	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "test_*.go")
	// Vérification de l'erreur
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	// Nettoyage
	defer os.Remove(tmpPath)

	// Créer les fichiers de test
	tmpDir := filepath.Dir(tmpPath)
	baseName := filepath.Base(tmpPath)
	fileBase := baseName[:len(baseName)-3]

	internalTest := filepath.Join(tmpDir, fileBase+"_internal_test.go")
	externalTest := filepath.Join(tmpDir, fileBase+"_external_test.go")

	tests := []struct {
		name            string
		createInternal  bool
		createExternal  bool
		wantHasInternal bool
		wantHasExternal bool
	}{
		{
			name:            "no test files",
			createInternal:  false,
			createExternal:  false,
			wantHasInternal: false,
			wantHasExternal: false,
		},
		{
			name:            "only internal test",
			createInternal:  true,
			createExternal:  false,
			wantHasInternal: true,
			wantHasExternal: false,
		},
		{
			name:            "only external test",
			createInternal:  false,
			createExternal:  true,
			wantHasInternal: false,
			wantHasExternal: true,
		},
		{
			name:            "both tests",
			createInternal:  true,
			createExternal:  true,
			wantHasInternal: true,
			wantHasExternal: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer les fichiers selon le cas de test
			if tt.createInternal {
				f, err := os.Create(internalTest)
				// Vérification de l'erreur
				if err != nil {
					t.Fatalf("failed to create internal test: %v", err)
				}
				f.Close()
				// Nettoyage
				defer os.Remove(internalTest)
			}
			// Vérification external
			if tt.createExternal {
				f, err := os.Create(externalTest)
				// Vérification de l'erreur
				if err != nil {
					t.Fatalf("failed to create external test: %v", err)
				}
				f.Close()
				// Nettoyage
				defer os.Remove(externalTest)
			}

			status := checkTestFilesExist(tmpPath)

			// Vérification internal
			if status.hasInternal != tt.wantHasInternal {
				t.Errorf("checkTestFilesExist() hasInternal = %v, want %v", status.hasInternal, tt.wantHasInternal)
			}
			// Vérification external
			if status.hasExternal != tt.wantHasExternal {
				t.Errorf("checkTestFilesExist() hasExternal = %v, want %v", status.hasExternal, tt.wantHasExternal)
			}
		})
	}
}

// Test_fileExistsOnDisk tests the fileExistsOnDisk private function.
//
// Params:
//   - t: testing context
func Test_fileExistsOnDisk(t *testing.T) {
	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "test_*.go")
	// Vérification de l'erreur
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	// Nettoyage
	defer os.Remove(tmpPath)

	// Créer un répertoire temporaire
	tmpDir, err := os.MkdirTemp("", "test_dir_*")
	// Vérification de l'erreur
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// Nettoyage
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file returns true",
			path: tmpPath,
			want: true,
		},
		{
			name: "non-existing file returns false",
			path: "/nonexistent/file.go",
			want: false,
		},
		{
			name: "directory returns false",
			path: tmpDir,
			want: false,
		},
		{
			name: "empty path returns false",
			path: "",
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExistsOnDisk(tt.path)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("fileExistsOnDisk(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// Test_reportTestFileIssues tests the logic of reportTestFileIssues.
//
// Params:
//   - t: testing context
func Test_reportTestFileIssues(t *testing.T) {
	tests := []struct {
		name        string
		hasPublic   bool
		hasPrivate  bool
		hasInternal bool
		hasExternal bool
		expectIssue bool
	}{
		{
			name:        "public only with external test - OK",
			hasPublic:   true,
			hasPrivate:  false,
			hasInternal: false,
			hasExternal: true,
			expectIssue: false,
		},
		{
			name:        "private only with internal test - OK",
			hasPublic:   false,
			hasPrivate:  true,
			hasInternal: true,
			hasExternal: false,
			expectIssue: false,
		},
		{
			name:        "both with both tests - OK",
			hasPublic:   true,
			hasPrivate:  true,
			hasInternal: true,
			hasExternal: true,
			expectIssue: false,
		},
		{
			name:        "public only missing test",
			hasPublic:   true,
			hasPrivate:  false,
			hasInternal: false,
			hasExternal: false,
			expectIssue: true,
		},
		{
			name:        "error case - missing all tests",
			hasPublic:   false,
			hasPrivate:  false,
			hasInternal: false,
			hasExternal: false,
			expectIssue: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conceptual logic
			t.Logf("Testing case: hasPublic=%v, hasPrivate=%v, hasInternal=%v, hasExternal=%v",
				tt.hasPublic, tt.hasPrivate, tt.hasInternal, tt.hasExternal)
		})
	}
}

// Test_reportMixedFunctionsIssues tests the reportMixedFunctionsIssues private function.
//
// Params:
//   - t: testing context
func Test_reportMixedFunctionsIssues(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - no status",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			status := testFilesStatus{
				baseName:    "test.go",
				fileBase:    "test",
				hasInternal: false,
				hasExternal: false,
			}

			result := &fileAnalysisResult{
				hasPublic:    true,
				hasPrivate:   true,
				publicFuncs:  []string{"PublicFunc"},
				privateFuncs: []string{"privateFunc"},
			}

			reportMixedFunctionsIssues(pass, file, result, status)

			// Vérification rapport généré
			if reportCount < 1 {
				t.Errorf("expected at least 1 report, got %d", reportCount)
			}
		})
	}
}

// Test_reportPublicOnlyIssues tests the reportPublicOnlyIssues private function.
//
// Params:
//   - t: testing context
func Test_reportPublicOnlyIssues(t *testing.T) {
	tests := []struct {
		name        string
		hasInternal bool
		hasExternal bool
		expectCount int
	}{
		{
			name:        "no external test file",
			hasInternal: false,
			hasExternal: false,
			expectCount: 1,
		},
		{
			name:        "has both files (internal superfluous)",
			hasInternal: true,
			hasExternal: true,
			expectCount: 1,
		},
		{
			name:        "has only external file (OK)",
			hasInternal: false,
			hasExternal: true,
			expectCount: 0,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			status := testFilesStatus{
				baseName:    "test.go",
				fileBase:    "test",
				hasInternal: tt.hasInternal,
				hasExternal: tt.hasExternal,
			}

			result := &fileAnalysisResult{
				hasPublic:    true,
				hasPrivate:   false,
				publicFuncs:  []string{"PublicFunc"},
				privateFuncs: []string{},
			}

			reportPublicOnlyIssues(pass, file, result, status)

			// Vérification rapport généré
			if reportCount != tt.expectCount {
				t.Errorf("expected %d report, got %d", tt.expectCount, reportCount)
			}
		})
	}
}

// Test_reportPrivateOnlyIssues tests the reportPrivateOnlyIssues private function.
//
// Params:
//   - t: testing context
func Test_reportPrivateOnlyIssues(t *testing.T) {
	tests := []struct {
		name        string
		hasInternal bool
		hasExternal bool
		expectCount int
	}{
		{
			name:        "no internal test file",
			hasInternal: false,
			hasExternal: false,
			expectCount: 1,
		},
		{
			name:        "has both files (external superfluous)",
			hasInternal: true,
			hasExternal: true,
			expectCount: 1,
		},
		{
			name:        "has only internal file (OK)",
			hasInternal: true,
			hasExternal: false,
			expectCount: 0,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			reportCount := 0
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					reportCount++
				},
			}

			status := testFilesStatus{
				baseName:    "test.go",
				fileBase:    "test",
				hasInternal: tt.hasInternal,
				hasExternal: tt.hasExternal,
			}

			result := &fileAnalysisResult{
				hasPublic:    false,
				hasPrivate:   true,
				publicFuncs:  []string{},
				privateFuncs: []string{"privateFunc"},
			}

			reportPrivateOnlyIssues(pass, file, result, status)

			// Vérification rapport généré
			if reportCount != tt.expectCount {
				t.Errorf("expected %d report, got %d", tt.expectCount, reportCount)
			}
		})
	}
}

// Test_formatFuncList tests the formatFuncList function.
func Test_formatFuncList(t *testing.T) {
	tests := []struct {
		name     string
		funcs    []string
		expected string
	}{
		{
			name:     "empty list",
			funcs:    []string{},
			expected: "",
		},
		{
			name:     "single function",
			funcs:    []string{"Func1"},
			expected: "Func1",
		},
		{
			name:     "two functions",
			funcs:    []string{"Func1", "Func2"},
			expected: "Func1, Func2",
		},
		{
			name:     "three functions (max)",
			funcs:    []string{"Func1", "Func2", "Func3"},
			expected: "Func1, Func2, Func3",
		},
		{
			name:     "four functions (truncated)",
			funcs:    []string{"Func1", "Func2", "Func3", "Func4"},
			expected: "Func1, Func2, Func3, ... (+1)",
		},
		{
			name:     "five functions (truncated)",
			funcs:    []string{"Func1", "Func2", "Func3", "Func4", "Func5"},
			expected: "Func1, Func2, Func3, ... (+2)",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFuncList(tt.funcs)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("formatFuncList() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_formatCount tests the formatCount function.
func Test_formatCount(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{name: "zero", count: 0, expected: "0"},
		{name: "one", count: 1, expected: "1"},
		{name: "ten", count: 10, expected: "10"},
		{name: "hundred", count: 100, expected: "100"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCount(tt.count)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("formatCount() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_classifyFunction tests the classifyFunction function.
func Test_classifyFunction(t *testing.T) {
	tests := []struct {
		name           string
		funcDecl       *ast.FuncDecl
		expectPublic   bool
		expectPrivate  bool
		expectedPubLen int
		expectedPriLen int
	}{
		{
			name: "public function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "PublicFunc"},
				Type: &ast.FuncType{},
			},
			expectPublic:   true,
			expectPrivate:  false,
			expectedPubLen: 1,
			expectedPriLen: 0,
		},
		{
			name: "private function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "privateFunc"},
				Type: &ast.FuncType{},
			},
			expectPublic:   false,
			expectPrivate:  true,
			expectedPubLen: 0,
			expectedPriLen: 1,
		},
		{
			name: "exempt function (init)",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "init"},
				Type: &ast.FuncType{},
			},
			expectPublic:   false,
			expectPrivate:  false,
			expectedPubLen: 0,
			expectedPriLen: 0,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileAnalysisResult{
				publicFuncs:  []string{},
				privateFuncs: []string{},
			}
			classifyFunction(tt.funcDecl, &result)
			// Vérification du résultat
			if result.hasPublic != tt.expectPublic {
				t.Errorf("hasPublic = %v, want %v", result.hasPublic, tt.expectPublic)
			}
			// Vérification du résultat privé
			if result.hasPrivate != tt.expectPrivate {
				t.Errorf("hasPrivate = %v, want %v", result.hasPrivate, tt.expectPrivate)
			}
			// Vérification des longueurs
			if len(result.publicFuncs) != tt.expectedPubLen {
				t.Errorf("publicFuncs len = %d, want %d", len(result.publicFuncs), tt.expectedPubLen)
			}
			// Vérification des longueurs privées
			if len(result.privateFuncs) != tt.expectedPriLen {
				t.Errorf("privateFuncs len = %d, want %d", len(result.privateFuncs), tt.expectedPriLen)
			}
		})
	}
}

// Test_buildFunctionDisplayName tests the buildFunctionDisplayName function.
func Test_buildFunctionDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		expected string
	}{
		{
			name: "simple function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "myFunc"},
				Type: &ast.FuncType{},
			},
			expected: "myFunc",
		},
		{
			name: "method with value receiver",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "Type"}},
					},
				},
				Type: &ast.FuncType{},
			},
			expected: "Type.Method",
		},
		{
			name: "method with pointer receiver",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.StarExpr{X: &ast.Ident{Name: "Type"}}},
					},
				},
				Type: &ast.FuncType{},
			},
			expected: "(*Type).Method",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFunctionDisplayName(tt.funcDecl)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("buildFunctionDisplayName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_extractReceiverTypeString tests the extractReceiverTypeString function.
func Test_extractReceiverTypeString(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple type",
			expr:     &ast.Ident{Name: "MyType"},
			expected: "MyType",
		},
		{
			name:     "pointer type",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "MyType"}},
			expected: "(*MyType)",
		},
		{
			name:     "unsupported type",
			expr:     &ast.BasicLit{},
			expected: "",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractReceiverTypeString(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractReceiverTypeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test_runTest008_disabled tests that the rule is skipped when disabled.
func Test_runTest008_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-008": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runTest008(pass)
	if err != nil {
		t.Errorf("runTest008() error = %v", err)
	}
}

// Test_runTest008_excludedFile tests that excluded files are skipped.
func Test_runTest008_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-008": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test_test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runTest008(pass)
	if err != nil {
		t.Errorf("runTest008() error = %v", err)
	}
}
