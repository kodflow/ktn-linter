// Package formatter provides output formatting for lint diagnostics.
package formatter

// JSONLocation represents the location of a diagnostic.
// Specifies file path, line number, and column position.
type JSONLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}
