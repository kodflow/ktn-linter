// Package updater provides self-update functionality for ktn-linter binary.
package updater

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestUpdater_parseVersion tests version parsing for various formats.
func TestUpdater_parseVersion(t *testing.T) {
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
		{
			name:     "version with extra parts",
			version:  "v1.2.3.4",
			expected: [3]int{1, 2, 3},
		},
		{
			name:     "version with leading zeros",
			version:  "v01.02.03",
			expected: [3]int{1, 2, 3},
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := u.parseVersion(tt.version)
			// Check result matches expected
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, result, tt.expected)
			}
		})
	}
}

// TestUpdater_isNewer tests version comparison logic for various scenarios.
func TestUpdater_isNewer(t *testing.T) {
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
		{
			name:     "same major higher remote minor",
			current:  "v2.0.0",
			latest:   "v2.1.0",
			expected: true,
		},
		{
			name:     "same major minor higher remote patch",
			current:  "v2.1.0",
			latest:   "v2.1.5",
			expected: true,
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
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

// TestUpdater_getBinaryName tests platform-specific binary name generation.
func TestUpdater_getBinaryName(t *testing.T) {
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
		{
			name:          "empty version binary name",
			version:       "",
			wantContains:  "ktn-linter-",
			wantMinLength: 15,
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			u := NewUpdater(tt.version)
			name := u.getBinaryName()
			// Check name is not empty
			if name == "" {
				t.Error("getBinaryName() returned empty string")
			}
			// Check name contains expected prefix
			if !strings.Contains(name, tt.wantContains) {
				t.Errorf("getBinaryName() = %q, want to contain %q", name, tt.wantContains)
			}
			// Check minimum length
			if len(name) < tt.wantMinLength {
				t.Errorf("getBinaryName() length = %d, want >= %d", len(name), tt.wantMinLength)
			}
		})
	}
}

// TestUpdater_getLatestVersion tests fetching latest version from GitHub API.
func TestUpdater_getLatestVersion(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		responseStatus int
		wantVersion    string
		wantErr        bool
		wantErrContain string
	}{
		{
			name:           "successful response",
			responseBody:   `{"tag_name": "v1.2.3"}`,
			responseStatus: http.StatusOK,
			wantVersion:    "v1.2.3",
			wantErr:        false,
		},
		{
			name:           "not found status",
			responseBody:   `{"message": "Not Found"}`,
			responseStatus: http.StatusNotFound,
			wantVersion:    "",
			wantErr:        true,
			wantErrContain: "unexpected status",
		},
		{
			name:           "server error status",
			responseBody:   `{"message": "Internal Server Error"}`,
			responseStatus: http.StatusInternalServerError,
			wantVersion:    "",
			wantErr:        true,
			wantErrContain: "unexpected status",
		},
		{
			name:           "invalid json response",
			responseBody:   `{invalid json`,
			responseStatus: http.StatusOK,
			wantVersion:    "",
			wantErr:        true,
			wantErrContain: "parsing release info",
		},
		{
			name:           "empty tag name",
			responseBody:   `{"tag_name": ""}`,
			responseStatus: http.StatusOK,
			wantVersion:    "",
			wantErr:        false,
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create updater with custom client pointing to test server
			u := &Updater{
				version: "v1.0.0",
				client:  server.Client(),
			}

			// Override the API URL by using a custom transport
			originalClient := u.client
			u.client = &http.Client{
				Transport: &mockTransport{
					url:    server.URL,
					client: originalClient,
				},
			}

			version, err := u.getLatestVersion()

			// Check error expectation
			if tt.wantErr {
				// Expect error
				if err == nil {
					t.Errorf("getLatestVersion() error = nil, wantErr %v", tt.wantErr)
				} else if tt.wantErrContain != "" && !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("getLatestVersion() error = %v, want error containing %q", err, tt.wantErrContain)
				}
			} else {
				// Expect no error
				if err != nil {
					t.Errorf("getLatestVersion() unexpected error = %v", err)
				}
			}

			// Check version
			if version != tt.wantVersion {
				t.Errorf("getLatestVersion() = %q, want %q", version, tt.wantVersion)
			}
		})
	}
}

// mockTransport is a custom RoundTripper that redirects requests to test server.
type mockTransport struct {
	url    string
	client *http.Client
}

// RoundTrip implements http.RoundTripper interface.
//
// Params:
//   - req: the HTTP request
//
// Returns:
//   - *http.Response: the HTTP response
//   - error: any error during transport
func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m == nil || m.client == nil {
		return nil, errors.New("mockTransport: nil client")
	}
	// Guard against nil request
	if req == nil || req.URL == nil {
		return nil, errors.New("mockTransport: nil request")
	}

	// Parse test-server base URL
	base, err := url.Parse(m.url)
	if err != nil {
		return nil, err
	}

	// Preserve original path + query, redirect to test server host
	target := base.ResolveReference(&url.URL{
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	})

	// Clone the request preserving context
	newReq := req.Clone(req.Context())
	newReq.URL = target
	newReq.Host = target.Host
	newReq.RequestURI = "" // Clear RequestURI to prevent HTTP client errors
	newReq.Header = req.Header.Clone()

	// Get transport with nil fallback (and prevent self-recursion)
	rt := m.client.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	// Prevent self-recursion
	if _, ok := rt.(*mockTransport); ok {
		rt = http.DefaultTransport
	}
	// Execute request
	return rt.RoundTrip(newReq)
}

// TestUpdater_downloadAndReplace tests the download and replace functionality.
func TestUpdater_downloadAndReplace(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		wantErr        bool
		wantErrContain string
	}{
		{
			name:           "download not found",
			responseStatus: http.StatusNotFound,
			responseBody:   "Not Found",
			wantErr:        true,
			wantErrContain: "download failed",
		},
		{
			name:           "server error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Internal Server Error",
			wantErr:        true,
			wantErrContain: "download failed",
		},
		{
			name:           "forbidden access",
			responseStatus: http.StatusForbidden,
			responseBody:   "Forbidden",
			wantErr:        true,
			wantErrContain: "download failed",
		},
	}

	// Run each test case
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create test server that returns specified response
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create updater with custom client
			u := &Updater{
				version: "v1.0.0",
				client: &http.Client{
					Transport: &mockTransport{
						url:    server.URL,
						client: server.Client(),
					},
				},
			}

			err := u.downloadAndReplace("v1.1.0")

			// Check error expectation
			if tt.wantErr {
				// Expect error
				if err == nil {
					t.Errorf("downloadAndReplace() error = nil, wantErr %v", tt.wantErr)
				} else if tt.wantErrContain != "" && !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("downloadAndReplace() error = %v, want error containing %q", err, tt.wantErrContain)
				}
			} else {
				// Expect no error
				if err != nil {
					t.Errorf("downloadAndReplace() unexpected error = %v", err)
				}
			}
		})
	}
}
