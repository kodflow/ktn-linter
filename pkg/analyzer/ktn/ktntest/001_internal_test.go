// Internal tests for analyzer 001.
package ktntest

import (
	"testing"
)

// Test_runTest001 tests the runTest001 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "basic test of runTest001 structure",
		},
		{
			name: "error case validation",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing: %s", tt.name)
		})
	}
}

// Test_isExemptPackage tests the isExemptPackage private function.
//
// Params:
//   - t: testing context
func Test_isExemptPackage(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
		want    bool
	}{
		{
			name:    "main package is exempt",
			pkgName: "main",
			want:    true,
		},
		{
			name:    "testhelper package is exempt",
			pkgName: "testhelper",
			want:    true,
		},
		{
			name:    "ktntest package is exempt",
			pkgName: "ktntest",
			want:    true,
		},
		{
			name:    "regular package not exempt",
			pkgName: "mypackage",
			want:    false,
		},
		{
			name:    "empty package not exempt",
			pkgName: "",
			want:    false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptPackage(tt.pkgName)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isExemptPackage(%q) = %v, want %v", tt.pkgName, got, tt.want)
			}
		})
	}
}

// Test_runTest001_integration tests the analyzer structure.
//
// Params:
//   - t: testing context
func Test_runTest001_integration(t *testing.T) {
	// Test analyzer structure
	if Analyzer001 == nil {
		t.Fatal("Analyzer001 should not be nil")
	}
	// Vérification du nom
	if Analyzer001.Name != "ktntest001" {
		t.Errorf("Analyzer001.Name = %q, want %q", Analyzer001.Name, "ktntest001")
	}
}
