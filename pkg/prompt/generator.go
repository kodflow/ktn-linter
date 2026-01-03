// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

import (
	"io"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"github.com/kodflow/ktn-linter/pkg/rules"
)

// ktnPrefix is the standard prefix for rule codes.
const ktnPrefix string = "KTN-"

// Generator generates AI-optimized prompts from linter diagnostics.
// Coordinates orchestrator, rule metadata extraction, and phase classification.
type Generator struct {
	orch    *orchestrator.Orchestrator
	verbose bool
}

// NewGenerator creates a new prompt generator.
//
// Params:
//   - stderr: writer for error output
//   - verbose: enable verbose logging
//
// Returns:
//   - *Generator: new generator instance
func NewGenerator(stderr io.Writer, verbose bool) *Generator {
	// Create orchestrator
	orch := orchestrator.NewOrchestrator(stderr, verbose)

	// Return generator
	return &Generator{
		orch:    orch,
		verbose: verbose,
	}
}

// Generate runs the linter and generates a prompt output.
//
// Params:
//   - patterns: package patterns to analyze
//   - opts: orchestrator options
//
// Returns:
//   - *PromptOutput: generated prompt output
//   - error: generation error if any
func (g *Generator) Generate(patterns []string, opts orchestrator.Options) (*PromptOutput, error) {
	// Run linter pipeline
	diagnostics, err := g.runLinter(patterns, opts)
	// Check for error
	if err != nil {
		// Return nil for linter error
		return nil, err
	}

	// Group violations by rule
	ruleViolations := g.collectViolations(diagnostics)

	// Enrich with metadata
	enriched := g.enrichWithMetadata(ruleViolations)

	// Sort into phases
	phases := SortRulesByPhase(enriched)

	// Build output
	output := g.buildOutput(enriched, phases)

	// Return generated output
	return output, nil
}

// runLinter executes the linter and returns raw diagnostics.
//
// Params:
//   - patterns: package patterns to analyze
//   - opts: orchestrator options
//
// Returns:
//   - []orchestrator.DiagnosticResult: raw diagnostics
//   - error: linter error if any
func (g *Generator) runLinter(patterns []string, opts orchestrator.Options) ([]orchestrator.DiagnosticResult, error) {
	// Load packages
	pkgs, err := g.orch.LoadPackages(patterns)
	// Check for error
	if err != nil {
		var emptyDiags []orchestrator.DiagnosticResult
		// Return empty slice for load error
		return emptyDiags, err
	}

	// Select all analyzers (ignore filters for prompt)
	analyzers, err := g.orch.SelectAnalyzers(opts)
	// Check for error
	if err != nil {
		var emptyDiags []orchestrator.DiagnosticResult
		// Return empty slice for analyzer selection error
		return emptyDiags, err
	}

	// Run analyzers
	rawDiags := g.orch.RunAnalyzers(pkgs, analyzers)

	// Filter diagnostics
	filtered := g.orch.FilterDiagnostics(rawDiags)

	// Return filtered diagnostics
	return filtered, nil
}

// collectViolations groups diagnostics by rule code.
//
// Params:
//   - diagnostics: raw diagnostics from linter
//
// Returns:
//   - map[string]*RuleViolations: violations grouped by rule code
func (g *Generator) collectViolations(diagnostics []orchestrator.DiagnosticResult) map[string]*RuleViolations {
	// Preallocate map with estimated capacity
	result := make(map[string]*RuleViolations, len(diagnostics))

	// Process each diagnostic
	for i := range diagnostics {
		diag := &diagnostics[i]
		// Extract rule code from message
		code := extractRuleCode(diag.Diag.Message)
		// If no KTN code, use analyzer name for modernize suggestions
		if code == "" {
			code = buildModernizeCode(diag.AnalyzerName)
		}
		// Skip if still no code
		if code == "" {
			continue
		}

		// Get or create rule violations
		rv, exists := result[code]
		// Initialize if not exists
		if !exists {
			var emptyViolations []Violation
			rv = &RuleViolations{
				Code:       code,
				Violations: emptyViolations,
			}
			result[code] = rv
		}

		// Add violation
		pos := diag.Position()
		violation := Violation{
			FilePath: pos.Filename,
			Line:     pos.Line,
			Column:   pos.Column,
			Message:  extractMessage(diag.Diag.Message),
		}
		rv.Violations = append(rv.Violations, violation)
	}

	// Return collected violations
	return result
}

// enrichWithMetadata adds description and examples to violations.
//
// Params:
//   - violations: map of rule violations
//
// Returns:
//   - []RuleViolations: enriched violations slice
func (g *Generator) enrichWithMetadata(violations map[string]*RuleViolations) []RuleViolations {
	result := make([]RuleViolations, 0, len(violations))

	// Enrich each rule
	for code, rv := range violations {
		// Get rule info
		info := rules.GetRuleInfoByCode(code)
		// Enrich if found
		if info != nil {
			rv.Category = info.Category
			rv.Description = info.Description
			rv.GoodExample = rules.LoadGoodExample(code)
		}

		// Append to result
		result = append(result, *rv)
	}

	// Return enriched violations
	return result
}

