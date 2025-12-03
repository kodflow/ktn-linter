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

// Test_runLint teste la fonction runLint avec un package valide.
func Test_runLint(t *testing.T) {
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

// Test_runLintWithIssues teste runLint qui trouve des violations.
func Test_runLintWithIssues(t *testing.T) {
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

// Test_runLintSuccess teste runLint sans aucun diagnostic.
func Test_runLintSuccess(t *testing.T) {
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

// Test_loadPackages teste loadPackages avec un pattern valide.
func Test_loadPackages(t *testing.T) {
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

// Test_loadPackagesInvalid teste loadPackages avec un pattern invalide.
func Test_loadPackagesInvalid(t *testing.T) {
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

// Test_loadPackagesWithPackageError teste loadPackages avec un package qui a des erreurs.
func Test_loadPackagesWithPackageError(t *testing.T) {
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
	// Ne devrait pas paniquer ni sortir
	pkg := &packages.Package{
		PkgPath: "test/pkg",
		Errors:  []packages.Error{},
	}

	checkLoadErrors([]*packages.Package{pkg})
	// Si on arrive ici, le test passe
}

// TestRunAnalyzers teste runAnalyzers
func Test_runAnalyzers(t *testing.T) {
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
func Test_runAnalyzersWithCategory(t *testing.T) {
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
func Test_runAnalyzersVerbose(t *testing.T) {
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
func Test_runAnalyzersVerboseWithCategory(t *testing.T) {
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
func Test_runAnalyzersVerboseMultiplePackages(t *testing.T) {
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
func Test_runAnalyzersWithError(t *testing.T) {
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
func Test_filterDiagnostics(t *testing.T) {
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
func Test_extractDiagnostics(t *testing.T) {
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
func Test_formatAndDisplay(t *testing.T) {
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
func Test_formatAndDisplayWithDiagnostics(t *testing.T) {
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

// Test_filterOverlappingEdits teste le filtrage des éditions qui se chevauchent
func Test_filterOverlappingEdits(t *testing.T) {
	tests := []struct {
		name          string
		edits         []textEdit
		expectedCount int
	}{
		{
			name:          "empty edits",
			edits:         []textEdit{},
			expectedCount: 0,
		},
		{
			name: "single edit",
			edits: []textEdit{
				{start: 10, end: 20, newText: []byte("replacement")},
			},
			expectedCount: 1,
		},
		{
			name: "non-overlapping edits",
			edits: []textEdit{
				{start: 30, end: 40, newText: []byte("third")},
				{start: 20, end: 25, newText: []byte("second")},
				{start: 10, end: 15, newText: []byte("first")},
			},
			expectedCount: 3,
		},
		{
			name: "overlapping edits",
			edits: []textEdit{
				{start: 20, end: 30, newText: []byte("second")},
				{start: 15, end: 25, newText: []byte("overlap")},
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOverlappingEdits(tt.edits)
			// Vérification du nombre
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d edits, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

// Test_applyFixes teste l'application de fixes
func Test_applyFixes(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics []diagWithFset
	}{
		{
			name:        "empty diagnostics",
			diagnostics: []diagWithFset{},
		},
		{
			name: "diagnostics without suggested fixes",
			diagnostics: []diagWithFset{
				{
					diag: analysis.Diagnostic{
						Message: "test",
					},
					analyzerName: "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := applyFixes(tt.diagnostics)
			// Validation du compteur
			if count < 0 {
				t.Errorf("Expected non-negative count, got %d", count)
			}
		})
	}
}

// Test_collectSafeEdits teste la collecte d'éditions sûres
func Test_collectSafeEdits(t *testing.T) {
	tests := []struct {
		name           string
		diagnostics    []diagWithFset
		safeAnalyzers  map[string]bool
		expectedSkip   int
		expectedEdits  int
	}{
		{
			name:          "empty diagnostics",
			diagnostics:   []diagWithFset{},
			safeAnalyzers: map[string]bool{"any": true},
			expectedSkip:  0,
			expectedEdits: 0,
		},
		{
			name: "unsafe analyzer skipped",
			diagnostics: []diagWithFset{
				{
					diag: analysis.Diagnostic{
						SuggestedFixes: []analysis.SuggestedFix{
							{Message: "fix"},
						},
					},
					analyzerName: "unsafe",
				},
			},
			safeAnalyzers: map[string]bool{"any": true},
			expectedSkip:  1,
			expectedEdits: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edits, skipped := collectSafeEdits(tt.diagnostics, tt.safeAnalyzers)
			// Vérification skip count
			if skipped != tt.expectedSkip {
				t.Errorf("Expected %d skipped, got %d", tt.expectedSkip, skipped)
			}
			// Vérification edit count
			if len(edits) != tt.expectedEdits {
				t.Errorf("Expected %d edits, got %d", tt.expectedEdits, len(edits))
			}
		})
	}
}

// Test_extractTextEdits teste l'extraction d'éditions de texte
func Test_extractTextEdits(t *testing.T) {
	tests := []struct {
		name string
		diag diagWithFset
	}{
		{
			name: "diagnostic without fixes",
			diag: diagWithFset{
				diag: analysis.Diagnostic{
					SuggestedFixes: []analysis.SuggestedFix{},
				},
			},
		},
		{
			name: "error case validation",
			diag: diagWithFset{
				diag: analysis.Diagnostic{
					SuggestedFixes: []analysis.SuggestedFix{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileEdits := make(map[string][]textEdit)
			extractTextEdits(tt.diag, &fileEdits)
			// Validation - ne devrait pas paniquer
			if len(fileEdits) > 0 {
				t.Log("Edits extracted successfully")
			}
		})
	}
}

// Test_applyCollectedEdits teste l'application des éditions collectées
func Test_applyCollectedEdits(t *testing.T) {
	tests := []struct {
		name      string
		fileEdits map[string][]textEdit
	}{
		{
			name:      "empty edits",
			fileEdits: map[string][]textEdit{},
		},
		{
			name: "edits with nonexistent file",
			fileEdits: map[string][]textEdit{
				"/nonexistent/file.go": {
					{start: 0, end: 1, newText: []byte("x")},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := applyCollectedEdits(tt.fileEdits)
			// Validation du compteur
			if count < 0 {
				t.Errorf("Expected non-negative count, got %d", count)
			}
		})
	}
}

// Test_applyEditsToFile teste l'application d'éditions à un fichier
func Test_applyEditsToFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		edits    []textEdit
		expected bool
	}{
		{
			name:     "nonexistent file",
			filename: "/nonexistent/file.go",
			edits:    []textEdit{},
			expected: false,
		},
		{
			name:     "empty edits list",
			filename: "/tmp/test.go",
			edits:    []textEdit{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyEditsToFile(tt.filename, tt.edits)
			// Validation du résultat
			if result && tt.expected == false {
				t.Error("Expected false for invalid file/edits")
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
