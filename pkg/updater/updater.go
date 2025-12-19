// Package updater provides self-update functionality for ktn-linter binary.
// It checks GitHub releases for newer versions and downloads/replaces the binary.
package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Updater constants for GitHub API, versioning, and file operations.
const (
	// repoOwner is the GitHub repository owner.
	repoOwner string = "kodflow"
	// repoName is the GitHub repository name.
	repoName string = "ktn-linter"
	// apiURL is the GitHub releases API endpoint.
	apiURL string = "https://api.github.com/repos/%s/%s/releases/latest"
	// downloadURL is the release asset download URL pattern.
	downloadURL string = "https://github.com/%s/%s/releases/download/%s/%s"
	// httpTimeout is the timeout for HTTP requests.
	httpTimeout time.Duration = 30 * time.Second
	// semverMajorIdx is the index of major version component.
	semverMajorIdx int = 0
	// semverMinorIdx is the index of minor version component.
	semverMinorIdx int = 1
	// semverPatchIdx is the index of patch version component.
	semverPatchIdx int = 2
	// semverComponents is the number of semver components.
	semverComponents int = 3
	// executablePerm is the permission for executable files.
	executablePerm os.FileMode = 0755
)

// NewUpdater creates a new updater instance.
//
// Params:
//   - version: current binary version (empty means dev build)
//
// Returns:
//   - *Updater: configured updater instance
func NewUpdater(version string) *Updater {
	// Return configured updater with timeout
	return &Updater{
		version: version,
		client:  &http.Client{Timeout: httpTimeout},
	}
}

// CheckForUpdate checks if an update is available.
//
// Returns:
//   - UpdateInfo: information about available update
//   - error: any error during check
func (u *Updater) CheckForUpdate() (UpdateInfo, error) {
	// Check if this is a dev build
	if u.version == "" || u.version == "dev" {
		// Return error for dev builds
		return UpdateInfo{CurrentVersion: u.version}, fmt.Errorf("cannot check updates for dev build")
	}

	// Fetch latest version from GitHub API
	latest, err := u.getLatestVersion()
	// Handle API errors
	if err != nil {
		// Return error with current version info
		return UpdateInfo{CurrentVersion: u.version}, err
	}

	// Compare versions to determine if update is available
	available := u.isNewer(latest)

	// Return complete update info
	return UpdateInfo{
		Available:      available,
		CurrentVersion: u.version,
		LatestVersion:  latest,
	}, nil
}

// Upgrade downloads and applies the latest version.
//
// Returns:
//   - UpdateInfo: information about the update
//   - error: any download or replacement error
func (u *Updater) Upgrade() (UpdateInfo, error) {
	// Check for available updates first
	info, err := u.CheckForUpdate()
	// Handle check errors
	if err != nil {
		// Return error from check
		return info, err
	}

	// Verify update is available
	if !info.Available {
		// Return info indicating already up to date
		return info, nil
	}

	// Download and replace binary with new version
	err = u.downloadAndReplace(info.LatestVersion)
	// Handle download/replace errors
	if err != nil {
		// Return error from download
		return info, err
	}

	// Return success info
	return info, nil
}

// getLatestVersion fetches the latest release version from GitHub.
//
// Returns:
//   - string: latest version tag
//   - error: any API error
func (u *Updater) getLatestVersion() (string, error) {
	// Build API URL for latest release
	url := fmt.Sprintf(apiURL, repoOwner, repoName)
	// Make HTTP request to GitHub API
	resp, err := u.client.Get(url)
	// Handle HTTP request errors
	if err != nil {
		// Return wrapped HTTP error
		return "", fmt.Errorf("fetching release info: %w", err)
	}
	defer resp.Body.Close()

	// Validate response status code
	if resp.StatusCode != http.StatusOK {
		// Return error for non-200 status
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Parse JSON response
	var release releaseInfo
	// Decode JSON body
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		// Return wrapped decode error
		return "", fmt.Errorf("parsing release info: %w", err)
	}

	// Return the version tag
	return release.TagName, nil
}

