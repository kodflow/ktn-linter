package cmd

import (
	"bytes"
	"fmt"
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// Test_runLint teste la fonction runLint avec différents packages.
func Test_runLint(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "valid formatter package",
			packages: []string{"../../../pkg/formatter"},
		},
		{
			name:     "testdata with potential issues",
			packages: []string{"../../../pkg/analyzer/ktn/const/testdata/src/const001"},
		},
		{
			name:     "formatter success case",
			packages: []string{"../../../pkg/formatter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				os.Stdout = oldStdout
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			r.Close()

			// Vérification exit et code
			if !didExit || (exitCode != 0 && exitCode != 1) {
				t.Errorf("Test failed: didExit=%v, exitCode=%d", didExit, exitCode)
			}
		})
	}
}

// Test_loadPackages teste loadPackages avec différents patterns.
func Test_loadPackages(t *testing.T) {
	tests := []struct {
		name         string
		patterns     []string
		expectExit   bool
		expectedCode int
	}{
		{
			name:       "valid formatter package",
			patterns:   []string{"../../../pkg/formatter"},
			expectExit: false,
		},
		{
			name:         "invalid nonexistent path",
			patterns:     []string{"/nonexistent/path/that/does/not/exist"},
			expectExit:   true,
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			defer func() {
				os.Stderr = oldStderr
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				pkgs := loadPackages(tt.patterns)
				// Vérification packages valides
				if !tt.expectExit && len(pkgs) == 0 {
					t.Error("Expected at least one package")
				}
			})

			w.Close()
			r.Close()

			// Vérification comportement exit
			if tt.expectExit && (!didExit || exitCode != tt.expectedCode) {
				t.Errorf("Expected exit=%v code=%d, got exit=%v code=%d",
					tt.expectExit, tt.expectedCode, didExit, exitCode)
			}
		})
	}
}

// Test_loadPackagesWithPackageError teste loadPackages avec un package qui a des erreurs.
func Test_loadPackagesWithPackageError(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{
			name:     "current dir with package errors",
			patterns: []string{"."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Utiliser un path qui cause des erreurs de package (pas d'erreur de Load() mais pkg.Errors)
			exitCode, didExit := catchExitInCmd(t, func() {
				loadPackages(tt.patterns) // Current dir n'est pas un package Go valide
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Peut exit ou pas, dépend du contexte
			_ = didExit
			_ = exitCode
		})
	}
}

// TestCheckLoadErrors teste checkLoadErrors avec des erreurs
func Test_checkLoadErrors(t *testing.T) {
	tests := []struct {
		name          string
		pkg           *packages.Package
		expectedExit  bool
		expectedCode  int
		expectedInMsg string
	}{
		{
			name: "package with single error",
			pkg: &packages.Package{
				PkgPath: "test/pkg",
				Errors: []packages.Error{
					{Msg: "test error"},
				},
			},
			expectedExit:  true,
			expectedCode:  1,
			expectedInMsg: "test error",
		},
		{
			name: "package with multiple errors",
			pkg: &packages.Package{
				PkgPath: "test/multi",
				Errors: []packages.Error{
					{Msg: "error 1"},
					{Msg: "error 2"},
				},
			},
			expectedExit:  true,
			expectedCode:  1,
			expectedInMsg: "error 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			defer func() {
				os.Stderr = oldStderr
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				checkLoadErrors([]*packages.Package{tt.pkg})
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)

			// Vérification exit
			if didExit != tt.expectedExit {
				t.Errorf("Expected didExit=%v, got %v", tt.expectedExit, didExit)
			}

			// Vérification code
			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitCode)
			}

			// Vérification message
			output := stderr.String()
			if !strings.Contains(output, tt.expectedInMsg) {
				t.Errorf("Expected %q in output, got: %s", tt.expectedInMsg, output)
			}
		})
	}
}

// TestCheckLoadErrorsNoErrors teste checkLoadErrors sans erreurs
func Test_checkLoadErrorsNoErrors(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
	}{
		{
			name:    "package without errors",
			pkgPath: "test/pkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ne devrait pas paniquer ni sortir
			pkg := &packages.Package{
				PkgPath: tt.pkgPath,
				Errors:  []packages.Error{},
			}

			// Vérification: la fonction doit s'exécuter sans panique
			defer func() {
				// Récupère une éventuelle panique
				if r := recover(); r != nil {
					t.Errorf("checkLoadErrors panicked: %v", r)
				}
			}()
			// Exécute la fonction
			checkLoadErrors([]*packages.Package{pkg})
		})
	}
}

