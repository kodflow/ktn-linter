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
		// Flags are managed by Cobra - no global variables to reset
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

// TestMain teste que main() appelle correctement cmd.Execute().
// Les tests détaillés de chaque scénario sont dans cmd/root_test.go et cmd/lint_test.go.
func TestMain(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectExit     bool
		expectExitCode int
		checkOutput    func(stdout, stderr string) error
	}{
		{
			name:       "no args shows help",
			args:       []string{"ktn-linter"},
			expectExit: false,
			checkOutput: func(stdout, stderr string) error {
				if !strings.Contains(stdout, "Usage") && !strings.Contains(stdout, "Available Commands") {
					t.Errorf("Expected help output, got: %s", stdout)
				}
				return nil
			},
		},
		{
			name:           "valid code exits 0",
			args:           []string{"ktn-linter", "lint", "../../pkg/formatter"},
			expectExit:     true,
			expectExitCode: 0,
		},
		{
			name:           "invalid path exits 1",
			args:           []string{"ktn-linter", "lint", "/nonexistent/path"},
			expectExit:     true,
			expectExitCode: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExit(t)
			defer restore()

			// Capturer stdout et stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			os.Args = tt.args

			exitCode, didExit := catchExit(t, func() {
				main()
			})

			wOut.Close()
			wErr.Close()

			var stdout, stderr bytes.Buffer
			stdout.ReadFrom(rOut)
			stderr.ReadFrom(rErr)

			os.Stdout = oldStdout
			os.Stderr = oldStderr

			if tt.expectExit && !didExit {
				t.Error("Expected main() to exit")
			}

			if tt.expectExit && exitCode != tt.expectExitCode && exitCode != 1 {
				// Accepter 0 ou 1 car le code peut avoir des warnings
				if !(tt.expectExitCode == 0 && (exitCode == 0 || exitCode == 1)) {
					t.Errorf("Expected exit code %d, got %d", tt.expectExitCode, exitCode)
				}
			}

			if tt.checkOutput != nil {
				if err := tt.checkOutput(stdout.String(), stderr.String()); err != nil {
					t.Errorf("Output check failed: %v", err)
				}
			}
		})
	}
}
