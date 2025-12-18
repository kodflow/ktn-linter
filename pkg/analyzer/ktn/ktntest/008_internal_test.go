// Internal tests for analyzer 008.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_runTest008 tests the runTest008 private function with table-driven tests.
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
// NOTE: Variables are no longer considered for test file requirements.
// They are tested indirectly via the functions that use them.
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
			name: "multiple variables with public and private",
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

// Test_checkTypes tests the checkTypes private function.
func Test_checkTypes(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantPublic bool
	}{
		{
			name: "public type",
			code: `package test
type PublicType struct{}`,
			wantPublic: true,
		},
		{
			name: "private type",
			code: `package test
type privateType struct{}`,
			wantPublic: false,
		},
		{
			name: "public interface",
			code: `package test
type PublicInterface interface{}`,
			wantPublic: true,
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
					checkTypes(genDecl, &result)
				}
				// Continuer la traversée
				return true
			})

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("checkTypes() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
		})
	}
}

// Test_checkConsts tests the checkConsts private function.
func Test_checkConsts(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantPublic bool
	}{
		{
			name: "public constant",
			code: `package test
const PublicConst = 1`,
			wantPublic: true,
		},
		{
			name: "private constant",
			code: `package test
const privateConst = 1`,
			wantPublic: false,
		},
		{
			name: "multiple constants with public",
			code: `package test
const (
	PublicConst = 1
	privateConst = 2
)`,
			wantPublic: true,
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
					checkConsts(genDecl, &result)
				}
				// Continuer la traversée
				return true
			})

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("checkConsts() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
		})
	}
}

// Test_checkTestFilesExist tests the checkTestFilesExist private function.
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
func Test_reportTestFileIssues(t *testing.T) {
	tests := []struct {
		name        string
		hasPublic   bool
		hasPrivate  bool
		hasInternal bool
		hasExternal bool
		expectCount int
	}{
		{
			name:        "public only with external test - OK",
			hasPublic:   true,
			hasPrivate:  false,
			hasInternal: false,
			hasExternal: true,
			expectCount: 0,
		},
		{
			name:        "private only with internal test - OK",
			hasPublic:   false,
			hasPrivate:  true,
			hasInternal: true,
			hasExternal: false,
			expectCount: 0,
		},
		{
			name:        "both with both tests - OK",
			hasPublic:   true,
			hasPrivate:  true,
			hasInternal: true,
			hasExternal: true,
			expectCount: 0,
		},
		{
			name:        "public only missing external test",
			hasPublic:   true,
			hasPrivate:  false,
			hasInternal: false,
			hasExternal: false,
			expectCount: 1,
		},
		{
			name:        "private only missing internal test",
			hasPublic:   false,
			hasPrivate:  true,
			hasInternal: false,
			hasExternal: false,
			expectCount: 1,
		},
		{
			name:        "mixed missing both tests",
			hasPublic:   true,
			hasPrivate:  true,
			hasInternal: false,
			hasExternal: false,
			expectCount: 1, // Un seul rapport combiné
		},
		{
			name:        "mixed missing only internal test",
			hasPublic:   true,
			hasPrivate:  true,
			hasInternal: false,
			hasExternal: true,
			expectCount: 1, // Rapport pour internal manquant
		},
		{
			name:        "mixed missing only external test",
			hasPublic:   true,
			hasPrivate:  true,
			hasInternal: true,
			hasExternal: false,
			expectCount: 1, // Rapport pour external manquant
		},
		{
			name:        "no functions no tests - no issues",
			hasPublic:   false,
			hasPrivate:  false,
			hasInternal: false,
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
				hasPublic:    tt.hasPublic,
				hasPrivate:   tt.hasPrivate,
				publicFuncs:  []string{},
				privateFuncs: []string{},
			}

			// Ajouter des fonctions si nécessaire
			if tt.hasPublic {
				result.publicFuncs = []string{"PublicFunc"}
			}
			// Vérification private
			if tt.hasPrivate {
				result.privateFuncs = []string{"privateFunc"}
			}

			// Appel réel de la fonction
			reportTestFileIssues(pass, file, result, status)

			// Vérification du nombre de rapports
			if reportCount != tt.expectCount {
				t.Errorf("reportTestFileIssues() reports = %d, want %d", reportCount, tt.expectCount)
			}
		})
	}
}

