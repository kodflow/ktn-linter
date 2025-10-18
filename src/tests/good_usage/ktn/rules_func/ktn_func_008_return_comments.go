package rules_func

import "errors"

// SimpleReturnWithComment retourne une valeur simple avec commentaire.
//
// Returns:
//   - string: message de bienvenue
func SimpleReturnWithComment() string {
	// Retourne le message de bienvenue standard
	return "Hello"
}

// MultipleReturnsWithComments a plusieurs returns avec commentaires.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - int: la valeur traitée
//   - error: erreur éventuelle
func MultipleReturnsWithComments(value int) (int, error) {
	if value < 0 {
		// Erreur car valeur négative non autorisée
		return 0, errors.New("negative value")
	}

	if value == 0 {
		// Retourne 0 sans erreur pour une valeur nulle
		return 0, nil
	}

	// Retourne le double de la valeur
	return value * 2, nil
}

// ComplexConditionsWithComments a des conditions complexes avec commentaires.
//
// Params:
//   - x: première valeur
//   - y: deuxième valeur
//
// Returns:
//   - string: résultat de la comparaison
func ComplexConditionsWithComments(x, y int) string {
	if x > y {
		if x > 100 {
			// Retourne "very large" car x dépasse largement y et 100
			return "very large"
		}
		// Retourne "large" car x est supérieur à y
		return "large"
	}

	if x == y {
		// Retourne "equal" car les valeurs sont identiques
		return "equal"
	}

	// Retourne "small" car x est inférieur à y
	return "small"
}

// SwitchWithReturnComments utilise switch avec commentaires sur returns.
//
// Params:
//   - status: le code de statut
//
// Returns:
//   - string: message correspondant
func SwitchWithReturnComments(status int) string {
	switch status {
	case 200:
		// Retourne "OK" pour un statut HTTP 200
		return "OK"
	case 404:
		// Retourne "Not Found" pour un statut HTTP 404
		return "Not Found"
	case 500:
		// Retourne "Internal Error" pour un statut HTTP 500
		return "Internal Error"
	default:
		// Retourne "Unknown" pour tous les autres statuts
		return "Unknown"
	}
}

// EarlyReturnWithComment fait un early return avec commentaire.
//
// Params:
//   - data: les données à valider
//
// Returns:
//   - error: erreur de validation
func EarlyReturnWithComment(data []string) error {
	if len(data) == 0 {
		// Erreur car les données sont vides
		return errors.New("empty data")
	}

	for _, item := range data {
		if item == "" {
			// Erreur car un élément est vide
			return errors.New("empty item")
		}
	}

	// Succès car toutes les validations sont passées
	return nil
}

// NestedReturnsWithComments a des returns imbriqués avec commentaires.
//
// Params:
//   - a: première valeur
//   - b: deuxième valeur
//   - c: troisième valeur
//
// Returns:
//   - bool: résultat de la vérification
func NestedReturnsWithComments(a, b, c int) bool {
	if a > 0 {
		if b > 0 {
			if c > 0 {
				// Retourne true car toutes les valeurs sont positives
				return true
			}
			// Retourne false car c n'est pas positif
			return false
		}
		// Retourne false car b n'est pas positif
		return false
	}
	// Retourne false car a n'est pas positif
	return false
}

// ErrorHandlingWithComments gère des erreurs avec commentaires.
//
// Params:
//   - input: valeur d'entrée
//
// Returns:
//   - int: résultat du calcul
//   - error: erreur éventuelle
func ErrorHandlingWithComments(input int) (int, error) {
	if input < 0 {
		// Erreur car l'entrée ne peut pas être négative
		return 0, errors.New("negative input")
	}

	if input > 1000 {
		// Erreur car l'entrée dépasse la limite maximale
		return 0, errors.New("input too large")
	}

	result := input * 2

	if result > 500 {
		// Retourne la moitié du résultat pour éviter les valeurs trop grandes
		return result / 2, nil
	}

	// Retourne le résultat normal (double de l'entrée)
	return result, nil
}

// BooleanLogicWithComments utilise logique booléenne avec commentaires.
//
// Params:
//   - enabled: si activé
//   - ready: si prêt
//
// Returns:
//   - string: statut
func BooleanLogicWithComments(enabled, ready bool) string {
	if enabled && ready {
		// Retourne "active" car activé et prêt
		return "active"
	}

	if enabled && !ready {
		// Retourne "pending" car activé mais pas encore prêt
		return "pending"
	}

	if !enabled && ready {
		// Retourne "disabled" car prêt mais désactivé
		return "disabled"
	}

	// Retourne "inactive" car ni activé ni prêt
	return "inactive"
}

// DivideWithComments effectue une division avec gestion d'erreur commentée.
//
// Params:
//   - a: numérateur
//   - b: dénominateur
//
// Returns:
//   - float64: résultat de la division
//   - error: erreur si division par zéro
func DivideWithComments(a, b int) (float64, error) {
	if b == 0 {
		// Erreur car division par zéro impossible
		return 0, errors.New("division by zero")
	}

	// Retourne le résultat de la division
	return float64(a) / float64(b), nil
}

// ValidateRangeWithComments valide qu'une valeur est dans une plage.
//
// Params:
//   - value: valeur à valider
//   - min: minimum inclusif
//   - max: maximum inclusif
//
// Returns:
//   - bool: true si dans la plage
func ValidateRangeWithComments(value, min, max int) bool {
	if value < min {
		// Retourne false car inférieur au minimum
		return false
	}

	if value > max {
		// Retourne false car supérieur au maximum
		return false
	}

	// Retourne true car dans la plage valide
	return true
}
