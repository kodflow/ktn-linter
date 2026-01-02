// Package var031 contains test cases for KTN-VAR-031 rule.
package var031

const (
	// BadCapacityTen is capacity value 10
	BadCapacityTen int = 10
	// BadPiValue is approximate pi value
	BadPiValue float64 = 3.14
)

// badCloneStringMap demonstrates manual map cloning pattern.
//
// Params:
//   - original: map to clone
//
// Returns:
//   - map[string]int: cloned map
func badCloneStringMap(original map[string]int) map[string]int {
	// Manual clone pattern - should use maps.Clone
	clone := make(map[string]int, len(original))
	// Range over original and copy each key-value pair
	for k, v := range original {
		clone[k] = v
	}
	// Return cloned map
	return clone
}

// badCloneBoolMap demonstrates manual bool map cloning.
//
// Params:
//   - source: map to clone
//
// Returns:
//   - map[int]bool: cloned map
func badCloneBoolMap(source map[int]bool) map[int]bool {
	// Manual clone pattern - should use maps.Clone
	dest := make(map[int]bool, len(source))
	// Range over source and copy each key-value pair
	for k, v := range source {
		dest[k] = v
	}
	// Return cloned map
	return dest
}

// badCloneFloat64Map demonstrates manual float64 map cloning.
//
// Params:
//   - data: map to clone
//
// Returns:
//   - map[string]float64: cloned map
func badCloneFloat64Map(data map[string]float64) map[string]float64 {
	// Manual clone pattern - should use maps.Clone
	result := make(map[string]float64, len(data))
	// Range over data and copy each key-value pair
	for key, value := range data {
		result[key] = value
	}
	// Return cloned map
	return result
}

// init uses the bad functions to avoid unused warnings.
func init() {
	// Create test map with capacity
	testMap := make(map[string]int, BadCapacityTen)
	testMap["a"] = 1
	// Call badCloneStringMap
	_ = badCloneStringMap(testMap)
	// Create bool test map
	boolMap := make(map[int]bool, BadCapacityTen)
	boolMap[1] = true
	// Call badCloneBoolMap
	_ = badCloneBoolMap(boolMap)
	// Create float64 test map
	floatMap := make(map[string]float64, BadCapacityTen)
	floatMap["pi"] = BadPiValue
	// Call badCloneFloat64Map
	_ = badCloneFloat64Map(floatMap)
}
