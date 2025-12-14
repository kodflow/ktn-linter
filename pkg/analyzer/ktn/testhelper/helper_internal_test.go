package testhelper

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
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

// Test_runAnalyzerInternal teste le comportement interne de RunAnalyzer.
// Note: Public API tests are in helper_external_test.go
func Test_runAnalyzerInternal(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedDiagCount int
	}{
		{
			name:             "valid go file with no errors",
			content:          "package test\n\nfunc Example() {}\n",
			expectedDiagCount: 0,
		},
		{
			name:             "another valid go file",
			content:          "package test\n\nfunc AnotherExample() int { return 42 }\n",
			expectedDiagCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			tmpFile := createTempGoFile(t, tt.content)
			diags := RunAnalyzer(t, testAnalyzer, tmpFile)
			// Vérification que le slice est vide (peut être nil ou vide)
			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiagCount, len(diags))
			}
		})
	}
}

// TestRunAnalyzerWithDiagnostics teste RunAnalyzer qui génère des diagnostics.
func TestRunAnalyzerWithDiagnostics(t *testing.T) {
	tests := []struct {
		name            string
		expectedCount   int
		expectedMessage string
	}{
		{name: "generates diagnostic", expectedCount: 1, expectedMessage: "test diagnostic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAnalyzer := &analysis.Analyzer{
				Name: "test",
				Doc:  "Test analyzer",
				Run: func(pass *analysis.Pass) (any, error) {
					pass.Report(analysis.Diagnostic{
						Pos:     pass.Files[0].Package,
						Message: "test diagnostic",
					})
					return nil, nil
				},
			}

			tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
			diags := RunAnalyzer(t, testAnalyzer, tmpFile)
			// Vérification combinée
			if len(diags) != tt.expectedCount ||
				(len(diags) > 0 && diags[0].Message != tt.expectedMessage) {
				t.Errorf("Expected %d diagnostics with message %q", tt.expectedCount, tt.expectedMessage)
			}
		})
	}
}

// TestRunAnalyzerWithRequiredAnalyzer teste avec un analyzer requis.
func TestRunAnalyzerWithRequiredAnalyzer(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedDiagCount int
	}{
		{
			name:             "analyzer with inspect requirement",
			content:          "package test\n\nfunc Example() {}\n",
			expectedDiagCount: 0,
		},
		{
			name:             "analyzer with inspect requirement on complex code",
			content:          "package test\n\nfunc Complex() { if true { return } }\n",
			expectedDiagCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			tmpFile := createTempGoFile(t, tt.content)
			diags := RunAnalyzer(t, testAnalyzer, tmpFile)
			// Vérification qu'aucun diagnostic n'est généré
			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiagCount, len(diags))
			}
		})
	}
}

// TestRunAnalyzerError teste le cas où l'analyzer retourne une erreur.
func TestRunAnalyzerError(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		analyzerError    string
		expectFatalfCall bool
	}{
		{
			name:             "analyzer returns error",
			content:          "package test\n\nfunc Example() {}\n",
			analyzerError:    "analyzer error",
			expectFatalfCall: true,
		},
		{
			name:             "analyzer returns different error",
			content:          "package test\n\nfunc Test() {}\n",
			analyzerError:    "another error",
			expectFatalfCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Utilisation d'un mock qui ne fail pas vraiment
			mockT := &MockTestingT{}

			// Analyzer qui retourne une erreur
			errorAnalyzer := &analysis.Analyzer{
				Name: "test-error",
				Doc:  "Test analyzer that errors",
				Run: func(pass *analysis.Pass) (any, error) {
					// Retour avec erreur
					return nil, errors.New(tt.analyzerError)
				},
			}

			// Appel qui devrait trigger Fatalf
			tmpFile := createTempGoFile(t, tt.content)
			RunAnalyzer(mockT, errorAnalyzer, tmpFile)

			// Vérification que Fatalf a été appelé
			if tt.expectFatalfCall && !mockT.FatalfCalled {
				t.Error("Expected Fatalf to be called when analyzer returns error")
			}
		})
	}
}

