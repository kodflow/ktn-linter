// Package updater provides self-update functionality for ktn-linter binary.
package updater

import (
	"testing"
)

// TestParseVersion tests version parsing for various formats.
func TestParseVersion(t *testing.T) {
	u := NewUpdater("v1.0.0")

	tests := []struct {
		name     string
		version  string
		expected [3]int
	}{
		{
			name:     "standard version with v prefix",
			version:  "v1.2.3",
			expected: [3]int{1, 2, 3},
		},
		{
			name:     "version without v prefix",
			version:  "1.2.3",
			expected: [3]int{1, 2, 3},
		},
		{
			name:     "major version only",
			version:  "v2",
			expected: [3]int{2, 0, 0},
		},
		{
			name:     "major and minor only",
			version:  "v2.5",
			expected: [3]int{2, 5, 0},
		},
		{
			name:     "empty version string",
			version:  "",
			expected: [3]int{0, 0, 0},
		},
		{
			name:     "invalid version format",
			version:  "invalid",
			expected: [3]int{0, 0, 0},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := u.parseVersion(tt.version)
			// Check result matches expected
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, result, tt.expected)
			}
		})
	}
}

// TestIsNewer tests version comparison logic for various scenarios.
func TestIsNewer(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{
			name:     "major version upgrade available",
			current:  "v1.0.0",
			latest:   "v2.0.0",
			expected: true,
		},
		{
			name:     "minor version upgrade available",
			current:  "v1.0.0",
			latest:   "v1.1.0",
			expected: true,
		},
		{
			name:     "patch version upgrade available",
			current:  "v1.0.0",
			latest:   "v1.0.1",
			expected: true,
		},
		{
			name:     "same version no upgrade",
			current:  "v1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "older major version no upgrade",
			current:  "v2.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "older minor version no upgrade",
			current:  "v1.5.0",
			latest:   "v1.4.0",
			expected: false,
		},
		{
			name:     "older patch version no upgrade",
			current:  "v1.0.5",
			latest:   "v1.0.4",
			expected: false,
		},
		{
			name:     "major upgrade from high minor patch",
			current:  "v1.9.9",
			latest:   "v2.0.0",
			expected: true,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpdater(tt.current)
			result := u.isNewer(tt.latest)
			// Check result matches expected
			if result != tt.expected {
				t.Errorf("isNewer(%q) with current %q = %v, want %v",
					tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}

// TestGetBinaryName tests platform-specific binary name generation.
func TestGetBinaryName(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		wantContains  string
		wantMinLength int
	}{
		{
			name:          "standard version binary name",
			version:       "v1.0.0",
			wantContains:  "ktn-linter-",
			wantMinLength: 15,
		},
		{
			name:          "dev version binary name",
			version:       "dev",
			wantContains:  "ktn-linter-",
			wantMinLength: 15,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpdater(tt.version)
			name := u.getBinaryName()
			// Check name is not empty
			if name == "" {
				t.Error("getBinaryName() returned empty string")
			}
			// Check name contains expected prefix
			if !containsSubstring(name, tt.wantContains) {
				t.Errorf("getBinaryName() = %q, want to contain %q", name, tt.wantContains)
			}
			// Check minimum length
			if len(name) < tt.wantMinLength {
				t.Errorf("getBinaryName() length = %d, want >= %d", len(name), tt.wantMinLength)
			}
		})
	}
}

// containsSubstring checks if a string contains a substring.
//
// Params:
//   - s: string to search in
//   - substr: substring to find
//
// Returns:
//   - bool: true if substr found in s
func containsSubstring(s, substr string) bool {
	// Check each possible starting position
	for i := 0; i <= len(s)-len(substr); i++ {
		// Check if substring matches at position
		if s[i:i+len(substr)] == substr {
			// Match found
			return true
		}
	}
	// No match found
	return false
}
