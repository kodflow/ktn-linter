// Lint command implementation for analyzing Go code.
package cmd

import (
	"fmt"
	"go/token"
	"os"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/formatter"
	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// lintOrchestrator defines the interface for linting orchestration.
// Abstracts the orchestrator for testability.
type lintOrchestrator interface {
	LoadPackages(patterns []string) ([]*packages.Package, error)
	SelectAnalyzers(opts orchestrator.Options) ([]*analysis.Analyzer, error)
	RunAnalyzers(pkgs []*packages.Package, analyzers []*analysis.Analyzer) []orchestrator.DiagnosticResult
	FilterDiagnostics(diagnostics []orchestrator.DiagnosticResult) []orchestrator.DiagnosticResult
	ExtractDiagnostics(diagnostics []orchestrator.DiagnosticResult) []analysis.Diagnostic
}

// lintCmd represents the lint command.
var lintCmd *cobra.Command = &cobra.Command{
	Use:   "lint [packages...]",
	Short: "Lint Go packages using KTN rules",
	Long:  `Lint analyzes Go packages and reports issues based on KTN conventions.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runLint,
}

// init registers the lint command with root.
//
// Params: none
//
// Returns: none
func init() {
	rootCmd.AddCommand(lintCmd)
}

// runLint executes the linting analysis.
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: command line arguments
//
// Returns: none
func runLint(_ *cobra.Command, args []string) {
	opts := parseOptions()

	// Load configuration
	loadConfiguration(opts)

	// Propagate verbose flag to config
	config.Get().Verbose = opts.Verbose

	// Create orchestrator
	orch := orchestrator.NewOrchestrator(os.Stderr, opts.Verbose)

	// Run the linting pipeline
	diags, fset, err := runPipeline(orch, args, opts)
	// Check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		OsExit(1)
	}

	// Format and display results
	formatAndDisplay(diags, fset, opts.Verbose)

	// Exit with appropriate code
	if len(diags) > 0 {
		OsExit(1)
	}
	OsExit(0)
}

// parseOptions extracts options from Cobra flags.
//
// Returns:
//   - orchestrator.Options: extracted options
func parseOptions() orchestrator.Options {
	flags := rootCmd.PersistentFlags()

	verbose, _ := flags.GetBool(flagVerbose)
	category, _ := flags.GetString(flagCategory)
	onlyRule, _ := flags.GetString(flagOnlyRule)
	configPath, _ := flags.GetString(flagConfig)

	// Return parsed options
	return orchestrator.Options{
		Verbose:    verbose,
		Category:   category,
		OnlyRule:   onlyRule,
		ConfigPath: configPath,
	}
}

// loadConfiguration loads the linter configuration.
//
// Params:
//   - opts: lint options
//
// Returns: none
func loadConfiguration(opts orchestrator.Options) {
	// Check if config file specified
	if opts.ConfigPath != "" {
		// Load from specified file
		if err := config.LoadAndSet(opts.ConfigPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", opts.ConfigPath, err)
			OsExit(1)
		}
		// Log if verbose
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "Loaded configuration from %s\n", opts.ConfigPath)
		}
		// Return
		return
	}

	// Try default locations
	if err := config.LoadAndSet(""); err == nil {
		// Log if verbose
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "Loaded configuration from default location\n")
		}
	}
}

// runPipeline runs the complete linting pipeline.
//
// Params:
//   - orch: linting orchestrator interface
//   - args: package patterns
//   - opts: linting options
//
// Returns:
//   - []analysis.Diagnostic: found issues
//   - *token.FileSet: first fileset for formatting
//   - error: pipeline error if any
func runPipeline(orch lintOrchestrator, args []string, opts orchestrator.Options) ([]analysis.Diagnostic, *token.FileSet, error) {
	// Load packages
	pkgs, err := orch.LoadPackages(args)
	// Check for error
	if err != nil {
		// Return error
		return []analysis.Diagnostic{}, nil, err
	}

	// Select analyzers
	analyzers, err := orch.SelectAnalyzers(opts)
	// Check for error
	if err != nil {
		// Return error
		return []analysis.Diagnostic{}, nil, err
	}

	// Run analyzers
	rawDiags := orch.RunAnalyzers(pkgs, analyzers)

	// Filter diagnostics
	filtered := orch.FilterDiagnostics(rawDiags)

	// Get first fset for formatting
	var fset *token.FileSet
	// Check if diagnostics exist
	if len(filtered) > 0 {
		fset = filtered[0].Fset
	}

	// Extract and deduplicate
	diags := orch.ExtractDiagnostics(filtered)

	// Return results
	return diags, fset, nil
}

// formatAndDisplay formats and displays diagnostics.
//
// Params:
//   - diagnostics: diagnostics to display
//   - fset: fileset for positions
//   - verbose: enable verbose output
//
// Returns: none
func formatAndDisplay(diagnostics []analysis.Diagnostic, fset *token.FileSet, verbose bool) {
	fmtr := formatter.NewFormatter(os.Stdout, false, true, true, verbose)

	// Check if empty
	if len(diagnostics) == 0 {
		fmtr.Format(nil, nil)
		// Return
		return
	}

	fmtr.Format(fset, diagnostics)
}
