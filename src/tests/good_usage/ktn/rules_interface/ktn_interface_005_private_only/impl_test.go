package goodinterfaces

import (
	"testing"

)

// TestNewService teste TODO.
//
// Params:
//   - t: contexte de test
func TestNewService(t *testing.T) {
	svc := NewService("test-service")
	if svc == nil {
		t.Error("NewService() returned nil")
	}
}

// TestNewHelper teste TODO.
//
// Params:
//   - t: contexte de test
func TestNewHelper(t *testing.T) {
	helper := NewHelper()
	if helper == nil {
		t.Error("NewHelper() returned nil")
	}
}
