// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
	"go/token"
	"io"
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
	DiscoverModules(paths []string) ([]string, error)
	RunMultiModule(paths []string, opts orchestrator.Options) ([]orchestrator.DiagnosticResult, error)
}

// lintCmd represents the lint command.
var lintCmd *cobra.Command = &cobra.Command{
	Use:   "lint [packages...]",
	Short: "Lint Go packages using KTN rules",
	Long:  `Lint analyzes Go packages and reports issues based on KTN conventions.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runLint,
}

// Flag names for lint command.
const (
	// flagSarif is the flag name for SARIF output.
	flagSarif string = "sarif"
	// flagJSON is the flag name for JSON output.
	flagJSON string = "json"
)

// init registers the lint command with root.
//
// Params: none
//
// Returns: none
func init() {
	rootCmd.AddCommand(lintCmd)
	// Add lint-specific flags
	lintCmd.Flags().Bool(flagSarif, false, "Output in SARIF format (for IDE integration)")
	lintCmd.Flags().Bool(flagJSON, false, "Output in JSON format")
}

// runLint executes the linting analysis.
//
// Params:
//   - cmd: Cobra command (used to get flags)
//   - args: command line arguments
//
// Returns: none
func runLint(cmd *cobra.Command, args []string) {
	opts := parseOptions(cmd)

	// Load configuration
	loadConfiguration(opts.Options)

	// Propagate verbose flag to config
	config.Get().Verbose = opts.Verbose

	// Create orchestrator
	orch := orchestrator.NewOrchestrator(os.Stderr, opts.Verbose)

	// Run the linting pipeline
	diags, fset, err := runPipeline(orch, args, opts.Options)
	// Check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		OsExit(1)
	}

	// Format and display results
	formatAndDisplay(diags, fset, opts)

	// Exit with appropriate code
	if len(diags) > 0 {
		OsExit(1)
	}
	OsExit(0)
}

// lintOptions extends orchestrator options with CLI-specific settings.
type lintOptions struct {
	orchestrator.Options
	Format     formatter.OutputFormat
	OutputPath string
}

// parseOptions extracts options from Cobra flags.
//
// Params:
//   - cmd: Cobra command with flags
//
// Returns:
//   - lintOptions: extracted options including format and output
func parseOptions(cmd *cobra.Command) lintOptions {
	flags := rootCmd.PersistentFlags()

	verbose, _ := flags.GetBool(flagVerbose)
	category, _ := flags.GetString(flagCategory)
	onlyRule, _ := flags.GetString(flagOnlyRule)
	configPath, _ := flags.GetString(flagConfig)
	outputPath, _ := flags.GetString(flagOutput)

	// Check lint-specific format flags (--sarif, --json)
	sarifMode, _ := cmd.Flags().GetBool(flagSarif)
	jsonMode, _ := cmd.Flags().GetBool(flagJSON)

	// Determine output format
	outputFormat := formatter.FormatText
	// Check for SARIF format
	if sarifMode {
		outputFormat = formatter.FormatSARIF
	} else if jsonMode {
		// Check for JSON format
		outputFormat = formatter.FormatJSON
	}

	// Return parsed options
	return lintOptions{
		Options: orchestrator.Options{
			Verbose:    verbose,
			Category:   category,
			OnlyRule:   onlyRule,
			ConfigPath: configPath,
		},
		Format:     outputFormat,
		OutputPath: outputPath,
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
//   - args: package patterns or paths
//   - opts: linting options
//
// Returns:
//   - []analysis.Diagnostic: found issues
//   - *token.FileSet: first fileset for formatting
//   - error: pipeline error if any
func runPipeline(orch lintOrchestrator, args []string, opts orchestrator.Options) ([]analysis.Diagnostic, *token.FileSet, error) {
	// Check if we need multi-module discovery
	if needsModuleDiscovery(args) {
		// Use multi-module approach
		return runMultiModulePipeline(orch, args, opts)
	}

	// Use standard single-module approach
	return runSingleModulePipeline(orch, args, opts)
}

// needsModuleDiscovery checks if args require module discovery.
//
// Params:
//   - args: command line arguments
//
// Returns:
//   - bool: true if discovery needed
func needsModuleDiscovery(args []string) bool {
	// Check each arg
	for _, arg := range args {
		// Skip standard Go patterns
		if arg == "./..." || arg == "." {
			continue
		}
		// Check if path exists as directory
		info, err := os.Stat(arg)
		// Skip if not accessible
		if err != nil {
			continue
		}
		// Check if directory
		if info.IsDir() {
			// Directory found, discovery needed
			return true
		}
	}
	// No directory found, no discovery needed
	return false
}

// runMultiModulePipeline runs analysis across multiple modules.
//
// Params:
//   - orch: linting orchestrator interface
//   - args: paths to analyze
//   - opts: linting options
//
// Returns:
//   - []analysis.Diagnostic: found issues
//   - *token.FileSet: first fileset for formatting
//   - error: pipeline error if any
func runMultiModulePipeline(orch lintOrchestrator, args []string, opts orchestrator.Options) ([]analysis.Diagnostic, *token.FileSet, error) {
	// Run multi-module analysis
	rawDiags, err := orch.RunMultiModule(args, opts)
	// Check for error
	if err != nil {
		// Return error from multi-module analysis
		return []analysis.Diagnostic{}, nil, err
	}

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

// runSingleModulePipeline runs analysis for a single module.
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
func runSingleModulePipeline(orch lintOrchestrator, args []string, opts orchestrator.Options) ([]analysis.Diagnostic, *token.FileSet, error) {
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
//   - opts: lint options including format and output path
//
// Returns: none
func formatAndDisplay(diagnostics []analysis.Diagnostic, fset *token.FileSet, opts lintOptions) {
	// Get output writer
	writer, cleanup := getOutputWriter(opts.OutputPath)
	// Defer cleanup
	if cleanup != nil {
		defer cleanup()
	}

	// Create formatter options
	// SimpleMode désactivé : on affiche toujours le format complet
	// VerboseMode n'affecte plus les messages (toujours longs)
	fmtOpts := formatter.FormatterOptions{
		AIMode:      false,
		NoColor:     opts.OutputPath != "",
		SimpleMode:  false,
		VerboseMode: false,
	}

	// Create formatter based on format
	fmtr := formatter.NewFormatterByFormat(opts.Format, writer, fmtOpts)

	// Check if empty
	if len(diagnostics) == 0 {
		fmtr.Format(nil, nil)
		// Return
		return
	}

	fmtr.Format(fset, diagnostics)
}

// getOutputWriter returns the writer for output and optional cleanup function.
//
// Params:
//   - outputPath: path to output file (empty for stdout)
//
// Returns:
//   - io.Writer: output writer
//   - func(): cleanup function (may be nil)
func getOutputWriter(outputPath string) (io.Writer, func()) {
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
