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

// TestCheckLoadErrors teste checkLoadErrors avec des erreurs
func TestCheckLoadErrors(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Créer un package avec des erreurs
	pkg := &packages.Package{
		PkgPath: "test/pkg",
		Errors: []packages.Error{
			{Msg: "test error"},
		},
	}

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	exitCode, didExit := catchExitInCmd(t, func() {
		checkLoadErrors([]*packages.Package{pkg})
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)

	if !didExit {
		t.Error("Expected checkLoadErrors to exit with errors")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected error message in output, got: %s", output)
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
	pkgs := loadPackages([]string{"../../../pkg/formatter"})
	diagnostics := runAnalyzers(pkgs)

	// Les diagnostics peuvent être vides ou non selon le code
	// L'important est que la fonction ne panique pas
	_ = diagnostics
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
	if lintCmd.Use != "lint [packages...]" {
		t.Errorf("Expected Use='lint [packages...]', got '%s'", lintCmd.Use)
	}

	if lintCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if lintCmd.Long == "" {
		t.Error("Long description should not be empty")
	}

	if lintCmd.Args == nil {
		t.Error("Args validator should not be nil")
	}

	if lintCmd.Run == nil {
		t.Error("Run function should not be nil")
	}
}
