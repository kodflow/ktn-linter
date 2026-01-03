// Package updater provides self-update functionality for ktn-linter binary.
package updater

// UpdateInfo contains comprehensive information about available updates.
// It provides the current version, latest available version, and availability status.
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
}
