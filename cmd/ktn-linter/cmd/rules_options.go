// Package cmd implements the CLI commands for ktn-linter.
package cmd

// rulesOptions contains options for the rules command.
type rulesOptions struct {
	Format     string // Output format (markdown, json, text)
	NoExamples bool   // Whether to skip examples
}

// parseRulesOptions extracts options from command flags.
//
// Params:
//   - flags: flag getter for command flags
//
// Returns:
//   - rulesOptions: extracted options
func parseRulesOptions(flags flagGetter) rulesOptions {
	format, _ := flags.GetString(flagRulesFormat)
	noExamples, _ := flags.GetBool(flagRulesNoExamples)

	// Return parsed options
	return rulesOptions{
		Format:     format,
		NoExamples: noExamples,
	}
}
