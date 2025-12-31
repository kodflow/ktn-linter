// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

import (
	"fmt"
	"io"
)

// MarkdownFormatter formats prompt output as markdown.
// Writes structured output with phases, rules, and checkboxes.
type MarkdownFormatter struct {
	writer io.Writer
}

// NewMarkdownFormatter creates a new markdown formatter.
//
// Params:
//   - w: writer for output
//
// Returns:
//   - *MarkdownFormatter: new formatter instance
func NewMarkdownFormatter(w io.Writer) *MarkdownFormatter {
	// Return new formatter instance
	return &MarkdownFormatter{writer: w}
}

// Format writes the complete prompt output as markdown.
// Implements phase dependencies: only shows active phase, blocks subsequent phases.
//
// Params:
//   - output: prompt output to format
func (f *MarkdownFormatter) Format(output *PromptOutput) {
	// Guard against nil receiver, writer, or output
	if f == nil || f.writer == nil || output == nil {
		// Return early to avoid nil pointer dereference
		return
	}

	// Write header
	f.writeHeader(output)

	// Write instructions
	f.writeInstructions()

	// Track if previous phases have unresolved violations
	previousPhasesHaveViolations := false

	// Write each phase with dependency logic
	for i := range output.Phases {
		phase := &output.Phases[i]
		phaseNum := i + 1

		// Count violations in this phase
		phaseViolationCount := countPhaseViolations(phase)

		// Check if this phase should be blocked
		if previousPhasesHaveViolations {
			// Show blocked message for this phase
			f.writeBlockedPhase(phase, phaseNum, phaseViolationCount)
		} else {
			// Show full phase details
			f.writePhase(phase, phaseNum)
		}

		// Update tracking: if this phase has violations, block subsequent phases
		if phaseViolationCount > 0 {
			previousPhasesHaveViolations = true
		}
	}
}

// countPhaseViolations counts total violations in a phase.
//
// Params:
//   - phase: phase group to count
//
// Returns:
//   - int: total violation count
func countPhaseViolations(phase *PhaseGroup) int {
	// Guard against nil phase
	if phase == nil {
		// Return zero for nil phase
		return 0
	}

	count := 0
	// Sum violations from all rules
	for i := range phase.Rules {
		count += len(phase.Rules[i].Violations)
	}
	// Return total count
	return count
}

// writeHeader writes the markdown header with summary.
//
// Params:
//   - output: prompt output for summary stats
func (f *MarkdownFormatter) writeHeader(output *PromptOutput) {
	// Main title
	fmt.Fprintln(f.writer, "# KTN-Linter Correction Prompt")
	fmt.Fprintln(f.writer)

	// Summary line
	fmt.Fprintf(f.writer, "**Total**: %d violations across %d rules\n\n",
		output.TotalViolations, output.TotalRules)
}

// writeInstructions writes the instructions section.
func (f *MarkdownFormatter) writeInstructions() {
	fmt.Fprintln(f.writer, "## Instructions")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "Ce prompt guide la correction des violations KTN. **Suivez les phases dans l'ordre strict.**")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "### Workflow")
	fmt.Fprintln(f.writer, "1. Corrigez toutes les violations de la phase active")
	fmt.Fprintln(f.writer, "2. Re-executez le linter: `ktn-linter prompt ./...`")
	fmt.Fprintln(f.writer, "3. La phase suivante sera alors debloquee")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "### Pourquoi cet ordre?")
	fmt.Fprintln(f.writer, "- **Phase 1-2**: Peuvent creer/supprimer/deplacer des fichiers")
	fmt.Fprintln(f.writer, "- **Phase 3**: Modifications locales dans les fichiers existants")
	fmt.Fprintln(f.writer, "- **Phase 4**: Commentaires (en dernier, quand le code est finalise)")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "Les phases bloquees (⏸️) indiquent des violations en attente.")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "---")
	fmt.Fprintln(f.writer)
}

// writePhase writes a single phase group.
//
// Params:
//   - phase: phase group to write (pointer for efficiency)
//   - phaseNum: phase number for display
func (f *MarkdownFormatter) writePhase(phase *PhaseGroup, phaseNum int) {
	// Guard against nil receiver, writer, or phase
	if f == nil || f.writer == nil || phase == nil {
		// Return early to avoid nil pointer dereference
		return
	}

	// Phase header
	fmt.Fprintf(f.writer, "## Phase %d: %s\n\n", phaseNum, phase.Name)

	// Phase description
	fmt.Fprintf(f.writer, "%s\n\n", phase.Description)

	// Re-run warning if needed
	if phase.NeedsRerun {
		fmt.Fprintln(f.writer, "> **Re-executez le linter apres cette phase**")
		fmt.Fprintln(f.writer)
	}

	// Write each rule
	for i := range phase.Rules {
		f.writeRule(&phase.Rules[i])
	}

	// Phase separator
	fmt.Fprintln(f.writer, "---")
	fmt.Fprintln(f.writer)
}

