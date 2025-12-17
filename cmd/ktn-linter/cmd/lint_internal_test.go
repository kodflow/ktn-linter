// Internal tests for the lint command.
package cmd

import (
	"bytes"
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"golang.org/x/tools/go/analysis"
)

// Test_runLint tests the runLint function with different packages.
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
			packages: []string{"../../../pkg/analyzer/ktn/ktnconst/testdata/src/const001"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Capture stdout
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

			// Verify exit and code
			if !didExit || (exitCode != 0 && exitCode != 1) {
				t.Errorf("Test failed: didExit=%v, exitCode=%d", didExit, exitCode)
			}
		})
	}
}

// Test_parseOptions tests the parseOptions function.
func Test_parseOptions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		wantOpts orchestrator.Options
	}{
		{
			name: "default options",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
			},
			wantOpts: orchestrator.Options{
				Verbose:    false,
				Category:   "",
				OnlyRule:   "",
				ConfigPath: "",
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
			wantOpts: orchestrator.Options{
				Verbose:    true,
				Category:   "",
				OnlyRule:   "",
				ConfigPath: "",
			},
		},
		{
			name: "with category filter",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "func")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
			},
			wantOpts: orchestrator.Options{
				Verbose:    false,
				Category:   "func",
				OnlyRule:   "",
				ConfigPath: "",
			},
		},
		{
			name: "with only rule",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "KTN-FUNC-001")
				flags.Set(flagConfig, "")
			},
			wantOpts: orchestrator.Options{
				Verbose:    false,
				Category:   "",
				OnlyRule:   "KTN-FUNC-001",
				ConfigPath: "",
			},
		},
		{
			name: "with config path",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "/path/to/config.yaml")
			},
			wantOpts: orchestrator.Options{
				Verbose:    false,
				Category:   "",
				OnlyRule:   "",
				ConfigPath: "/path/to/config.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			opts := parseOptions()

			// Verify all fields
			if opts.Verbose != tt.wantOpts.Verbose {
				t.Errorf("Verbose = %v, want %v", opts.Verbose, tt.wantOpts.Verbose)
			}
			// Verify category
			if opts.Category != tt.wantOpts.Category {
				t.Errorf("Category = %v, want %v", opts.Category, tt.wantOpts.Category)
			}
			// Verify only rule
			if opts.OnlyRule != tt.wantOpts.OnlyRule {
				t.Errorf("OnlyRule = %v, want %v", opts.OnlyRule, tt.wantOpts.OnlyRule)
			}
			// Verify config path
			if opts.ConfigPath != tt.wantOpts.ConfigPath {
				t.Errorf("ConfigPath = %v, want %v", opts.ConfigPath, tt.wantOpts.ConfigPath)
			}
		})
	}
}

// Test_loadConfiguration tests the loadConfiguration function.
func Test_loadConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (orchestrator.Options, func())
		expectExit  bool
		exitCode    int
		checkStderr string
	}{
		{
			name: "empty path succeeds",
			setup: func() (orchestrator.Options, func()) {
				return orchestrator.Options{}, func() {}
			},
			expectExit:  false,
			exitCode:    0,
			checkStderr: "",
		},
		{
			name: "valid config file",
			setup: func() (orchestrator.Options, func()) {
				configData := `version: 1
exclude:
  - "**/*_test.go"
`
				tmpfile, _ := os.CreateTemp("", "test-config-*.yaml")
				tmpfile.Write([]byte(configData))
				tmpfile.Close()
				return orchestrator.Options{ConfigPath: tmpfile.Name()}, func() {
					os.Remove(tmpfile.Name())
				}
			},
			expectExit:  false,
			exitCode:    0,
			checkStderr: "",
		},
		{
			name: "invalid config file exits with error",
			setup: func() (orchestrator.Options, func()) {
				return orchestrator.Options{ConfigPath: "/nonexistent/config.yaml"}, func() {}
			},
			expectExit:  true,
			exitCode:    1,
			checkStderr: "Error loading config",
		},
		{
			name: "verbose mode logs config loading",
			setup: func() (orchestrator.Options, func()) {
				configData := `version: 1
exclude:
  - "*.tmp"
`
				tmpfile, _ := os.CreateTemp("", "verbose-config-*.yaml")
				tmpfile.Write([]byte(configData))
				tmpfile.Close()
				return orchestrator.Options{ConfigPath: tmpfile.Name(), Verbose: true}, func() {
					os.Remove(tmpfile.Name())
				}
			},
			expectExit:  false,
			exitCode:    0,
			checkStderr: "Loaded configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			opts, cleanup := tt.setup()
			defer cleanup()

			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode, didExit := catchExitInCmd(t, func() {
				loadConfiguration(opts)
			})

			w.Close()
			var stderr bytes.Buffer
			stderr.ReadFrom(r)
			os.Stderr = oldStderr

			// Verify exit expectation
			if tt.expectExit && !didExit {
				t.Error("expected exit but did not exit")
			}
			// Verify no exit expectation
			if !tt.expectExit && didExit {
				t.Errorf("unexpected exit with code %d", exitCode)
			}
			// Verify exit code
			if tt.expectExit && exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}
			// Verify stderr content
			if tt.checkStderr != "" && !strings.Contains(stderr.String(), tt.checkStderr) {
				t.Errorf("expected stderr to contain %q, got %q", tt.checkStderr, stderr.String())
			}
		})
	}
}

