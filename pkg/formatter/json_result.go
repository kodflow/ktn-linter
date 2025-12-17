// Package formatter provides output formatting for lint diagnostics.
package formatter

// JSONResult represents a single diagnostic result in JSON format.
// Contains rule identification, severity, message, and location.
type JSONResult struct {
	RuleID   string       `json:"ruleId"`
	Level    string       `json:"level"`
	Message  string       `json:"message"`
	Location JSONLocation `json:"location"`
}