// writeBlockedPhase writes a blocked phase with summary only.
// Shows that violations exist but details are hidden until previous phases complete.
//
// Params:
//   - phase: phase group to write (pointer for efficiency)
//   - phaseNum: phase number for display
//   - violationCount: number of violations in this phase
func (f *MarkdownFormatter) writeBlockedPhase(phase *PhaseGroup, phaseNum int, violationCount int) {
	// Guard against nil receiver, writer, or phase
	if f == nil || f.writer == nil || phase == nil {
		// Return early to avoid nil pointer dereference
		return
	}

	// Phase header with blocked indicator
	fmt.Fprintf(f.writer, "## Phase %d: %s ⏸️ BLOQUEE\n\n", phaseNum, phase.Name)

	// Show violation count without details
	if violationCount > 0 {
		fmt.Fprintf(f.writer, "**%d violations** en attente dans cette phase.\n\n", violationCount)
		// List rule codes only (no details)
		fmt.Fprintln(f.writer, "Regles concernees:")
		// Iterate over rules
		for i := range phase.Rules {
			rule := &phase.Rules[i]
			fmt.Fprintf(f.writer, "- %s (%d violations)\n", rule.Code, len(rule.Violations))
		}
		fmt.Fprintln(f.writer)
	}

	// Blocking message
	fmt.Fprintln(f.writer, "> ⚠️ **Completez les phases precedentes avant de traiter celle-ci.**")
	fmt.Fprintln(f.writer, "> Les fichiers peuvent changer suite aux corrections structurelles.")
	fmt.Fprintln(f.writer)

	// Phase separator
	fmt.Fprintln(f.writer, "---")
	fmt.Fprintln(f.writer)
}

// writeRule writes a single rule with its violations.
//
// Params:
//   - rule: rule violations to write (pointer for efficiency)
func (f *MarkdownFormatter) writeRule(rule *RuleViolations) {
	// Guard against nil receiver, writer, or rule
	if f == nil || f.writer == nil || rule == nil {
		// Return early to avoid nil pointer dereference
		return
	}

	// Rule header
	fmt.Fprintf(f.writer, "### %s\n\n", rule.Code)

	// Category
	fmt.Fprintf(f.writer, "**Category**: %s\n\n", rule.Category)

	// Description
	fmt.Fprintf(f.writer, "**Description**: %s\n\n", rule.Description)

	// Good example if available
	if rule.GoodExample != "" {
		f.writeGoodExample(rule.GoodExample)
	}

	// Violations list
	f.writeViolations(rule.Violations)
}

// writeGoodExample writes the good example code block.
//
// Params:
//   - example: example code content
func (f *MarkdownFormatter) writeGoodExample(example string) {
	fmt.Fprintln(f.writer, "#### Good Example")
	fmt.Fprintln(f.writer)
	fmt.Fprintln(f.writer, "```go")
	fmt.Fprint(f.writer, example)
	// Ensure newline at end
	if len(example) > 0 && example[len(example)-1] != '\n' {
		fmt.Fprintln(f.writer)
	}
	fmt.Fprintln(f.writer, "```")
	fmt.Fprintln(f.writer)
}

// writeViolations writes the violations as a checkbox list.
//
// Params:
//   - violations: list of violations
func (f *MarkdownFormatter) writeViolations(violations []Violation) {
	// Violations header
	fmt.Fprintf(f.writer, "#### Files to Fix (%d violations)\n\n", len(violations))

	// Write each violation as checkbox
	for _, v := range violations {
		f.writeViolation(v)
	}

	fmt.Fprintln(f.writer)
}

// writeViolation writes a single violation as a checkbox item.
//
// Params:
//   - v: violation to write
func (f *MarkdownFormatter) writeViolation(v Violation) {
	// Format: - [ ] `path/file.go:line` - message
	if v.Message != "" {
		fmt.Fprintf(f.writer, "- [ ] `%s:%d` - %s\n", v.FilePath, v.Line, v.Message)
	} else {
		// Handle empty message case
		fmt.Fprintf(f.writer, "- [ ] `%s:%d`\n", v.FilePath, v.Line)
	}
}
