// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flag names as constants for type safety.
const (
	// flagVerbose is the flag name for verbose mode.
	flagVerbose string = "verbose"
	// flagCategory is the flag name for category filter.
	flagCategory string = "category"
	// flagOnlyRule is the flag name for single rule filter.
	flagOnlyRule string = "only-rule"
	// flagConfig is the flag name for config path.
	flagConfig string = "config"
	// flagFormat is the flag name for output format.
	flagFormat string = "format"
	// flagOutput is the flag name for output file path.
	flagOutput string = "output"
)

// Global state for testing and version.
var (
	// version stocke la version du linter.
	version string = "dev"

	// OsExit est une variable pour permettre le mocking dans les tests.
	// Par défaut, elle pointe vers os.Exit, mais peut être remplacée par un mock.
	OsExit func(int) = os.Exit

	// rootCmd represents the base command when called without any subcommands.
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "ktn-linter lint [packages...]",
		Short: "KTN-Linter - Linter for Go code following KTN conventions",
		Long: `KTN-Linter ` + version + ` - Linter for Go code following KTN (Kodflow Typing Notation) conventions.

Examples:
  ktn-linter lint ./...                           Lint all packages
  ktn-linter lint --category=func ./...           Lint only function rules
  ktn-linter lint --only-rule=KTN-FUNC-001 .      Lint only a specific rule
  ktn-linter lint -v ./pkg/...                    Lint with verbose output`,
	}
)

// SetVersion définit la version du linter.
//
// Params:
//   - v: version à définir
//
// Returns: aucun
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
	rootCmd.Long = `KTN-Linter ` + v + ` - Linter for Go code following KTN (Kodflow Typing Notation) conventions.

Examples:
  ktn-linter lint ./...                           Lint all packages
  ktn-linter lint --category=func ./...           Lint only function rules
  ktn-linter lint --only-rule=KTN-FUNC-001 .      Lint only a specific rule
  ktn-linter lint -v ./pkg/...                    Lint with verbose output`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
//
// Returns: aucun
//
// Params: aucun
func Execute() {
	// Vérification de la condition
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		OsExit(1)
	}
}

// init configure les flags globaux de la commande root.
//
// Returns: aucun
//
// Params: aucun
func init() {
	// Disable default completion and help commands
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Define persistent flags (available to all subcommands)
	pf := rootCmd.PersistentFlags()
	pf.BoolP(flagVerbose, "v", false, "Verbose output")
	pf.String(flagCategory, "", "Run only rules from specific category (func, var, error, etc.)")
	pf.String(flagOnlyRule, "", "Run only a specific rule by code (e.g., KTN-FUNC-001)")
	pf.StringP(flagConfig, "c", "", "Path to configuration file (.ktn-linter.yaml)")
	pf.String(flagFormat, "text", "Output format: text, json, or sarif")
	pf.StringP(flagOutput, "o", "", "Output file path (default: stdout)")
}
