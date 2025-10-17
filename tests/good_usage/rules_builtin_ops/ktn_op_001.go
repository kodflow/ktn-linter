package rules_builtin_ops

import "errors"

// ✅ GOOD: vérifier avant division
func safeDiv(x, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("division by zero")
	}
	return x / y, nil
}

func safeMod(x, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("modulo by zero")
	}
	return x % y, nil
}
