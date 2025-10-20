package testhelper

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// createTempGoFile crée un fichier Go temporaire pour les tests.
func createTempGoFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

// TestRunAnalyzer teste la fonction RunAnalyzer.
func TestRunAnalyzer(t *testing.T) {
	// Création d'un analyzer de test simple
	testAnalyzer := &analysis.Analyzer{
		Name: "test",
		Doc:  "Test analyzer",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour de la fonction
			return nil, nil
		},
	}

	// Test avec un fichier Go valide
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	diags := RunAnalyzer(t, testAnalyzer, tmpFile)
	// Vérification que le slice est vide (peut être nil ou vide)
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics, got %d", len(diags))
	}
}

// TestRunAnalyzerWithDiagnostics teste RunAnalyzer qui génère des diagnostics.
func TestRunAnalyzerWithDiagnostics(t *testing.T) {
	// Création d'un analyzer qui génère un diagnostic
	testAnalyzer := &analysis.Analyzer{
		Name: "test",
		Doc:  "Test analyzer",
		Run: func(pass *analysis.Pass) (any, error) {
			// Génération d'un diagnostic
			pass.Report(analysis.Diagnostic{
				Pos:     pass.Files[0].Package,
				Message: "test diagnostic",
			})
			// Retour de la fonction
			return nil, nil
		},
	}

	// Test avec un fichier Go valide
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	diags := RunAnalyzer(t, testAnalyzer, tmpFile)
	// Vérification qu'un diagnostic a été généré
	if len(diags) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diags))
	}
	// Vérification du message
	if len(diags) > 0 && diags[0].Message != "test diagnostic" {
		t.Errorf("Expected message 'test diagnostic', got '%s'", diags[0].Message)
	}
}

// TestRunAnalyzerWithRequiredAnalyzer teste avec un analyzer requis.
func TestRunAnalyzerWithRequiredAnalyzer(t *testing.T) {
	// Analyzer qui nécessite inspect
	testAnalyzer := &analysis.Analyzer{
		Name:     "test-with-require",
		Doc:      "Test analyzer with requires",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			// Utilisation de l'inspector
			_ = pass.ResultOf[inspect.Analyzer]
			// Retour de la fonction
			return nil, nil
		},
	}

	// Test avec un fichier Go valide
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	diags := RunAnalyzer(t, testAnalyzer, tmpFile)
	// Vérification qu'aucun diagnostic n'est généré
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics, got %d", len(diags))
	}
}

// TestRunAnalyzerError teste le cas où l'analyzer retourne une erreur.
func TestRunAnalyzerError(t *testing.T) {
	// Utilisation d'un mock qui ne fail pas vraiment
	mockT := &MockTestingT{}

	// Analyzer qui retourne une erreur
	errorAnalyzer := &analysis.Analyzer{
		Name: "test-error",
		Doc:  "Test analyzer that errors",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour avec erreur
			return nil, errors.New("analyzer error")
		},
	}

	// Appel qui devrait trigger Fatalf
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	RunAnalyzer(mockT, errorAnalyzer, tmpFile)

	// Vérification que Fatalf a été appelé
	if !mockT.FatalfCalled {
		t.Error("Expected Fatalf to be called when analyzer returns error")
	}
}

// TestRunAnalyzerRequiredError teste l'erreur d'un analyzer requis.
func TestRunAnalyzerRequiredError(t *testing.T) {
	mockT := &MockTestingT{}

	// Analyzer requis qui fail
	failingRequired := &analysis.Analyzer{
		Name: "failing-required",
		Doc:  "Failing required analyzer",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour avec erreur
			return nil, errors.New("required analyzer failed")
		},
	}

	// Analyzer qui dépend du failing
	testAnalyzer := &analysis.Analyzer{
		Name:     "test-with-failing-require",
		Doc:      "Test with failing require",
		Requires: []*analysis.Analyzer{failingRequired},
		Run: func(pass *analysis.Pass) (any, error) {
			// Ne devrait jamais être appelé
			return nil, nil
		},
	}

	// Appel qui devrait trigger Fatalf
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	RunAnalyzer(mockT, testAnalyzer, tmpFile)

	// Vérification que Fatalf a été appelé
	if !mockT.FatalfCalled {
		t.Error("Expected Fatalf to be called when required analyzer fails")
	}
}

// TestRunAnalyzerInvalidFile teste avec un fichier invalide.
func TestRunAnalyzerInvalidFile(t *testing.T) {
	mockT := &MockTestingT{}

	testAnalyzer := &analysis.Analyzer{
		Name: "test",
		Doc:  "Test analyzer",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour de la fonction
			return nil, nil
		},
	}

	// Fichier inexistant
	RunAnalyzer(mockT, testAnalyzer, "nonexistent_file.go")

	// Vérification que Fatalf a été appelé
	if !mockT.FatalfCalled {
		t.Error("Expected Fatalf to be called for nonexistent file")
	}
}

// createTestDataStructure crée une structure testdata/ temporaire pour TestGoodBad.
func createTestDataStructure(t *testing.T, testDir, goodContent, badContent string) string {
	tmpDir := t.TempDir()
	testdataPath := filepath.Join(tmpDir, "testdata", "src", testDir)
	err := os.MkdirAll(testdataPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create testdata structure: %v", err)
	}

	goodFile := filepath.Join(testdataPath, "good.go")
	badFile := filepath.Join(testdataPath, "bad.go")

	err = os.WriteFile(goodFile, []byte(goodContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write good.go: %v", err)
	}

	err = os.WriteFile(badFile, []byte(badContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write bad.go: %v", err)
	}

	// Change to tmpDir so TestGoodBad can find testdata/
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(oldWd) })

	return tmpDir
}

