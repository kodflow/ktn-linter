package rules_func

// Violations intentionnelles pour tester les helpers d'interfaces

// bad_function fonction mal nommée (snake_case)
func bad_function() string {
	return "test"
}

// MissingDoc fonction sans documentation
func MissingDoc(x int) int {
	return x * 2
}