// TestRunAnalyzerRequiredError teste l'erreur d'un analyzer requis.
func TestRunAnalyzerRequiredError(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		errorMessage     string
		expectFatalfCall bool
	}{
		{
			name:             "required analyzer fails",
			content:          "package test\n\nfunc Example() {}\n",
			errorMessage:     "required analyzer failed",
			expectFatalfCall: true,
		},
		{
			name:             "required analyzer fails with different error",
			content:          "package test\n\nfunc Test() {}\n",
			errorMessage:     "another failure",
			expectFatalfCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &MockTestingT{}

			// Analyzer requis qui fail
			failingRequired := &analysis.Analyzer{
				Name: "failing-required",
				Doc:  "Failing required analyzer",
				Run: func(pass *analysis.Pass) (any, error) {
					// Retour avec erreur
					return nil, errors.New(tt.errorMessage)
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
			tmpFile := createTempGoFile(t, tt.content)
			RunAnalyzer(mockT, testAnalyzer, tmpFile)

			// Vérification que Fatalf a été appelé
			if tt.expectFatalfCall && !mockT.FatalfCalled {
				t.Error("Expected Fatalf to be called when required analyzer fails")
			}
		})
	}
}

// TestRunAnalyzerInvalidFile teste avec un fichier invalide.
func TestRunAnalyzerInvalidFile(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		expectFatalfCall bool
	}{
		{
			name:             "nonexistent file",
			filename:         "nonexistent_file.go",
			expectFatalfCall: true,
		},
		{
			name:             "another nonexistent file",
			filename:         "/invalid/path/file.go",
			expectFatalfCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			RunAnalyzer(mockT, testAnalyzer, tt.filename)

			// Vérification que Fatalf a été appelé
			if tt.expectFatalfCall && !mockT.FatalfCalled {
				t.Error("Expected Fatalf to be called for nonexistent file")
			}
		})
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
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to tmpDir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	})

	return tmpDir
}

// Test_testGoodBadInternal teste le comportement interne de TestGoodBad.
// Note: Public API tests are in helper_external_test.go
func Test_testGoodBadInternal(t *testing.T) {
	tests := []struct {
		name            string
		testDir         string
		goodContent     string
		badContent      string
		expectedErrors  int
		expectErrorCall bool
	}{
		{
			name:            "valid good and bad files",
			testDir:         "test001",
			goodContent:     "package test\n\nfunc Good() {}\n",
			badContent:      "package test\n\nfunc Bad() {}\n",
			expectedErrors:  0,
			expectErrorCall: false,
		},
		{
			name:            "another valid set of files",
			testDir:         "test002",
			goodContent:     "package test\n\nfunc AnotherGood() int { return 1 }\n",
			badContent:      "package test\n\nfunc AnotherBad() int { return 2 }\n",
			expectedErrors:  0,
			expectErrorCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			createTestDataStructure(t, tt.testDir, tt.goodContent, tt.badContent)

			// Test normal - ne devrait pas fail
			TestGoodBad(mockT, goodAnalyzer, tt.testDir, tt.expectedErrors)

			// Vérification qu'aucune erreur n'a été appelée
			if tt.expectErrorCall && !mockT.ErrorfCalled {
				t.Error("Expected Errorf to be called")
			}
			if !tt.expectErrorCall && mockT.ErrorfCalled {
				t.Error("Expected no errors for valid good.go and bad.go")
			}
		})
	}
}

// TestTestGoodBadWithErrors teste TestGoodBad avec mauvais nombre d'erreurs.
func TestTestGoodBadWithErrors(t *testing.T) {
	tests := []struct {
		name             string
		testDir          string
		goodContent      string
		badContent       string
		expectedErrors   int
		expectErrorCall  bool
	}{
		{
			name:             "good file has unexpected errors",
			testDir:          "test002",
			goodContent:      "package test\n\nfunc Good() {}\n",
			badContent:       "package test\n\nfunc Bad() {}\n",
			expectedErrors:   0,
			expectErrorCall:  true,
		},
		{
			name:             "another error detection case",
			testDir:          "test003",
			goodContent:      "package test\n\nfunc AnotherGood() {}\n",
			badContent:       "package test\n\nfunc AnotherBad() {}\n",
			expectedErrors:   0,
			expectErrorCall:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			createTestDataStructure(t, tt.testDir, tt.goodContent, tt.badContent)

			// Test qui devrait détecter des erreurs sur good.go
			TestGoodBad(mockT, badAnalyzer, tt.testDir, tt.expectedErrors)

			// Vérification qu'Errorf a été appelé
			if tt.expectErrorCall && !mockT.ErrorfCalled {
				t.Error("Expected Errorf to be called when good.go has unexpected errors")
			}
		})
	}
}