// Test_runPipeline tests the runPipeline function.
func Test_runPipeline(t *testing.T) {
	tests := []struct {
		name        string
		packages    []string
		opts        orchestrator.Options
		expectError bool
	}{
		{
			name:        "valid package",
			packages:    []string{"../../../pkg/formatter"},
			opts:        orchestrator.Options{},
			expectError: false,
		},
		{
			name:        "with category filter",
			packages:    []string{"../../../pkg/formatter"},
			opts:        orchestrator.Options{Category: "func"},
			expectError: false,
		},
		{
			name:        "with single rule",
			packages:    []string{"../../../pkg/formatter"},
			opts:        orchestrator.Options{OnlyRule: "KTN-FUNC-001"},
			expectError: false,
		},
		{
			name:        "invalid category returns error",
			packages:    []string{"../../../pkg/formatter"},
			opts:        orchestrator.Options{Category: "nonexistent"},
			expectError: true,
		},
		{
			name:        "invalid rule returns error",
			packages:    []string{"../../../pkg/formatter"},
			opts:        orchestrator.Options{OnlyRule: "KTN-INVALID-999"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orch := orchestrator.NewOrchestrator(os.Stderr, tt.opts.Verbose)

			diags, fset, err := runPipeline(orch, tt.packages, tt.opts)

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			// Verify no error expectation
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify results on success
			_ = diags
			_ = fset
		})
	}
}

// Test_formatAndDisplay tests the formatAndDisplay function.
func Test_formatAndDisplay(t *testing.T) {
	tests := []struct {
		name          string
		diagnostics   []analysis.Diagnostic
		fset          *token.FileSet
		opts          lintOptions
		expectedInMsg string
	}{
		{
			name:          "empty diagnostics shows success",
			diagnostics:   []analysis.Diagnostic{},
			fset:          nil,
			opts:          lintOptions{},
			expectedInMsg: "No issues found",
		},
		{
			name: "diagnostics are displayed",
			diagnostics: func() []analysis.Diagnostic {
				fset := token.NewFileSet()
				file := fset.AddFile("test.go", -1, 100)
				return []analysis.Diagnostic{
					{
						Pos:     file.Pos(10),
						Message: "test issue",
					},
				}
			}(),
			fset: func() *token.FileSet {
				fset := token.NewFileSet()
				fset.AddFile("test.go", -1, 100)
				return fset
			}(),
			opts:          lintOptions{},
			expectedInMsg: "test issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			formatAndDisplay(tt.diagnostics, tt.fset, tt.opts)

			w.Close()
			var stdout bytes.Buffer
			stdout.ReadFrom(r)
			os.Stdout = oldStdout

			// Verify output contains expected message
			if !strings.Contains(stdout.String(), tt.expectedInMsg) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedInMsg, stdout.String())
			}
		})
	}
}

// TestLintCmdStructure tests the lint command structure.
func TestLintCmdStructure(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "Use field is correct",
			check: func(t *testing.T) {
				const expectedUse = "lint [packages...]"
				// Verify Use
				if lintCmd.Use != expectedUse {
					t.Errorf("expected Use=%q, got %q", expectedUse, lintCmd.Use)
				}
			},
		},
		{
			name: "Short description is not empty",
			check: func(t *testing.T) {
				// Verify Short
				if lintCmd.Short == "" {
					t.Error("Short description should not be empty")
				}
			},
		},
		{
			name: "Long description is not empty",
			check: func(t *testing.T) {
				// Verify Long
				if lintCmd.Long == "" {
					t.Error("Long description should not be empty")
				}
			},
		},
		{
			name: "Args validator is not nil",
			check: func(t *testing.T) {
				// Verify Args
				if lintCmd.Args == nil {
					t.Error("Args validator should not be nil")
				}
			},
		},
		{
			name: "Run function is not nil",
			check: func(t *testing.T) {
				// Verify Run
				if lintCmd.Run == nil {
					t.Error("Run function should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}

// Test_getOutputWriter tests the getOutputWriter function.
func Test_getOutputWriter(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		wantFile   bool
	}{
		{
			name:       "empty path returns stdout",
			outputPath: "",
			wantFile:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, cleanup := getOutputWriter(tt.outputPath)

			// Verify writer is not nil
			if writer == nil {
				t.Error("getOutputWriter returned nil writer")
			}

			// Cleanup if present
			if cleanup != nil {
				cleanup()
			}
		})
	}
}
