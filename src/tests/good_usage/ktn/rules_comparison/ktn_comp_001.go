package rules_comparison

// ✅ GOOD: comparaisons directes
func cleanCheck(isValid bool) bool {
	if isValid { // ✅ direct
		// Continue inspection/processing.
		return true
	}
	// Stop inspection/processing.
	return false
}

func negation(isActive bool) bool {
	if !isActive { // ✅ négation claire
		// Stop inspection/processing.
		return false
	}
	// Continue inspection/processing.
	return true
}

func directReturn(enabled bool) bool {
	// Early return from function.
	return enabled // ✅ simple
}