// buildOutput constructs the final PromptOutput.
//
// Params:
//   - rules: enriched rule violations
//   - phases: sorted phase groups
//
// Returns:
//   - *PromptOutput: complete output structure
func (g *Generator) buildOutput(rulesList []RuleViolations, phases []PhaseGroup) *PromptOutput {
	// Count total violations
	totalViolations := 0
	// Sum violations from all rules
	for _, rv := range rulesList {
		totalViolations += len(rv.Violations)
	}

	// Return complete output structure
	return &PromptOutput{
		TotalViolations: totalViolations,
		TotalRules:      len(rulesList),
		Phases:          phases,
	}
}

// buildModernizeCode converts an analyzer name to a KTN-MODERNIZE code.
//
// Params:
//   - analyzerName: name of the analyzer
//
// Returns:
//   - string: KTN-MODERNIZE-XXX code or empty if not a modernize analyzer
func buildModernizeCode(analyzerName string) string {
	// Map of known modernize analyzer names to codes
	modernizeCodes := map[string]string{
		"stringscut":       "KTN-MODERNIZE-001", // strings.Cut
		"stringscutprefix": "KTN-MODERNIZE-002", // CutPrefix/CutSuffix
		"any":              "KTN-MODERNIZE-003", // interface{} → any
		"minmax":           "KTN-MODERNIZE-004", // min/max builtins
		"slicescontains":   "KTN-MODERNIZE-005", // slices.Contains
		"slicessort":       "KTN-MODERNIZE-006", // slices.Sort
		"fmtappendf":       "KTN-MODERNIZE-007", // fmt.Appendf
		"rangeint":         "KTN-MODERNIZE-008", // range over int
		"bloop":            "KTN-MODERNIZE-009", // b.Loop
		"waitgroup":        "KTN-MODERNIZE-010", // wg.Go
		"omitzero":         "KTN-MODERNIZE-011", // omitzero tag
		"testingcontext":   "KTN-MODERNIZE-012", // t.Context
		"stringsbuilder":   "KTN-MODERNIZE-013", // strings.Builder
		"mapsloop":         "KTN-MODERNIZE-014", // maps package
		"stditerators":     "KTN-MODERNIZE-015", // iterators
		"stringsseq":       "KTN-MODERNIZE-016", // SplitSeq/FieldsSeq
		"unsafefuncs":      "KTN-MODERNIZE-017", // unsafe functions
		"reflecttypefor":   "KTN-MODERNIZE-018", // TypeFor
		"plusbuild":        "KTN-MODERNIZE-019", // obsolete build comments
		"forvar":           "KTN-MODERNIZE-020", // loop variables
	}

	// Check if it's a known modernize analyzer
	if code, ok := modernizeCodes[analyzerName]; ok {
		// Return the modernize code
		return code
	}

	// Return empty for unknown analyzers
	return ""
}

// extractRuleCode extracts KTN-XXX-YYY from a diagnostic message.
//
// Params:
//   - message: diagnostic message
//
// Returns:
//   - string: rule code or empty if not found
func extractRuleCode(message string) string {
	// Check for bracketed format [KTN-XXX-YYY]
	if strings.HasPrefix(message, "["+ktnPrefix) {
		endIdx := strings.Index(message, "]")
		// Check if bracket found
		if endIdx > 0 {
			// Return extracted rule code from brackets
			return message[1:endIdx]
		}
	}

	// Check for colon format KTN-XXX-YYY:
	if strings.HasPrefix(message, ktnPrefix) {
		// Try colon separator first
		if code, _, found := strings.Cut(message, ":"); found {
			// Return rule code before colon
			return code
		}

		// Try space separator
		if code, _, found := strings.Cut(message, " "); found {
			// Return rule code before space
			return code
		}
	}

	// Return empty string when no rule code found
	return ""
}

// extractMessage extracts the message part after the rule code.
//
// Params:
//   - message: full diagnostic message
//
// Returns:
//   - string: message without rule code prefix
func extractMessage(message string) string {
	// Check for bracketed format [KTN-XXX-YYY]
	if strings.HasPrefix(message, "["+ktnPrefix) {
		endIdx := strings.Index(message, "]")
		// Check if bracket found
		if endIdx > 0 && len(message) > endIdx+1 {
			// Return message after bracket
			return strings.TrimSpace(message[endIdx+1:])
		}
	}

	// Check for colon format KTN-XXX-YYY: message
	if strings.HasPrefix(message, ktnPrefix) {
		colonIdx := strings.Index(message, ":")
		// Check if colon found
		if colonIdx > 0 && len(message) > colonIdx+1 {
			// Return message after colon
			return strings.TrimSpace(message[colonIdx+1:])
		}
	}

	// Return original message when no prefix found
	return message
}
