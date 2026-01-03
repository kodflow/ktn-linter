// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"fmt"

	"github.com/kodflow/ktn-linter/pkg/updater"
	"github.com/spf13/cobra"
)

// Flag for check-only mode.
const (
	// flagCheck is the flag name for check-only mode.
	flagCheck string = "check"
)

// upgradeCmd represents the upgrade command.
var upgradeCmd *cobra.Command = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade ktn-linter to the latest version",
	Long: `Upgrade ktn-linter to the latest version from GitHub releases.

This command checks for a newer version and downloads/replaces the
current binary if an update is available.

Examples:
  ktn-linter upgrade          Check and upgrade to latest version
  ktn-linter upgrade --check  Only check for updates without upgrading`,
	Run: runUpgrade,
}

// init registers the upgrade command with root.
func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().Bool(flagCheck, false, "Only check for updates without upgrading")
}

// runUpgrade executes the upgrade command.
//
// Params:
//   - cmd: cobra command
//   - _args: command arguments (unused - upgrade takes no args)
func runUpgrade(cmd *cobra.Command, _args []string) {
	// Create updater with current version
	upd := updater.NewUpdater(version)

	// Delegate to internal function with interfaces
	runUpgradeWithDeps(cmd.Flags(), upd)
}

// runUpgradeWithDeps executes the upgrade with injected dependencies.
//
// Params:
//   - flags: flag getter for command flags
//   - upd: updater service
func runUpgradeWithDeps(flags flagGetter, upd updaterService) {
	// Check if check-only mode
	checkOnly, _ := flags.GetBool(flagCheck)

	// Handle check-only mode
	if checkOnly {
		handleCheckOnly(upd)
		// Exit after check-only mode
		return
	}

	// Perform upgrade
	handleUpgrade(upd)
}

// handleCheckOnly checks for updates without upgrading.
//
// Params:
//   - upd: updater service
func handleCheckOnly(upd updaterService) {
	fmt.Println("Checking for updates...")

	info, err := upd.CheckForUpdate()
	// Check for errors
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		OsExit(1)
		// Exit after error
		return
	}

	// Display result
	if info.Available {
		fmt.Printf("Update available: %s → %s\n", info.CurrentVersion, info.LatestVersion)
		fmt.Println("Run 'ktn-linter upgrade' to update.")
		// Update is available, handled above
	} else {
		// Already up to date
		fmt.Printf("Already up to date: %s\n", info.CurrentVersion)
	}
}

// handleUpgrade performs the actual upgrade.
//
// Params:
//   - upd: updater service
func handleUpgrade(upd updaterService) {
	fmt.Println("Checking for updates...")

	info, err := upd.Upgrade()
	// Check for errors
	if err != nil {
		fmt.Printf("Error upgrading: %v\n", err)
		OsExit(1)
		// Exit after error
		return
	}

	// Display result
	if info.Available {
		fmt.Printf("Successfully upgraded: %s → %s\n", info.CurrentVersion, info.LatestVersion)
		// Upgrade completed successfully
	} else {
		// Already up to date
		fmt.Printf("Already up to date: %s\n", info.CurrentVersion)
	}
}