// TestTestGoodBadWrongErrorCount teste avec mauvais nombre d'erreurs dans bad.go.
func TestTestGoodBadWrongErrorCount(t *testing.T) {
	tests := []struct {
		name             string
		testDir          string
		goodContent      string
		badContent       string
		expectedErrors   int
		expectErrorCall  bool
	}{
		{
			name:             "bad file has wrong error count",
			testDir:          "test003",
			goodContent:      "package test\n\nfunc Good() {}\n",
			badContent:       "package test\n\nfunc Bad() {}\n",
			expectedErrors:   10,
			expectErrorCall:  true,
		},
		{
			name:             "another wrong error count case",
			testDir:          "test004",
			goodContent:      "package test\n\nfunc GoodFunc() {}\n",
			badContent:       "package test\n\nfunc BadFunc() {}\n",
			expectedErrors:   5,
			expectErrorCall:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			createTestDataStructure(t, tt.testDir, tt.goodContent, tt.badContent)

			// Test qui attend des erreurs mais n'en aura 0
			TestGoodBad(mockT, noErrorAnalyzer, tt.testDir, tt.expectedErrors)

			// Vérification qu'Errorf a été appelé
			if tt.expectErrorCall && !mockT.ErrorfCalled {
				t.Error("Expected Errorf to be called when bad.go has wrong error count")
			}
		})
	}
}

// TestRunAnalyzerWithReadFile teste que ReadFile fonctionne.
func TestRunAnalyzerWithReadFile(t *testing.T) {
	const EXPECTED_DIAG_COUNT int = 0

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "ReadFile succeeds",
			check: func(t *testing.T) {
				var readErr error
				// Analyzer qui utilise ReadFile
				readFileAnalyzer := &analysis.Analyzer{
					Name: "test-readfile",
					Doc:  "Test analyzer that uses ReadFile",
					Run: func(pass *analysis.Pass) (any, error) {
						// Lecture du fichier
						_, readErr = pass.ReadFile(pass.Fset.Position(pass.Files[0].Pos()).Filename)
						// Retour de la fonction
						return nil, nil
					},
				}

				// Test avec un fichier Go valide
				tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
				RunAnalyzer(t, readFileAnalyzer, tmpFile)
				// Vérification pas d'erreur
				if readErr != nil {
					t.Errorf("ReadFile failed: %v", readErr)
				}
			},
		},
		{
			name: "ReadFile returns non-empty content",
			check: func(t *testing.T) {
				var content []byte
				// Analyzer qui utilise ReadFile
				readFileAnalyzer := &analysis.Analyzer{
					Name: "test-readfile-content",
					Doc:  "Test analyzer that checks content",
					Run: func(pass *analysis.Pass) (any, error) {
						// Lecture du fichier
						content, _ = pass.ReadFile(pass.Fset.Position(pass.Files[0].Pos()).Filename)
						// Retour de la fonction
						return nil, nil
					},
				}

				// Test avec un fichier Go valide
				tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
				RunAnalyzer(t, readFileAnalyzer, tmpFile)
				// Vérification contenu non vide
				if len(content) == 0 {
					t.Error("ReadFile returned empty content")
				}
			},
		},
		{
			name: "No diagnostics generated",
			check: func(t *testing.T) {
				// Analyzer qui utilise ReadFile
				readFileAnalyzer := &analysis.Analyzer{
					Name: "test-readfile-no-diag",
					Doc:  "Test analyzer with no diagnostics",
					Run: func(pass *analysis.Pass) (any, error) {
						// Lecture du fichier
						_, _ = pass.ReadFile(pass.Fset.Position(pass.Files[0].Pos()).Filename)
						// Retour de la fonction
						return nil, nil
					},
				}

				// Test avec un fichier Go valide
				tmpFile := createTempGoFile(t, "package test\n\nfunc Example() {}\n")
				diags := RunAnalyzer(t, readFileAnalyzer, tmpFile)
				// Vérification qu'aucun diagnostic n'est généré
				if len(diags) != EXPECTED_DIAG_COUNT {
					t.Errorf("Expected %d diagnostics, got %d", EXPECTED_DIAG_COUNT, len(diags))
				}
			},
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}

