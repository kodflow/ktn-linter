package func012

const (
	// MULTIPLIER represents the multiplication factor
	MULTIPLIER int = 2
	// LOOP_MAX represents the maximum loop iterations
	LOOP_MAX int = 10
	// THRESHOLD represents the threshold value
	THRESHOLD int = 10
)

// badCheckPositive vérifie si un nombre est positif avec else inutile.
//
// Params:
//   - x: nombre à vérifier
//
// Returns:
//   - string: "positive" ou "negative"
func badCheckPositive(x int) string {
	// Check if number is positive
	if x > 0 {
		// Return positive case
		return "positive"
	} else {
		// Return negative case
		return "negative"
	}
}

// badProcessValue traite une valeur avec else inutile après return.
//
// Params:
//   - val: valeur à traiter
//
// Returns:
//   - int: 0 si négatif, sinon val doublée
func badProcessValue(val int) int {
	// Check if value is negative
	if val < 0 {
		// Return zero for negative values
		return 0
	} else {
		// Return multiplied value
		return val * MULTIPLIER
	}
}

// badFindMax trouve le maximum avec else inutile après return.
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: le maximum des deux
func badFindMax(a, b int) int {
	// Check if a is greater than b
	if a > b {
		// Return a if it's larger
		return a
	} else {
		// Return b otherwise
		return b
	}
}

// badLoopExample demonstrates a loop with unnecessary else after continue.
// This function violates KTN-FUNC-012 by using else after continue.
func badLoopExample() {
	// Iterate from 0 to LOOP_MAX
	for i := 0; i < LOOP_MAX; i++ {
		// Check if i is even
		if i%MULTIPLIER == 0 {
			// Skip even numbers
			continue
		} else {
			// Process odd numbers
			_ = i
		}
	}
}

// badSwitchExample illustre else inutile après break.
//
// Params:
//   - x: valeur à traiter
func badSwitchExample(x int) {
	// Loop until x exceeds threshold
	for {
		// Check if threshold exceeded
		if x > THRESHOLD {
			// Exit loop when threshold exceeded
			break
		} else {
			// Increment x otherwise
			x++
		}
	}
}

// badValidateInput valide une entrée avec else inutile après return.
//
// Params:
//   - input: chaîne à valider
//
// Returns:
//   - error: erreur ou nil
func badValidateInput(input string) error {
	// Check if input is empty
	if input == "" {
		// Return nil for empty input
		return nil
	} else {
		// Return nil for non-empty input
		return nil
	}
}
