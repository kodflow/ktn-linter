// Package orchestrator coordinates the linting pipeline.
package orchestrator

// Options holds all options for the linting orchestration.
// Contains verbose, category, rule, and config path settings.
type Options struct {
	Verbose    bool
	Category   string
	OnlyRule   string
	ConfigPath string
}
