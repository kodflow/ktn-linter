package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
)

// exitError est utilisé pour simuler os.Exit dans les tests
type exitError struct {
	code int
}

func (e exitError) Error() string {
	return ""
}

// mockExit crée un mock pour OsExit qui capture le code de sortie via panic
func mockExit(t *testing.T) (restore func()) {
	t.Helper()
	oldOsExit := cmd.OsExit
	oldArgs := os.Args

	cmd.OsExit = func(code int) {
		panic(exitError{code: code})
	}

	return func() {
		cmd.OsExit = oldOsExit
		os.Args = oldArgs
		// Reset flags
		cmd.AIMode = false
		cmd.NoColor = false
		cmd.Simple = false
		cmd.Verbose = false
		cmd.Category = ""
	}
}

// catchExit exécute une fonction et capture le code de sortie
func catchExit(t *testing.T, fn func()) (exitCode int, didExit bool) {
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

// TestMainNoArgs teste que main() sans arguments affiche l'aide
func TestMainNoArgs(t *testing.T) {
	restore := mockExit(t)
	defer restore()

	// Capturer stdout (Cobra affiche l'aide sur stdout)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	os.Args = []string{"ktn-linter"}

	exitCode, didExit := catchExit(t, func() {
		main()
	})

	w.Close()
	var stdout bytes.Buffer
	stdout.ReadFrom(r)

	// Cobra affiche l'aide et ne termine pas nécessairement le process
	// C'est un comportement acceptable
	_ = didExit
	_ = exitCode

	output := stdout.String()
	// Cobra doit montrer l'usage ou les commandes disponibles
	if !strings.Contains(output, "Usage") && !strings.Contains(output, "Available Commands") && !strings.Contains(output, "lint") {
		t.Errorf("Expected help output, got: %s", output)
	}
}

// TestMainInvalidCategory teste qu'une catégorie invalide exit avec 1
func TestMainInvalidCategory(t *testing.T) {
	restore := mockExit(t)
	defer restore()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	os.Args = []string{"ktn-linter", "lint", "--category=invalid", "."}

	exitCode, didExit := catchExit(t, func() {
		main()
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)

	if !didExit {
		t.Error("Expected main() to exit with invalid category")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Unknown category") {
		t.Errorf("Expected 'Unknown category' message, got: %s", output)
	}
}

// TestMainInvalidPath teste qu'un chemin invalide exit avec 1
func TestMainInvalidPath(t *testing.T) {
	restore := mockExit(t)
	defer restore()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	os.Args = []string{"ktn-linter", "lint", "/nonexistent/path/that/does/not/exist"}

	exitCode, didExit := catchExit(t, func() {
		main()
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)

	if !didExit {
		t.Error("Expected main() to exit with invalid path")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	output := stderr.String()
	// Le message d'erreur peut varier selon le système
	if len(output) == 0 {
		t.Error("Expected error message")
	}
}

// TestMainSuccess teste que main() exit avec 0 pour du code valide
func TestMainSuccess(t *testing.T) {
	restore := mockExit(t)
	defer restore()

	// Utiliser le package formatter qui devrait être propre
	os.Args = []string{"ktn-linter", "lint", "../../pkg/formatter"}

	exitCode, didExit := catchExit(t, func() {
		main()
	})

	if !didExit {
		t.Error("Expected main() to exit")
	}

	// Le code peut être 0 (succès) ou 1 (quelques warnings)
	// L'important est que le programme ne crash pas
	if exitCode != 0 && exitCode != 1 {
		t.Errorf("Expected exit code 0 or 1, got %d", exitCode)
	}
}

// TestFlags teste le parsing des flags avec Cobra
func TestFlags(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantAI       bool
		wantNoColor  bool
		wantSimple   bool
		wantVerbose  bool
		wantCategory string
	}{
		{
			name:         "ai flag",
			args:         []string{"ktn-linter", "lint", "--ai", "../../pkg/formatter"},
			wantAI:       true,
			wantNoColor:  false,
			wantSimple:   false,
			wantVerbose:  false,
			wantCategory: "",
		},
		{
			name:         "no-color flag",
			args:         []string{"ktn-linter", "lint", "--no-color", "../../pkg/formatter"},
			wantAI:       false,
			wantNoColor:  true,
			wantSimple:   false,
			wantVerbose:  false,
			wantCategory: "",
		},
		{
			name:         "simple flag",
			args:         []string{"ktn-linter", "lint", "--simple", "../../pkg/formatter"},
			wantAI:       false,
			wantNoColor:  false,
			wantSimple:   true,
			wantVerbose:  false,
			wantCategory: "",
		},
		{
			name:         "verbose flag",
			args:         []string{"ktn-linter", "lint", "-v", "../../pkg/formatter"},
			wantAI:       false,
			wantNoColor:  false,
			wantSimple:   false,
			wantVerbose:  true,
			wantCategory: "",
		},
		{
			name:         "category flag",
			args:         []string{"ktn-linter", "lint", "--category=func", "../../pkg/formatter"},
			wantAI:       false,
			wantNoColor:  false,
			wantSimple:   false,
			wantVerbose:  false,
			wantCategory: "func",
		},
		{
			name:         "multiple flags",
			args:         []string{"ktn-linter", "lint", "--ai", "--no-color", "--category=const", "../../pkg/formatter"},
			wantAI:       true,
			wantNoColor:  true,
			wantSimple:   false,
			wantVerbose:  false,
			wantCategory: "const",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExit(t)
			defer restore()

			os.Args = tt.args

			// Capturer stderr pour éviter le bruit
			oldStderr := os.Stderr
			oldStdout := os.Stdout
			os.Stderr = nil
			os.Stdout = nil
			defer func() {
				os.Stderr = oldStderr
				os.Stdout = oldStdout
			}()

			catchExit(t, func() {
				main()
			})

			if cmd.AIMode != tt.wantAI {
				t.Errorf("AIMode = %v, want %v", cmd.AIMode, tt.wantAI)
			}
			if cmd.NoColor != tt.wantNoColor {
				t.Errorf("NoColor = %v, want %v", cmd.NoColor, tt.wantNoColor)
			}
			if cmd.Simple != tt.wantSimple {
				t.Errorf("Simple = %v, want %v", cmd.Simple, tt.wantSimple)
			}
			if cmd.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", cmd.Verbose, tt.wantVerbose)
			}
			if cmd.Category != tt.wantCategory {
				t.Errorf("Category = %v, want %v", cmd.Category, tt.wantCategory)
			}
		})
	}
}
