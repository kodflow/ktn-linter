package func003

// Edge case: main function should be exempt
func main() {
	// Main function is exempt from verb-prefix rule
	_ = 0
}

// Edge case: init function should be exempt
func init() {
	// Init function is exempt from verb-prefix rule
	_ = 0
}

// Edge case: function with empty name (tested via extractFirstWord with empty string)
// This case is indirectly tested when a function name has special characteristics
