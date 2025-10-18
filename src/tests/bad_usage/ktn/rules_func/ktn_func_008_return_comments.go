package rules_func

import "errors"

// simpleReturnWithoutComment retourne une valeur simple sans commentaire.
//
// Returns:
//   - string: message de bienvenue
func simpleReturnWithoutComment() string {
	// Early return from function.
	return "Hello" // ❌ KTN-FUNC-008: Pas de commentaire au-dessus du return
}

// multipleReturnsNoComments a plusieurs returns sans commentaires.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - int: la valeur traitée
//   - error: erreur éventuelle
func multipleReturnsNoComments(value int) (int, error) {
	if value < 0 {
		// Early return from function.
		return 0, errors.New("negative value") // ❌ KTN-FUNC-008: Pas de commentaire
	}

	if value == 0 {
		// Early return from function.
		return 0, nil // ❌ KTN-FUNC-008: Pas de commentaire
	}

	// Early return from function.
	return value * 2, nil // ❌ KTN-FUNC-008: Pas de commentaire
}

// complexConditionsWithoutComments a des conditions complexes sans commentaires.
//
// Params:
//   - x: première valeur
//   - y: deuxième valeur
//
// Returns:
//   - string: résultat de la comparaison
func complexConditionsWithoutComments(x, y int) string {
	if x > y {
		if x > 100 {
			// Early return from function.
			return "very large" // ❌ KTN-FUNC-008: Pas de commentaire
		}
		// Early return from function.
		return "large" // ❌ KTN-FUNC-008: Pas de commentaire
	}

	if x == y {
		// Early return from function.
		return "equal" // ❌ KTN-FUNC-008: Pas de commentaire
	}

	// Early return from function.
	return "small" // ❌ KTN-FUNC-008: Pas de commentaire
}

// switchWithoutReturnComments utilise switch sans commentaires sur returns.
//
// Params:
//   - status: le code de statut
//
// Returns:
//   - string: message correspondant
func switchWithoutReturnComments(status int) string {
	switch status {
	case 200:
		// Early return from function.
		return "OK" // ❌ KTN-FUNC-008: Pas de commentaire
	case 404:
		// Early return from function.
		return "Not Found" // ❌ KTN-FUNC-008: Pas de commentaire
	case 500:
		// Early return from function.
		return "Internal Error" // ❌ KTN-FUNC-008: Pas de commentaire
	default:
		// Early return from function.
		return "Unknown" // ❌ KTN-FUNC-008: Pas de commentaire
	}
}

// earlyReturnWithoutComment fait un early return sans commentaire.
//
// Params:
//   - data: les données à valider
//
// Returns:
//   - error: erreur de validation
func earlyReturnWithoutComment(data []string) error {
	if len(data) == 0 {
		// Return error to caller.
		return errors.New("empty data") // ❌ KTN-FUNC-008: Pas de commentaire
	}

	for _, item := range data {
		if item == "" {
			// Return error to caller.
			return errors.New("empty item") // ❌ KTN-FUNC-008: Pas de commentaire
		}
	}

	// Early return from function.
	return nil // ❌ KTN-FUNC-008: Pas de commentaire
}

// nestedReturnsWithoutComments a des returns imbriqués sans commentaires.
//
// Params:
//   - a: première valeur
//   - b: deuxième valeur
//   - c: troisième valeur
//
// Returns:
//   - bool: résultat de la vérification
func nestedReturnsWithoutComments(a, b, c int) bool {
	if a > 0 {
		if b > 0 {
			if c > 0 {
				// Continue inspection/processing.
				return true // ❌ KTN-FUNC-008: Pas de commentaire
			}
			// Stop inspection/processing.
			return false // ❌ KTN-FUNC-008: Pas de commentaire
		}
		// Stop inspection/processing.
		return false // ❌ KTN-FUNC-008: Pas de commentaire
	}
	// Stop inspection/processing.
	return false // ❌ KTN-FUNC-008: Pas de commentaire
}

// errorHandlingWithoutComments gère des erreurs sans commentaires.
//
// Params:
//   - input: valeur d'entrée
//
// Returns:
//   - int: résultat du calcul
//   - error: erreur éventuelle
func errorHandlingWithoutComments(input int) (int, error) {
	if input < 0 {
		// Early return from function.
		return 0, errors.New("negative input") // ❌ KTN-FUNC-008: Pas de commentaire
	}

	if input > 1000 {
		// Early return from function.
		return 0, errors.New("input too large") // ❌ KTN-FUNC-008: Pas de commentaire
	}

	result := input * 2

	if result > 500 {
		// Early return from function.
		return result / 2, nil // ❌ KTN-FUNC-008: Pas de commentaire
	}

	// Early return from function.
	return result, nil // ❌ KTN-FUNC-008: Pas de commentaire
}

// booleanLogicWithoutComments utilise logique booléenne sans commentaires.
//
// Params:
//   - enabled: si activé
//   - ready: si prêt
//
// Returns:
//   - string: statut
func booleanLogicWithoutComments(enabled, ready bool) string {
	if enabled && ready {
		// Early return from function.
		return "active" // ❌ KTN-FUNC-008: Pas de commentaire
	}

	if enabled && !ready {
		// Early return from function.
		return "pending" // ❌ KTN-FUNC-008: Pas de commentaire
	}

	if !enabled && ready {
		// Early return from function.
		return "disabled" // ❌ KTN-FUNC-008: Pas de commentaire
	}

	// Early return from function.
	return "inactive" // ❌ KTN-FUNC-008: Pas de commentaire
}
