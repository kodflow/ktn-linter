// Root command configuration for the ktn-linter CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flag names as constants for type safety.
const (
	// flagAI is the flag name for AI mode.
	flagAI string = "ai"
	// flagNoColor is the flag name for no-color mode.
	flagNoColor string = "no-color"
	// flagSimple is the flag name for simple mode.
	flagSimple string = "simple"
	// flagVerbose is the flag name for verbose mode.
	flagVerbose string = "verbose"
	// flagCategory is the flag name for category filter.
	flagCategory string = "category"
	// flagOnlyRule is the flag name for single rule filter.
	flagOnlyRule string = "only-rule"
	// flagFix is the flag name for auto-fix mode.
	flagFix string = "fix"
	// flagConfig is the flag name for config path.
	flagConfig string = "config"
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
		Use:   "ktn-linter",
		Short: "KTN-Linter - Linter for Go code following KTN conventions",
		Long: `KTN-Linter is a specialized linter that enforces naming conventions and code quality standards for Go projects.

It analyzes Go code to ensure compliance with KTN (Kodflow Typing Notation) standards.`,
		Version: version,
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
	// Define persistent flags (available to all subcommands)
	pf := rootCmd.PersistentFlags()
	pf.Bool(flagAI, false, "Enable AI-friendly output format")
	pf.Bool(flagNoColor, false, "Disable colored output")
	pf.Bool(flagSimple, false, "Simple one-line format for IDE integration")
	pf.BoolP(flagVerbose, "v", false, "Verbose output")
	pf.String(flagCategory, "", "Run only rules from specific category (func, var, error, etc.)")
	pf.String(flagOnlyRule, "", "Run only a specific rule by code (e.g., KTN-FUNC-001)")
	pf.Bool(flagFix, false, "Automatically apply suggested fixes from modernize analyzers")
	pf.StringP(flagConfig, "c", "", "Path to configuration file (.ktn-linter.yaml)")
}
