// Package formatter provides output formatting for lint diagnostics.
package formatter

// JSONSummary represents the summary section of JSON report.
// Provides counts of issues by severity level.
type JSONSummary struct {
	TotalIssues int            `json:"totalIssues"`
	ByLevel     map[string]int `json:"byLevel"`
}
