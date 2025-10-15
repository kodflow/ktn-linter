package goodinterfaces_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_interface/ktn_interface_007_private_only"
)

// TestNewService teste TODO.
//
// Params:
//   - t: contexte de test
func TestNewService(t *testing.T) {
	svc := goodinterfaces.NewService("test-service")
	if svc == nil {
		t.Error("NewService() returned nil")
	}
}

// TestNewHelper teste TODO.
//
// Params:
//   - t: contexte de test
func TestNewHelper(t *testing.T) {
	helper := goodinterfaces.NewHelper()
	if helper == nil {
		t.Error("NewHelper() returned nil")
	}
}
