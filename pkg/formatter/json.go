// Package formatter provides output formatting for lint diagnostics.
package formatter

import (
	"encoding/json"
	"go/token"
	"io"

	"github.com/kodflow/ktn-linter/pkg/severity"
	"golang.org/x/tools/go/analysis"
)

// jsonFormatter implements JSON output formatting.
// Provides structured output compatible with CI/CD tools.
type jsonFormatter struct {
	writer  io.Writer
	verbose bool
}

// NewJSONFormatter creates a new JSON formatter.
//
// Params:
//   - w: writer for output
//   - verbose: enable verbose messages
//
// Returns:
//   - Formatter: JSON formatter instance
func NewJSONFormatter(w io.Writer, verbose bool) Formatter {
	// Return new JSON formatter
	return &jsonFormatter{
		writer:  w,
		verbose: verbose,
	}
}

// Format outputs diagnostics in JSON format.
//
// Params:
//   - fset: fileset for position information
//   - diagnostics: list of diagnostics to format
func (f *jsonFormatter) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	// Create report structure
	report := f.buildReport(fset, diagnostics)

	// Encode as JSON with indentation
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")

	// Write JSON output
	_ = encoder.Encode(report)
}

// buildReport builds the JSON report from diagnostics.
//
// Params:
//   - fset: fileset for position information
//   - diagnostics: list of diagnostics
//
// Returns:
//   - JSONReport: constructed report
func (f *jsonFormatter) buildReport(fset *token.FileSet, diagnostics []analysis.Diagnostic) JSONReport {
	// Initialize level counts
	byLevel := map[string]int{
		"error":   0,
		"warning": 0,
		"info":    0,
	}

	// Build results and count by level
	results := make([]JSONResult, 0, len(diagnostics))

	// Iterate over diagnostics
	for _, diag := range diagnostics {
		result := f.buildResult(fset, diag)
		results = append(results, result)

		// Count by level
		byLevel[result.Level]++
	}

	// Return complete report
	return JSONReport{
		Schema:  "https://ktn-linter.dev/schema/v1/report.json",
		Version: "1.0.0",
		Tool: JSONTool{
			Name:    "ktn-linter",
			Version: "1.0.0",
		},
		Summary: JSONSummary{
			TotalIssues: len(diagnostics),
			ByLevel:     byLevel,
		},
		Results: results,
	}
}

// buildResult builds a single JSON result from a diagnostic.
//
// Params:
//   - fset: fileset for position information
//   - diag: diagnostic to convert
//
// Returns:
//   - JSONResult: converted result
func (f *jsonFormatter) buildResult(fset *token.FileSet, diag analysis.Diagnostic) JSONResult {
	// Get position from fileset
	pos := fset.Position(diag.Pos)

	// Extract rule code from message
	code := extractCode(diag.Message)

	// Get severity level
	level := severity.GetSeverity(code)
	levelStr := f.severityToLevel(level)

	// Extract message without code prefix
	message := extractMessageWithOptions(diag.Message, !f.verbose)

	// Return constructed result
	return JSONResult{
		RuleID:  code,
		Level:   levelStr,
		Message: message,
		Location: JSONLocation{
			File:   pos.Filename,
			Line:   pos.Line,
			Column: pos.Column,
		},
	}
}

// severityToLevel converts severity level to JSON level string.
//
// Params:
//   - level: severity level
//
// Returns:
//   - string: level string (error, warning, info)
func (f *jsonFormatter) severityToLevel(level severity.Level) string {
	// Map severity to JSON level
	switch level {
	// Error level case
	case severity.SeverityError:
		// Return error string
		return "error"
	// Warning level case
	case severity.SeverityWarning:
		// Return warning string
		return "warning"
	// Info level case
	case severity.SeverityInfo:
		// Return info string
		return "info"
	}

	// Default to warning
	return "warning"
}
