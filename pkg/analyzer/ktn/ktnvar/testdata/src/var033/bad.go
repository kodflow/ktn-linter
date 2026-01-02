// Package var033 contains test cases for KTN-VAR-033 rule.
package var033

const (
	// BadDefaultPort represents the default port for testing.
	BadDefaultPort int = 8080
	// BadDefaultHost represents the default host for testing.
	BadDefaultHost string = "localhost"
)

// badGetPort returns port or default if zero using manual pattern.
//
// Params:
//   - port: port to check.
//
// Returns:
//   - int: port or default value.
func badGetPort(port int) int {
	// Check if port is non-zero
	if port != 0 { // want "KTN-VAR-033"
		// Return the provided port
		return port
	}
	// Return default port
	return BadDefaultPort
}

// badGetHost returns host or default if empty using manual pattern.
//
// Params:
//   - host: host to check.
//
// Returns:
//   - string: host or default value.
func badGetHost(host string) string {
	// Check if host is non-empty
	if host != "" { // want "KTN-VAR-033"
		// Return the provided host
		return host
	}
	// Return default host
	return BadDefaultHost
}

// badGetPointer returns ptr or default if nil using manual pattern.
//
// Params:
//   - ptr: pointer to check.
//   - defaultVal: default value if nil.
//
// Returns:
//   - *int: pointer or default value.
func badGetPointer(ptr *int, defaultVal *int) *int {
	// Check if ptr is not nil
	if ptr != nil { // want "KTN-VAR-033"
		// Return the provided pointer
		return ptr
	}
	// Return default pointer
	return defaultVal
}

// badGetSlice returns slice or default if nil using manual pattern.
//
// Params:
//   - slice: slice to check.
//   - defaultVal: default value if nil.
//
// Returns:
//   - []int: slice or default value.
func badGetSlice(slice []int, defaultVal []int) []int {
	// Check if slice is not nil
	if slice != nil { // want "KTN-VAR-033"
		// Return the provided slice
		return slice
	}
	// Return default slice
	return defaultVal
}

// init uses all the bad functions to avoid unused warnings.
func init() {
	// Test badGetPort function
	_ = badGetPort(0)

	// Test badGetHost function
	_ = badGetHost("")

	// Create test values
	val := 42
	defaultVal := 0

	// Test badGetPointer function
	_ = badGetPointer(&val, &defaultVal)

	// Test badGetSlice function
	_ = badGetSlice(nil, []int{1, 2, 3})
}
