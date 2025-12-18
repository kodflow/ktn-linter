// Mock helper for testing.
package test013

// MockHelper is a mock helper.
//
// Params:
//   - input: input value
//
// Returns:
//   - string: output value
//   - error: any error
func MockHelper(input string) (string, error) {
	// Return mock result
	return "mock:" + input, nil
}
