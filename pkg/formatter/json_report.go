// Package formatter provides output formatting for lint diagnostics.
package formatter

// JSONReport represents the top-level JSON report structure.
// Contains metadata, summary, and list of diagnostic results.
type JSONReport struct {
	Schema  string       `json:"$schema"`
	Version string       `json:"version"`
	Tool    JSONTool     `json:"tool"`
	Summary JSONSummary  `json:"summary"`
	Results []JSONResult `json:"results"`
}
