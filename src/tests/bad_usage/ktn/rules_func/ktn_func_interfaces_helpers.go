package rules_func

// Violations intentionnelles pour tester les helpers d'interfaces

// bad_function fonction mal nommée (snake_case)
func bad_function() string {
	// Early return from function.
	return "test"
}

// MissingDoc fonction sans documentation
func MissingDoc(x int) int {
	// Early return from function.
	return x * 2
}
