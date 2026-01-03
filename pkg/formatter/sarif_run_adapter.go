// Package formatter provides output formatting for lint diagnostics.
package formatter

import sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"

// sarifRunAdapter wraps *sarif.Run to implement sarifRunWriter.
type sarifRunAdapter struct {
	run *sarif.Run
}

// AddResult adds a result to the wrapped SARIF run.
//
// Params:
//   - result: rule result to add
//
// Returns:
//   - *sarif.Run: the underlying run
func (a *sarifRunAdapter) AddResult(result *sarif.Result) *sarif.Run {
	// Delegate to underlying run
	return a.run.AddResult(result)
}

// GetTool returns the tool from the wrapped SARIF run.
//
// Returns:
//   - *sarif.Tool: tool reference
func (a *sarifRunAdapter) GetTool() *sarif.Tool {
	// Return tool from run
	return a.run.Tool
}
