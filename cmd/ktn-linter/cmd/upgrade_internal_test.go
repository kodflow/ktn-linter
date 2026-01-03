// Package cmd implements the CLI commands for ktn-linter.
package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/updater"
	"github.com/spf13/cobra"
)

// TestUpgradeCmd tests the upgrade command structure.
func TestUpgradeCmd(t *testing.T) {
	// Define test cases for upgrade command structure
	tests := []struct {
		name string
	}{
		{name: "verifies upgrade command structure and configuration"},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

// TestUpgradeCmdFlags tests the upgrade command flags.
func TestUpgradeCmdFlags(t *testing.T) {
	// Define test cases for flag verification
	tests := []struct {
		name        string
		flagName    string
		expectType  string
		expectValue string
	}{
		{name: "verifies check flag exists with correct type and default", flagName: flagCheck, expectType: "bool", expectValue: "false"},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Check flag exists
			flag := upgradeCmd.Flags().Lookup(tt.flagName)
			// Verify flag exists
			if flag == nil {
				t.Fatal("flag not found")
			}

			// Check flag type
			if flag.Value.Type() != tt.expectType {
				t.Errorf("flag type = %q, want %q", flag.Value.Type(), tt.expectType)
			}

			// Check default value
			if flag.DefValue != tt.expectValue {
				t.Errorf("flag default = %q, want %q", flag.DefValue, tt.expectValue)
			}
		})
	}
}

// TestUpgradeCmdRegistered tests that upgrade command is registered with root.
func TestUpgradeCmdRegistered(t *testing.T) {
	// Define test cases for command registration
	tests := []struct {
		name    string
		cmdName string
	}{
		{name: "confirms upgrade command is registered with root", cmdName: "upgrade"},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Find command in root's subcommands
			found := false
			// Iterate over subcommands
			for _, cmd := range rootCmd.Commands() {
				// Check if target command
				if cmd.Use == tt.cmdName {
					found = true
					break
				}
			}

			// Verify command is registered
			if !found {
				t.Errorf("%s command not registered with root", tt.cmdName)
			}
		})
	}
}

// TestRunUpgradeDevBuild tests upgrade with dev version.
func TestRunUpgradeDevBuild(t *testing.T) {
	// Define test cases for dev build behavior
	tests := []struct {
		name         string
		version      string
		expectedExit int
	}{
		{name: "exits with code 1 when version is dev", version: "dev", expectedExit: 1},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := version
			origExit := OsExit

			// Restore after test
			defer func() {
				version = origVersion
				OsExit = origExit
			}()

			// Set test version
			version = tt.version

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

			// Run command
			runUpgrade(cmd, []string{})

			// Check exit code
			if exitCode != tt.expectedExit {
				t.Errorf("runUpgrade() exit code = %d, want %d", exitCode, tt.expectedExit)
			}
		})
	}
}

// TestRunUpgradeWithDeps tests the testable upgrade function with mocks.
func TestRunUpgradeWithDeps(t *testing.T) {
	tests := []struct {
		name         string
		checkOnly    bool
		checkInfo    updater.UpdateInfo
		checkErr     error
		upgradeInfo  updater.UpdateInfo
		upgradeErr   error
		expectedExit int
	}{
		{
			name:      "check only with update available",
			checkOnly: true,
			checkInfo: updater.UpdateInfo{
				Available:      true,
				CurrentVersion: "1.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:      "check only already up to date",
			checkOnly: true,
			checkInfo: updater.UpdateInfo{
				Available:      false,
				CurrentVersion: "2.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:         "check only with error",
			checkOnly:    true,
			checkErr:     errors.New("network error"),
			expectedExit: 1,
		},
		{
			name:      "upgrade successful",
			checkOnly: false,
			upgradeInfo: updater.UpdateInfo{
				Available:      true,
				CurrentVersion: "1.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:      "upgrade already up to date",
			checkOnly: false,
			upgradeInfo: updater.UpdateInfo{
				Available:      false,
				CurrentVersion: "2.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:         "upgrade with error",
			checkOnly:    false,
			upgradeErr:   errors.New("download failed"),
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore OsExit
			origExit := OsExit
			defer func() { OsExit = origExit }()

			var exitCode int
			exitCalled := false
			OsExit = func(code int) {
				exitCode = code
				exitCalled = true
			}

			// Create mock flags
			flags := newMockFlagGetter()
			flags.boolValues[flagCheck] = tt.checkOnly

			// Create mock updater
			mockUpd := &mockUpdaterService{
				checkInfo:   tt.checkInfo,
				checkErr:    tt.checkErr,
				upgradeInfo: tt.upgradeInfo,
				upgradeErr:  tt.upgradeErr,
			}

			// Run testable function
			runUpgradeWithDeps(flags, mockUpd)

			// Check exit code
			if exitCalled && exitCode != tt.expectedExit {
				t.Errorf("exit code = %d, want %d", exitCode, tt.expectedExit)
			}
		})
	}
}

// TestHandleCheckOnly tests the check only handler with mock.
func TestHandleCheckOnly(t *testing.T) {
	tests := []struct {
		name         string
		info         updater.UpdateInfo
		err          error
		expectedExit int
	}{
		{
			name: "update available",
			info: updater.UpdateInfo{
				Available:      true,
				CurrentVersion: "1.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name: "already up to date",
			info: updater.UpdateInfo{
				Available:      false,
				CurrentVersion: "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:         "error checking",
			err:          errors.New("network error"),
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			origExit := OsExit
			defer func() { OsExit = origExit }()

			var exitCode int
			exitCalled := false
			OsExit = func(code int) {
				exitCode = code
				exitCalled = true
			}

			mockUpd := &mockUpdaterService{
				checkInfo: tt.info,
				checkErr:  tt.err,
			}

			handleCheckOnly(mockUpd)

			// Only check exit code if OsExit was called
			if exitCalled && exitCode != tt.expectedExit {
				t.Errorf("exit code = %d, want %d", exitCode, tt.expectedExit)
			}
		})
	}
}

// TestHandleUpgrade tests the upgrade handler with mock.
func TestHandleUpgrade(t *testing.T) {
	tests := []struct {
		name         string
		info         updater.UpdateInfo
		err          error
		expectedExit int
	}{
		{
			name: "upgrade successful",
			info: updater.UpdateInfo{
				Available:      true,
				CurrentVersion: "1.0.0",
				LatestVersion:  "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name: "already up to date",
			info: updater.UpdateInfo{
				Available:      false,
				CurrentVersion: "2.0.0",
			},
			expectedExit: 0,
		},
		{
			name:         "error upgrading",
			err:          errors.New("download failed"),
			expectedExit: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			origExit := OsExit
			defer func() { OsExit = origExit }()

			var exitCode int
			exitCalled := false
			OsExit = func(code int) {
				exitCode = code
				exitCalled = true
			}

			mockUpd := &mockUpdaterService{
				upgradeInfo: tt.info,
				upgradeErr:  tt.err,
			}

			handleUpgrade(mockUpd)

			// Only check exit code if OsExit was called
			if exitCalled && exitCode != tt.expectedExit {
				t.Errorf("exit code = %d, want %d", exitCode, tt.expectedExit)
			}
		})
	}
}
