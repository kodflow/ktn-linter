// Package formatter provides output formatting for lint diagnostics.
package formatter

import "io"

// sarifFormatter implements SARIF output formatting.
type sarifFormatter struct {
	writer  io.Writer
	verbose bool
}
