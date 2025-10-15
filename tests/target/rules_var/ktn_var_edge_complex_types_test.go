package rules_var_test

import "testing"

func TestNestedMapGood(t *testing.T) {
	if nestedMapGood == nil {
		t.Error("nestedMapGood is nil")
	}
	if len(nestedMapGood) != 1 {
		t.Errorf("nestedMapGood length = %v, want 1", len(nestedMapGood))
	}
}

func TestSliceOfStructsGood(t *testing.T) {
	if sliceOfStructsGood == nil {
		t.Error("sliceOfStructsGood is nil")
	}
	if len(sliceOfStructsGood) != 2 {
		t.Errorf("sliceOfStructsGood length = %v, want 2", len(sliceOfStructsGood))
	}
}

func TestChannelMapGood(t *testing.T) {
	if channelMapGood == nil {
		t.Error("channelMapGood is nil")
	}
	// Test capacity
	if cap(channelMapGood) != 10 {
		t.Errorf("channelMapGood capacity = %v, want 10", cap(channelMapGood))
	}
}
