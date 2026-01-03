// Package updater provides self-update functionality for ktn-linter binary.
package updater

// releaseInfo represents the GitHub release API response structure.
// It contains the tag name which is used to determine the latest version.
type releaseInfo struct {
	TagName string `json:"tag_name"`
}
