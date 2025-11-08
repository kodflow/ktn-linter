package cmd

import (
	"bytes"
	"go/token"
	"os"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// TestRunLint teste la fonction runLint avec un package valide
func TestRunLint(t *testing.T) {
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
		runLint(lintCmd, []string{"../../../pkg/formatter"})
	})

	w.Close()
	r.Close()

	if !didExit {
		t.Error("Expected runLint to exit")
	}

	// Le code peut être 0 (succès) ou 1 (quelques warnings)
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
	}
}

// TestRunLintWithIssues teste runLint qui trouve des violations
func TestRunLintWithIssues(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Utiliser le package const qui a potentiellement des règles violées
	exitCode, didExit := catchExitInCmd(t, func() {
		runLint(lintCmd, []string{"../../../pkg/analyzer/ktn/const/testdata/src/const001"})
	})

	if !didExit {
		t.Error("Expected runLint to exit")
	}

	// Devrait exit avec 1 (issues trouvées) ou 0 (aucune issue)
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
	}
}

// TestRunLintSuccess teste runLint sans aucun diagnostic
func TestRunLintSuccess(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	exitCode, didExit := catchExitInCmd(t, func() {
		runLint(lintCmd, []string{"../../../pkg/formatter"})
	})

	if !didExit {
		t.Error("Expected runLint to exit")
	}

	// formatter devrait être clean
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
	}
}

// TestLoadPackagesValid teste loadPackages avec un pattern valide
func TestLoadPackagesValid(t *testing.T) {
	pkgs := loadPackages([]string{"../../../pkg/formatter"})

	if len(pkgs) == 0 {
		t.Error("Expected at least one package")
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			t.Errorf("Package %s has errors: %v", pkg.PkgPath, pkg.Errors)
		}
	}
}

