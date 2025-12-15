// Package func012 contains test cases for KTN rules.
package func012

const (
	// TWO_COUNT représente la valeur 2 pour les compteurs
	TWO_COUNT int = 2
	// THIRTY_AGE représente l'âge 30
	THIRTY_AGE int = 30
	// NINETY_FIVE_SCORE représente le score 95.5
	NINETY_FIVE_SCORE float64 = 95.5
)

// FourUnnamedReturns demonstrates a function with 4 unnamed return values (violates KTN-FUNC-012).
//
// Returns:
//   - int: identifier
//   - string: message
//   - bool: status
//   - error: operation error
func FourUnnamedReturns() (int, string, bool, error) {
	// Return successful test data with no error
	return 1, "test", true, nil
}

// FiveUnnamedReturns demonstrates a function with 5 unnamed return values (violates KTN-FUNC-012).
//
// Returns:
//   - int: identifier
//   - int: count value
//   - string: message
//   - bool: status
//   - error: operation error
func FiveUnnamedReturns() (int, int, string, bool, error) {
	// Return successful test data with counts and no error
	return 1, TWO_COUNT, "test", true, nil
}

// ManyUnnamedReturns demonstrates a function with 6 unnamed return values (violates KTN-FUNC-012).
//
// Returns:
//   - int: identifier
//   - string: message
//   - int: age value
//   - bool: status
//   - float64: score value
//   - error: operation error
func ManyUnnamedReturns() (int, string, int, bool, float64, error) {
	// Return complete test data including score with no error
	return 1, "test", THIRTY_AGE, true, NINETY_FIVE_SCORE, nil
}

// SevenUnnamedReturns demonstrates a function with 7 unnamed return values (violates KTN-FUNC-012).
//
// Returns:
//   - int: identifier
//   - string: message
//   - int: age value
//   - bool: status
//   - float64: score value
//   - string: extra information
//   - error: operation error
func SevenUnnamedReturns() (int, string, int, bool, float64, string, error) {
	// Return full test data with extra field and no error
	return 1, "test", THIRTY_AGE, true, NINETY_FIVE_SCORE, "extra", nil
}
