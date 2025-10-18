package rules_return

// ✅ GOOD: return explicite dans fonction longue
func complexProcessGood() (int, error) {
	// ... beaucoup de code ...
	result := 0
	for i := 0; i < 100; i++ {
		if i == 50 {
			result = i
		}
	}
	// ... encore plus de code ...
	if !someConditionGood() {
		// Early return from function.
		return 0, nil
	}
	// Early return from function.
	return result, nil // ✅ explicite et clair!
}

// ✅ GOOD: naked return OK dans fonction courte
func shortFunc() (x int) {
	x = 42
	// Early return from function.
	return // ✅ OK: fonction courte
}

func someConditionGood() bool { return true }
