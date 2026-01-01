package ktnvar

import (
	"testing"
)

// TestLoopVars004 tests the loopVars004 map.
func TestLoopVars004(t *testing.T) {
	expectedLoopVars := []string{"i", "j", "k", "n", "x", "y", "z", "v"}

	for _, v := range expectedLoopVars {
		if !loopVars004[v] {
			t.Errorf("loopVars004 should contain %q", v)
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
