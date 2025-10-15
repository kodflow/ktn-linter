package rules_var_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_var"
)

func TestNestedMapGood(t *testing.T) {
	if rules_var.NestedMapGood == nil {
		t.Error("NestedMapGood is nil")
	}
	if len(rules_var.NestedMapGood) != 1 {
		t.Errorf("NestedMapGood length = %v, want 1", len(rules_var.NestedMapGood))
	}
}

func TestSliceOfStructsGood(t *testing.T) {
	if rules_var.SliceOfStructsGood == nil {
		t.Error("SliceOfStructsGood is nil")
	}
	if len(rules_var.SliceOfStructsGood) != 2 {
		t.Errorf("SliceOfStructsGood length = %v, want 2", len(rules_var.SliceOfStructsGood))
	}
}

func TestChannelMapGood(t *testing.T) {
	if rules_var.ChannelMapGood == nil {
		t.Error("ChannelMapGood is nil")
	}
	// Test capacity
	if cap(rules_var.ChannelMapGood) != 10 {
		t.Errorf("ChannelMapGood capacity = %v, want 10", cap(rules_var.ChannelMapGood))
	}
}
