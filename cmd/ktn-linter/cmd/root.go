// Root command configuration for the ktn-linter CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global flags and commands.
var (
	// AIMode enables AI-friendly output format.
	AIMode bool = false
	// NoColor disables colored output.
	NoColor bool = false
	// Simple enables simple one-line format for IDE integration.
	Simple bool = false
	// Verbose enables verbose output.
	Verbose bool = false
	// Category filters rules by specific category.
	Category string = ""
	// OnlyRule runs only a specific rule by code (e.g., KTN-FUNC-001).
	OnlyRule string = ""
	// Fix enables automatic fix application for modernize analyzers.
	Fix bool = false
	// ConfigPath is the path to the configuration file.
	ConfigPath string = ""
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
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&AIMode, "ai", false, "Enable AI-friendly output format")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&Simple, "simple", false, "Simple one-line format for IDE integration")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&Category, "category", "", "Run only rules from specific category (func, var, error, etc.)")
	rootCmd.PersistentFlags().StringVar(&OnlyRule, "only-rule", "", "Run only a specific rule by code (e.g., KTN-FUNC-001)")
	rootCmd.PersistentFlags().BoolVar(&Fix, "fix", false, "Automatically apply suggested fixes from modernize analyzers")
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "Path to configuration file (.ktn-linter.yaml)")
}
