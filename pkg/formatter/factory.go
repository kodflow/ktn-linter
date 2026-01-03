// Package formatter provides output formatting for lint diagnostics.
package formatter

import (
	"io"
)

// FormatterOptions contains options for formatter creation.
// It provides configuration for all formatter types.
type FormatterOptions struct {
	AIMode      bool
	NoColor     bool
	SimpleMode  bool
	VerboseMode bool
}

// NewFormatterByFormat creates a formatter based on output format.
//
// Params:
//   - format: output format (text, json, sarif)
//   - w: writer for output
//   - opts: formatter options
//
// Returns:
//   - Formater: appropriate formatter for the format
func NewFormatterByFormat(format OutputFormat, w io.Writer, opts FormatterOptions) Formater {
	// Select formatter based on format
	switch format {
	// JSON format case
	case FormatJSON:
		// Return JSON formatter
		return NewJSONFormatter(w, opts.VerboseMode)
	// SARIF format case
	case FormatSARIF:
		// Return SARIF formatter
		return NewSARIFFormatter(w, opts.VerboseMode)
	// Default case
	default:
		// Return default text formatter
		return NewFormatter(w, opts.AIMode, opts.NoColor, opts.SimpleMode, opts.VerboseMode)
	}
}
