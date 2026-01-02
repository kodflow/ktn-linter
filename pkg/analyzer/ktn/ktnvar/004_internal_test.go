package ktnvar

import (
	"testing"
)

// TestIdiomaticOneChar004 tests the idiomaticOneChar004 map.
func TestIdiomaticOneChar004(t *testing.T) {
	expectedOneChar := []string{"i", "j", "k", "n", "b", "c", "f", "m", "r", "s", "t", "w", "_"}

	for _, v := range expectedOneChar {
		if !idiomaticOneChar004[v] {
			t.Errorf("idiomaticOneChar004 should contain %q", v)
		}
	}
}

// TestIdiomaticShort004 tests the idiomaticShort004 map.
func TestIdiomaticShort004(t *testing.T) {
	expectedIdioms := []string{"ok"}

	for _, v := range expectedIdioms {
		if !idiomaticShort004[v] {
			t.Errorf("idiomaticShort004 should contain %q", v)
		}
	}
}

// TestMinVarNameLength004 tests the minimum name length constant.
func TestMinVarNameLength004(t *testing.T) {
	if minVarNameLength004 != 2 {
		t.Errorf("minVarNameLength004 = %d, want 2", minVarNameLength004)
	}
}
