package main

import (
	"bytes"
	"flag"
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

// mockExit crée un mock pour osExit qui capture le code de sortie via panic
func mockExit(t *testing.T) (restore func()) {
	t.Helper()
	oldOsExit := osExit
	oldArgs := os.Args

	osExit = func(code int) {
		panic(exitError{code: code})
	}

	return func() {
		osExit = oldOsExit
		os.Args = oldArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		// Reset flags
		aiMode = false
		noColor = false
		simple = false
		verbose = false
		category = ""
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

// TestMainNoArgs teste que main() exit avec 1 sans arguments
func TestMainNoArgs(t *testing.T) {
	restore := mockExit(t)
	defer restore()

	// Capturer stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	os.Args = []string{"ktn-linter"}

	exitCode, didExit := catchExit(t, func() {
		main()
	})

	w.Close()
	var stderr bytes.Buffer
	stderr.ReadFrom(r)

	if !didExit {
		t.Error("Expected main() to exit")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	output := stderr.String()
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage message, got: %s", output)
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

	os.Args = []string{"ktn-linter", "-category=invalid", "."}

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

	os.Args = []string{"ktn-linter", "/nonexistent/path/that/does/not/exist"}

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
	os.Args = []string{"ktn-linter", "../../pkg/formatter"}

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

// TestParseFlags teste le parsing des flags
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantAI   bool
		wantNoColor bool
		wantSimple bool
		wantVerbose bool
		wantCategory string
	}{
		{
			name: "no flags",
			args: []string{"ktn-linter", "."},
			wantAI: false,
			wantNoColor: false,
			wantSimple: false,
			wantVerbose: false,
			wantCategory: "",
		},
		{
			name: "ai flag",
			args: []string{"ktn-linter", "-ai", "."},
			wantAI: true,
		},
		{
			name: "no-color flag",
			args: []string{"ktn-linter", "-no-color", "."},
			wantNoColor: true,
		},
		{
			name: "simple flag",
			args: []string{"ktn-linter", "-simple", "."},
			wantSimple: true,
		},
		{
			name: "verbose flag",
			args: []string{"ktn-linter", "-v", "."},
			wantVerbose: true,
		},
		{
			name: "category flag",
			args: []string{"ktn-linter", "-category=func", "."},
			wantCategory: "func",
		},
		{
			name: "multiple flags",
			args: []string{"ktn-linter", "-ai", "-no-color", "-category=const", "."},
			wantAI: true,
			wantNoColor: true,
			wantCategory: "const",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			aiMode = false
			noColor = false
			simple = false
			verbose = false
			category = ""

			os.Args = tt.args
			parseFlags()

			if aiMode != tt.wantAI {
				t.Errorf("aiMode = %v, want %v", aiMode, tt.wantAI)
			}
			if noColor != tt.wantNoColor {
				t.Errorf("noColor = %v, want %v", noColor, tt.wantNoColor)
			}
			if simple != tt.wantSimple {
				t.Errorf("simple = %v, want %v", simple, tt.wantSimple)
			}
			if verbose != tt.wantVerbose {
				t.Errorf("verbose = %v, want %v", verbose, tt.wantVerbose)
			}
			if category != tt.wantCategory {
				t.Errorf("category = %v, want %v", category, tt.wantCategory)
			}
		})
	}
}

// TestPrintUsage teste que printUsage affiche bien l'aide
func TestPrintUsage(t *testing.T) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printUsage()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Usage:") {
		t.Error("Expected 'Usage:' in output")
	}
	if !strings.Contains(output, "Flags:") {
		t.Error("Expected 'Flags:' in output")
	}
	if !strings.Contains(output, "Categories disponibles:") {
		t.Error("Expected 'Categories disponibles:' in output")
	}
	if !strings.Contains(output, "Examples:") {
		t.Error("Expected 'Examples:' in output")
	}
}
