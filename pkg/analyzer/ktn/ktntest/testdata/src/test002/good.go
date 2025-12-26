// Package test003 provides item processing utilities.
package test003

import "errors"

const (
	// MIN_LENGTH longueur minimale
	MIN_LENGTH int = 1
	// MAX_COUNT compteur maximal
	MAX_COUNT int = 1000
)

// ProcessItem traite un item.
//
// Params:
//   - item: item à traiter
//
// Returns:
//   - string: résultat
//   - error: erreur si invalide
func ProcessItem(item string) (string, error) {
	// Vérification longueur
	if len(item) < MIN_LENGTH {
		// Retour erreur
		return "", errors.New("item too short")
	}
	// Retour résultat
	return "processed:" + item, nil
}

// CountItems compte des items.
//
// Params:
//   - items: liste d'items
//
// Returns:
//   - int: nombre d'items
//   - error: erreur si trop d'items
func CountItems(items []string) (int, error) {
	count := len(items)
	// Vérification dépassement
	if count > MAX_COUNT {
		// Retour erreur
		return 0, errors.New("too many items")
	}
	// Retour compteur
	return count, nil
}
