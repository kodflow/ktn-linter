// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

// TestUpgradeCmd tests the upgrade command structure.
func TestUpgradeCmd(t *testing.T) {
	// Check command exists
	if upgradeCmd == nil {
		t.Fatal("upgradeCmd is nil")
	}

	// Check command name
	if upgradeCmd.Use != "upgrade" {
		t.Errorf("upgradeCmd.Use = %q, want %q", upgradeCmd.Use, "upgrade")
	}

	// Check short description
	if upgradeCmd.Short == "" {
		t.Error("upgradeCmd.Short is empty")
	}

	// Check run function is set
	if upgradeCmd.Run == nil {
		t.Error("upgradeCmd.Run is nil")
	}
}

// TestUpgradeCmdFlags tests the upgrade command flags.
func TestUpgradeCmdFlags(t *testing.T) {
	// Check --check flag exists
	flag := upgradeCmd.Flags().Lookup(flagCheck)
	// Verify flag exists
	if flag == nil {
		t.Fatal("--check flag not found")
	}

	// Check flag type
	if flag.Value.Type() != "bool" {
		t.Errorf("--check flag type = %q, want %q", flag.Value.Type(), "bool")
	}

	// Check default value
	if flag.DefValue != "false" {
		t.Errorf("--check flag default = %q, want %q", flag.DefValue, "false")
	}
}

// TestUpgradeCmdRegistered tests that upgrade command is registered with root.
func TestUpgradeCmdRegistered(t *testing.T) {
	// Find upgrade command in root's subcommands
	found := false
	// Iterate over subcommands
	for _, cmd := range rootCmd.Commands() {
		// Check if upgrade command
		if cmd.Use == "upgrade" {
			found = true
			break
		}
	}

	// Verify command is registered
	if !found {
		t.Error("upgrade command not registered with root")
	}
}

// TestRunUpgradeDevBuild tests upgrade with dev version.
func TestRunUpgradeDevBuild(t *testing.T) {
	// Save original values
	origVersion := version
	origExit := OsExit

	// Restore after test
	defer func() {
		version = origVersion
		OsExit = origExit
	}()

	// Set dev version
	version = "dev"

	// Track exit code
	var exitCode int
	OsExit = func(code int) {
		exitCode = code
	}

	// Create test command
	cmd := &cobra.Command{}
	cmd.Flags().Bool(flagCheck, true, "")

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Run command (should fail for dev build)
	runUpgrade(cmd, []string{})

	// Check exit code
	if exitCode != 1 {
		t.Errorf("runUpgrade() exit code = %d, want 1 for dev build", exitCode)
	}
}
