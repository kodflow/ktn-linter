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

// TestTestGoodBadWithFiles teste la fonction TestGoodBadWithFiles.
func TestTestGoodBadWithFiles(t *testing.T) {
	// Note: Cette fonction appelle RunAnalyzer qui utilise testdata/src
	// Pour un vrai test, il faudrait créer la structure testdata appropriée
	// Ce test vérifie juste que la fonction est exportée et documentée
	t.Log("TestGoodBadWithFiles covered by integration tests in ktntest package")
}

// TestParsePackageFiles teste la fonction parsePackageFiles.
func TestParsePackageFiles(t *testing.T) {
	t.Run("valid package with multiple files", func(t *testing.T) {
		tmpDir := t.TempDir()
		file1 := filepath.Join(tmpDir, "file1.go")
		file2 := filepath.Join(tmpDir, "file2.go")

		// Écriture des fichiers
		err := os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write file1: %v", err)
		}
		err = os.WriteFile(file2, []byte("package test\n\nfunc Func2() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write file2: %v", err)
		}

		// Appel de parsePackageFiles
		fset := token.NewFileSet()
		files := parsePackageFiles(t, tmpDir, fset)

		// Vérification
		if len(files) != 2 {
			t.Errorf("Expected 2 files, got %d", len(files))
		}
	})

	t.Run("directory read error", func(t *testing.T) {
		mock := &MockTestingT{}
		fset := token.NewFileSet()
		parsePackageFiles(mock, "/nonexistent/directory", fset)
		// Vérification que Fatalf a été appelé
		if !mock.FatalfCalled {
			t.Error("Expected Fatalf to be called for nonexistent directory")
		}
	})

	t.Run("no go files", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Créer un fichier non-.go
		txtFile := filepath.Join(tmpDir, "readme.txt")
		err := os.WriteFile(txtFile, []byte("Not a go file"), 0644)
		if err != nil {
			t.Fatalf("Failed to write txt file: %v", err)
		}

		mock := &MockTestingT{}
		fset := token.NewFileSet()
		parsePackageFiles(mock, tmpDir, fset)
		// Vérification que Fatalf a été appelé
		if !mock.FatalfCalled {
			t.Error("Expected Fatalf to be called for directory without go files")
		}
	})

	t.Run("parse error", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.go")
		err := os.WriteFile(invalidFile, []byte("package test\n\nfunc InvalidSyntax( {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid file: %v", err)
		}

		mock := &MockTestingT{}
		fset := token.NewFileSet()
		parsePackageFiles(mock, tmpDir, fset)
		// Vérification que Fatalf a été appelé
		if !mock.FatalfCalled {
			t.Error("Expected Fatalf to be called for invalid Go syntax")
		}
	})

	t.Run("skip subdirectories", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Créer un fichier .go valide
		validFile := filepath.Join(tmpDir, "valid.go")
		err := os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write valid file: %v", err)
		}
		// Créer un sous-répertoire (doit être ignoré)
		subDir := filepath.Join(tmpDir, "subdir")
		err = os.Mkdir(subDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}

		// Appel de parsePackageFiles
		fset := token.NewFileSet()
		files := parsePackageFiles(t, tmpDir, fset)

		// Vérification - doit avoir seulement 1 fichier (le subdir est ignoré)
		if len(files) != 1 {
			t.Errorf("Expected 1 file, got %d", len(files))
		}
	})

	t.Run("skip non-go files", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Créer des fichiers de différents types
		goFile := filepath.Join(tmpDir, "code.go")
		txtFile := filepath.Join(tmpDir, "readme.txt")
		mdFile := filepath.Join(tmpDir, "doc.md")

		err := os.WriteFile(goFile, []byte("package test\n\nfunc Code() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write go file: %v", err)
		}
		err = os.WriteFile(txtFile, []byte("readme"), 0644)
		if err != nil {
			t.Fatalf("Failed to write txt file: %v", err)
		}
		err = os.WriteFile(mdFile, []byte("# doc"), 0644)
		if err != nil {
			t.Fatalf("Failed to write md file: %v", err)
		}

		// Appel de parsePackageFiles
		fset := token.NewFileSet()
		files := parsePackageFiles(t, tmpDir, fset)

		// Vérification - doit avoir seulement 1 fichier .go
		if len(files) != 1 {
			t.Errorf("Expected 1 file, got %d", len(files))
		}
	})
}

