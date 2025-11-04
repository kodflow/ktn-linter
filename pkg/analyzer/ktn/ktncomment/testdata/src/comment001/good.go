package comment001

// goodMeaningfulComment explains why, not what.
func goodMeaningfulComment() int {
	// Using 42 as default timeout in seconds
	return 42
}

// goodExplainsReason provides context.
func goodExplainsReason() {
	// Preallocate to avoid reallocation during append
	x := 10
	_ = x
}

// goodNoComment is self-explanatory.
func goodNoComment() error {
	return nil
}

// goodContextualComment adds value.
func goodContextualComment() {
	i := 0
	// Skip first element as it's the header
	i++
	_ = i
}

// goodComplexLogic explains complex behavior.
func goodComplexLogic() {
	// Process only even numbers for performance optimization
	for i := 0; i < 10; i++ {
		_ = i
	}
}

// goodWhyNotWhat explains rationale.
func goodWhyNotWhat(x int) int {
	// Multiply by 2 to match API requirement
	return x * 2
}
