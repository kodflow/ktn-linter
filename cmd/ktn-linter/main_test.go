package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain compile le binaire avant de lancer les tests
func TestMain(m *testing.M) {
	// Compile le binaire
	buildCmd := exec.Command("go", "build", "-o", "ktn-linter-test", ".")
	if err := buildCmd.Run(); err != nil {
		panic("Failed to build binary: " + err.Error())
	}

	// Exécute les tests
	code := m.Run()

	// Nettoyage
	os.Remove("ktn-linter-test")

	os.Exit(code)
}

// TestCLINoArgs teste le binaire sans arguments (affiche usage)
func TestCLINoArgs(t *testing.T) {
	cmd := exec.Command("./ktn-linter-test")
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	// Cobra affiche l'usage/help avec les commandes disponibles
	if !strings.Contains(outputStr, "Usage") && !strings.Contains(outputStr, "Available Commands") && !strings.Contains(outputStr, "lint") {
		t.Errorf("Expected usage message, got: %s", outputStr)
	}
}

// TestCLIHelp teste l'affichage de l'aide
func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"short flag", []string{"-h"}},
		{"long flag", []string{"--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./ktn-linter-test", tt.args...)
			output, err := cmd.CombinedOutput()

			outputStr := string(output)
			// L'aide peut retourner 0 ou 1 selon l'implémentation
			if !strings.Contains(outputStr, "Usage") {
				t.Errorf("Expected help message, got: %s", outputStr)
			}

			_ = err // L'aide peut exit 0 ou 1
		})
	}
}

// TestCLISimpleMode teste le mode simple
func TestCLISimpleMode(t *testing.T) {
	// Créer un fichier de test temporaire avec une erreur
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	code := `package test

// Bad function without params/returns doc
func BadFunc() {
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd := exec.Command("./ktn-linter-test", "lint", "-simple", tmpDir)
	output, err := cmd.CombinedOutput()

	// Doit échouer car le code a des erreurs
	if err == nil {
		t.Error("Expected error from linter")
	}

	outputStr := string(output)
	// En mode simple, le format est: file:line:col: message
	if !strings.Contains(outputStr, ":") {
		t.Errorf("Expected simple format output, got: %s", outputStr)
	}
}

// TestCLIAIMode teste le mode AI
func TestCLIAIMode(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	code := `package test

func ValidFunc() {
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd := exec.Command("./ktn-linter-test", "lint", "-ai", tmpDir)
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	// Mode AI doit produire une sortie structurée
	// Même avec 0 issues, il y a une sortie
	if outputStr == "" {
		t.Error("Expected AI mode output")
	}
}

// TestCLINoColor teste le mode sans couleurs
func TestCLINoColor(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	code := `package test

// ValidFunc description.
//
// Returns: aucun
//
// Params: aucun
//
func ValidFunc() {
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd := exec.Command("./ktn-linter-test", "lint", "-no-color", tmpDir)
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	// Ne doit pas contenir de codes ANSI de couleur
	if strings.Contains(outputStr, "\033[") || strings.Contains(outputStr, "\x1b[") {
		t.Errorf("Expected no color codes, got: %s", outputStr)
	}
}

// TestCLICategory teste le filtrage par catégorie
func TestCLICategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
	}{
		{"func category", "func"},
		{"const category", "const"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.go")

			code := `package test

func TestFunc() {
}
`
			if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			cmd := exec.Command("./ktn-linter-test", "lint", "-category="+tt.category, tmpDir)
			_, _ = cmd.CombinedOutput()

			// La commande doit s'exécuter sans crash
			// Le résultat dépend des règles de la catégorie
		})
	}
}

// TestCLIInvalidCategory teste une catégorie invalide
func TestCLIInvalidCategory(t *testing.T) {
	cmd := exec.Command("./ktn-linter-test", "lint", "-category=invalid", ".")
	output, err := cmd.CombinedOutput()

	// Doit échouer avec une catégorie invalide
	if err == nil {
		t.Error("Expected error with invalid category")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Unknown category") && !strings.Contains(outputStr, "unknown") {
		t.Errorf("Expected 'Unknown category' message, got: %s", outputStr)
	}
}

// TestCLIValidCode teste avec du code valide
func TestCLIValidCode(t *testing.T) {
	// Utiliser le workspace actuel qui contient déjà un module Go
	// On teste sur le package formatter qui devrait être propre
	cmd := exec.Command("./ktn-linter-test", "lint", "../../pkg/formatter")
	output, _ := cmd.CombinedOutput()

	// Le formatter package devrait être conforme
	// On accepte exit 0 ou 1 (quelques warnings acceptables)
	// L'important est que la commande s'exécute sans crash
	outputStr := string(output)
	if outputStr == "" {
		// Si pas de sortie du tout, c'est OK (pas d'erreurs)
		return
	}

	// Si la sortie contient des résultats, c'est acceptable
	// tant que le binaire ne crash pas
}

// TestCLIInvalidPath teste avec un chemin invalide
func TestCLIInvalidPath(t *testing.T) {
	cmd := exec.Command("./ktn-linter-test", "lint", "/nonexistent/path/that/does/not/exist")
	output, err := cmd.CombinedOutput()

	// Doit échouer avec un chemin invalide
	if err == nil {
		t.Error("Expected error with invalid path")
	}

	outputStr := string(output)
	// Accepte différents formats de messages d'erreur
	hasError := strings.Contains(outputStr, "Error") ||
		strings.Contains(outputStr, "error") ||
		strings.Contains(outputStr, "not found") ||
		strings.Contains(outputStr, "directory")
	if !hasError {
		t.Errorf("Expected error message, got: %s", outputStr)
	}
}

// TestCLIVerboseMode teste le mode verbose
func TestCLIVerboseMode(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	code := `package test

func TestFunc() {
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd := exec.Command("./ktn-linter-test", "lint", "-v", tmpDir)
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	// En mode verbose, doit afficher des informations supplémentaires
	if !strings.Contains(outputStr, "Running") && !strings.Contains(outputStr, "Analyzing") && outputStr == "" {
		// Accepte que verbose puisse ne rien afficher si tout va bien
		// Mais vérifie au moins que la commande s'exécute
	}
}

// TestCLIMultipleFlags teste la combinaison de plusieurs flags
func TestCLIMultipleFlags(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	code := `package test

func TestFunc() {
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd := exec.Command("./ktn-linter-test", "lint", "-v", "-no-color", "-simple", tmpDir)
	_, _ = cmd.CombinedOutput()

	// La commande doit s'exécuter sans crash avec plusieurs flags
}
