package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// exitError est utilisé pour simuler os.Exit dans les tests
type exitError struct {
	code int
}

func (e exitError) Error() string {
	return ""
}

// mockExit crée un mock pour OsExit qui capture le code de sortie via panic
func mockExitInCmd(t *testing.T) (restore func()) {
	t.Helper()
	oldOsExit := OsExit

	OsExit = func(code int) {
		panic(exitError{code: code})
	}

	return func() {
		OsExit = oldOsExit
		// Reset flags
		AIMode = false
		NoColor = false
		Simple = false
		Verbose = false
		Category = ""
	}
}

// catchExit exécute une fonction et capture le code de sortie
func catchExitInCmd(t *testing.T, fn func()) (exitCode int, didExit bool) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exitError); ok {
				exitCode = e.code
				didExit = true
			} else {
				panic(r) // Re-panic si ce n'est pas notre exitError
			}
		}
	}()

	fn()
	return 0, false
}

// TestExecuteSuccess teste l'exécution réussie de Execute
func TestExecuteSuccess(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Capturer stdout pour éviter le bruit
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Simuler des arguments valides avec un chemin qui existe
	oldArgs := os.Args
	os.Args = []string{"ktn-linter", "lint", "../../pkg/formatter"}
	defer func() {
		os.Args = oldArgs
	}()

	exitCode, didExit := catchExitInCmd(t, func() {
		Execute()
	})

	w.Close()
	r.Close()

	// Vérification de la condition
	if !didExit {
		t.Error("Expected Execute() to call OsExit")
	}

	// Le code peut être 0 (succès) ou 1 (quelques warnings)
	// L'important est que le programme ne crash pas
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
	}
}

// TestExecuteError teste l'exécution avec une erreur
func TestExecuteError(t *testing.T) {
	restore := mockExitInCmd(t)
	defer restore()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	// Simuler des arguments invalides
	oldArgs := os.Args
	os.Args = []string{"ktn-linter", "invalid-command"}
	defer func() {
		os.Args = oldArgs
	}()

	exitCode, didExit := catchExitInCmd(t, func() {
		Execute()
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)

	// Vérification de la condition
	if !didExit {
		t.Error("Expected Execute() to exit on invalid command")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

// TestInitFlags teste que les flags sont correctement initialisés
func TestInitFlags(t *testing.T) {
	// Les flags sont initialisés dans init()
	// Vérifier qu'ils existent dans rootCmd
	if rootCmd.PersistentFlags().Lookup("ai") == nil {
		t.Error("Flag 'ai' not found")
	}
	if rootCmd.PersistentFlags().Lookup("no-color") == nil {
		t.Error("Flag 'no-color' not found")
	}
	if rootCmd.PersistentFlags().Lookup("simple") == nil {
		t.Error("Flag 'simple' not found")
	}
	if rootCmd.PersistentFlags().Lookup("verbose") == nil {
		t.Error("Flag 'verbose' not found")
	}
	if rootCmd.PersistentFlags().Lookup("category") == nil {
		t.Error("Flag 'category' not found")
	}
}

// TestRootCmdStructure teste la structure de la commande root
func TestRootCmdStructure(t *testing.T) {
	if rootCmd.Use != "ktn-linter" {
		t.Errorf("Expected Use='ktn-linter', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}