// TestRunAnalyzerWithTypeError teste que Error callback est appelé.
func TestRunAnalyzerWithTypeError(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedDiagCount int
	}{
		{
			name: "type error ignored",
			content: `package test

const VALID_CONST string = "valid"

// Type error: cannot use string as int
var wrongType int = "this is a string"
`,
			expectedDiagCount: 0,
		},
		{
			name: "another type error ignored",
			content: `package test

var anotherWrongType bool = 42
`,
			expectedDiagCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			tmpFile := createTempGoFile(t, tt.content)
			diags := RunAnalyzer(t, testAnalyzer, tmpFile)
			// Le test ne devrait pas planter même avec une erreur de type
			// car Error callback ignore les erreurs
			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics (type errors are ignored), got %d", tt.expectedDiagCount, len(diags))
			}
		})
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
func (m *MockTestingT) Fatalf(format string, args ...any) {
	m.FatalfCalled = true
}

// Errorf enregistre l'appel.
func (m *MockTestingT) Errorf(format string, args ...any) {
	m.ErrorfCalled = true
}

// Logf enregistre l'appel.
func (m *MockTestingT) Logf(format string, args ...any) {
	m.LogfCalled = true
}

// Test_testGoodBadWithFilesInternal teste le comportement interne de TestGoodBadWithFiles.
// Note: Public API tests are in helper_external_test.go
func Test_testGoodBadWithFilesInternal(t *testing.T) {
	tests := []struct {
		name         string
		goodContent  string
		badContent   string
		expectedBad  int
		shouldFail   bool
		errorMessage string
	}{
		{
			name:         "valid good and bad files",
			goodContent:  "package test\n\nfunc Good() {}\n",
			badContent:   "package test\n\nfunc Bad() {}\n",
			expectedBad:  0,
			shouldFail:   false,
			errorMessage: "",
		},
		{
			name:         "good file with errors",
			goodContent:  "package test\n\nfunc Good() {}\n",
			badContent:   "package test\n\nfunc Bad() {}\n",
			expectedBad:  0,
			shouldFail:   true,
			errorMessage: "good.go should have no errors",
		},
		{
			name:         "bad file with wrong error count",
			goodContent:  "package test\n\nfunc Good() {}\n",
			badContent:   "package test\n\nfunc Bad() {}\n",
			expectedBad:  5,
			shouldFail:   true,
			errorMessage: "bad.go should have expected errors",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			mockT := &MockTestingT{}

			// Créer analyzer approprié pour le cas de test
			var testAnalyzer *analysis.Analyzer
			// Vérification du cas de test
			if tt.errorMessage == "good.go should have no errors" {
				// Analyzer qui génère toujours des erreurs
				testAnalyzer = &analysis.Analyzer{
					Name: "test-errors",
					Doc:  "Test analyzer with errors",
					Run: func(pass *analysis.Pass) (any, error) {
						// Génération d'un diagnostic
						pass.Report(analysis.Diagnostic{
							Pos:     pass.Files[0].Package,
							Message: "test error",
						})
						// Retour de la fonction
						return nil, nil
					},
				}
			} else {
				// Analyzer sans erreurs
				testAnalyzer = &analysis.Analyzer{
					Name: "test-no-errors",
					Doc:  "Test analyzer without errors",
					Run: func(pass *analysis.Pass) (any, error) {
						// Retour de la fonction
						return nil, nil
					},
				}
			}

			// Créer structure testdata
			createTestDataStructure(t, "testfiles", tt.goodContent, tt.badContent)

			// Appel fonction avec bons paramètres (testDir, goodFilename, badFilename, expectedBad)
			TestGoodBadWithFiles(mockT, testAnalyzer, "testfiles", "good.go", "bad.go", tt.expectedBad)

			// Vérification résultat
			if tt.shouldFail && !mockT.ErrorfCalled {
				t.Errorf("Expected Errorf to be called: %s", tt.errorMessage)
			}
			// Vérification pas d'erreur si succès attendu
			if !tt.shouldFail && mockT.ErrorfCalled {
				t.Error("Expected no errors but Errorf was called")
			}
		})
	}
}

