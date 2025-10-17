package rules_var_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/good_usage/rules_var"
)

// TestNestedMapGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestNestedMapGood(t *testing.T) {
	if rules_var.NestedMapGood == nil {
		t.Error("NestedMapGood is nil")
	}
	if len(rules_var.NestedMapGood) != 1 {
		t.Errorf("NestedMapGood length = %v, want 1", len(rules_var.NestedMapGood))
	}
}

// TestSliceOfStructsGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestSliceOfStructsGood(t *testing.T) {
	if rules_var.SliceOfStructsGood == nil {
		t.Error("SliceOfStructsGood is nil")
	}
	if len(rules_var.SliceOfStructsGood) != 2 {
		t.Errorf("SliceOfStructsGood length = %v, want 2", len(rules_var.SliceOfStructsGood))
	}
}

// TestChannelMapGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestChannelMapGood(t *testing.T) {
	if rules_var.ChannelMapGood == nil {
		t.Error("ChannelMapGood is nil")
	}
	// Test capacity
	if cap(rules_var.ChannelMapGood) != 10 {
		t.Errorf("ChannelMapGood capacity = %v, want 10", cap(rules_var.ChannelMapGood))
	}
}
