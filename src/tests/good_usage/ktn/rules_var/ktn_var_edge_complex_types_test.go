package rules_var

import (
	"testing"
)

// TestNestedMapGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestNestedMapGood(t *testing.T) {
	if NestedMapGood == nil {
		t.Error("NestedMapGood is nil")
	}
	if len(NestedMapGood) != 1 {
		t.Errorf("NestedMapGood length = %v, want 1", len(NestedMapGood))
	}
}

// TestSliceOfStructsGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestSliceOfStructsGood(t *testing.T) {
	if SliceOfStructsGood == nil {
		t.Error("SliceOfStructsGood is nil")
	}
	if len(SliceOfStructsGood) != 2 {
		t.Errorf("SliceOfStructsGood length = %v, want 2", len(SliceOfStructsGood))
	}
}

// TestChannelMapGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestChannelMapGood(t *testing.T) {
	if ChannelMapGood == nil {
		t.Error("ChannelMapGood is nil")
	}
	// Test capacity
	if cap(ChannelMapGood) != 10 {
		t.Errorf("ChannelMapGood capacity = %v, want 10", cap(ChannelMapGood))
	}
}
