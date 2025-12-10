// Package func004 demonstrates violations of KTN-FUNC-004 rule.
// Contains examples of unused private functions that should trigger warnings.
package func004

// validateTagName is dead code created to bypass KTN-TEST-008.
// This is a private helper function to justify internal testing.
//
// Params:
//   - name: tag name to validate
//
// Returns:
//   - bool: true if valid, false otherwise
func validateTagName(name string) bool {
	// Cette fonction n'est JAMAIS appelée dans le code de production!
	return len(name) > 0
}

// unusedHelper is dead code created to bypass linting.
//
// Returns:
//   - string: message
func unusedHelper() string {
	// Jamais appelée!
	return "unused"
}

// formatData is dead code.
//
// Params:
//   - data: données
//
// Returns:
//   - string: formaté
func formatData(data string) string {
	// Jamais appelée!
	return "[" + data + "]"
}
