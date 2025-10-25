package testhelper

import (
	"errors"
	"go/ast"
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

// TestTestGoodBadWithFiles teste la fonction TestGoodBadWithFiles.
func TestTestGoodBadWithFiles(t *testing.T) {
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

// TestRunAnalyzerOnPackage teste la fonction RunAnalyzerOnPackage.
func TestRunAnalyzerOnPackage(t *testing.T) {
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

// TestTestGoodBadPackage teste la fonction TestGoodBadPackage.
func TestTestGoodBadPackage(t *testing.T) {
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
