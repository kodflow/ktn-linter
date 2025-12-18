// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"github.com/kodflow/ktn-linter/pkg/prompt"
	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command.
var promptCmd *cobra.Command = &cobra.Command{
	Use:   "prompt [packages...]",
	Short: "Generate AI-optimized prompt for fixing violations",
	Long: `Generate a markdown prompt that groups violations by rule with examples and phases.

The output is organized into phases:
  1. Structural Changes - Rules that may create/move/delete files
  2. Test Organization - Test file naming and placement rules
  3. Local Fixes - Code modifications within existing files
  4. Comments & Documentation - Applied last after code is finalized

Each rule includes:
  - Description and category
  - Good example from testdata
  - List of all files/lines to fix as checkboxes`,
	Args: cobra.MinimumNArgs(1),
	Run:  runPrompt,
}

// init registers the prompt command with root.
//
// Params: none
//
// Returns: none
func init() {
	rootCmd.AddCommand(promptCmd)
}

// runPrompt executes the prompt generation.
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: package patterns to analyze
//
// Returns: none
func runPrompt(_ *cobra.Command, args []string) {
	opts := parsePromptOptions()

	// Load configuration
	loadPromptConfiguration(opts.Options)

	// Propagate verbose flag to config
	config.Get().Verbose = opts.Verbose

	// Create prompt generator
	gen := prompt.NewGenerator(os.Stderr, opts.Verbose)

	// Generate prompt
	output, err := gen.Generate(args, opts.Options)
	// Check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		OsExit(1)
	}

	// Get output writer
	writer, cleanup := getPromptOutputWriter(opts.OutputPath)
	// Defer cleanup
	if cleanup != nil {
		defer cleanup()
	}

	// Format and display
	formatter := prompt.NewMarkdownFormatter(writer)
	formatter.Format(output)

	// Exit with appropriate code
	if output.TotalViolations > 0 {
		OsExit(1)
	}
	OsExit(0)
}

// promptOptions extends orchestrator options with output path.
type promptOptions struct {
	orchestrator.Options
	OutputPath string
}

// parsePromptOptions extracts options from Cobra flags.
//
// Returns:
//   - promptOptions: extracted options
func parsePromptOptions() promptOptions {
	flags := rootCmd.PersistentFlags()

	verbose, _ := flags.GetBool(flagVerbose)
	category, _ := flags.GetString(flagCategory)
	onlyRule, _ := flags.GetString(flagOnlyRule)
	configPath, _ := flags.GetString(flagConfig)
	outputPath, _ := flags.GetString(flagOutput)

	// Return parsed options
	return promptOptions{
		Options: orchestrator.Options{
			Verbose:    verbose,
			Category:   category,
			OnlyRule:   onlyRule,
			ConfigPath: configPath,
		},
		OutputPath: outputPath,
	}
}

// loadPromptConfiguration loads the linter configuration.
//
// Params:
//   - opts: orchestrator options
//
// Returns: none
func loadPromptConfiguration(opts orchestrator.Options) {
	// Check if config file specified
	if opts.ConfigPath != "" {
		// Load from specified file
		if err := config.LoadAndSet(opts.ConfigPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", opts.ConfigPath, err)
			OsExit(1)
		}
		// Return
		return
	}

	// Try default locations (ignore errors for defaults)
	_ = config.LoadAndSet("")
}

// getPromptOutputWriter returns the writer for output and optional cleanup.
//
// Params:
//   - outputPath: path to output file (empty for stdout)
//
// Returns:
//   - io.Writer: output writer
//   - func(): cleanup function (may be nil)
func getPromptOutputWriter(outputPath string) (io.Writer, func()) {
	// Check if output to file
	if outputPath != "" {
		// Open file for writing
		file, err := os.Create(outputPath)
		// Check for error
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			OsExit(1)
		}
		// Return file and cleanup
		return file, func() { file.Close() }
	}

	// Default to stdout
	return os.Stdout, nil
}