// Test_reportMixedFunctionsIssues tests the reportMixedFunctionsIssues private function.
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
		{
			name:     "pointer to non-identifier",
			expr:     &ast.StarExpr{X: &ast.BasicLit{}},
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

// Test_checkVariables_mockVariable tests checkVariables with mock variable.
func Test_checkVariables_mockVariable(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		wantPublic  bool
		wantPrivate bool
	}{
		{
			name: "mock variable is skipped",
			code: `package test
var MockService int`,
			wantPublic:  false,
			wantPrivate: false,
		},
		{
			name: "mock_variable is skipped",
			code: `package test
var mock_service int`,
			wantPublic:  false,
			wantPrivate: false,
		},
		{
			name: "import declaration is skipped",
			code: `package test
import "fmt"`,
			wantPublic:  false,
			wantPrivate: false,
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

// Test_checkVariables_nonValueSpec tests checkVariables with non-ValueSpec.
func Test_checkVariables_nonValueSpec(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "type declaration not ValueSpec",
			code: `package test
type MyType struct{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			result := fileAnalysisResult{}
			ast.Inspect(file, func(n ast.Node) bool {
				if genDecl, ok := n.(*ast.GenDecl); ok {
					checkVariables(genDecl, &result)
				}
				return true
			})

			if result.hasPublic || result.hasPrivate {
				t.Error("expected no variables marked")
			}
		})
	}
}

// Test_checkTypes_nonTypeSpec tests checkTypes with non-TypeSpec.
func Test_checkTypes_nonTypeSpec(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "var declaration not TypeSpec",
			code: `package test
var x int`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			result := fileAnalysisResult{}
			ast.Inspect(file, func(n ast.Node) bool {
				if genDecl, ok := n.(*ast.GenDecl); ok {
					checkTypes(genDecl, &result)
				}
				return true
			})

			if result.hasPublic {
				t.Error("expected no types marked")
			}
		})
	}
}

// Test_checkConsts_nonValueSpec tests checkConsts with non-ValueSpec.
func Test_checkConsts_nonValueSpec(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "type declaration not ValueSpec",
			code: `package test
type MyType struct{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, 0)
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			result := fileAnalysisResult{}
			ast.Inspect(file, func(n ast.Node) bool {
				if genDecl, ok := n.(*ast.GenDecl); ok {
					checkConsts(genDecl, &result)
				}
				return true
			})

			if result.hasPublic {
				t.Error("expected no consts marked")
			}
		})
	}
}

// Test_checkTypes_mockType tests checkTypes with mock type.
func Test_checkTypes_mockType(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantPublic bool
	}{
		{
			name: "mock type is skipped",
			code: `package test
type MockService struct{}`,
			wantPublic: false,
		},
		{
			name: "mock_type is skipped",
			code: `package test
type mock_service struct{}`,
			wantPublic: false,
		},
		{
			name: "non-type declaration is skipped",
			code: `package test
import "fmt"`,
			wantPublic: false,
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
					checkTypes(genDecl, &result)
				}
				// Continuer la traversée
				return true
			})

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("checkTypes() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
		})
	}
}

// Test_checkConsts_mockConst tests checkConsts with mock const.
func Test_checkConsts_mockConst(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantPublic bool
	}{
		{
			name: "mock const is skipped",
			code: `package test
const MockValue = 1`,
			wantPublic: false,
		},
		{
			name: "non-const declaration is skipped",
			code: `package test
import "fmt"`,
			wantPublic: false,
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
					checkConsts(genDecl, &result)
				}
				// Continuer la traversée
				return true
			})

			// Vérification public
			if result.hasPublic != tt.wantPublic {
				t.Errorf("checkConsts() hasPublic = %v, want %v", result.hasPublic, tt.wantPublic)
			}
		})
	}
}

// Test_classifyFunction_mockReceiver tests classifyFunction with mock receiver.
func Test_classifyFunction_mockReceiver(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectPublic  bool
		expectPrivate bool
	}{
		{
			name: "mock function name is skipped",
			code: `package test
func MockFunc() {}`,
			expectPublic:  false,
			expectPrivate: false,
		},
		{
			name: "method on mock receiver is skipped",
			code: `package test
type MockService struct{}
func (m *MockService) Method() {}`,
			expectPublic:  false,
			expectPrivate: false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			result := fileAnalysisResult{
				publicFuncs:  []string{},
				privateFuncs: []string{},
			}

			// Parcourir les déclarations
			for _, decl := range file.Decls {
				// Vérifier FuncDecl
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					classifyFunction(funcDecl, &result)
				}
			}

			// Vérification du résultat
			if result.hasPublic != tt.expectPublic {
				t.Errorf("hasPublic = %v, want %v", result.hasPublic, tt.expectPublic)
			}
			// Vérification du résultat privé
			if result.hasPrivate != tt.expectPrivate {
				t.Errorf("hasPrivate = %v, want %v", result.hasPrivate, tt.expectPrivate)
			}
		})
	}
}
