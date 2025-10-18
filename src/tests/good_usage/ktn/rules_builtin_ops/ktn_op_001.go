package rules_builtin_ops

import "errors"

// ✅ GOOD: vérifier avant division
func safeDiv(x, y int) (int, error) {
	if y == 0 {
		// Early return from function.
		return 0, errors.New("division by zero")
	}
	// Early return from function.
	return x / y, nil
}

func safeMod(x, y int) (int, error) {
	if y == 0 {
		// Early return from function.
		return 0, errors.New("modulo by zero")
	}
	// Early return from function.
	return x % y, nil
}
