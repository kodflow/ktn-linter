// Package updater provides self-update functionality for ktn-linter binary.
package updater

import (
	"net/http"
)

// releaseInfo represents the GitHub release API response structure.
// It contains the tag name which is used to determine the latest version.
type releaseInfo struct {
	TagName string `json:"tag_name"`
}

// UpdateInfo contains comprehensive information about available updates.
// It provides the current version, latest available version, and availability status.
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
}

// Updater handles self-update logic for the ktn-linter binary.
// It manages version checking via GitHub API and binary replacement.
type Updater struct {
	version string
	client  *http.Client
}