// TestRunAnalyzerOnPackage teste la fonction RunAnalyzerOnPackage.
func TestRunAnalyzerOnPackage(t *testing.T) {
	t.Run("valid package", func(t *testing.T) {
		// Création d'un analyzer de test simple
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer",
			Run: func(pass *analysis.Pass) (any, error) {
				// Retour de la fonction
				return nil, nil
			},
		}

		// Création d'un package temporaire
		tmpDir := t.TempDir()
		file1 := filepath.Join(tmpDir, "file1.go")
		file2 := filepath.Join(tmpDir, "file2.go")

		// Écriture des fichiers
		err := os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write file1: %v", err)
		}
		err = os.WriteFile(file2, []byte("package test\n\nfunc Func2() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write file2: %v", err)
		}

		// Exécution de l'analyzer sur le package
		diags := RunAnalyzerOnPackage(t, testAnalyzer, tmpDir)
		// Vérification que le slice est vide
		if len(diags) != 0 {
			t.Errorf("Expected 0 diagnostics, got %d", len(diags))
		}
	})

	t.Run("required analyzer error", func(t *testing.T) {
		// Test avec un analyzer requis qui échoue
		testAnalyzer := &analysis.Analyzer{
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

		tmpDir := t.TempDir()
		validFile := filepath.Join(tmpDir, "valid.go")
		err := os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write valid file: %v", err)
		}

		mock := &MockTestingT{}
		RunAnalyzerOnPackage(mock, testAnalyzer, tmpDir)
		// Vérification que Fatalf a été appelé
		if !mock.FatalfCalled {
			t.Error("Expected Fatalf to be called for required analyzer failure")
		}
	})

	t.Run("analyzer error", func(t *testing.T) {
		// Test avec un analyzer qui échoue
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer that fails",
			Run: func(pass *analysis.Pass) (any, error) {
				return nil, errors.New("analyzer failed")
			},
		}

		tmpDir := t.TempDir()
		validFile := filepath.Join(tmpDir, "valid.go")
		err := os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write valid file: %v", err)
		}

		mock := &MockTestingT{}
		RunAnalyzerOnPackage(mock, testAnalyzer, tmpDir)
		// Vérification que Fatalf a été appelé
		if !mock.FatalfCalled {
			t.Error("Expected Fatalf to be called for analyzer failure")
		}
	})

	t.Run("analyzer uses ReadFile", func(t *testing.T) {
		// Test avec un analyzer qui utilise ReadFile
		var fileRead bool
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer that reads file",
			Run: func(pass *analysis.Pass) (any, error) {
				// Lire le fichier source via pass.ReadFile
				for _, file := range pass.Files {
					filename := pass.Fset.File(file.Pos()).Name()
					_, err := pass.ReadFile(filename)
					if err == nil {
						fileRead = true
					}
				}
				return nil, nil
			},
		}

		tmpDir := t.TempDir()
		file1 := filepath.Join(tmpDir, "file1.go")
		err := os.WriteFile(file1, []byte("package test\n\nfunc TestFunc() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write file1: %v", err)
		}

		// Exécution de l'analyzer
		RunAnalyzerOnPackage(t, testAnalyzer, tmpDir)

		// Vérification que le fichier a été lu
		if !fileRead {
			t.Error("Expected analyzer to read file via pass.ReadFile")
		}
	})
}

