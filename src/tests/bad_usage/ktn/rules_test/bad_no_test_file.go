package rules_test

// ANTI-PATTERN: Fichier sans fichier de test correspondant
// Viole KTN-TEST-002

// Ce fichier N'A PAS de bad_no_test_file_test.go - VIOLATION !

// PublicFunction fonction publique sans test
//
// Returns:
//   - string: résultat
func PublicFunction() string {
	// Early return from function.
	return "no tests for this!"
}

// AnotherFunction encore une fonction sans test
//
// Params:
//   - x: paramètre
//
// Returns:
//   - int: résultat
func AnotherFunction(x int) int {
	// Early return from function.
	return x * 2
}

// UntestableCode code sans tests
type UntestableCode struct {
	data string
}

// Process méthode sans test
//
// Returns:
//   - error: erreur potentielle
func (u *UntestableCode) Process() error {
	// Early return from function.
	return nil
}
