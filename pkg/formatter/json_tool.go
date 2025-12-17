// Package formatter provides output formatting for lint diagnostics.
package formatter

// JSONTool represents the tool information in JSON report.
// Includes the tool name and version.
type JSONTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