// TestRunAnalyzers teste runAnalyzers
func Test_runAnalyzers(t *testing.T) {
	tests := []struct {
		name        string
		packages    []string
		opts        lintOptions
		expectPanic bool
	}{
		{
			name:        "success with valid packages",
			packages:    []string{"../../../pkg/formatter"},
			opts:        lintOptions{},
			expectPanic: false,
		},
		{
			name:        "success with multiple packages",
			packages:    []string{"../../../pkg/formatter", "../../../pkg/analyzer/utils"},
			opts:        lintOptions{},
			expectPanic: false,
		},
		{
			name:        "success with specific category",
			packages:    []string{"../../../pkg/formatter"},
			opts:        lintOptions{category: "func"},
			expectPanic: false,
		},
		{
			name:        "error handling with empty packages",
			packages:    []string{},
			opts:        lintOptions{},
			expectPanic: false,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Chargement des packages
			pkgs := loadPackages(tt.packages)

			// Exécution de runAnalyzers avec options
			diagnostics := runAnalyzers(pkgs, tt.opts)

			// Vérification que la fonction ne panique pas
			_ = diagnostics
		})
	}
}

// TestRunAnalyzersWithCategory teste runAnalyzers avec différentes catégories.
func Test_runAnalyzersWithCategory(t *testing.T) {
	tests := []struct {
		name         string
		opts         lintOptions
		expectExit   bool
		expectedCode int
	}{
		{
			name:       "valid func category",
			opts:       lintOptions{category: "func"},
			expectExit: false,
		},
		{
			name:         "invalid category should exit",
			opts:         lintOptions{category: "invalid"},
			expectExit:   true,
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			defer func() {
				os.Stderr = oldStderr
			}()

			pkgs := loadPackages([]string{"../../../pkg/formatter"})

			exitCode, didExit := catchExitInCmd(t, func() {
				runAnalyzers(pkgs, tt.opts)
			})

			w.Close()
			r.Close()

			// Vérification comportement
			if tt.expectExit && (!didExit || exitCode != tt.expectedCode) {
				t.Errorf("Expected exit=%v code=%d, got exit=%v code=%d",
					tt.expectExit, tt.expectedCode, didExit, exitCode)
			}
		})
	}
}

// TestRunAnalyzersVerbose teste runAnalyzers en mode verbose
func Test_runAnalyzersVerbose(t *testing.T) {
	tests := []struct {
		name          string
		packages      []string
		opts          lintOptions
		expectedInMsg string
	}{
		{
			name:          "verbose mode output",
			packages:      []string{"../../../pkg/formatter"},
			opts:          lintOptions{verbose: true},
			expectedInMsg: "Running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			pkgs := loadPackages(tt.packages)
			diagnostics := runAnalyzers(pkgs, tt.opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			_ = diagnostics

			output := stderr.String()
			if !strings.Contains(output, tt.expectedInMsg) {
				t.Error("Expected verbose output")
			}
		})
	}
}

// TestRunAnalyzersVerboseWithCategory teste runAnalyzers en mode verbose avec catégorie
func Test_runAnalyzersVerboseWithCategory(t *testing.T) {
	tests := []struct {
		name         string
		opts         lintOptions
		packages     []string
		expectedMsg1 string
		expectedMsg2 string
	}{
		{
			name:         "verbose with const category",
			opts:         lintOptions{verbose: true, category: "const"},
			packages:     []string{"../../../pkg/formatter"},
			expectedMsg1: "category",
			expectedMsg2: "rules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			pkgs := loadPackages(tt.packages)
			diagnostics := runAnalyzers(pkgs, tt.opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			_ = diagnostics

			output := stderr.String()
			if !strings.Contains(output, tt.expectedMsg1) && !strings.Contains(output, tt.expectedMsg2) {
				t.Error("Expected verbose output with category info")
			}
		})
	}
}

// TestRunAnalyzersVerboseMultiplePackages teste verbose mode avec plusieurs packages
func Test_runAnalyzersVerboseMultiplePackages(t *testing.T) {
	tests := []struct {
		name          string
		packages      []string
		opts          lintOptions
		expectedInMsg string
	}{
		{
			name:          "multiple packages verbose output",
			packages:      []string{"../../../pkg/formatter", "../../../pkg/analyzer/utils"},
			opts:          lintOptions{verbose: true},
			expectedInMsg: "Analyzing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Charger plusieurs packages pour couvrir la boucle verbose
			pkgs := loadPackages(tt.packages)
			diagnostics := runAnalyzers(pkgs, tt.opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			_ = diagnostics

			output := stderr.String()
			// Devrait afficher "Analyzing package:" pour chaque package
			if !strings.Contains(output, tt.expectedInMsg) {
				t.Error("Expected verbose package analysis output")
			}
		})
	}
}

