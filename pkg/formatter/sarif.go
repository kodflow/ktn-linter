// Package formatter provides output formatting for lint diagnostics.
package formatter

import (
	"go/token"
	"io"

	"github.com/kodflow/ktn-linter/pkg/severity"
	sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"golang.org/x/tools/go/analysis"
)

// sarifFormatter implements SARIF output formatting.
type sarifFormatter struct {
	writer  io.Writer
	verbose bool
}

// NewSARIFFormatter creates a new SARIF formatter.
//
// Params:
//   - w: writer for output
//   - verbose: enable verbose messages
//
// Returns:
//   - Formatter: SARIF formatter instance
func NewSARIFFormatter(w io.Writer, verbose bool) Formatter {
	// Return new SARIF formatter
	return &sarifFormatter{
		writer:  w,
		verbose: verbose,
	}
}

// Format outputs diagnostics in SARIF format.
//
// Params:
//   - fset: fileset for position information
//   - diagnostics: list of diagnostics to format
func (f *sarifFormatter) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	// Create new SARIF report
	report := sarif.NewReport()

	// Create run with tool information
	run := sarif.NewRunWithInformationURI("ktn-linter", "https://github.com/kodflow/ktn-linter")

	// Add rules and results
	f.addResults(run, fset, diagnostics)

	// Add run to report
	report.AddRun(run)

	// Write SARIF to output
	_ = report.Write(f.writer)
}

// addResults adds all diagnostic results to the SARIF run.
//
// Params:
//   - run: SARIF run to add results to
//   - fset: fileset for position information
//   - diagnostics: list of diagnostics
func (f *sarifFormatter) addResults(run *sarif.Run, fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	// Track seen rules for deduplication
	seenRules := make(map[string]bool, len(diagnostics))

	// Iterate over diagnostics
	for _, diag := range diagnostics {
		// Get position
		pos := fset.Position(diag.Pos)

		// Extract rule code
		code := extractCode(diag.Message)

		// Add rule if not seen
		if !seenRules[code] {
			f.addRule(run, code)
			seenRules[code] = true
		}

		// Get severity level
		level := severity.GetSeverity(code)
		sarifLevel := f.severityToSARIF(level)

		// Extract message
		message := extractMessageWithOptions(diag.Message, !f.verbose)

		// Create result
		result := sarif.NewRuleResult(code)
		result.Level = sarifLevel
		result.Message = sarif.NewTextMessage(message)

		// Create location
		location := sarif.NewLocation()
		physicalLocation := sarif.NewPhysicalLocation()
		physicalLocation.ArtifactLocation = sarif.NewSimpleArtifactLocation(pos.Filename)
		physicalLocation.Region = sarif.NewRegion()
		physicalLocation.Region.StartLine = &pos.Line
		physicalLocation.Region.StartColumn = &pos.Column

		location.PhysicalLocation = physicalLocation
		result.Locations = append(result.Locations, location)

		// Add result to run
		run.AddResult(result)
	}
}

// addRule adds a rule definition to the SARIF run.
//
// Params:
//   - run: SARIF run to add rule to
//   - code: rule code (e.g., "KTN-FUNC-001")
func (f *sarifFormatter) addRule(run *sarif.Run, code string) {
	// Get severity for rule
	level := severity.GetSeverity(code)

	// Create rule
	rule := sarif.NewRule(code)
	rule.ShortDescription = sarif.NewMultiformatMessageString()
	rule.ShortDescription.Text = &code
	rule.DefaultConfiguration = sarif.NewReportingConfiguration()
	rule.DefaultConfiguration.Level = f.severityToSARIF(level)

	// Add rule to driver
	run.Tool.Driver.Rules = append(run.Tool.Driver.Rules, rule)
}

// severityToSARIF converts severity level to SARIF level.
//
// Params:
//   - level: severity level
//
// Returns:
//   - string: SARIF level (error, warning, note)
func (f *sarifFormatter) severityToSARIF(level severity.Level) string {
	// Map severity to SARIF level
	switch level {
	// Error case
	case severity.SeverityError:
		// Return error level
		return "error"
	// Warning case
	case severity.SeverityWarning:
		// Return warning level
		return "warning"
	// Info case
	case severity.SeverityInfo:
		// Return note level
		return "note"
	}

	// Default to warning
	return "warning"
}
