package testhelper_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"golang.org/x/tools/go/analysis"
)

// mockTestingT is a test implementation of TestingT for external tests.
type mockTestingT struct {
	fatalfCalled bool
	errorfCalled bool
	logfCalled   bool
}

// Fatalf implements TestingT.
func (m *mockTestingT) Fatalf(format string, args ...any) {
	m.fatalfCalled = true
}

// Errorf implements TestingT.
func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errorfCalled = true
}

// Logf implements TestingT.
func (m *mockTestingT) Logf(format string, args ...any) {
	m.logfCalled = true
}

// createTempGoFileForTest creates a temp Go file for testing.
func createTempGoFileForTest(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	// Vérification erreur
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Retour du chemin
	return tmpFile
}

// TestRunAnalyzer tests the RunAnalyzer function.
func TestRunAnalyzer(t *testing.T) {
	const EXPECTED_ZERO int = 0

	tests := []struct {
		name         string
		setupFunc    func(*testing.T) (*analysis.Analyzer, string)
		expectDiag   int
		expectFatalf bool
	}{
		{
			name: "valid file with no diagnostics",
			setupFunc: func(t *testing.T) (*analysis.Analyzer, string) {
				analyzer := &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
				tmpFile := createTempGoFileForTest(t, "package test\n\nfunc Example() {}\n")
				return analyzer, tmpFile
			},
			expectDiag:   EXPECTED_ZERO,
			expectFatalf: false,
		},
		{
			name: "valid file with one diagnostic",
			setupFunc: func(t *testing.T) (*analysis.Analyzer, string) {
				analyzer := &analysis.Analyzer{
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
				tmpFile := createTempGoFileForTest(t, "package test\n\nfunc Example() {}\n")
				return analyzer, tmpFile
			},
			expectDiag:   1,
			expectFatalf: false,
		},
		{
			name: "nonexistent file triggers fatalf",
			setupFunc: func(t *testing.T) (*analysis.Analyzer, string) {
				analyzer := &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
				return analyzer, "nonexistent_file.go"
			},
			expectDiag:   EXPECTED_ZERO,
			expectFatalf: true,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			analyzer, filename := tt.setupFunc(t)
			// Vérification si fatalf attendu
			if tt.expectFatalf {
				mock := &mockTestingT{}
				testhelper.RunAnalyzer(mock, analyzer, filename)
				// Vérification Fatalf appelé
				if !mock.fatalfCalled {
					t.Error("Expected Fatalf to be called")
				}
			} else {
				diags := testhelper.RunAnalyzer(t, analyzer, filename)
				// Vérification nombre diagnostics
				if len(diags) != tt.expectDiag {
					t.Errorf("Expected %d diagnostics, got %d", tt.expectDiag, len(diags))
				}
			}
		})
	}
}

// createTestDataStructureForTest creates a testdata structure for tests.
func createTestDataStructureForTest(t *testing.T, testDir, goodContent, badContent string) {
	tmpDir := t.TempDir()
	testdataPath := filepath.Join(tmpDir, "testdata", "src", testDir)
	err := os.MkdirAll(testdataPath, 0755)
	// Vérification erreur
	if err != nil {
		t.Fatalf("Failed to create testdata structure: %v", err)
	}

	goodFile := filepath.Join(testdataPath, "good.go")
	badFile := filepath.Join(testdataPath, "bad.go")

	err = os.WriteFile(goodFile, []byte(goodContent), 0644)
	// Vérification erreur
	if err != nil {
		t.Fatalf("Failed to write good.go: %v", err)
	}

	err = os.WriteFile(badFile, []byte(badContent), 0644)
	// Vérification erreur
	if err != nil {
		t.Fatalf("Failed to write bad.go: %v", err)
	}

	// Change to tmpDir so TestGoodBad can find testdata/
	oldWd, err := os.Getwd()
	// Vérification erreur
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	// Vérification erreur
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to tmpDir: %v", err)
	}
	t.Cleanup(func() {
		// Restoration du répertoire
		_ = os.Chdir(oldWd)
	})
}