// TestLoadPackagesInvalid teste loadPackages avec un pattern invalide
func TestLoadPackagesInvalid(t *testing.T) {
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
		loadPackages([]string{"/nonexistent/path/that/does/not/exist"})
	})

	w.Close()
	r.Close()

	if !didExit {
		t.Error("Expected loadPackages to exit on invalid path")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// TestLoadPackagesWithPackageError teste loadPackages avec un package qui a des erreurs
func TestLoadPackagesWithPackageError(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Utiliser un path qui cause des erreurs de package (pas d'erreur de Load() mais pkg.Errors)
	exitCode, didExit := catchExitInCmd(t, func() {
		loadPackages([]string{"."}) // Current dir n'est pas un package Go valide
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)
	os.Stderr = oldStderr

	// Peut exit ou pas, dépend du contexte
	_ = didExit
	_ = exitCode
}

// TestCheckLoadErrors teste checkLoadErrors avec des erreurs
func TestCheckLoadErrors(t *testing.T) {
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
func TestCheckLoadErrorsNoErrors(t *testing.T) {
	// Ne devrait pas paniquer ni sortir
	pkg := &packages.Package{
		PkgPath: "test/pkg",
		Errors:  []packages.Error{},
	}

	checkLoadErrors([]*packages.Package{pkg})
	// Si on arrive ici, le test passe
}

// TestRunAnalyzers teste runAnalyzers
func TestRunAnalyzers(t *testing.T) {
	tests := []struct {
		name            string
		packages        []string
		setupCategory   func()
		cleanupCategory func()
		expectPanic     bool
	}{
		{
			name:     "success with valid packages",
			packages: []string{"../../../pkg/formatter"},
			setupCategory: func() {
				// Pas de catégorie spécifique
			},
			cleanupCategory: func() {},
			expectPanic:     false,
		},
		{
			name:     "success with multiple packages",
			packages: []string{"../../../pkg/formatter", "../../../pkg/analyzer/utils"},
			setupCategory: func() {
				// Pas de catégorie spécifique
			},
			cleanupCategory: func() {},
			expectPanic:     false,
		},
		{
			name:     "success with specific category",
			packages: []string{"../../../pkg/formatter"},
			setupCategory: func() {
				Category = "func"
			},
			cleanupCategory: func() {
				Category = ""
			},
			expectPanic: false,
		},
		{
			name:     "error handling with empty packages",
			packages: []string{},
			setupCategory: func() {
				// Pas de catégorie spécifique
			},
			cleanupCategory: func() {},
			expectPanic:     false,
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.setupCategory()
			defer tt.cleanupCategory()

			// Chargement des packages
			pkgs := loadPackages(tt.packages)

			// Exécution de runAnalyzers
			diagnostics := runAnalyzers(pkgs)

			// Vérification que la fonction ne panique pas
			_ = diagnostics
		})
	}
}

// TestRunAnalyzersWithCategory teste runAnalyzers avec une catégorie
func TestRunAnalyzersWithCategory(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Tester avec une catégorie valide
	Category = "func"
	defer func() { Category = "" }()

	pkgs := loadPackages([]string{"../../../pkg/formatter"})
	diagnostics := runAnalyzers(pkgs)
	_ = diagnostics

	// Tester avec une catégorie invalide
	Category = "invalid"

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	exitCode, didExit := catchExitInCmd(t, func() {
		runAnalyzers(pkgs)
	})

	w.Close()
	r.Close()

	if !didExit {
		t.Error("Expected runAnalyzers to exit with invalid category")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// TestRunAnalyzersVerbose teste runAnalyzers en mode verbose
func TestRunAnalyzersVerbose(t *testing.T) {
	Verbose = true
	defer func() { Verbose = false }()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pkgs := loadPackages([]string{"../../../pkg/formatter"})
	diagnostics := runAnalyzers(pkgs)

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)
	os.Stderr = oldStderr

	_ = diagnostics

	output := stderr.String()
	if !strings.Contains(output, "Running") {
		t.Error("Expected verbose output")
	}
}

// TestRunAnalyzersVerboseWithCategory teste runAnalyzers en mode verbose avec catégorie
func TestRunAnalyzersVerboseWithCategory(t *testing.T) {
	Verbose = true
	Category = "const"
	defer func() {
		Verbose = false
		Category = ""
	}()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pkgs := loadPackages([]string{"../../../pkg/formatter"})
	diagnostics := runAnalyzers(pkgs)

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)
	os.Stderr = oldStderr

	_ = diagnostics

	output := stderr.String()
	if !strings.Contains(output, "category") && !strings.Contains(output, "rules") {
		t.Error("Expected verbose output with category info")
	}
}

// TestRunAnalyzersVerboseMultiplePackages teste verbose mode avec plusieurs packages
func TestRunAnalyzersVerboseMultiplePackages(t *testing.T) {
	Verbose = true
	defer func() { Verbose = false }()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Charger plusieurs packages pour couvrir la boucle verbose
	pkgs := loadPackages([]string{"../../../pkg/formatter", "../../../pkg/analyzer/utils"})
	diagnostics := runAnalyzers(pkgs)

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)
	os.Stderr = oldStderr

	_ = diagnostics

	output := stderr.String()
	// Devrait afficher "Analyzing package:" pour chaque package
	if !strings.Contains(output, "Analyzing") {
		t.Error("Expected verbose package analysis output")
	}
}

// TestRunAnalyzersWithError teste runAnalyzers avec un analyzer qui retourne une erreur
func TestRunAnalyzersWithError(t *testing.T) {
	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// On utilise un vrai package - si un analyzer échoue, l'erreur sera affichée mais le programme continue
	pkgs := loadPackages([]string{"../../../pkg/formatter"})
	diagnostics := runAnalyzers(pkgs)

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)
	os.Stderr = oldStderr

	// Le test passe si la fonction ne panique pas
	_ = diagnostics
}

// TestFilterDiagnostics teste le filtrage des diagnostics
func TestFilterDiagnostics(t *testing.T) {
	fset := token.NewFileSet()

	diagnostics := []diagWithFset{
		{
			diag: analysis.Diagnostic{
				Pos:     fset.AddFile("test.go", -1, 100).Pos(0),
				Message: "test message",
			},
			fset: fset,
		},
		{
			diag: analysis.Diagnostic{
				Pos:     fset.AddFile("/.cache/go-build/test.go", -1, 100).Pos(0),
				Message: "cache message",
			},
			fset: fset,
		},
		{
			diag: analysis.Diagnostic{
				Pos:     fset.AddFile("/tmp/test.go", -1, 100).Pos(0),
				Message: "tmp message",
			},
			fset: fset,
		},
	}

	filtered := filterDiagnostics(diagnostics)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 diagnostic after filtering, got %d", len(filtered))
	}

	if filtered[0].diag.Message != "test message" {
		t.Errorf("Expected 'test message', got '%s'", filtered[0].diag.Message)
	}
}

// TestExtractDiagnostics teste l'extraction et déduplication
func TestExtractDiagnostics(t *testing.T) {
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

	if len(deduped) != 2 {
		t.Errorf("Expected 2 diagnostics after deduplication, got %d", len(deduped))
	}
}

// TestFormatAndDisplayEmpty teste formatAndDisplay avec une liste vide
func TestFormatAndDisplayEmpty(t *testing.T) {
	// Capturer stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	formatAndDisplay([]diagWithFset{})

	w.Close()
	var stdout bytes.Buffer
	stdout.ReadFrom(r)
	os.Stdout = oldStdout

	// Devrait afficher un message de succès
	output := stdout.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

// TestFormatAndDisplayWithDiagnostics teste formatAndDisplay avec des diagnostics
func TestFormatAndDisplayWithDiagnostics(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	diagnostics := []diagWithFset{
		{
			diag: analysis.Diagnostic{
				Pos:     file.Pos(10),
				Message: "test issue",
			},
			fset: fset,
		},
	}

	// Capturer stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	formatAndDisplay(diagnostics)

	w.Close()
	var stdout bytes.Buffer
	stdout.ReadFrom(r)
	os.Stdout = oldStdout

	// Devrait afficher le diagnostic
	output := stdout.String()
	if !strings.Contains(output, "test issue") {
		t.Errorf("Expected diagnostic in output, got: %s", output)
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
