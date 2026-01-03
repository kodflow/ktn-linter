// Internal tests for sarif_run_adapter.
package formatter

import (
	"testing"

	sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

// TestSarifRunAdapterInternal tests internal adapter behavior.
func TestSarifRunAdapterInternal(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "adapter stores run reference",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a real SARIF run
			run := sarif.NewRunWithInformationURI("test", "http://test")

			// Create adapter with the run
			adapter := &sarifRunAdapter{run: run}

			// Verify adapter stores reference correctly
			if adapter.run != run {
				t.Error("expected adapter to store run reference")
			}

			// Verify GetTool returns correct reference
			if adapter.GetTool() != run.Tool {
				t.Error("expected GetTool to return run's tool")
			}
		})
	}
}

// TestSarifRunAdapterAddResult tests the AddResult method.
func TestSarifRunAdapterAddResult(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "add result returns run",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a real SARIF run
			run := sarif.NewRunWithInformationURI("test", "http://test")

			// Create adapter
			adapter := &sarifRunAdapter{run: run}

			// Add a result
			result := sarif.NewRuleResult("test-rule")
			returnedRun := adapter.AddResult(result)

			// Verify it returns the underlying run
			if returnedRun != run {
				t.Error("expected AddResult to return the underlying run")
			}
		})
	}
}
