package rules_func

import (
	"context"
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// FONCTIONS PARFAITES
// ════════════════════════════════════════════════════════════════════════════
// Ce fichier contient UNIQUEMENT des exemples de fonctions parfaites
// qui respectent TOUTES les règles KTN-FUNC-001 à KTN-FUNC-007
// ════════════════════════════════════════════════════════════════════════════

// ✅ FONCTION PARFAITE 1: Simple, bien nommée, bien documentée
// ProcessData traite les données fournies et retourne une erreur si échec
func ProcessData(data string) error {
	if data == "" {
		return errors.New("data vide")
	}
	return nil
}

// ✅ FONCTION PARFAITE 2: Fonction avec Context en paramètre
// FetchUserData récupère les données utilisateur depuis une source externe
// Retourne une erreur si l'utilisateur n'existe pas ou si le contexte expire
func FetchUserData(ctx context.Context, userID int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if userID <= 0 {
			return "", errors.New("userID invalide")
		}
		return "user data", nil
	}
}

// ✅ FONCTION PARFAITE 3: Nom en mixedCaps (fonction privée)
// calculateTotal est une fonction privée bien documentée
// Elle calcule le total en additionnant les valeurs
func calculateTotal(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}

// ✅ FONCTION PARFAITE 4: Avec 5 paramètres (limite OK)
// CreateUser crée un nouvel utilisateur avec les paramètres fournis
// Les paramètres ctx, name, email, age et active contrôlent la création
// Retourne l'ID du nouvel utilisateur ou une erreur si les données sont invalides
func CreateUser(ctx context.Context, name string, email string, age int, active bool) (int, error) {
	_ = ctx
	if name == "" {
		return 0, errors.New("name requis")
	}
	if email == "" {
		return 0, errors.New("email requis")
	}
	if age < 0 {
		return 0, errors.New("age invalide")
	}
	return 1, nil
}

// ✅ FONCTION PARFAITE 5: Complexité faible (< 10)
// ValidateInput valide les données d'entrée selon plusieurs critères
// Retourne une erreur si validation échoue
func ValidateInput(input string) error {
	if input == "" {
		return errors.New("input vide")
	}
	if len(input) < 3 {
		return errors.New("input trop court")
	}
	if len(input) > 100 {
		return errors.New("input trop long")
	}
	return nil
}

// ✅ FONCTION PARFAITE 6: Longueur OK (< 35 lignes)
// ProcessOrder traite une commande et effectue plusieurs validations
// Les paramètres ctx, orderID et userID contrôlent le traitement
// Retourne une erreur si la commande est invalide
func ProcessOrder(ctx context.Context, orderID int, userID int) error {
	if orderID <= 0 {
		return errors.New("orderID invalide")
	}

	if userID <= 0 {
		return errors.New("userID invalide")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Traitement OK
		return nil
	}
}

// ✅ FONCTION PARFAITE 7: Fonction avec 2 paramètres (pas besoin de tout documenter)
// ValidateEmail valide le format d'un email
func ValidateEmail(email string) bool {
	return len(email) > 0 && len(email) < 100
}

// ✅ FONCTION PARFAITE 8: Fonction avec Context en paramètre
// SendNotification envoie une notification avec timeout optionnel
// Le paramètre message contient le texte et ctx permet l'annulation
// Retourne une erreur si l'envoi échoue
func SendNotification(message string, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if message == "" {
			return errors.New("message vide")
		}
		return nil
	}
}

// ✅ FONCTION PARFAITE 9: Fonction avec un seul retour et 3 params
// FormatMessage formate un message avec les paramètres fournis
// Les paramètres template, name et value sont utilisés pour le formatage
func FormatMessage(template string, name string, value int) string {
	if template == "" {
		return ""
	}
	return template + name
}