// isNewer checks if the given version is newer than current.
//
// Params:
//   - latest: version to compare against
//
// Returns:
//   - bool: true if latest is newer
func (u *Updater) isNewer(latest string) bool {
	// Parse both versions into components
	current := u.parseVersion(u.version)
	remote := u.parseVersion(latest)

	// Compare major version first
	if remote[semverMajorIdx] > current[semverMajorIdx] {
		// Remote has higher major version
		return true
	}
	// Check if current major is higher
	if remote[semverMajorIdx] < current[semverMajorIdx] {
		// Current has higher major version
		return false
	}

	// Major versions equal, compare minor
	if remote[semverMinorIdx] > current[semverMinorIdx] {
		// Remote has higher minor version
		return true
	}
	// Check if current minor is higher
	if remote[semverMinorIdx] < current[semverMinorIdx] {
		// Current has higher minor version
		return false
	}

	// Major and minor equal, compare patch
	return remote[semverPatchIdx] > current[semverPatchIdx]
}

// parseVersion parses a semver string into components.
//
// Params:
//   - v: version string (e.g., "v1.2.3")
//
// Returns:
//   - [3]int: major, minor, patch as integers
func (u *Updater) parseVersion(v string) [semverComponents]int {
	// Remove optional 'v' prefix
	v = strings.TrimPrefix(v, "v")
	// Split by dots
	parts := strings.Split(v, ".")

	// Initialize result array
	var result [semverComponents]int
	// Iterate over version parts
	for i := 0; i < semverComponents && i < len(parts); i++ {
		// Parse integer, defaulting to 0 on error
		result[i], _ = strconv.Atoi(parts[i])
	}
	// Return parsed version
	return result
}

// downloadAndReplace downloads the new binary and replaces the current one.
//
// Params:
//   - version: version to download
//
// Returns:
//   - error: any download or replacement error
func (u *Updater) downloadAndReplace(version string) error {
	// Get platform-specific binary name
	binaryName := u.getBinaryName()
	// Build download URL
	url := fmt.Sprintf(downloadURL, repoOwner, repoName, version, binaryName)

	// Make HTTP request to download binary
	resp, err := u.client.Get(url)
	// Handle download request errors
	if err != nil {
		// Return wrapped download error
		return fmt.Errorf("downloading binary: %w", err)
	}
	defer resp.Body.Close()

	// Validate download response status
	if resp.StatusCode != http.StatusOK {
		// Return error for failed download
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// Get path to current executable
	execPath, err := os.Executable()
	// Handle path resolution errors
	if err != nil {
		// Return wrapped path error
		return fmt.Errorf("getting executable path: %w", err)
	}
	// Resolve any symlinks in path
	execPath, err = filepath.EvalSymlinks(execPath)
	// Handle symlink resolution errors
	if err != nil {
		// Return wrapped symlink error
		return fmt.Errorf("resolving symlinks: %w", err)
	}

	// Create temporary file for download
	tmpFile, err := os.CreateTemp(filepath.Dir(execPath), "ktn-linter-update-*")
	// Handle temp file creation errors
	if err != nil {
		// Return wrapped temp file error
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Copy downloaded content to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	// Handle write errors
	if err != nil {
		os.Remove(tmpPath)
		// Return wrapped write error
		return fmt.Errorf("writing temp file: %w", err)
	}

	// Set executable permissions on temp file
	err = os.Chmod(tmpPath, executablePerm)
	// Handle chmod errors
	if err != nil {
		os.Remove(tmpPath)
		// Return wrapped chmod error
		return fmt.Errorf("setting permissions: %w", err)
	}

	// Atomically replace old binary with new one
	err = os.Rename(tmpPath, execPath)
	// Handle rename errors
	if err != nil {
		os.Remove(tmpPath)
		// Return wrapped rename error
		return fmt.Errorf("replacing binary: %w", err)
	}

	// Return nil on success
	return nil
}

// getBinaryName returns the binary name for the current platform.
//
// Returns:
//   - string: binary name with platform suffix
func (u *Updater) getBinaryName() string {
	// Build base name with OS and architecture
	name := fmt.Sprintf("ktn-linter-%s-%s", runtime.GOOS, runtime.GOARCH)
	// Check if Windows platform
	if runtime.GOOS == "windows" {
		// Add .exe extension for Windows
		name += ".exe"
	}
	// Return platform-specific binary name
	return name
}
