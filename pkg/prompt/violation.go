// Package prompt provides AI-optimized prompt generation for KTN linter violations.
package prompt

// Violation represents a single rule violation with its location.
// Contains file path, line, column, and the diagnostic message.
type Violation struct {
	// FilePath is the absolute path to the file.
	FilePath string
	// Line is the line number where the violation occurs.
	Line int
	// Column is the column position of the violation.
	Column int
	// Message is the diagnostic message without the rule code prefix.
	Message string
}