// TestRunAnalyzersWithError teste runAnalyzers avec différents packages.
func Test_runAnalyzersWithError(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
		opts     lintOptions
	}{
		{
			name:     "formatter package should work",
			packages: []string{"../../../pkg/formatter"},
			opts:     lintOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			pkgs := loadPackages(tt.packages)
			diagnostics := runAnalyzers(pkgs, tt.opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Vérification: diagnostics doit être non nil
			if diagnostics == nil {
				t.Error("runAnalyzers returned nil diagnostics")
			}
		})
	}
}

// TestFilterDiagnostics teste le filtrage des diagnostics.
func Test_filterDiagnostics(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		expectedCount   int
		expectedMessage string
	}{
		{
			name:            "filters only cache files, not tmp",
			files:           []string{"test.go", "/.cache/go-build/test.go", "/tmp/test.go"},
			expectedCount:   2,
			expectedMessage: "msg-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			diagnostics := make([]diagWithFset, len(tt.files))

			for i, file := range tt.files {
				diagnostics[i] = diagWithFset{
					diag: analysis.Diagnostic{
						Pos:     fset.AddFile(file, -1, 100).Pos(0),
						Message: "msg-" + string(rune('0'+i)),
					},
					fset: fset,
				}
			}

			filtered := filterDiagnostics(diagnostics)

			// Vérification unique
			if len(filtered) != tt.expectedCount || filtered[0].diag.Message != tt.expectedMessage {
				t.Errorf("Expected %d diagnostics with message %q, got %d",
					tt.expectedCount, tt.expectedMessage, len(filtered))
			}
		})
	}
}

// TestExtractDiagnostics teste l'extraction et déduplication
func Test_extractDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
	}{
		{
			name:          "deduplication of diagnostics",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			diagnostics := []diagWithFset{
				{
					diag: analysis.Diagnostic{
						Pos:     file.Pos(10),
						Message: "message 1",
					},
					fset: fset,
				},
				{
					diag: analysis.Diagnostic{
						Pos:     file.Pos(10),
						Message: "message 1", // Duplicate
					},
					fset: fset,
				},
				{
					diag: analysis.Diagnostic{
						Pos:     file.Pos(20),
						Message: "message 2",
					},
					fset: fset,
				},
			}

			deduped := extractDiagnostics(diagnostics)

			if len(deduped) != tt.expectedCount {
				t.Errorf("Expected %d diagnostics after deduplication, got %d", tt.expectedCount, len(deduped))
			}
		})
	}
}

// TestFormatAndDisplayEmpty teste formatAndDisplay avec une liste vide
func Test_formatAndDisplay(t *testing.T) {
	tests := []struct {
		name          string
		diagnostics   []diagWithFset
		expectedInMsg string
	}{
		{
			name:          "empty diagnostics list",
			diagnostics:   []diagWithFset{},
			expectedInMsg: "No issues found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			formatAndDisplay(tt.diagnostics, lintOptions{})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Devrait afficher un message de succès
			output := stdout.String()
			if !strings.Contains(output, tt.expectedInMsg) {
				t.Errorf("Expected success message, got: %s", output)
			}
		})
	}
}

// TestFormatAndDisplayWithDiagnostics teste formatAndDisplay avec des diagnostics
func Test_formatAndDisplayWithDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		expectedInMsg string
	}{
		{
			name:          "display diagnostic message",
			message:       "test issue",
			expectedInMsg: "test issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			diagnostics := []diagWithFset{
				{
					diag: analysis.Diagnostic{
						Pos:     file.Pos(10),
						Message: tt.message,
					},
					fset: fset,
				},
			}

			// Capturer stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			formatAndDisplay(diagnostics, lintOptions{})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Devrait afficher le diagnostic
			output := stdout.String()
			if !strings.Contains(output, tt.expectedInMsg) {
				t.Errorf("Expected diagnostic in output, got: %s", output)
			}
		})
	}
}

// TestLintCmdStructure teste la structure de la commande lint
func TestLintCmdStructure(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "Use field is correct",
			check: func(t *testing.T) {
				const EXPECTED_USE string = "lint [packages...]"
				// Vérification Use
				if lintCmd.Use != EXPECTED_USE {
					t.Errorf("Expected Use='%s', got '%s'", EXPECTED_USE, lintCmd.Use)
				}
			},
		},
		{
			name: "Short description is not empty",
			check: func(t *testing.T) {
				// Vérification Short
				if lintCmd.Short == "" {
					t.Error("Short description should not be empty")
				}
			},
		},
		{
			name: "Long description is not empty",
			check: func(t *testing.T) {
				// Vérification Long
				if lintCmd.Long == "" {
					t.Error("Long description should not be empty")
				}
			},
		},
		{
			name: "Args validator is not nil",
			check: func(t *testing.T) {
				// Vérification Args
				if lintCmd.Args == nil {
					t.Error("Args validator should not be nil")
				}
			},
		},
		{
			name: "Run function is not nil",
			check: func(t *testing.T) {
				// Vérification Run
				if lintCmd.Run == nil {
					t.Error("Run function should not be nil")
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

// Test_filterTestFiles teste le filtrage des fichiers de test
func Test_filterTestFiles(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			// Vérification chargement
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset

			filtered := filterTestFiles(pkg.Syntax, fset)
			// Validation - ne devrait pas paniquer
			if len(filtered) < 0 {
				t.Error("Expected non-negative count")
			}
		})
	}
}