// TestParsePackageFiles teste la fonction parsePackageFiles.
func TestParsePackageFiles(t *testing.T) {
	const EXPECTED_FILE_COUNT_TWO int = 2
	const EXPECTED_FILE_COUNT_ONE int = 1

	tests := []struct {
		name       string
		setupFunc  func(*testing.T) string
		useMock    bool
		expectFail bool
		fileCount  int
	}{
		{
			name: "valid package with multiple files",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.go")
				file2 := filepath.Join(tmpDir, "file2.go")
				// Écriture des fichiers
				_ = os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)
				_ = os.WriteFile(file2, []byte("package test\n\nfunc Func2() {}\n"), 0644)
				return tmpDir
			},
			useMock:    false,
			expectFail: false,
			fileCount:  EXPECTED_FILE_COUNT_TWO,
		},
		{
			name: "directory read error",
			setupFunc: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			useMock:    true,
			expectFail: true,
			fileCount:  0,
		},
		{
			name: "no go files",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				txtFile := filepath.Join(tmpDir, "readme.txt")
				_ = os.WriteFile(txtFile, []byte("Not a go file"), 0644)
				return tmpDir
			},
			useMock:    true,
			expectFail: true,
			fileCount:  0,
		},
		{
			name: "parse error",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				invalidFile := filepath.Join(tmpDir, "invalid.go")
				_ = os.WriteFile(invalidFile, []byte("package test\n\nfunc InvalidSyntax( {}\n"), 0644)
				return tmpDir
			},
			useMock:    true,
			expectFail: true,
			fileCount:  0,
		},
		{
			name: "skip subdirectories",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				validFile := filepath.Join(tmpDir, "valid.go")
				_ = os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
				subDir := filepath.Join(tmpDir, "subdir")
				_ = os.Mkdir(subDir, 0755)
				return tmpDir
			},
			useMock:    false,
			expectFail: false,
			fileCount:  EXPECTED_FILE_COUNT_ONE,
		},
		{
			name: "skip non-go files",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				goFile := filepath.Join(tmpDir, "code.go")
				txtFile := filepath.Join(tmpDir, "readme.txt")
				mdFile := filepath.Join(tmpDir, "doc.md")
				_ = os.WriteFile(goFile, []byte("package test\n\nfunc Code() {}\n"), 0644)
				_ = os.WriteFile(txtFile, []byte("readme"), 0644)
				_ = os.WriteFile(mdFile, []byte("# doc"), 0644)
				return tmpDir
			},
			useMock:    false,
			expectFail: false,
			fileCount:  EXPECTED_FILE_COUNT_ONE,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc(t)
			fset := token.NewFileSet()

			// Vérification si mock nécessaire
			if tt.useMock {
				mock := &MockTestingT{}
				parsePackageFiles(mock, dir, fset)
				// Vérification que Fatalf a été appelé
				if tt.expectFail && !mock.FatalfCalled {
					t.Error("Expected Fatalf to be called")
				}
			} else {
				files := parsePackageFiles(t, dir, fset)
				// Vérification du nombre de fichiers
				if len(files) != tt.fileCount {
					t.Errorf("Expected %d files, got %d", tt.fileCount, len(files))
				}
			}
		})
	}
}

