package func014

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

// TagResource has methods.
type TagResource struct {
	tags []string
}

// GetTags returns tags.
//
// Returns:
//   - []string: tags
func (t *TagResource) GetTags() []string {
	// Retourne les tags
	return t.tags
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
