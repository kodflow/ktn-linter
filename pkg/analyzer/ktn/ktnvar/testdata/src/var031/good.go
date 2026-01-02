// Package var031 provides good test cases.
package var031

const (
	// CapacityTen is capacity value 10
	CapacityTen int = 10
	// MultiplierTwo is multiplier value 2
	MultiplierTwo int = 2
)

// init demonstrates correct usage patterns (compliant with KTN-VAR-031)
func init() {
	// Good: Transformation with multiplication (not a simple clone)
	original := make(map[string]int, CapacityTen)
	original["key"] = 1
	transformed := make(map[string]int, len(original))
	// This transforms the value, so it's not a simple clone
	for k, v := range original {
		transformed[k] = v * MultiplierTwo
	}
	_ = transformed

	// Good: Filter operation (not a simple clone)
	filterResult := make(map[string]int, len(original))
	// This filters values, so it's not a simple clone
	for k, v := range original {
		// Only include values greater than zero
		if v > 0 {
			filterResult[k] = v
		}
	}
	_ = filterResult

	// Good: Key transformation (not a simple clone)
	keyTransformed := make(map[string]int, len(original))
	// This transforms the key, so it's not a simple clone
	for k, v := range original {
		keyTransformed["prefix_"+k] = v
	}
	_ = keyTransformed

	// Good: Multiple statements in loop body (not a simple clone)
	multiStmt := make(map[string]int, len(original))
	// Multiple statements means it's doing more than just cloning
	for k, v := range original {
		multiStmt[k] = v
		// This extra statement prevents detection as simple clone
		_ = k
	}
	_ = multiStmt
}