// Test_runAnalyzerOnPackageInternal teste le comportement interne de RunAnalyzerOnPackage.
// Note: Public API tests are in helper_external_test.go
func Test_runAnalyzerOnPackageInternal(t *testing.T) {
	const EXPECTED_DIAG_COUNT int = 0

	tests := []struct {
		name           string
		createAnalyzer func(*bool) *analysis.Analyzer
		setupFiles     func(*testing.T) string
		useMock        bool
		expectFatal    bool
		expectDiagZero bool
		fileReadFlag   *bool
	}{
		{
			name: "valid package",
			createAnalyzer: func(_ *bool) *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
			},
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.go")
				file2 := filepath.Join(tmpDir, "file2.go")
				_ = os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)
				_ = os.WriteFile(file2, []byte("package test\n\nfunc Func2() {}\n"), 0644)
				return tmpDir
			},
			useMock:        false,
			expectFatal:    false,
			expectDiagZero: true,
			fileReadFlag:   nil,
		},
		{
			name: "required analyzer error",
			createAnalyzer: func(_ *bool) *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
					Requires: []*analysis.Analyzer{
						{
							Name: "required",
							Doc:  "Required analyzer that fails",
							Run: func(pass *analysis.Pass) (any, error) {
								return nil, errors.New("required analyzer failed")
							},
						},
					},
				}
			},
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				validFile := filepath.Join(tmpDir, "valid.go")
				_ = os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
				return tmpDir
			},
			useMock:        true,
			expectFatal:    true,
			expectDiagZero: false,
			fileReadFlag:   nil,
		},
		{
			name: "analyzer error",
			createAnalyzer: func(_ *bool) *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer that fails",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, errors.New("analyzer failed")
					},
				}
			},
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				validFile := filepath.Join(tmpDir, "valid.go")
				_ = os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
				return tmpDir
			},
			useMock:        true,
			expectFatal:    true,
			expectDiagZero: false,
			fileReadFlag:   nil,
		},
		{
			name: "analyzer uses ReadFile",
			createAnalyzer: func(fileRead *bool) *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer that reads file",
					Run: func(pass *analysis.Pass) (any, error) {
						// Lire le fichier source via pass.ReadFile
						for _, file := range pass.Files {
							filename := pass.Fset.File(file.Pos()).Name()
							_, err := pass.ReadFile(filename)
							// Vérification de l'erreur
							if err == nil {
								*fileRead = true
							}
						}
						return nil, nil
					},
				}
			},
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.go")
				_ = os.WriteFile(file1, []byte("package test\n\nfunc TestFunc() {}\n"), 0644)
				return tmpDir
			},
			useMock:        false,
			expectFatal:    false,
			expectDiagZero: true,
			fileReadFlag:   new(bool),
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := tt.setupFiles(t)
			analyzer := tt.createAnalyzer(tt.fileReadFlag)

			// Vérification si mock nécessaire
			if tt.useMock {
				mock := &MockTestingT{}
				RunAnalyzerOnPackage(mock, analyzer, tmpDir)
				// Vérification que Fatalf a été appelé
				if tt.expectFatal && !mock.FatalfCalled {
					t.Error("Expected Fatalf to be called")
				}
			} else {
				diags := RunAnalyzerOnPackage(t, analyzer, tmpDir)
				// Vérification du nombre de diagnostics
				if tt.expectDiagZero && len(diags) != EXPECTED_DIAG_COUNT {
					t.Errorf("Expected %d diagnostics, got %d", EXPECTED_DIAG_COUNT, len(diags))
				}
				// Vérification du flag ReadFile si applicable
				if tt.fileReadFlag != nil && !*tt.fileReadFlag {
					t.Error("Expected analyzer to read file via pass.ReadFile")
				}
			}
		})
	}
}

