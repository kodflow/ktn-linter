// Package formatter provides output formatting for lint diagnostics.
package formatter

// OutputFormat represents the output format for lint diagnostics.
type OutputFormat string

// Output format constants.
const (
	// FormatText represents plain text output format.
	FormatText OutputFormat = "text"
	// FormatJSON represents JSON output format.
	FormatJSON OutputFormat = "json"
	// FormatSARIF represents SARIF output format.
	FormatSARIF OutputFormat = "sarif"
)

// ParseOutputFormat parses a string to an OutputFormat.
// If the format is unknown, it defaults to FormatText.
//
// Params:
//   - s: the format string to parse
//
// Returns:
//   - OutputFormat: the parsed format, or FormatText if unknown
func ParseOutputFormat(s string) OutputFormat {
	// Convert string to OutputFormat
	format := OutputFormat(s)

	// Check if format is valid
	if format.IsValid() {
		// Return the valid format
		return format
	}

	// Default to FormatText for unknown formats
	return FormatText
}

// IsValid checks if the OutputFormat is a valid format.
//
// Returns:
//   - bool: true if the format is valid, false otherwise
func (f OutputFormat) IsValid() bool {
	// Check against all valid formats
	switch f {
	// Match any of the valid format constants
	case FormatText, FormatJSON, FormatSARIF:
		// Format is valid
		return true
	}

	// Format is invalid
	return false
}
