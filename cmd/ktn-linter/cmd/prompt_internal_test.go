// Internal tests for the prompt command.
package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// Test_runPrompt tests the runPrompt function with different scenarios.
//
// Params:
//   - t: testing object
func Test_runPrompt(t *testing.T) {
	tests := []struct {
		name        string
		packages    []string
		expectExit  bool
		exitCode    int
		checkStderr string
	}{
		{
			name:       "valid formatter package",
			packages:   []string{"github.com/kodflow/ktn-linter/pkg/formatter"},
			expectExit: true,
			exitCode:   0,
		},
		{
			name:        "testdata with potential issues",
			packages:    []string{"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst/testdata/src/const001"},
			expectExit:  true,
			exitCode:    1,
			checkStderr: "",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := mockExitInCmd(t)
			defer restore()

			// Reset all flags to defaults
			flags := rootCmd.PersistentFlags()
			flags.Set(flagVerbose, "false")
			flags.Set(flagCategory, "")
			flags.Set(flagOnlyRule, "")
			flags.Set(flagConfig, "")
			flags.Set(flagOutput, "")

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() {
				os.Stdout = oldStdout
			}()

			exitCode, didExit := catchExitInCmd(t, func() {
				runPrompt(promptCmd, tt.packages)
			})

			w.Close()
			r.Close()

			// Verify exit expectation
			if tt.expectExit && !didExit {
				t.Error("expected exit but did not exit")
			}

			// Verify exit code
			if tt.expectExit && exitCode != tt.exitCode {
				t.Errorf("expected exit code %d, got %d", tt.exitCode, exitCode)
			}
		})
	}
}

// Test_parsePromptOptions tests the parsePromptOptions function.
//
// Params:
//   - t: testing object
func Test_parsePromptOptions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		wantOpts promptOptions
	}{
		{
			name: "default options",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
				flags.Set(flagOutput, "")
			},
			wantOpts: promptOptions{
				Options: orchestrator.Options{
					Verbose:    false,
					Category:   "",
					OnlyRule:   "",
					ConfigPath: "",
				},
				OutputPath: "",
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
				flags.Set(flagOutput, "")
			},
			wantOpts: promptOptions{
				Options: orchestrator.Options{
					Verbose:    true,
					Category:   "",
					OnlyRule:   "",
					ConfigPath: "",
				},
				OutputPath: "",
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
				flags.Set(flagOutput, "")
			},
			wantOpts: promptOptions{
				Options: orchestrator.Options{
					Verbose:    false,
					Category:   "func",
					OnlyRule:   "",
					ConfigPath: "",
				},
				OutputPath: "",
			},
		},
		{
			name: "with output path",
			setup: func() {
				flags := rootCmd.PersistentFlags()
				flags.Set(flagVerbose, "false")
				flags.Set(flagCategory, "")
				flags.Set(flagOnlyRule, "")
				flags.Set(flagConfig, "")
				flags.Set(flagOutput, "/tmp/output.md")
			},
			wantOpts: promptOptions{
				Options: orchestrator.Options{
					Verbose:    false,
					Category:   "",
					OnlyRule:   "",
					ConfigPath: "",
				},
				OutputPath: "/tmp/output.md",
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			opts := parsePromptOptions()

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

			// Verify output path
			if opts.OutputPath != tt.wantOpts.OutputPath {
				t.Errorf("OutputPath = %v, want %v", opts.OutputPath, tt.wantOpts.OutputPath)
			}
		})
	}
}

// Test_loadPromptConfiguration tests the loadPromptConfiguration function.
//
// Params:
//   - t: testing object
func Test_loadPromptConfiguration(t *testing.T) {
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
	}

	// Run tests
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
				loadPromptConfiguration(opts)
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

// Test_getPromptOutputWriter tests the getPromptOutputWriter function.
//
// Params:
//   - t: testing object
func Test_getPromptOutputWriter(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (string, func())
		expectStdout bool
		expectFile   bool
	}{
		{
			name: "empty path returns stdout",
			setup: func() (string, func()) {
				return "", func() {}
			},
			expectStdout: true,
			expectFile:   false,
		},
		{
			name: "valid path returns file writer",
			setup: func() (string, func()) {
				tmpfile, _ := os.CreateTemp("", "test-output-*.md")
				tmpfile.Close()
				return tmpfile.Name(), func() {
					os.Remove(tmpfile.Name())
				}
			},
			expectStdout: false,
			expectFile:   true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath, cleanup := tt.setup()
			defer cleanup()

			writer, writerCleanup := getPromptOutputWriter(outputPath)

			// Verify writer is not nil
			if writer == nil {
				t.Error("getPromptOutputWriter returned nil writer")
			}

			// Verify stdout expectation
			if tt.expectStdout && writer != os.Stdout {
				t.Error("expected stdout writer")
			}

			// Verify file expectation
			if tt.expectFile && writer == os.Stdout {
				t.Error("expected file writer, got stdout")
			}

			// Cleanup if present
			if writerCleanup != nil {
				writerCleanup()
			}
		})
	}
}

// TestPromptCmdStructure tests the prompt command structure.
//
// Params:
//   - t: testing object
func TestPromptCmdStructure(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "Use field is correct",
			check: func(t *testing.T) {
				const expectedUse = "prompt [packages...]"
				// Verify Use
				if promptCmd.Use != expectedUse {
					t.Errorf("expected Use=%q, got %q", expectedUse, promptCmd.Use)
				}
			},
		},
		{
			name: "Short description is not empty",
			check: func(t *testing.T) {
				// Verify Short
				if promptCmd.Short == "" {
					t.Error("Short description should not be empty")
				}
			},
		},
		{
			name: "Long description is not empty",
			check: func(t *testing.T) {
				// Verify Long
				if promptCmd.Long == "" {
					t.Error("Long description should not be empty")
				}
			},
		},
		{
			name: "Args validator is not nil",
			check: func(t *testing.T) {
				// Verify Args
				if promptCmd.Args == nil {
					t.Error("Args validator should not be nil")
				}
			},
		},
		{
			name: "Run function is not nil",
			check: func(t *testing.T) {
				// Verify Run
				if promptCmd.Run == nil {
					t.Error("Run function should not be nil")
				}
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
