package ktntest

import (
	"testing"
)

// TestIsExemptPackage tests isExemptPackage from 001.go
func TestIsExemptPackage(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
		want    bool
	}{
		{"main package", "main", true},
		{"testhelper package", "testhelper", true},
		{"cmd package", "cmd", true},
		{"utils package", "utils", true},
		{"formatter package", "formatter", true},
		{"ktn package", "ktn", true},
		{"ktnconst package", "ktnconst", true},
		{"ktnfunc package", "ktnfunc", true},
		{"ktntest package", "ktntest", true},
		{"ktnvar package", "ktnvar", true},
		{"other package", "mypackage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptPackage(tt.pkgName)
			if got != tt.want {
				t.Errorf("isExemptPackage(%q) = %v, want %v", tt.pkgName, got, tt.want)
			}
		})
	}
}