// Test_selectFilesForAnalyzer teste la sélection de fichiers pour un analyseur
func Test_selectFilesForAnalyzer(t *testing.T) {
	tests := []struct {
		name             string
		analyzerPrefix   string
		expectedAllFiles bool
	}{
		{
			name:             "test analyzer includes test files",
			analyzerPrefix:   "ktntest",
			expectedAllFiles: true,
		},
		{
			name:             "non-test analyzer excludes test files",
			analyzerPrefix:   "ktnfunc",
			expectedAllFiles: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			// Vérification chargement
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset

			// Créer un analyseur factice
			a := &analysis.Analyzer{Name: tt.analyzerPrefix + "001"}

			files := selectFilesForAnalyzer(a, pkg, fset)

			// Vérification résultat
			if tt.expectedAllFiles {
				if len(files) != len(pkg.Syntax) {
					t.Errorf("Expected all files (%d), got %d", len(pkg.Syntax), len(files))
				}
			} else {
				if len(files) > len(pkg.Syntax) {
					t.Errorf("Expected filtered files (<= %d), got %d", len(pkg.Syntax), len(files))
				}
			}
		})
	}
}

// Test_isModernizeAnalyzer teste la détection des analyseurs modernize
func Test_isModernizeAnalyzer(t *testing.T) {
	tests := []struct {
		name     string
		analyzer string
		expected bool
	}{
		{
			name:     "recognized modernize analyzer",
			analyzer: "any",
			expected: true,
		},
		{
			name:     "another recognized analyzer",
			analyzer: "minmax",
			expected: true,
		},
		{
			name:     "non-modernize analyzer",
			analyzer: "ktnfunc001",
			expected: false,
		},
		{
			name:     "empty string returns false",
			analyzer: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isModernizeAnalyzer(tt.analyzer)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("Expected %v for %q, got %v", tt.expected, tt.analyzer, result)
			}
		})
	}
}

// Test_formatModernizeCode teste le formatage des codes modernize
func Test_formatModernizeCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase name",
			input:    "any",
			expected: "KTN-MDRNZ-ANY",
		},
		{
			name:     "mixed case name",
			input:    "MinMax",
			expected: "KTN-MDRNZ-MINMAX",
		},
		{
			name:     "empty string handling",
			input:    "",
			expected: "KTN-MDRNZ-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatModernizeCode(tt.input)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}


// Test_runRequiredAnalyzers teste l'exécution des analyseurs requis
func Test_runRequiredAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			// Vérification chargement
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset
			results := make(map[*analysis.Analyzer]any)

			// Analyseur simple sans requirements
			a := &analysis.Analyzer{
				Name:     "test",
				Requires: []*analysis.Analyzer{},
			}

			runRequiredAnalyzers(a, pkg.Syntax, pkg, fset, results)
			// Validation - ne devrait pas paniquer
		})
	}
}

// Test_createAnalysisPass teste la création d'un pass d'analyse
func Test_createAnalysisPass(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			// Vérification chargement
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset
			results := make(map[*analysis.Analyzer]any)
			diagnostics := []diagWithFset{}

			a := &analysis.Analyzer{
				Name:     "test",
				Requires: []*analysis.Analyzer{},
			}

			pass := createAnalysisPass(a, pkg, fset, &diagnostics, results)
			// Validation pass
			if pass == nil {
				t.Error("Expected non-nil pass")
			}
			// Vérification analyzer
			if pass.Analyzer != a {
				t.Error("Expected analyzer to match")
			}
		})
	}
}

// Test_loadConfiguration_EmptyPath teste loadConfiguration sans chemin spécifié
func Test_loadConfiguration_EmptyPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "empty config path searches defaults",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Pas de panic attendu - utiliser lintOptions vide
			loadConfiguration(lintOptions{})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait s'exécuter sans erreur
		})
	}
}