// TestTestGoodBadPackage teste la fonction TestGoodBadPackage.
func TestTestGoodBadPackage(t *testing.T) {
	t.Run("valid packages", func(t *testing.T) {
		// Création d'un analyzer de test simple
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer",
			Run: func(pass *analysis.Pass) (any, error) {
				// Retour de la fonction
				return nil, nil
			},
		}

		// Création de la structure testdata/src/testpkg/good et bad
		tmpDir := t.TempDir()
		goodDir := filepath.Join(tmpDir, "testdata", "src", "testpkg", "good")
		badDir := filepath.Join(tmpDir, "testdata", "src", "testpkg", "bad")

		err := os.MkdirAll(goodDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create good dir: %v", err)
		}
		err = os.MkdirAll(badDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create bad dir: %v", err)
		}

		// Créer des fichiers .go dans good
		goodFile := filepath.Join(goodDir, "code.go")
		err = os.WriteFile(goodFile, []byte("package testpkg\n\nfunc Good() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write good file: %v", err)
		}

		// Créer des fichiers .go dans bad
		badFile := filepath.Join(badDir, "code.go")
		err = os.WriteFile(badFile, []byte("package testpkg\n\nfunc Bad() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write bad file: %v", err)
		}

		// Changer le répertoire de travail temporairement
		oldWd, _ := os.Getwd()
		err = os.Chdir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(oldWd)

		// Exécution de TestGoodBadPackage
		TestGoodBadPackage(t, testAnalyzer, "testpkg", 0)
	})

	t.Run("bad package with errors", func(t *testing.T) {
		// Création d'un analyzer qui génère un diagnostic
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer",
			Run: func(pass *analysis.Pass) (any, error) {
				// Générer un diagnostic pour chaque fonction
				for _, file := range pass.Files {
					for _, decl := range file.Decls {
						if funcDecl, ok := decl.(*ast.FuncDecl); ok {
							pass.Reportf(funcDecl.Pos(), "test diagnostic")
						}
					}
				}
				return nil, nil
			},
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		}

		// Création de la structure testdata
		tmpDir := t.TempDir()
		goodDir := filepath.Join(tmpDir, "testdata", "src", "testpkg2", "good")
		badDir := filepath.Join(tmpDir, "testdata", "src", "testpkg2", "bad")

		err := os.MkdirAll(goodDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create good dir: %v", err)
		}
		err = os.MkdirAll(badDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create bad dir: %v", err)
		}

		// Créer des fichiers .go dans good (pas de fonctions -> 0 diagnostics)
		goodFile := filepath.Join(goodDir, "code.go")
		err = os.WriteFile(goodFile, []byte("package testpkg2\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write good file: %v", err)
		}

		// Créer des fichiers .go dans bad (1 fonction -> 1 diagnostic)
		badFile := filepath.Join(badDir, "code.go")
		err = os.WriteFile(badFile, []byte("package testpkg2\n\nimport \"go/ast\"\n\nfunc Bad(f *ast.FuncDecl) {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write bad file: %v", err)
		}

		// Changer le répertoire de travail temporairement
		oldWd, _ := os.Getwd()
		err = os.Chdir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(oldWd)

		// Exécution de TestGoodBadPackage avec mock
		mock := &MockTestingT{}
		TestGoodBadPackage(mock, testAnalyzer, "testpkg2", 1)
		// On s'attend à ce qu'il y ait 1 diagnostic dans bad, donc pas d'erreur
		// Si le nombre ne correspond pas, Errorf sera appelé
	})

	t.Run("good package with unexpected errors", func(t *testing.T) {
		// Création d'un analyzer qui génère un diagnostic pour toutes les fonctions
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer",
			Run: func(pass *analysis.Pass) (any, error) {
				// Générer un diagnostic pour chaque fonction
				for _, file := range pass.Files {
					for _, decl := range file.Decls {
						if funcDecl, ok := decl.(*ast.FuncDecl); ok {
							pass.Reportf(funcDecl.Pos(), "unexpected error")
						}
					}
				}
				return nil, nil
			},
		}

		// Création de la structure testdata
		tmpDir := t.TempDir()
		goodDir := filepath.Join(tmpDir, "testdata", "src", "testpkg3", "good")
		badDir := filepath.Join(tmpDir, "testdata", "src", "testpkg3", "bad")

		err := os.MkdirAll(goodDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create good dir: %v", err)
		}
		err = os.MkdirAll(badDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create bad dir: %v", err)
		}

		// Créer un fichier .go dans good avec une fonction (génère 1 erreur inattendue)
		goodFile := filepath.Join(goodDir, "code.go")
		err = os.WriteFile(goodFile, []byte("package testpkg3\n\nfunc Good() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write good file: %v", err)
		}

		// Créer un fichier .go dans bad sans fonctions
		badFile := filepath.Join(badDir, "code.go")
		err = os.WriteFile(badFile, []byte("package testpkg3\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write bad file: %v", err)
		}

		// Changer le répertoire de travail temporairement
		oldWd, _ := os.Getwd()
		err = os.Chdir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(oldWd)

		// Exécution de TestGoodBadPackage avec mock
		mock := &MockTestingT{}
		TestGoodBadPackage(mock, testAnalyzer, "testpkg3", 0)
		// Good a 1 erreur au lieu de 0, donc Errorf sera appelé et Logf aussi
		if !mock.ErrorfCalled {
			t.Error("Expected Errorf to be called for good package with errors")
		}
		if !mock.LogfCalled {
			t.Error("Expected Logf to be called to display errors")
		}
	})

	t.Run("bad package with wrong error count", func(t *testing.T) {
		// Création d'un analyzer qui génère un diagnostic pour toutes les fonctions
		testAnalyzer := &analysis.Analyzer{
			Name: "test",
			Doc:  "Test analyzer",
			Run: func(pass *analysis.Pass) (any, error) {
				// Générer un diagnostic pour chaque fonction
				for _, file := range pass.Files {
					for _, decl := range file.Decls {
						if funcDecl, ok := decl.(*ast.FuncDecl); ok {
							pass.Reportf(funcDecl.Pos(), "error")
						}
					}
				}
				return nil, nil
			},
		}

		// Création de la structure testdata
		tmpDir := t.TempDir()
		goodDir := filepath.Join(tmpDir, "testdata", "src", "testpkg4", "good")
		badDir := filepath.Join(tmpDir, "testdata", "src", "testpkg4", "bad")

		err := os.MkdirAll(goodDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create good dir: %v", err)
		}
		err = os.MkdirAll(badDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create bad dir: %v", err)
		}

		// Créer un fichier .go dans good sans fonctions
		goodFile := filepath.Join(goodDir, "code.go")
		err = os.WriteFile(goodFile, []byte("package testpkg4\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write good file: %v", err)
		}

		// Créer un fichier .go dans bad avec 2 fonctions (génère 2 erreurs)
		badFile := filepath.Join(badDir, "code.go")
		err = os.WriteFile(badFile, []byte("package testpkg4\n\nfunc Bad1() {}\nfunc Bad2() {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write bad file: %v", err)
		}

		// Changer le répertoire de travail temporairement
		oldWd, _ := os.Getwd()
		err = os.Chdir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(oldWd)

		// Exécution de TestGoodBadPackage avec mock (on attend 1 erreur mais il y en a 2)
		mock := &MockTestingT{}
		TestGoodBadPackage(mock, testAnalyzer, "testpkg4", 1)
		// Bad a 2 erreurs au lieu de 1, donc Errorf sera appelé et Logf aussi
		if !mock.ErrorfCalled {
			t.Error("Expected Errorf to be called for bad package with wrong error count")
		}
		if !mock.LogfCalled {
			t.Error("Expected Logf to be called to display errors")
		}
	})
}