// TestTestGoodBad tests the TestGoodBad function.
func TestTestGoodBad(t *testing.T) {
	tests := []struct {
		name          string
		analyzerFunc  func() *analysis.Analyzer
		expectedBad   int
		expectErrorf  bool
		goodContent   string
		badContent    string
		testDirSuffix string
	}{
		{
			name: "valid good and bad - no errors",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test-good",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
			},
			expectedBad:   0,
			expectErrorf:  false,
			goodContent:   "package test\n\nfunc Good() {}\n",
			badContent:    "package test\n\nfunc Bad() {}\n",
			testDirSuffix: "001",
		},
		{
			name: "good with errors triggers errorf",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test-bad",
					Doc:  "Test analyzer with diagnostics",
					Run: func(pass *analysis.Pass) (any, error) {
						pass.Report(analysis.Diagnostic{
							Pos:     pass.Files[0].Package,
							Message: "test error",
						})
						return nil, nil
					},
				}
			},
			expectedBad:   0,
			expectErrorf:  true,
			goodContent:   "package test\n\nfunc Good() {}\n",
			badContent:    "package test\n\nfunc Bad() {}\n",
			testDirSuffix: "002",
		},
		{
			name: "bad with wrong count triggers errorf",
			analyzerFunc: func() *analysis.Analyzer {
				return &analysis.Analyzer{
					Name: "test-no-error",
					Doc:  "Test analyzer with no errors",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
			},
			expectedBad:   10,
			expectErrorf:  true,
			goodContent:   "package test\n\nfunc Good() {}\n",
			badContent:    "package test\n\nfunc Bad() {}\n",
			testDirSuffix: "003",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			testDir := "testgoodbad" + tt.testDirSuffix
			createTestDataStructureForTest(t, testDir, tt.goodContent, tt.badContent)

			mock := &mockTestingT{}
			testhelper.TestGoodBad(mock, tt.analyzerFunc(), testDir, tt.expectedBad)

			// Vérification résultat
			if tt.expectErrorf && !mock.errorfCalled {
				t.Error("Expected Errorf to be called")
			}
			// Vérification pas d'erreur si non attendu
			if !tt.expectErrorf && mock.errorfCalled {
				t.Error("Unexpected Errorf call")
			}
		})
	}
}

// TestTestGoodBadWithFiles tests the TestGoodBadWithFiles function.
func TestTestGoodBadWithFiles(t *testing.T) {
	tests := []struct {
		name         string
		goodContent  string
		badContent   string
		expectedBad  int
		expectErrorf bool
		testDirName  string
	}{
		{
			name:         "valid good and bad files",
			goodContent:  "package test\n\nfunc Good() {}\n",
			badContent:   "package test\n\nfunc Bad() {}\n",
			expectedBad:  0,
			expectErrorf: false,
			testDirName:  "testfiles001",
		},
		{
			name:         "bad file with wrong error count",
			goodContent:  "package test\n\nfunc Good() {}\n",
			badContent:   "package test\n\nfunc Bad() {}\n",
			expectedBad:  5,
			expectErrorf: true,
			testDirName:  "testfiles002",
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			createTestDataStructureForTest(t, tt.testDirName, tt.goodContent, tt.badContent)

			mock := &mockTestingT{}
			testAnalyzer := &analysis.Analyzer{
				Name: "test",
				Doc:  "Test analyzer",
				Run: func(pass *analysis.Pass) (any, error) {
					return nil, nil
				},
			}

			testhelper.TestGoodBadWithFiles(mock, testAnalyzer, tt.testDirName, "good.go", "bad.go", tt.expectedBad)

			// Vérification résultat
			if tt.expectErrorf && !mock.errorfCalled {
				t.Error("Expected Errorf to be called")
			}
			// Vérification pas d'erreur si non attendu
			if !tt.expectErrorf && mock.errorfCalled {
				t.Error("Unexpected Errorf call")
			}
		})
	}
}

