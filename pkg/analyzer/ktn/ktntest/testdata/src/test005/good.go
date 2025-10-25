package test005

import (
	"errors"
	"strings"
)

const (
	// MIN_EMAIL_LENGTH longueur minimale d'un email
	MIN_EMAIL_LENGTH int = 3
	// EMAIL_PARTS_COUNT nombre de parties d'un email (avant/après @)
	EMAIL_PARTS_COUNT int = 2
	// DECIMAL_BASE base décimale pour conversion
	DECIMAL_BASE int = 10
	// FACTORIAL_LOOP_START début de la boucle factorielle
	FACTORIAL_LOOP_START int = 2
)

// Calculator effectue des opérations arithmétiques.
//
// Params:
//   - op: opération à effectuer (+, -, *, /)
//   - a: premier opérande
//   - b: deuxième opérande
//
// Returns:
//   - int: résultat de l'opération
//   - error: erreur si opération invalide ou division par zéro
func Calculator(op string, a, b int) (int, error) {
	// Sélection selon l'opération
	switch op {
	// Addition
	case "+":
		// Retour du résultat
		return a + b, nil
	// Soustraction
	case "-":
		// Retour du résultat
		return a - b, nil
	// Multiplication
	case "*":
		// Retour du résultat
		return a * b, nil
	// Division
	case "/":
		// Vérification division par zéro
		if b == 0 {
			// Retour erreur
			return 0, errors.New("division par zéro")
		}
		// Retour du résultat
		return a / b, nil
	// Opération inconnue
	default:
		// Retour erreur
		return 0, errors.New("opération invalide")
	}
}

// ValidateEmail vérifie si une adresse email est valide.
//
// Params:
//   - email: adresse email à valider
//
// Returns:
//   - bool: true si valide, false sinon
func ValidateEmail(email string) bool {
	// Vérification longueur minimale
	if len(email) < MIN_EMAIL_LENGTH {
		// Email trop court
		return false
	}

	// Vérification présence @
	if !strings.Contains(email, "@") {
		// Pas de @
		return false
	}

	parts := strings.Split(email, "@")
	// Vérification exactement 2 parties
	if len(parts) != EMAIL_PARTS_COUNT {
		// Format invalide
		return false
	}

	// Vérification parties non vides
	if parts[0] == "" || parts[1] == "" {
		// Parties vides
		return false
	}

	// Vérification domaine
	if !strings.Contains(parts[1], ".") {
		// Pas de point dans le domaine
		return false
	}

	// Email valide
	return true
}

// ParseInt convertit une string en entier.
//
// Params:
//   - s: string à convertir
//
// Returns:
//   - int: valeur entière
//   - error: erreur si conversion impossible
func ParseInt(s string) (int, error) {
	// Vérification string vide
	if s == "" {
		// Retour erreur
		return 0, errors.New("string vide")
	}

	result := 0
	negative := false

	// Vérification signe négatif
	if s[0] == '-' {
		negative = true
		s = s[1:]
	}

	// Parcours des caractères
	for _, c := range s {
		// Vérification caractère numérique
		if c < '0' || c > '9' {
			// Retour erreur
			return 0, errors.New("caractère non numérique")
		}
		result = result*DECIMAL_BASE + int(c-'0')
	}

	// Application du signe
	if negative {
		result = -result
	}

	// Retour du résultat
	return result, nil
}

// Factorial calcule la factorielle d'un nombre.
//
// Params:
//   - n: nombre dont calculer la factorielle
//
// Returns:
//   - int: factorielle de n
//   - error: erreur si n négatif
func Factorial(n int) (int, error) {
	// Vérification n négatif
	if n < 0 {
		// Retour erreur
		return 0, errors.New("nombre négatif")
	}

	// Cas de base
	if n == 0 || n == 1 {
		// Retour 1
		return 1, nil
	}

	result := 1
	// Calcul de la factorielle
	for i := FACTORIAL_LOOP_START; i <= n; i++ {
		result *= i
	}

	// Retour du résultat
	return result, nil
}
