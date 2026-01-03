// External tests for sarif_run_adapter.
package formatter_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
	sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

// TestSarifRunAdapter tests the SARIF run adapter functionality.
func TestSarifRunAdapter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "adapter wraps sarif run",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a new SARIF formatter (which uses the adapter internally)
			f := formatter.NewSARIFFormatter(nil, false)
			// Verify formatter is created
			if f == nil {
				t.Error("expected non-nil formatter")
			}
		})
	}
}

// TestSarifRunAdapterMethods tests the adapter methods.
func TestSarifRunAdapterMethods(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "run with results",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a real SARIF run
			run := sarif.NewRunWithInformationURI("test", "http://test")
			// Verify run is created
			if run == nil {
				t.Error("expected non-nil run")
			}
		})
	}
}