// createPackageTestDataForTest creates package testdata structure.
func createPackageTestDataForTest(t *testing.T, pkgName, goodContent, badContent string) {
	tmpDir := t.TempDir()
	goodDir := filepath.Join(tmpDir, "testdata", "src", pkgName, "good")
	badDir := filepath.Join(tmpDir, "testdata", "src", pkgName, "bad")

	// Création des répertoires
	_ = os.MkdirAll(goodDir, 0755)
	_ = os.MkdirAll(badDir, 0755)

	// Création des fichiers
	goodFile := filepath.Join(goodDir, "code.go")
	badFile := filepath.Join(badDir, "code.go")
	_ = os.WriteFile(goodFile, []byte(goodContent), 0644)
	_ = os.WriteFile(badFile, []byte(badContent), 0644)

	// Changer le répertoire de travail temporairement
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	t.Cleanup(func() {
		// Restoration du répertoire
		_ = os.Chdir(oldWd)
	})
}

// TestRunAnalyzerOnPackage tests the RunAnalyzerOnPackage function.
func TestRunAnalyzerOnPackage(t *testing.T) {
	const EXPECTED_ZERO int = 0

	tests := []struct {
		name       string
		setupFunc  func(*testing.T) (*analysis.Analyzer, string)
		expectDiag int
	}{
		{
			name: "valid package",
			setupFunc: func(t *testing.T) (*analysis.Analyzer, string) {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.go")
				file2 := filepath.Join(tmpDir, "file2.go")
				_ = os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)
				_ = os.WriteFile(file2, []byte("package test\n\nfunc Func2() {}\n"), 0644)

				analyzer := &analysis.Analyzer{
					Name: "test",
					Doc:  "Test analyzer",
					Run: func(pass *analysis.Pass) (any, error) {
						return nil, nil
					},
				}
				return analyzer, tmpDir
			},
			expectDiag: EXPECTED_ZERO,
		},
		{
			name: "package with diagnostics",
			setupFunc: func(t *testing.T) (*analysis.Analyzer, string) {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "file1.go")
				_ = os.WriteFile(file1, []byte("package test\n\nfunc Func1() {}\n"), 0644)

				analyzer := &analysis.Analyzer{
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
				return analyzer, tmpDir
			},
			expectDiag: 1,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			analyzer, dir := tt.setupFunc(t)
			diags := testhelper.RunAnalyzerOnPackage(t, analyzer, dir)
			// Vérification nombre diagnostics
			if len(diags) != tt.expectDiag {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectDiag, len(diags))
			}
		})
	}
}

// TestTestGoodBadPackage tests the TestGoodBadPackage function.
func TestTestGoodBadPackage(t *testing.T) {
	tests := []struct {
		name           string
		pkgName        string
		goodContent    string
		badContent     string
		expectedErrors int
		expectErrorf   bool
	}{
		{
			name:           "valid packages",
			pkgName:        "testpkg1",
			goodContent:    "package testpkg1\n\nfunc Good() {}\n",
			badContent:     "package testpkg1\n\nfunc Bad() {}\n",
			expectedErrors: 0,
			expectErrorf:   false,
		},
		{
			name:           "bad package with wrong error count",
			pkgName:        "testpkg2",
			goodContent:    "package testpkg2\n\nfunc Good() {}\n",
			badContent:     "package testpkg2\n\nfunc Bad() {}\n",
			expectedErrors: 5,
			expectErrorf:   true,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			createPackageTestDataForTest(t, tt.pkgName, tt.goodContent, tt.badContent)

			mock := &mockTestingT{}
			testAnalyzer := &analysis.Analyzer{
				Name: "test",
				Doc:  "Test analyzer",
				Run: func(pass *analysis.Pass) (any, error) {
					return nil, nil
				},
			}

			testhelper.TestGoodBadPackage(mock, testAnalyzer, tt.pkgName, tt.expectedErrors)

			// Vérification résultat
			if tt.expectErrorf && !mock.errorfCalled {
				t.Error("Expected Errorf to be called")
			}
			// Vérification pas d'erreur si non attendu
			if !tt.expectErrorf && mock.errorfCalled {
				t.Error("Unexpected Errorf call")
			}
		})
	}
}
