// Package func003 contains test cases for KTN rules.
package func003

const (
	// Multiplier represents the multiplication factor
	Multiplier int = 2
	// LoopMax represents the maximum loop iterations
	LoopMax int = 10
	// Threshold represents the threshold value
	Threshold int = 10
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
		return val * Multiplier
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
// This function violates KTN-FUNC-003 by using else after continue.
func badLoopExample() {
	// Iterate from 0 to LoopMax
	for i := range LoopMax {
		// Check if i is even
		if i%Multiplier == 0 {
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
		if x > Threshold {
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

// badPanicExample illustre else inutile après panic.
//
// Params:
//   - x: valeur à vérifier
func badPanicExample(x int) {
	// Check if x is negative
	if x < 0 {
		// Panic for negative values
		panic("negative value")
	} else { // want "KTN-FUNC-003: else inutile après panic, utiliser early return"
		// Process positive values
		_ = x
	}
}

// badElseIfExample illustre else if inutile après return.
//
// Params:
//   - x: valeur à classifier
//
// Returns:
//   - string: catégorie de la valeur
func badElseIfExample(x int) string {
	// Check if negative
	if x < 0 {
		// Return negative category
		return "negative"
		// else if checks if zero
	} else if x == 0 {
		// Return zero category
		return "zero"
		// else handles positive
	} else {
		// Return positive category
		return "positive"
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badCheckPositive
	_ = badCheckPositive(0)
	// Appel de badProcessValue
	_ = badProcessValue(0)
	// Appel de badFindMax
	_ = badFindMax(1, 0)
	// Appel de badLoopExample
	badLoopExample()
	// Appel de badSwitchExample
	badSwitchExample(0)
	// Appel de badValidateInput
	_ = badValidateInput("")
	// Appel de badPanicExample
	badPanicExample(1)
	// Appel de badElseIfExample
	_ = badElseIfExample(0)
}
