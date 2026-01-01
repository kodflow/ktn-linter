// Package var033 contains test cases for KTN-VAR-033 rule.
package var033

import "cmp"

const (
	// DefaultPort represents the default port for testing.
	DefaultPort int = 8080
	// DefaultHost represents the default host for testing.
	DefaultHost string = "localhost"
	// SampleThreshold represents a sample threshold for testing.
	SampleThreshold int = 10
)

// goodGetPort returns port or default using cmp.Or.
//
// Params:
//   - port: port to check.
//
// Returns:
//   - int: port or default value.
func goodGetPort(port int) int {
	// Use cmp.Or for idiomatic zero-value check
	return cmp.Or(port, DefaultPort)
}

// goodGetHost returns host or default using cmp.Or.
//
// Params:
//   - host: host to check.
//
// Returns:
//   - string: host or default value.
func goodGetHost(host string) string {
	// Use cmp.Or for idiomatic empty-string check
	return cmp.Or(host, DefaultHost)
}

// goodComplexCondition performs a complex check (not zero value).
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: val or default based on complex condition.
func goodComplexCondition(val int) int {
	// Check if val is greater than threshold (not zero value check)
	if val > SampleThreshold {
		// Return the value
		return val
	}
	// Return default value
	return 5
}

// goodWithElse has an else clause (not simple cmp.Or pattern).
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: result based on condition.
func goodWithElse(val int) int {
	// Check if val is non-zero with else clause
	if val != 0 {
		// Return the value
		return val
	} else {
		// Else clause makes it different from cmp.Or pattern
		return val + 1
	}
}

// goodEqualCheck uses equality instead of not-equal.
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: result based on condition.
func goodEqualCheck(val int) int {
	// Check if val equals zero (equality, not != pattern)
	if val == 0 {
		// Return default
		return DefaultPort
	}
	// Return the value
	return val
}

// goodDifferentReturn returns different value than checked variable.
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: result based on condition.
func goodDifferentReturn(val int) int {
	// Check if val is non-zero but return different value
	if val != 0 {
		// Return transformed value (not the same variable)
		return val * 2
	}
	// Return default
	return DefaultPort
}

// goodWithInit has init statement in if (too complex).
//
// Params:
//   - vals: slice of values.
//
// Returns:
//   - int: result based on condition.
func goodWithInit(vals []int) int {
	// If with init statement (too complex for cmp.Or)
	if val := getFirst(vals); val != 0 {
		// Return the value
		return val
	}
	// Return default
	return DefaultPort
}

// getFirst returns first element or zero.
//
// Params:
//   - vals: slice of values.
//
// Returns:
//   - int: first element or zero.
func getFirst(vals []int) int {
	// Check if slice is not empty
	if len(vals) > 0 {
		// Return first element
		return vals[0]
	}
	// Return zero for empty slice
	return 0
}

// goodMultipleStatements has multiple statements in if body.
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: result based on condition.
func goodMultipleStatements(val int) int {
	// Check if val is non-zero
	if val != 0 {
		// Multiple statements in body (not simple cmp.Or pattern)
		_ = val
		// Return the value
		return val
	}
	// Return default
	return DefaultPort
}

// goodNoReturnAfter has different statement after if.
//
// Params:
//   - val: value to check.
//
// Returns:
//   - int: always returns default.
func goodNoReturnAfter(val int) int {
	// Check if val is non-zero
	if val != 0 {
		// Return the value
		return val
	}
	// Multiple statements before return (not simple pattern)
	_ = val
	// Return default
	return DefaultPort
}

// init uses all the good functions to avoid unused warnings.
func init() {
	// Test goodGetPort function
	_ = goodGetPort(0)

	// Test goodGetHost function
	_ = goodGetHost("")

	// Test goodComplexCondition function
	_ = goodComplexCondition(5)

	// Test goodWithElse function
	_ = goodWithElse(0)

	// Test goodEqualCheck function
	_ = goodEqualCheck(0)

	// Test goodDifferentReturn function
	_ = goodDifferentReturn(0)

	// Test goodWithInit function
	_ = goodWithInit([]int{1, 2, 3})

	// Test goodMultipleStatements function
	_ = goodMultipleStatements(0)

	// Test goodNoReturnAfter function
	_ = goodNoReturnAfter(0)
}