// TestTestGoodBad teste la fonction TestGoodBad.
func TestTestGoodBad(t *testing.T) {
	mockT := &MockTestingT{}

	// Analyzer simple qui ne génère pas de diagnostic
	goodAnalyzer := &analysis.Analyzer{
		Name: "test-good",
		Doc:  "Test analyzer",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour de la fonction
			return nil, nil
		},
	}

	// Création de la structure testdata
	createTestDataStructure(t, "test001",
		"package test\n\nfunc Good() {}\n",
		"package test\n\nfunc Bad() {}\n")

	// Test normal - ne devrait pas fail
	TestGoodBad(mockT, goodAnalyzer, "test001", 0)

	// Vérification qu'aucune erreur n'a été appelée
	if mockT.ErrorfCalled {
		t.Error("Expected no errors for valid good.go and bad.go")
	}
}

// TestTestGoodBadWithErrors teste TestGoodBad avec mauvais nombre d'erreurs.
func TestTestGoodBadWithErrors(t *testing.T) {
	mockT := &MockTestingT{}

	// Analyzer qui génère toujours 1 diagnostic
	badAnalyzer := &analysis.Analyzer{
		Name: "test-bad",
		Doc:  "Test analyzer with diagnostics",
		Run: func(pass *analysis.Pass) (any, error) {
			// Génération d'un diagnostic sur tous les fichiers
			pass.Report(analysis.Diagnostic{
				Pos:     pass.Files[0].Package,
				Message: "test error",
			})
			// Retour de la fonction
			return nil, nil
		},
	}

	// Création de la structure testdata
	createTestDataStructure(t, "test002",
		"package test\n\nfunc Good() {}\n",
		"package test\n\nfunc Bad() {}\n")

	// Test qui devrait détecter des erreurs sur good.go
	TestGoodBad(mockT, badAnalyzer, "test002", 0)

	// Vérification qu'Errorf a été appelé
	if !mockT.ErrorfCalled {
		t.Error("Expected Errorf to be called when good.go has unexpected errors")
	}
}

// TestTestGoodBadWrongErrorCount teste avec mauvais nombre d'erreurs dans bad.go.
func TestTestGoodBadWrongErrorCount(t *testing.T) {
	mockT := &MockTestingT{}

	// Analyzer qui ne génère jamais de diagnostic
	noErrorAnalyzer := &analysis.Analyzer{
		Name: "test-no-error",
		Doc:  "Test analyzer with no errors",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour de la fonction sans diagnostic
			return nil, nil
		},
	}

	// Création de la structure testdata
	createTestDataStructure(t, "test003",
		"package test\n\nfunc Good() {}\n",
		"package test\n\nfunc Bad() {}\n")

	// Test qui attend 10 erreurs mais n'en aura 0
	TestGoodBad(mockT, noErrorAnalyzer, "test003", 10)

	// Vérification qu'Errorf a été appelé
	if !mockT.ErrorfCalled {
		t.Error("Expected Errorf to be called when bad.go has wrong error count")
	}
}

// TestRunAnalyzerWithReadFile teste que ReadFile fonctionne.
func TestRunAnalyzerWithReadFile(t *testing.T) {
	// Analyzer qui utilise ReadFile
	readFileAnalyzer := &analysis.Analyzer{
		Name: "test-readfile",
		Doc:  "Test analyzer that uses ReadFile",
		Run: func(pass *analysis.Pass) (any, error) {
			// Lecture du fichier
			content, err := pass.ReadFile(pass.Fset.Position(pass.Files[0].Pos()).Filename)
			if err != nil {
				t.Errorf("ReadFile failed: %v", err)
			}
			if len(content) == 0 {
				t.Error("ReadFile returned empty content")
			}
			// Retour de la fonction
			return nil, nil
		},
	}

	// Test avec un fichier Go valide
	tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
	diags := RunAnalyzer(t, readFileAnalyzer, tmpFile)
	// Vérification qu'aucun diagnostic n'est généré
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics, got %d", len(diags))
	}
}

// TestRunAnalyzerWithTypeError teste que Error callback est appelé.
func TestRunAnalyzerWithTypeError(t *testing.T) {
	// Simple analyzer - le Error callback sera appelé lors du type checking
	testAnalyzer := &analysis.Analyzer{
		Name: "test-type-error",
		Doc:  "Test with type error",
		Run: func(pass *analysis.Pass) (any, error) {
			// Retour de la fonction
			return nil, nil
		},
	}

	// Test avec un fichier qui a une erreur de type
	// Le Error callback sera appelé pendant types.Config.Check()
	tmpFile := createTempGoFile(t, `package test

const VALID_CONST string = "valid"

// Type error: cannot use string as int
var wrongType int = "this is a string"
`)
	diags := RunAnalyzer(t, testAnalyzer, tmpFile)
	// Le test ne devrait pas planter même avec une erreur de type
	// car Error callback ignore les erreurs
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics (type errors are ignored), got %d", len(diags))
	}
}

// MockTestingT est un mock de testing.T.
type MockTestingT struct {
	FatalfCalled bool
	ErrorfCalled bool
	LogfCalled   bool
	Messages     []string
}

// Fatalf enregistre l'appel.
//
// Params:
//   - format: format du message
//   - args: arguments du message
func (m *MockTestingT) Fatalf(format string, args ...any) {
	m.FatalfCalled = true
}

// Errorf enregistre l'appel.
//
// Params:
//   - format: format du message
//   - args: arguments du message
func (m *MockTestingT) Errorf(format string, args ...any) {
	m.ErrorfCalled = true
}

// Logf enregistre l'appel.
//
// Params:
//   - format: format du message
//   - args: arguments du message
func (m *MockTestingT) Logf(format string, args ...any) {
	m.LogfCalled = true
}