// Test_testGoodBadPackageInternal teste le comportement interne de TestGoodBadPackage.
// Note: Public API tests are in helper_external_test.go
func Test_testGoodBadPackageInternal(t *testing.T) {
	tests := []struct {
		name           string
		analyzerFunc   func() *analysis.Analyzer
		pkgName        string
		goodContent    string
		badContent     string
		expectedErrors int
		useMock        bool
		expectErrorf   bool
		expectLogf     bool
	}{
		{
			name: "valid packages",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
			},
			pkgName:        "testpkg",
			goodContent:    "package testpkg\n\nfunc Good() {}\n",
			badContent:     "package testpkg\n\nfunc Bad() {}\n",
			expectedErrors: 0,
			useMock:        false,
			expectErrorf:   false,
			expectLogf:     false,
		},
		{
			name: "bad package with errors",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						// Générer un diagnostic pour chaque fonction
						for _, file := range pass.Files {
							// Itération sur les déclarations
							for _, decl := range file.Decls {
								// Vérification du type
								if funcDecl, ok := decl.(*ast.FuncDecl); ok {
									pass.Reportf(funcDecl.Pos(), "test diagnostic")
								}
							}
						}
						return nil, nil
					},
					Requires: []*analysis.Analyzer{inspect.Analyzer},
				}
			},
			pkgName:        "testpkg2",
			goodContent:    "package testpkg2\n",
			badContent:     "package testpkg2\n\nimport \"go/ast\"\n\nfunc Bad(f *ast.FuncDecl) {}\n",
			expectedErrors: 1,
			useMock:        true,
			expectErrorf:   false,
			expectLogf:     false,
		},
		{
			name: "good package with unexpected errors",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						// Générer un diagnostic pour chaque fonction
						for _, file := range pass.Files {
							// Itération sur les déclarations
							for _, decl := range file.Decls {
								// Vérification du type
								if funcDecl, ok := decl.(*ast.FuncDecl); ok {
									pass.Reportf(funcDecl.Pos(), "unexpected error")
								}
							}
						}
						return nil, nil
					},
				}
			},
			pkgName:        "testpkg3",
			goodContent:    "package testpkg3\n\nfunc Good() {}\n",
			badContent:     "package testpkg3\n",
			expectedErrors: 0,
			useMock:        true,
			expectErrorf:   true,
			expectLogf:     true,
		},
		{
			name: "bad package with wrong error count",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						// Générer un diagnostic pour chaque fonction
						for _, file := range pass.Files {
							// Itération sur les déclarations
							for _, decl := range file.Decls {
								// Vérification du type
								if funcDecl, ok := decl.(*ast.FuncDecl); ok {
									pass.Reportf(funcDecl.Pos(), "error")
								}
							}
						}
						return nil, nil
					},
				}
			},
			pkgName:        "testpkg4",
			goodContent:    "package testpkg4\n",
			badContent:     "package testpkg4\n\nfunc Bad1() {}\nfunc Bad2() {}\n",
			expectedErrors: 1,
			useMock:        true,
			expectErrorf:   true,
			expectLogf:     true,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			goodDir := filepath.Join(tmpDir, "testdata", "src", tt.pkgName, "good")
			badDir := filepath.Join(tmpDir, "testdata", "src", tt.pkgName, "bad")

			// Création des répertoires
			_ = os.MkdirAll(goodDir, 0755)
			_ = os.MkdirAll(badDir, 0755)

			// Création des fichiers
			goodFile := filepath.Join(goodDir, "code.go")
			badFile := filepath.Join(badDir, "code.go")
			_ = os.WriteFile(goodFile, []byte(tt.goodContent), 0644)
			_ = os.WriteFile(badFile, []byte(tt.badContent), 0644)

			// Changer le répertoire de travail temporairement
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			// Création de l'analyzer
			testAnalyzer := tt.analyzerFunc()

			// Exécution du test
			if tt.useMock {
				mock := &MockTestingT{}
				TestGoodBadPackage(mock, testAnalyzer, tt.pkgName, tt.expectedErrors)
				// Vérifications
				if tt.expectErrorf && !mock.ErrorfCalled {
					t.Error("Expected Errorf to be called")
				}
				// Vérification pas d'Errorf si non attendu
				if !tt.expectErrorf && mock.ErrorfCalled {
					t.Error("Unexpected Errorf call")
				}
				// Vérification Logf
				if tt.expectLogf && !mock.LogfCalled {
					t.Error("Expected Logf to be called")
				}
			} else {
				TestGoodBadPackage(t, testAnalyzer, tt.pkgName, tt.expectedErrors)
			}
		})
	}
}

// Test_createTypeInfo tests the createTypeInfo private function.
func Test_createTypeInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - basic creation",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			info := createTypeInfo()
			// Vérification que l'info n'est pas nil
			if info == nil {
				t.Error("createTypeInfo() returned nil")
			}
			// Vérification que les maps sont initialisées
			if info.Types == nil {
				t.Error("Types map is nil")
			}
			// Vérification Defs
			if info.Defs == nil {
				t.Error("Defs map is nil")
			}
			// Vérification Uses
			if info.Uses == nil {
				t.Error("Uses map is nil")
			}
		})
	}
}

// Test_createPass tests the createPass private function.
func Test_createPass(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - basic pass creation",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			info := createTypeInfo()
			var diagnostics []analysis.Diagnostic
			pass := createPass(fset, file, nil, info, &diagnostics)

			// Vérification que le pass n'est pas nil
			if pass == nil {
				t.Error("createPass() returned nil")
			}
			// Vérification Fset
			if pass.Fset != fset {
				t.Error("Pass Fset not set correctly")
			}
			// Vérification Files
			if len(pass.Files) != 1 {
				t.Errorf("expected 1 file, got %d", len(pass.Files))
			}
		})
	}
}

// Test_createPassForPackage tests the createPassForPackage private function.
func Test_createPassForPackage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - basic package pass creation",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test", 0)
			// Vérification de l'erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			files := []*ast.File{file}
			var diagnostics []analysis.Diagnostic
			pass := createPassForPackage(fset, files, &diagnostics)

			// Vérification que le pass n'est pas nil
			if pass == nil {
				t.Error("createPassForPackage() returned nil")
			}
			// Vérification Fset
			if pass.Fset != fset {
				t.Error("Pass Fset not set correctly")
			}
			// Vérification Files
			if len(pass.Files) != 1 {
				t.Errorf("expected 1 file, got %d", len(pass.Files))
			}
		})
	}
}