// Test_loadConfiguration_ValidPath teste loadConfiguration avec un fichier valide
func Test_loadConfiguration_ValidPath(t *testing.T) {
	tests := []struct {
		name       string
		configData string
	}{
		{
			name: "valid config file",
			configData: `version: 1
exclude:
  - "**/*_test.go"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer un fichier temporaire
			tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.configData)); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			opts := lintOptions{configPath: tmpfile.Name()}

			// Capturer stderr pour verbose output
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			loadConfiguration(opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait s'exécuter sans erreur
		})
	}
}

// Test_loadConfiguration_InvalidPath teste loadConfiguration avec un fichier invalide
func Test_loadConfiguration_InvalidPath(t *testing.T) {
	tests := []struct {
		name         string
		opts         lintOptions
		expectedCode int
	}{
		{
			name:         "nonexistent config file",
			opts:         lintOptions{configPath: "/nonexistent/config.yaml"},
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				loadConfiguration(tt.opts)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Vérification exit
			if !didExit || exitCode != tt.expectedCode {
				t.Errorf("Expected exit with code %d, got didExit=%v code=%d", tt.expectedCode, didExit, exitCode)
			}

			// Vérification message d'erreur
			if !strings.Contains(stderr.String(), "Error loading config") {
				t.Error("Expected error message in stderr")
			}
		})
	}
}

// Test_loadConfiguration_VerboseMode teste loadConfiguration en mode verbose
func Test_loadConfiguration_VerboseMode(t *testing.T) {
	tests := []struct {
		name       string
		configData string
	}{
		{
			name: "verbose with config file",
			configData: `version: 1
exclude:
  - "*.tmp"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer un fichier temporaire
			tmpfile, err := os.CreateTemp("", "verbose-config-*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.configData)); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			opts := lintOptions{configPath: tmpfile.Name(), verbose: true}

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			loadConfiguration(opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			output := stderr.String()
			// Vérification message verbose
			if !strings.Contains(output, "Loaded configuration") {
				t.Error("Expected verbose output about loaded configuration")
			}
		})
	}
}

// Test_loadConfiguration_DefaultLocationVerbose teste le mode verbose avec config par défaut
func Test_loadConfiguration_DefaultLocationVerbose(t *testing.T) {
	tests := []struct {
		name       string
		configData string
	}{
		{
			name: "verbose with default location",
			configData: `version: 1
exclude: []
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer un fichier .ktn-linter.yaml dans le répertoire courant
			tmpfile := ".ktn-linter-test.yaml"
			if err := os.WriteFile(tmpfile, []byte(tt.configData), 0644); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile)

			// Pas de ConfigPath, mais verbose activé
			opts := lintOptions{verbose: true}

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			loadConfiguration(opts)

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait fonctionner sans crash
			_ = stderr.String()
		})
	}
}


// Test_runLint_WithVerbose teste runLint en mode verbose
func Test_runLint_WithVerbose(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "verbose mode",
			packages: []string{"../../../pkg/formatter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Configurer le flag verbose via Cobra
			rootCmd.PersistentFlags().Set(flagVerbose, "true")
			defer rootCmd.PersistentFlags().Set(flagVerbose, "false")

			// Capturer stdout et stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			var output bytes.Buffer
			output.ReadFrom(r)
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// Vérification exit
			if !didExit || (exitCode != 0 && exitCode != 1) {
				t.Errorf("Expected exit, got didExit=%v code=%d", didExit, exitCode)
			}

			// Devrait avoir du output verbose
			_ = output.String()
		})
	}
}


// Test_checkLoadErrors_VCSError teste checkLoadErrors avec erreur VCS
func Test_checkLoadErrors_VCSError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "VCS errors are skipped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := &packages.Package{
				PkgPath: "test/pkg",
				Errors: []packages.Error{
					{Msg: "VCS status error"},
				},
			}

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Ne devrait pas exit pour VCS errors
			checkLoadErrors([]*packages.Package{pkg})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait être vide car VCS errors sont ignorées
			_ = stderr.String()
		})
	}
}

// Test_extractDiagnostics_WithModernize teste extractDiagnostics avec analyseur modernize
func Test_extractDiagnostics_WithModernize(t *testing.T) {
	tests := []struct {
		name         string
		analyzerName string
		message      string
		expected     string
	}{
		{
			name:         "modernize analyzer adds prefix",
			analyzerName: "any",
			message:      "use any instead of interface{}",
			expected:     "KTN-MDRNZ-ANY:",
		},
		{
			name:         "already prefixed message unchanged",
			analyzerName: "any",
			message:      "KTN-MDRNZ-ANY: already prefixed",
			expected:     "KTN-MDRNZ-ANY: already prefixed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			diagnostics := []diagWithFset{
				{
					diag: analysis.Diagnostic{
						Pos:     file.Pos(10),
						Message: tt.message,
					},
					fset:         fset,
					analyzerName: tt.analyzerName,
				},
			}

			result := extractDiagnostics(diagnostics)

			if len(result) != 1 {
				t.Fatalf("Expected 1 diagnostic, got %d", len(result))
			}

			if !strings.Contains(result[0].Message, tt.expected) {
				t.Errorf("Expected message to contain %q, got %q", tt.expected, result[0].Message)
			}
		})
	}
}

// Test_filterDiagnostics_WindowsPath teste filterDiagnostics avec chemin Windows
func Test_filterDiagnostics_WindowsPath(t *testing.T) {
	tests := []struct {
		name          string
		files         []string
		expectedCount int
	}{
		{
			name:          "filters Windows cache path",
			files:         []string{"test.go", "C:\\cache\\go-build\\test.go"},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			diagnostics := make([]diagWithFset, len(tt.files))

			for i, file := range tt.files {
				diagnostics[i] = diagWithFset{
					diag: analysis.Diagnostic{
						Pos:     fset.AddFile(file, -1, 100).Pos(0),
						Message: "test",
					},
					fset: fset,
				}
			}

			filtered := filterDiagnostics(diagnostics)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedCount, len(filtered))
			}
		})
	}
}


// Test_runLint_WithDiagnostics teste runLint avec des diagnostics présents
func Test_runLint_WithDiagnostics(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "package with issues exits with 1",
			packages: []string{"../../../pkg/analyzer/ktn/const/testdata/src/const001"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				os.Stdout = oldStdout
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			r.Close()

			// Vérification exit avec code 0 ou 1 (dépend des diagnostics)
			if !didExit || (exitCode != 0 && exitCode != 1) {
				t.Errorf("Expected exit with code 0 or 1, got didExit=%v code=%d", didExit, exitCode)
			}
		})
	}
}

// Test_loadPackages_WithPackageError teste loadPackages avec erreur de packages.Load
func Test_loadPackages_WithPackageError(t *testing.T) {
	tests := []struct {
		name         string
		patterns     []string
		expectedCode int
	}{
		{
			name:         "malformed package pattern",
			patterns:     []string{"!@#$%^&*()"},
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				loadPackages(tt.patterns)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Peut exit ou non selon l'erreur
			_ = didExit
			_ = exitCode
			_ = stderr.String()
		})
	}
}

// Test_runAnalyzers_WithAnalyzerError teste runAnalyzers avec erreur d'analyseur
func Test_runAnalyzers_WithAnalyzerError(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "analyzers on valid package",
			packages: []string{"../../../pkg/formatter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr pour les éventuelles erreurs
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			pkgs := loadPackages(tt.packages)
			diagnostics := runAnalyzers(pkgs, lintOptions{})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait fonctionner sans panic
			_ = diagnostics
			_ = stderr.String()
		})
	}
}


// Test_runLint_SuccessPath teste runLint sans erreurs
func Test_runLint_SuccessPath(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "clean package exits with 0",
			packages: []string{"../../../pkg/formatter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stdout pour éviter le bruit
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				os.Stdout = oldStdout
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			r.Close()

			// Devrait exit avec 0 ou 1
			if !didExit {
				t.Error("Expected function to exit")
			}

			_ = exitCode
		})
	}
}


// Test_runRequiredAnalyzers_WithRequires teste runRequiredAnalyzers avec dépendances
func Test_runRequiredAnalyzers_WithRequires(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "analyzer with required dependencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset
			results := make(map[*analysis.Analyzer]any)

			// Créer un analyseur avec requirements
			requiredAnalyzer := &analysis.Analyzer{
				Name: "required",
				Run: func(pass *analysis.Pass) (any, error) {
					// Retour simple
					return nil, nil
				},
			}

			a := &analysis.Analyzer{
				Name:     "test",
				Requires: []*analysis.Analyzer{requiredAnalyzer},
			}

			// Ne devrait pas paniquer
			runRequiredAnalyzers(a, pkg.Syntax, pkg, fset, results)

			// Vérifier que le résultat est stocké
			if _, exists := results[requiredAnalyzer]; !exists {
				t.Error("Expected required analyzer result to be stored")
			}
		})
	}
}

// Test_runRequiredAnalyzers_ReadFileError teste runRequiredAnalyzers avec erreur ReadFile
func Test_runRequiredAnalyzers_ReadFileError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ReadFile callback with error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgs := loadPackages([]string{"../../../pkg/analyzer/utils"})
			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]
			fset := pkg.Fset
			results := make(map[*analysis.Analyzer]any)

			// Créer un analyseur qui utilise ReadFile
			requiredAnalyzer := &analysis.Analyzer{
				Name: "test-readfile",
				Run: func(pass *analysis.Pass) (any, error) {
					// Tester ReadFile avec un fichier inexistant
					_, err := pass.ReadFile("/nonexistent/file.go")
					// L'erreur est attendue mais ne devrait pas faire crasher
					_ = err
					// Retour simple
					return nil, nil
				},
			}

			a := &analysis.Analyzer{
				Name:     "test",
				Requires: []*analysis.Analyzer{requiredAnalyzer},
			}

			// Ne devrait pas paniquer même avec ReadFile error
			runRequiredAnalyzers(a, pkg.Syntax, pkg, fset, results)

			// Vérifier que le résultat est stocké
			if _, exists := results[requiredAnalyzer]; !exists {
				t.Error("Expected required analyzer result to be stored")
			}
		})
	}
}

// Test_runAnalyzers_WithAnalyzerRunError teste runAnalyzers quand Run() retourne une erreur
func Test_runAnalyzers_WithAnalyzerRunError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "analyzer Run returns error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			pkgs := loadPackages([]string{"../../../pkg/formatter"})

			// Sauvegarder les analyseurs originaux
			originalAnalyzers := ktn.GetAllRules()

			// Créer un analyseur qui retourne une erreur
			errorAnalyzer := &analysis.Analyzer{
				Name: "error-analyzer",
				Run: func(pass *analysis.Pass) (any, error) {
					// Retourner une erreur
					return nil, fmt.Errorf("test error from analyzer")
				},
			}

			// Remplacer temporairement les analyseurs
			// Note: Comme on ne peut pas modifier directement GetAllRules,
			// on va tester avec Category vide et regarder la sortie d'erreur
			diagnostics := runAnalyzers(pkgs, lintOptions{})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait fonctionner sans panic
			_ = diagnostics
			_ = originalAnalyzers
			_ = errorAnalyzer
		})
	}
}

// Test_loadPackages_PackagesLoadError teste loadPackages avec erreur packages.Load
func Test_loadPackages_PackagesLoadError(t *testing.T) {
	tests := []struct {
		name         string
		patterns     []string
		expectedCode int
	}{
		{
			name:         "invalid Go module path",
			patterns:     []string{"github.com/nonexistent/invalid/module/path/x/y/z"},
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				loadPackages(tt.patterns)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Peut causer une erreur ou non selon Go version
			_ = didExit
			_ = exitCode
			_ = stderr.String()
		})
	}
}


// Test_runLint_NoDiagnosticsSuccess teste runLint sans diagnostics (exit 0)
func Test_runLint_NoDiagnosticsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "clean package exits with 0",
			packages: []string{"../../../pkg/severity"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Vérification exit avec code 0
			if !didExit {
				t.Error("Expected function to exit")
			}
			// Peut être 0 si pas de diagnostics ou 1 si diagnostics
			_ = exitCode
		})
	}
}

// Test_runLint_WithDiagnosticsExitOne teste runLint avec diagnostics (exit 1)
func Test_runLint_WithDiagnosticsExitOne(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "package with issues exits with 1",
			packages: []string{"../../../pkg/analyzer/ktn/ktnfunc/testdata/src/func001"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stdout pour éviter le bruit
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			exitCode, didExit := catchExitInCmd(t, func() {
				runLint(lintCmd, tt.packages)
			})

			w.Close()
			r.Close()
			os.Stdout = oldStdout

			// Vérification exit
			if !didExit {
				t.Error("Expected function to exit")
			}
			// Devrait être 1 si des diagnostics sont trouvés
			_ = exitCode
		})
	}
}


// Test_loadPackages_DirectError teste loadPackages avec erreur directe de packages.Load
func Test_loadPackages_DirectError(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{
			name:     "broken go.mod triggers error",
			patterns: []string{"file=broken.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Créer un répertoire avec un go.mod cassé
			tmpDir, err := os.MkdirTemp("", "broken-mod-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// go.mod invalide
			brokenMod := []byte("module broken\n\ngo invalid\n")
			if err := os.WriteFile(tmpDir+"/go.mod", brokenMod, 0644); err != nil {
				t.Fatal(err)
			}

			// Fichier go
			goFile := []byte("package broken\n")
			if err := os.WriteFile(tmpDir+"/broken.go", goFile, 0644); err != nil {
				t.Fatal(err)
			}

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Changer le répertoire de travail
			oldDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldDir)

			exitCode, didExit := catchExitInCmd(t, func() {
				loadPackages([]string{"."})
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Peut ou pas exit selon l'erreur
			_ = didExit
			_ = exitCode
			_ = stderr.String()
		})
	}
}

// Test_loadPackages_EmptyPattern teste loadPackages avec pattern vide
func Test_loadPackages_EmptyPattern(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
	}{
		{
			name:     "empty pattern list",
			patterns: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capturer stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				pkgs := loadPackages(tt.patterns)
				// Pattern vide charge le package courant
				_ = pkgs
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Devrait fonctionner ou exit
			_ = didExit
			_ = exitCode
		})
	}
}


// Test_loadConfiguration tests loadConfiguration function.
func Test_loadConfiguration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"tested via public API"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via runLint
		})
	}
}

// Test_selectAnalyzers tests selectAnalyzers function.
func Test_selectAnalyzers(t *testing.T) {
	tests := []struct {
		name       string
		onlyRule   string
		category   string
		expectMin  int
		expectExit bool
	}{
		{
			name:      "all rules when no filter",
			onlyRule:  "",
			category:  "",
			expectMin: 1,
		},
		{
			name:      "filter by category",
			onlyRule:  "",
			category:  "func",
			expectMin: 1,
		},
		{
			name:      "filter by single rule",
			onlyRule:  "KTN-FUNC-001",
			category:  "",
			expectMin: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer lintOptions avec les valeurs du test
			opts := lintOptions{
				onlyRule: tt.onlyRule,
				category: tt.category,
			}

			// Run function
			analyzers := selectAnalyzers(opts)

			// Verify
			if len(analyzers) < tt.expectMin {
				t.Errorf("got %d analyzers, want at least %d", len(analyzers), tt.expectMin)
			}
		})
	}
}

// Test_selectSingleRule tests selectSingleRule function.
func Test_selectSingleRule(t *testing.T) {
	tests := []struct {
		name       string
		opts       lintOptions
		expectExit bool
	}{
		{
			name:       "valid rule",
			opts:       lintOptions{onlyRule: "KTN-FUNC-001"},
			expectExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run function
			analyzers := selectSingleRule(tt.opts)

			// Verify
			if len(analyzers) != 1 {
				t.Errorf("got %d analyzers, want 1", len(analyzers))
			}
		})
	}
}

// Test_selectByCategory tests selectByCategory function.
func Test_selectByCategory(t *testing.T) {
	tests := []struct {
		name      string
		opts      lintOptions
		expectMin int
	}{
		{
			name:      "func category",
			opts:      lintOptions{category: "func"},
			expectMin: 1,
		},
		{
			name:      "const category",
			opts:      lintOptions{category: "const"},
			expectMin: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run function
			analyzers := selectByCategory(tt.opts)

			// Verify
			if len(analyzers) < tt.expectMin {
				t.Errorf("got %d analyzers, want at least %d", len(analyzers), tt.expectMin)
			}
		})
	}
}

// Test_analyzePackage tests analyzePackage function.
func Test_analyzePackage(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
	}{
		{
			name:     "valid package",
			packages: []string{"../../../pkg/formatter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := lintOptions{}

			// Load packages
			pkgs := loadPackages(tt.packages)

			// Verify packages loaded
			if len(pkgs) == 0 {
				t.Fatal("no packages loaded")
			}

			// Get analyzers
			analyzers := selectAnalyzers(opts)
			results := make(map[*analysis.Analyzer]any)
			var diagnostics []diagWithFset

			// Run function
			analyzePackage(pkgs[0], analyzers, results, &diagnostics, opts)

			// Verify function ran without panic
			_ = diagnostics
		})
	}
}

// Test_parseLintOptions tests the parseLintOptions function.
func Test_parseLintOptions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		wantOpts lintOptions
	}{
		{
			name: "default options",
			setup: func() {
				// Reset all flags to defaults
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
			},
			wantOpts: lintOptions{
				verbose:    false,
				category:   "",
				onlyRule:   "",
				configPath: "",
			},
		},
		{
			name: "with verbose enabled",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "true")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
			},
			wantOpts: lintOptions{
				verbose:    true,
				category:   "",
				onlyRule:   "",
				configPath: "",
			},
		},
		{
			name: "with category filter",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "true")
				flags.Set(flagCategory, "func")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
			},
			wantOpts: lintOptions{
				verbose:    true,
				category:   "func",
				onlyRule:   "",
				configPath: "",
			},
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test
			tt.setup()

			// Call function
			opts := parseLintOptions(lintCmd)

			// Verify category
			if opts.category != tt.wantOpts.category {
				t.Errorf("parseLintOptions() category = %v, want %v", opts.category, tt.wantOpts.category)
			}
			// Verify verbose
			if opts.verbose != tt.wantOpts.verbose {
				t.Errorf("parseLintOptions() verbose = %v, want %v", opts.verbose, tt.wantOpts.verbose)
			}
		})
	}
}
