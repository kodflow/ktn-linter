package rules_comparison

// ✅ GOOD: comparaisons directes
func cleanCheck(isValid bool) bool {
	if isValid { // ✅ direct
		return true
	}
	return false
}

func negation(isActive bool) bool {
	if !isActive { // ✅ négation claire
		return false
	}
	return true
}

func directReturn(enabled bool) bool {
	return enabled // ✅ simple
}
