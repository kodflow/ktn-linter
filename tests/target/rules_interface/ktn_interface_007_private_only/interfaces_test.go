package goodinterfaces_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_interface/ktn_interface_007_private_only"
)

// TestServiceInterface teste TODO.
//
// Params:
//   - t: contexte de test
func TestServiceInterface(t *testing.T) {
	svc := goodinterfaces.NewService("test")

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "valid data",
			data:    "test data",
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Process(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Test Close
	if err := svc.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestHelperInterface teste TODO.
//
// Params:
//   - t: contexte de test
func TestHelperInterface(t *testing.T) {
	helper := goodinterfaces.NewHelper()

	msg := helper.Help()
	if msg == "" {
		t.Error("Help() returned empty string")
	}
}
